package app

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Nerzal/gocloak/v13"
	"github.com/bopoh24/ma_1/app/internal/model"
	"github.com/bopoh24/ma_1/app/internal/repository"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"strconv"
)

func (a *App) login(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(NewError(http.StatusBadRequest, err.Error()).JSON())
		return
	}

	client := gocloak.NewClient(a.cfg.Keycloak.URL)
	token, err := client.Login(context.Background(), a.cfg.Keycloak.ClientID, a.cfg.Keycloak.ClientSecret,
		a.cfg.Keycloak.Realm, payload.Username, payload.Password)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(NewError(http.StatusInternalServerError, err.Error()).JSON())
		return
	}
	// write token to response body
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(token); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(NewError(http.StatusInternalServerError, err.Error()).JSON())
		return
	}
}

func (a *App) logout(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		RefreshToken string `json:"refreshToken"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(NewError(http.StatusBadRequest, err.Error()).JSON())
		slog.Error("Error decoding payload: ", err)
		return
	}

	client := gocloak.NewClient(a.cfg.Keycloak.URL)
	err = client.Logout(context.Background(), a.cfg.Keycloak.ClientID, a.cfg.Keycloak.ClientSecret,
		a.cfg.Keycloak.Realm, payload.RefreshToken)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(NewError(http.StatusInternalServerError, err.Error()).JSON())
		slog.Error("Error logging out: ", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte{})
}

func (a *App) refresh(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		RefreshToken string `json:"refreshToken"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(NewError(http.StatusBadRequest, err.Error()).JSON())
		slog.Error("Error decoding payload: ", err)
		return
	}

	client := gocloak.NewClient(a.cfg.Keycloak.URL)
	token, err := client.RefreshToken(context.Background(), payload.RefreshToken, a.cfg.Keycloak.ClientID,
		a.cfg.Keycloak.ClientSecret, a.cfg.Keycloak.Realm)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(NewError(http.StatusInternalServerError, err.Error()).JSON())
		slog.Error("Error refreshing token: ", err)
		return
	}
	// write token to response body
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(token); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(NewError(http.StatusInternalServerError, err.Error()).JSON())
		slog.Error("Error encoding token: ", err)
		return
	}
}

func (a *App) register(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(NewError(http.StatusBadRequest, err.Error()).JSON())
		slog.Error("Error decoding payload: ", err)
		return
	}
	client := gocloak.NewClient(a.cfg.Keycloak.URL)
	token, err := client.LoginAdmin(context.Background(),
		a.cfg.Keycloak.Admin, a.cfg.Keycloak.Password, "master")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(NewError(http.StatusInternalServerError, err.Error()).JSON())
		slog.Error("Error logging in admin: ", err)
		return
	}
	_, err = client.CreateUser(context.Background(), token.AccessToken, a.cfg.Keycloak.Realm, gocloak.User{
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
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(NewError(http.StatusInternalServerError, err.Error()).JSON())
		slog.Error("Error creating user: ", err)
		return
	}
	token, err = client.Login(context.Background(), a.cfg.Keycloak.ClientID, a.cfg.Keycloak.ClientSecret, a.cfg.Keycloak.Realm,
		payload.Email, payload.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(NewError(http.StatusInternalServerError, err.Error()).JSON())
		slog.Error("Error logging in as new user: ", err)
		return
	}
	// return token
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(token); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(NewError(http.StatusInternalServerError, err.Error()).JSON())
		slog.Error("Error encoding token: ", err)
		return
	}
}

// UserCreate creates a new user
func (a *App) userCreate(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(NewError(http.StatusBadRequest, err.Error()).JSON())
		return
	}
	if err = a.userSrv.Create(&user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(NewError(http.StatusInternalServerError, err.Error()).JSON())
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (a *App) getUserIdFromUrl(r *http.Request) (int64, error) {
	id := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		return 0, err
	}
	return int64(userID), nil
}

func (a *App) userByID(w http.ResponseWriter, r *http.Request) {
	userID, err := a.getUserIdFromUrl(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(NewError(http.StatusBadRequest, err.Error()).JSON())
		return
	}
	// get user by id
	user, err := a.userSrv.ByID(userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			w.WriteHeader(http.StatusNotFound)
			w.Write(NewError(http.StatusNotFound, err.Error()).JSON())
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(NewError(http.StatusInternalServerError, err.Error()).JSON())
		return
	}
	// write user to response body
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(NewError(http.StatusInternalServerError, err.Error()).JSON())
		return
	}
}

// UserUpdate updates a user
func (a *App) userUpdate(w http.ResponseWriter, r *http.Request) {
	userID, err := a.getUserIdFromUrl(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(NewError(http.StatusBadRequest, err.Error()).JSON())
		return
	}
	var user model.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(NewError(http.StatusBadRequest, err.Error()).JSON())
		return
	}
	user.ID = userID
	err = a.userSrv.Update(&user)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			w.WriteHeader(http.StatusNotFound)
			w.Write(NewError(http.StatusNotFound, err.Error()).JSON())
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(NewError(http.StatusInternalServerError, err.Error()).JSON())
		return
	}
}

func (a *App) userDelete(w http.ResponseWriter, r *http.Request) {
	userID, err := a.getUserIdFromUrl(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(NewError(http.StatusBadRequest, err.Error()).JSON())
		return
	}
	err = a.userSrv.Delete(userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			w.WriteHeader(http.StatusNotFound)
			w.Write(NewError(http.StatusNotFound, err.Error()).JSON())
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(NewError(http.StatusInternalServerError, err.Error()).JSON())
		return
	}
}

