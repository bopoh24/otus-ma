package service

import (
	"encoding/json"
	"errors"
	"github.com/bopoh24/ma_1/internal/config"
	"github.com/bopoh24/ma_1/internal/model"
	"github.com/bopoh24/ma_1/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	r.Route("/user", func(r chi.Router) {
		r.Use(mw.Middleware)
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

	return http.ListenAndServe(":8000", r)
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
