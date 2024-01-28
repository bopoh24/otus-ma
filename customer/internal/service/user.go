package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Nerzal/gocloak/v13"
	"github.com/bopoh24/ma_1/customer/internal/config"
	"github.com/bopoh24/ma_1/customer/internal/model"
	"github.com/bopoh24/ma_1/customer/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log/slog"
	"net/http"
	"strconv"
)

type UserService struct {
	repo repository.Repository
	conf *config.Config
}

// NewUserService returns a new UserService instance
func NewUserService(cfg *config.Config, repo repository.Repository) *UserService {

	return &UserService{
		conf: cfg,
		repo: repo,
	}
}

// Run runs the UserService
func (s *UserService) Run() error {
	mw := NewMetricsMiddleware(newMetrics())
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "ok"}`))
	})

	r.Route("/auth", func(r chi.Router) {
		r.Use(mw.Middleware)
		r.Post("/login", s.login)
		r.Post("/logout", s.logout)
		r.Post("/register", s.register)
		r.Post("/refresh", s.refresh)
	})

	r.Route("/user", func(r chi.Router) {
		r.Use(mw.Middleware)
		r.Get("/me", s.userProfile)
		r.Put("/me", s.updateUserProfile)
		r.Post("/", s.userCreate)
		r.Get("/{id}", s.userByID)
		r.Put("/{id}", s.userUpdate)
		r.Delete("/{id}", s.userDelete)
	})

	// metrics handler
	r.Handle("/metrics", promhttp.Handler())

	// Readiness and liveness probes
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	r.Get("/readyz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return http.ListenAndServe(":80", r)
}

// UserCreate creates a new user
func (s *UserService) userCreate(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(NewError(http.StatusBadRequest, err.Error()).JSON())
		return
	}
	err = s.repo.UserCreate(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(NewError(http.StatusInternalServerError, err.Error()).JSON())
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (s *UserService) getUserIdFromUrl(r *http.Request) (int64, error) {
	id := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		return 0, err
	}
	return int64(userID), nil
}

// UserByID returns a user by id
func (s *UserService) userByID(w http.ResponseWriter, r *http.Request) {
	userID, err := s.getUserIdFromUrl(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(NewError(http.StatusBadRequest, err.Error()).JSON())
		return
	}
	// get user by id
	user, err := s.repo.UserByID(userID)
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
func (s *UserService) userUpdate(w http.ResponseWriter, r *http.Request) {
	userID, err := s.getUserIdFromUrl(r)
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
	err = s.repo.UserUpdate(&user)
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

// UserDelete deletes a user
func (s *UserService) userDelete(w http.ResponseWriter, r *http.Request) {
	userID, err := s.getUserIdFromUrl(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(NewError(http.StatusBadRequest, err.Error()).JSON())
		return
	}
	err = s.repo.UserDelete(userID)
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

func (s *UserService) userProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User")
	user, err := s.repo.UserByExternalID(userID)
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
		err = s.repo.UserCreate(&user)
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

func (s *UserService) updateUserProfile(w http.ResponseWriter, r *http.Request) {
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
	err = s.repo.UserUpdateByExternalID(&user)
	if err != nil {
		if !errors.Is(err, repository.ErrUserNotFound) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(NewError(http.StatusInternalServerError, err.Error()).JSON())
			slog.Error("Error updating user: ", err)
			return
		}
		// create new user
		user.Email = r.Header.Get("X-Email")
		err = s.repo.UserCreate(&user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(NewError(http.StatusInternalServerError, err.Error()).JSON())
			slog.Error("Error creating user: ", err)
			return
		}
	}
	user, err = s.repo.UserByExternalID(user.ExternalID)
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

func (s *UserService) login(w http.ResponseWriter, r *http.Request) {
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

	client := gocloak.NewClient(s.conf.Keycloak.URL)
	token, err := client.Login(context.Background(), s.conf.Keycloak.ClientID, s.conf.Keycloak.ClientSecret,
		s.conf.Keycloak.Realm, payload.Username, payload.Password)

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

func (s *UserService) logout(w http.ResponseWriter, r *http.Request) {
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

	client := gocloak.NewClient(s.conf.Keycloak.URL)
	err = client.Logout(context.Background(), s.conf.Keycloak.ClientID, s.conf.Keycloak.ClientSecret,
		s.conf.Keycloak.Realm, payload.RefreshToken)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(NewError(http.StatusInternalServerError, err.Error()).JSON())
		slog.Error("Error logging out: ", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte{})
}

func (s *UserService) refresh(w http.ResponseWriter, r *http.Request) {
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

	client := gocloak.NewClient(s.conf.Keycloak.URL)
	token, err := client.RefreshToken(context.Background(), payload.RefreshToken, s.conf.Keycloak.ClientID,
		s.conf.Keycloak.ClientSecret, s.conf.Keycloak.Realm)

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

func (s *UserService) register(w http.ResponseWriter, r *http.Request) {
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
	client := gocloak.NewClient(s.conf.Keycloak.URL)
	token, err := client.LoginAdmin(context.Background(),
		s.conf.Keycloak.Admin, s.conf.Keycloak.Password, "master")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(NewError(http.StatusInternalServerError, err.Error()).JSON())
		slog.Error("Error logging in admin: ", err)
		return
	}
	_, err = client.CreateUser(context.Background(), token.AccessToken, s.conf.Keycloak.Realm, gocloak.User{
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
	token, err = client.Login(context.Background(), s.conf.Keycloak.ClientID, s.conf.Keycloak.ClientSecret, s.conf.Keycloak.Realm,
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