func (a *App) userProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User")
	user, err := a.userSrv.ByExternalID(userID)
	if err != nil {
		if !errors.Is(err, repository.ErrUserNotFound) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(NewError(http.StatusInternalServerError, err.Error()).JSON())
			slog.Error("Error getting user by external id: ", err)
			return
		}
		// create new user
		user = model.User{
			ExternalID: userID,
			Email:      r.Header.Get("X-Email"),
			Username:   r.Header.Get("X-Preferred-Username"),
			FirstName:  r.Header.Get("X-Given-Name"),
			LastName:   r.Header.Get("X-Family-Name"),
		}
		err = a.userSrv.Create(&user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(NewError(http.StatusInternalServerError, err.Error()).JSON())
			slog.Error("Error creating user: ", err)
			return
		}
	}
	// write user to response body
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(NewError(http.StatusInternalServerError, err.Error()).JSON())
		slog.Error("Error encoding user: ", err)
		return
	}
}

func (a *App) updateUserProfile(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		Username    string `json:"username"`
		FirstName   string `json:"firstName"`
		LastName    string `json:"lastName"`
		Phone       string `json:"phone"`
		Description string `json:"description"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(NewError(http.StatusBadRequest, err.Error()).JSON())
		return
	}
	user := model.User{
		ExternalID:  r.Header.Get("X-User"),
		Username:    payload.Username,
		FirstName:   payload.FirstName,
		LastName:    payload.LastName,
		Phone:       payload.Phone,
		Description: payload.Description,
	}
	// update user
	err = a.userSrv.UpdateByExternalID(&user)
	if err != nil {
		if !errors.Is(err, repository.ErrUserNotFound) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(NewError(http.StatusInternalServerError, err.Error()).JSON())
			slog.Error("Error updating user: ", err)
			return
		}
		// create new user
		user.Email = r.Header.Get("X-Email")
		err = a.userSrv.Create(&user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(NewError(http.StatusInternalServerError, err.Error()).JSON())
			slog.Error("Error creating user: ", err)
			return
		}
	}
	user, err = a.userSrv.ByExternalID(user.ExternalID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(NewError(http.StatusInternalServerError, err.Error()).JSON())
		slog.Error("Error getting user by external id: ", err)
		return
	}
	// write user to response body
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(NewError(http.StatusInternalServerError, err.Error()).JSON())
		slog.Error("Error encoding user: ", err)
		return
	}
}

func (a *App) orderCreate(w http.ResponseWriter, r *http.Request) {
	var order model.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(NewError(http.StatusBadRequest, err.Error()).JSON())
		return
	}

	idempotencyKey := r.Header.Get("X-Idempotency-Key")
	if idempotencyKey == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(NewError(http.StatusBadRequest, "idempotency key is empty").JSON())
		return
	}

	order.ID, err = a.orderSrv.Create(&order, idempotencyKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(NewError(http.StatusInternalServerError, err.Error()).JSON())
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"order_id":` + strconv.FormatInt(order.ID, 10) + `}`))
}

func (a *App) getOrderIdFromUrl(r *http.Request) (int64, error) {
	id := chi.URLParam(r, "id")
	orderID, err := strconv.Atoi(id)
	if err != nil {
		return 0, err
	}
	return int64(orderID), nil
}

func (a *App) orderByID(w http.ResponseWriter, r *http.Request) {
	orderID, err := a.getOrderIdFromUrl(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(NewError(http.StatusBadRequest, err.Error()).JSON())
		return
	}
	// get order by id
	order, err := a.orderSrv.ByID(orderID)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			w.WriteHeader(http.StatusNotFound)
			w.Write(NewError(http.StatusNotFound, err.Error()).JSON())
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(NewError(http.StatusInternalServerError, err.Error()).JSON())
		return
	}
	// write order to response body
	err = json.NewEncoder(w).Encode(order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(NewError(http.StatusInternalServerError, err.Error()).JSON())
		return
	}
}

func (a *App) orderByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := a.getUserIdFromUrl(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(NewError(http.StatusBadRequest, err.Error()).JSON())
		return
	}
	// get orders by user id
	orders, err := a.orderSrv.ByUserID(userID)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			w.WriteHeader(http.StatusNotFound)
			w.Write(NewError(http.StatusNotFound, err.Error()).JSON())
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(NewError(http.StatusInternalServerError, err.Error()).JSON())
		return
	}
	// write orders to response body
	err = json.NewEncoder(w).Encode(orders)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(NewError(http.StatusInternalServerError, err.Error()).JSON())
		return
	}
}

func (a *App) orderUpdate(w http.ResponseWriter, r *http.Request) {
	orderID, err := a.getOrderIdFromUrl(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(NewError(http.StatusBadRequest, err.Error()).JSON())
		return
	}
	var order model.Order
	err = json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(NewError(http.StatusBadRequest, err.Error()).JSON())
		return
	}
	order.ID = orderID
	err = a.orderSrv.Update(&order)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			w.WriteHeader(http.StatusNotFound)
			w.Write(NewError(http.StatusNotFound, err.Error()).JSON())
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(NewError(http.StatusInternalServerError, err.Error()).JSON())
		return
	}
}

func (a *App) orderDelete(w http.ResponseWriter, r *http.Request) {
	orderID, err := a.getOrderIdFromUrl(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(NewError(http.StatusBadRequest, err.Error()).JSON())
		return
	}
	err = a.orderSrv.Delete(orderID)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			w.WriteHeader(http.StatusNotFound)
			w.Write(NewError(http.StatusNotFound, err.Error()).JSON())
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(NewError(http.StatusInternalServerError, err.Error()).JSON())
		return
	}
}
