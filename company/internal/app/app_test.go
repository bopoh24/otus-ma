package app

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/bopoh24/ma_1/company/internal/model"
	"github.com/bopoh24/ma_1/company/internal/repository"
	"github.com/bopoh24/ma_1/company/internal/service"
	mock "github.com/bopoh24/ma_1/company/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func createMockApp(repo service.Repository) *App {
	return &App{
		service: service.New(repo),
	}
}

func TestHandlerCompanyDetails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := httptest.NewRequest("GET", "/company/{id}", nil)
	rctx := chi.NewRouteContext()
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	rctx.URLParams.Add("id", "2")

	t.Run("Company not found", func(t *testing.T) {
		w := httptest.NewRecorder()

		repo := mock.NewMockRepository(ctrl)
		a := createMockApp(repo)

		repo.EXPECT().CompanyByID(gomock.Any(), int64(2)).Return(model.Company{}, repository.ErrCompanyNotFound).Times(1)
		a.handlerCompanyDetails(w, r)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("OK", func(t *testing.T) {
		w := httptest.NewRecorder()
		repo := mock.NewMockRepository(ctrl)
		a := createMockApp(repo)
		now := time.Now()
		repo.EXPECT().CompanyByID(gomock.Any(), int64(2)).Return(model.Company{
			ID:      2,
			Name:    "Company",
			Address: "Come city, some street",
			Email:   "some@email.com",
			Created: &now,
			Updated: &now,
		}, nil).Times(1)
		a.handlerCompanyDetails(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestHandlerUpdateCompany(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Bad request", func(t *testing.T) {
		r := httptest.NewRequest("PUT", "/company/{id}", nil)
		w := httptest.NewRecorder()
		rctx := chi.NewRouteContext()
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		rctx.URLParams.Add("id", "2")

		repo := mock.NewMockRepository(ctrl)
		a := createMockApp(repo)
		a.handlerUpdateCompany(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Not found", func(t *testing.T) {
		body := model.Company{
			ID:   22,
			Name: "Company",
		}
		reqBody, err := json.Marshal(body)
		assert.NoError(t, err)
		r := httptest.NewRequest("PUT", "/company/{id}", bytes.NewReader(reqBody))
		w := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		rctx.URLParams.Add("id", "2")

		repo := mock.NewMockRepository(ctrl)
		a := createMockApp(repo)
		repo.EXPECT().CompanyUpdate(gomock.Any(), gomock.Any()).Return(repository.ErrCompanyNotFound).Times(1)
		a.handlerUpdateCompany(w, r)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("OK", func(t *testing.T) {
		body := model.Company{
			ID:   22,
			Name: "Company",
		}
		reqBody, err := json.Marshal(body)
		assert.NoError(t, err)
		r := httptest.NewRequest("PUT", "/company/{id}", bytes.NewReader(reqBody))
		w := httptest.NewRecorder()
		rctx := chi.NewRouteContext()
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		rctx.URLParams.Add("id", "2")

		repo := mock.NewMockRepository(ctrl)
		a := createMockApp(repo)
		repo.EXPECT().CompanyUpdate(gomock.Any(), gomock.Any()).Return(nil).Times(1)
		a.handlerUpdateCompany(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestHandlerUpdateLogo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Bad request", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/company/{id}/logo", nil)
		w := httptest.NewRecorder()
		rctx := chi.NewRouteContext()
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		rctx.URLParams.Add("id", "2")

		repo := mock.NewMockRepository(ctrl)
		a := createMockApp(repo)
		a.handlerUpdateLogo(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Not found", func(t *testing.T) {
		body := model.Company{
			ID:   22,
			Logo: "logo",
		}
		reqBody, err := json.Marshal(body)
		assert.NoError(t, err)
		r := httptest.NewRequest("POST", "/company/{id}/logo", bytes.NewReader(reqBody))
		w := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		rctx.URLParams.Add("id", "2")

		repo := mock.NewMockRepository(ctrl)
		a := createMockApp(repo)
		repo.EXPECT().CompanyUpdateLogo(gomock.Any(), gomock.Any(), gomock.Any()).Return(repository.ErrCompanyNotFound).Times(1)
		a.handlerUpdateLogo(w, r)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("OK", func(t *testing.T) {
		body := model.Company{
			ID:   22,
			Logo: "logo",
		}
		reqBody, err := json.Marshal(body)
		assert.NoError(t, err)
		r := httptest.NewRequest("POST", "/company/{id}/logo", bytes.NewReader(reqBody))
		w := httptest.NewRecorder()
		rctx := chi.NewRouteContext()
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		rctx.URLParams.Add("id", "2")

		repo := mock.NewMockRepository(ctrl)
		a := createMockApp(repo)
		repo.EXPECT().CompanyUpdateLogo(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
		a.handlerUpdateLogo(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestHandlerUpdateLocation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Bad request", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/company/{id}/location", nil)
		w := httptest.NewRecorder()
		rctx := chi.NewRouteContext()
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		rctx.URLParams.Add("id", "2")

		repo := mock.NewMockRepository(ctrl)
		a := createMockApp(repo)
		a.handlerUpdateLocation(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Not found", func(t *testing.T) {
		body := struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		}{
			Lat: 22.22,
			Lng: 33.33,
		}
		reqBody, err := json.Marshal(body)
		assert.NoError(t, err)
		r := httptest.NewRequest("POST", "/company/{id}/location", bytes.NewReader(reqBody))
		w := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		rctx.URLParams.Add("id", "2")

		repo := mock.NewMockRepository(ctrl)
		a := createMockApp(repo)
		repo.EXPECT().CompanyUpdateLocation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(repository.ErrCompanyNotFound).Times(1)
		a.handlerUpdateLocation(w, r)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("OK", func(t *testing.T) {
		body := struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		}{
			Lat: 22.22,
			Lng: 33.33,
		}
		reqBody, err := json.Marshal(body)
		assert.NoError(t, err)
		r := httptest.NewRequest("POST", "/company/{id}/location", bytes.NewReader(reqBody))
		w := httptest.NewRecorder()
		rctx := chi.NewRouteContext()
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		rctx.URLParams.Add("id", "2")

		repo := mock.NewMockRepository(ctrl)
		a := createMockApp(repo)
		repo.EXPECT().CompanyUpdateLocation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
		a.handlerUpdateLocation(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestHandlerActivateDeactivate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Bad request", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/company/{id}/activate", nil)
		w := httptest.NewRecorder()
		repo := mock.NewMockRepository(ctrl)
		repo.EXPECT().CompanyActivateDeactivate(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(0)
		a := createMockApp(repo)
		a.handlerActivateDeactivate(true)(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Not found", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/company/{id}/activate", nil)
		w := httptest.NewRecorder()
		rctx := chi.NewRouteContext()
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		rctx.URLParams.Add("id", "2")

		repo := mock.NewMockRepository(ctrl)
		a := createMockApp(repo)
		repo.EXPECT().CompanyActivateDeactivate(gomock.Any(), gomock.Any(), gomock.Any()).Return(repository.ErrCompanyNotFound).Times(1)
		a.handlerActivateDeactivate(true)(w, r)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("OK", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/company/{id}/activate", nil)
		w := httptest.NewRecorder()
		rctx := chi.NewRouteContext()
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		rctx.URLParams.Add("id", "2")

		repo := mock.NewMockRepository(ctrl)
		a := createMockApp(repo)
		repo.EXPECT().CompanyActivateDeactivate(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
		a.handlerActivateDeactivate(true)(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestHandlerGetManagers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Not found", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/company/{id}/managers", nil)
		w := httptest.NewRecorder()
		rctx := chi.NewRouteContext()
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		rctx.URLParams.Add("id", "2")

		repo := mock.NewMockRepository(ctrl)
		a := createMockApp(repo)
		repo.EXPECT().CompanyManagers(gomock.Any(), gomock.Any()).Return(nil, repository.ErrCompanyNotFound).Times(1)
		a.handlerGetManagers(w, r)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("OK", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/company/{id}/managers", nil)
		w := httptest.NewRecorder()
		rctx := chi.NewRouteContext()
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		rctx.URLParams.Add("id", "2")

		repo := mock.NewMockRepository(ctrl)
		a := createMockApp(repo)
		repo.EXPECT().CompanyManagers(gomock.Any(), gomock.Any()).Return([]model.Manager{
			{
				ID:        1,
				CompanyID: 2,
				UserID:    "231232332",
				Email:     "test@mail.com",
				Role:      model.MangerRoleAdmin,
				Active:    true,
			},
		}, nil).Times(1)
		a.handlerGetManagers(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "test@mail.com")
		assert.Contains(t, w.Body.String(), model.MangerRoleAdmin)
	})
}

func TestHandlerCreateCompany(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Bad request", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/company", nil)
		w := httptest.NewRecorder()
		repo := mock.NewMockRepository(ctrl)
		a := createMockApp(repo)
		a.handlerCreateCompany(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("OK", func(t *testing.T) {
		body := model.Company{
			Name: "Company",
		}
		reqBody, err := json.Marshal(body)
		assert.NoError(t, err)
		r := httptest.NewRequest("POST", "/company", bytes.NewReader(reqBody))
		r.Header.Set("X-User", "123")
		r.Header.Set("X-Email", "mail@mail.com")
		r.Header.Set("X-Given-Name", "John")
		r.Header.Set("X-Family-Name", "Doe")

		w := httptest.NewRecorder()

		repo := mock.NewMockRepository(ctrl)
		a := createMockApp(repo)
		repo.EXPECT().CompanyCreate(gomock.Any(), "123", "mail@mail.com", "John", "Doe", body).Return(nil).Times(1)
		a.handlerCreateCompany(w, r)
		assert.Equal(t, http.StatusCreated, w.Code)
	})
}
