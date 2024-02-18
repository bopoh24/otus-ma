package app

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Nerzal/gocloak/v13"
	"github.com/bopoh24/ma_1/customer/internal/model"
	"github.com/bopoh24/ma_1/customer/internal/repository"
	"github.com/bopoh24/ma_1/pkg/http/helper"
	"github.com/bopoh24/ma_1/pkg/verifier/phone"
	"net/http"
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
	helper.JSONResponse(w, http.StatusOK, token)
}

func (a *App) hanlderLogout(w http.ResponseWriter, r *http.Request) {
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
	helper.JSONResponse(w, http.StatusOK, token)
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
		Email:         gocloak.StringP(payload.Email),
		EmailVerified: gocloak.BoolP(true),
		FirstName:     gocloak.StringP(payload.FirstName),
		LastName:      gocloak.StringP(payload.LastName),
		Enabled:       gocloak.BoolP(true),
		Username:      gocloak.StringP(payload.Email),
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
	helper.JSONResponse(w, http.StatusOK, token)
}

// byID returns a user by id
func (a *App) handlerCustomerByID(w http.ResponseWriter, r *http.Request) {
	claims, err := helper.ExtractClaims(r)
	if err != nil {
		helper.ErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	// get user by id
	customer, err := a.service.CustomerByID(r.Context(), claims.Id)
	if err != nil {
		if errors.Is(err, repository.ErrCustomerNotFound) {
			helper.ErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		a.log.Error("Get customer by id", "err", err)
		return
	}

	// write user to response body
	err = json.NewEncoder(w).Encode(customer)
	if err != nil {
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (a *App) handlerProfile(w http.ResponseWriter, r *http.Request) {
	claims, err := helper.ExtractClaims(r)
	if err != nil {
		helper.ErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	customer, err := a.service.CustomerByID(r.Context(), claims.Id)
	if err != nil {
		if !errors.Is(err, repository.ErrCustomerNotFound) {
			helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
			a.log.Error("Get customer by id", "err", err)
			return
		}
		// create new user
		customer = model.Customer{
			ID:        claims.Id,
			Email:     claims.Email,
			FirstName: claims.FirstName,
			LastName:  claims.LastName,
		}

		err = a.service.CreateCustomerProfile(r.Context(), customer)
		if err != nil {
			helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
			a.log.Error("Error creating customer", "err", err)
			return
		}
	}
	// write user to response body
	helper.JSONResponse(w, http.StatusOK, customer)
}

func (a *App) handlerProfileUpdate(w http.ResponseWriter, r *http.Request) {
	claims, err := helper.ExtractClaims(r)
	if err != nil {
		helper.ErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	var payload model.Customer
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	payload.ID = claims.Id
	// update customer profile
	err = a.service.UpdateCustomerProfile(r.Context(), payload)
	if err != nil {
		if !errors.Is(err, repository.ErrCustomerNotFound) {
			helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
			a.log.Error("Error updating customer profile", "err", err)
			return
		}
		// create new profile
		payload.Email = claims.Email
		err = a.service.CreateCustomerProfile(r.Context(), payload)
		if err != nil {
			helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
			a.log.Error("Error creating customer profile", "err", err)
			return
		}
	}

	customer, err := a.service.CustomerByID(r.Context(), payload.ID)
	if err != nil {
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		a.log.Error("Error getting user by external id", "err", err)
		return
	}
	// write user to response body
	helper.JSONResponse(w, http.StatusOK, customer)
}

func (a *App) handlerRequestPhoneVerification(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Phone string `json:"phone"`
	}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if payload.Phone == "" {
		helper.ErrorResponse(w, http.StatusBadRequest, "phone is required")
		return
	}

	err = a.service.RequestPhoneVerification(r.Context(), payload.Phone)
	if err != nil {
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		a.log.Error("Error requesting phone verification", "err", err)
		return
	}
	helper.JSONResponse(w, http.StatusOK, map[string]string{"result": "verification code sent"})
}

func (a *App) handlerVerifyPhone(w http.ResponseWriter, r *http.Request) {
	claims, err := helper.ExtractClaims(r)
	if err != nil {
		helper.ErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	var payload struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if payload.Phone == "" || payload.Code == "" {
		helper.ErrorResponse(w, http.StatusBadRequest, "phone and code are required")
		return
	}
	err = a.service.VerifyPhone(r.Context(), claims.Id, payload.Phone, payload.Code)
	if err != nil {
		if errors.Is(err, phone.ErrIncorrectVerificationCode) {
			helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		a.log.Error("Error verifying phone", "err", err)
		return
	}
	helper.JSONResponse(w, http.StatusOK, map[string]string{"result": "phone verified"})
}
