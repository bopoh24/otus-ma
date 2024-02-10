package app

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Nerzal/gocloak/v13"
	"github.com/bopoh24/ma_1/company/internal/model"
	"github.com/bopoh24/ma_1/company/internal/repository"
	"github.com/bopoh24/ma_1/pkg/http/helper"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (a *App) handlerLogin(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := a.keycloakClient.Login(context.Background(), a.conf.Keycloak.ClientID, a.conf.Keycloak.ClientSecret,
		a.conf.Keycloak.Realm, payload.Email, payload.Password)

	if err != nil {
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	// write token to response body
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(token); err != nil {
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (a *App) handlerLogout(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		RefreshToken string `json:"refreshToken"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = a.keycloakClient.Logout(context.Background(), a.conf.Keycloak.ClientID, a.conf.Keycloak.ClientSecret,
		a.conf.Keycloak.Realm, payload.RefreshToken)

	if err != nil {
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *App) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		RefreshToken string `json:"refreshToken"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := a.keycloakClient.RefreshToken(context.Background(), payload.RefreshToken, a.conf.Keycloak.ClientID,
		a.conf.Keycloak.ClientSecret, a.conf.Keycloak.Realm)

	if err != nil {
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	// write token to response body
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(token); err != nil {
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (a *App) handlerRegister(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		a.log.Error("Decoding register payload", "err", err)
		return
	}
	token, err := a.keycloakClient.LoginAdmin(context.Background(),
		a.conf.Keycloak.Admin, a.conf.Keycloak.Password, "master")
	if err != nil {
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		a.log.Error("Logging in as admin", "err", err)
		return
	}
	_, err = a.keycloakClient.CreateUser(context.Background(), token.AccessToken, a.conf.Keycloak.Realm, gocloak.User{
		Email:     gocloak.StringP(payload.Email),
		FirstName: gocloak.StringP(payload.FirstName),
		LastName:  gocloak.StringP(payload.LastName),
		Enabled:   gocloak.BoolP(true),
		Username:  gocloak.StringP(payload.Email),
		Credentials: &[]gocloak.CredentialRepresentation{
			{
				Temporary: gocloak.BoolP(false),
				Type:      gocloak.StringP("password"),
				Value:     gocloak.StringP(payload.Password),
			},
		},
	})
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		a.log.Error("Creating user", "err", err)
		return
	}
	token, err = a.keycloakClient.Login(context.Background(), a.conf.Keycloak.ClientID, a.conf.Keycloak.ClientSecret, a.conf.Keycloak.Realm,
		payload.Email, payload.Password)
	if err != nil {
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		a.log.Error("Logging in as new user", "err", err)
		return
	}
	// return token
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(token); err != nil {
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		a.log.Error("Encoding token", "err", err)
		return
	}
}

func (a *App) handlerCompanyDetails(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	company, err := a.service.CompanyByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrCompanyNotFound) {
			helper.ErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	// return company
	if err := json.NewEncoder(w).Encode(company); err != nil {
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (a *App) handlerUpdateCompany(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	var company model.Company
	err = json.NewDecoder(r.Body).Decode(&company)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	company.ID = id
	err = a.service.UpdateCompany(r.Context(), company)
	if err != nil {
		if errors.Is(err, repository.ErrCompanyNotFound) {
			helper.ErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *App) handlerUpdateLogo(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	payload := struct {
		Logo string `json:"logo"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	err = a.service.UpdateCompanyLogo(r.Context(), id, payload.Logo)
	if err != nil {
		if errors.Is(err, repository.ErrCompanyNotFound) {
			helper.ErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *App) handlerUpdateLocation(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	payload := struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	err = a.service.UpdateCompanyLocation(r.Context(), id, payload.Lat, payload.Lng)
	if err != nil {
		if errors.Is(err, repository.ErrCompanyNotFound) {
			helper.ErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *App) handlerActivateDeactivate(active bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = a.service.ActivateDeactivateCompany(r.Context(), id, active)
		if err != nil {
			if errors.Is(err, repository.ErrCompanyNotFound) {
				helper.ErrorResponse(w, http.StatusNotFound, err.Error())
				return
			}
			helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (a *App) handlerCreateCompany(w http.ResponseWriter, r *http.Request) {

	userID := r.Header.Get("X-User")
	email := r.Header.Get("X-Email")
	firstName := r.Header.Get("X-Given-Name")
	lastName := r.Header.Get("X-Family-Name")

	var company model.Company
	err := json.NewDecoder(r.Body).Decode(&company)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	err = a.service.CreateCompany(r.Context(), userID, email, firstName, lastName, company)
	if err != nil {
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (a *App) handlerGetManagers(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	managers, err := a.service.CompanyManagers(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrCompanyNotFound) {
			helper.ErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	// return managers
	if err := json.NewEncoder(w).Encode(managers); err != nil {
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}
