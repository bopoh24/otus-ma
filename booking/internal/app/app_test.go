package app

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/bopoh24/ma_1/booking/internal/repository"
	"github.com/bopoh24/ma_1/booking/internal/service"
	mock "github.com/bopoh24/ma_1/booking/mocks"
	"github.com/bopoh24/ma_1/booking/pkg/model"
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

func TestHandlerGetServices(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock.NewMockRepository(ctrl)
	a := createMockApp(repo)
	r := httptest.NewRequest("GET", "/booking/services", nil)
	w := httptest.NewRecorder()

	repo.EXPECT().Services(gomock.Any()).Return([]model.Service{
		{
			ID:       1,
			ParentID: 0,
			Name:     "test",
		},
		{
			ID:       2,
			ParentID: 1,
			Name:     "SubTest1",
		},
	}, nil)
	a.handlerGetServices(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "SubTest1")
}

func TestHandlerAddService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock.NewMockRepository(ctrl)
	a := createMockApp(repo)

	t.Run("bad request", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/booking/services", nil)
		w := httptest.NewRecorder()

		a.handlerAddService(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("root service", func(t *testing.T) {
		svc := model.Service{
			Name: "My service",
		}
		reqBody, err := json.Marshal(svc)
		assert.NoError(t, err)

		r := httptest.NewRequest("POST", "/booking/services", bytes.NewReader(reqBody))
		w := httptest.NewRecorder()

		repo.EXPECT().ServiceAdd(gomock.Any(), svc).Return(nil).Times(1)
		repo.EXPECT().Services(gomock.Any()).Return([]model.Service{}, nil).Times(1)
		a.handlerAddService(w, r)
		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("sub service", func(t *testing.T) {
		svc := model.Service{
			ParentID: 1,
			Name:     "My sub service",
		}
		reqBody, err := json.Marshal(svc)
		assert.NoError(t, err)

		r := httptest.NewRequest("POST", "/booking/services", bytes.NewReader(reqBody))
		w := httptest.NewRecorder()

		repo.EXPECT().ServiceAdd(gomock.Any(), svc).Return(nil).Times(1)
		repo.EXPECT().Services(gomock.Any()).Return([]model.Service{}, nil).Times(1)
		a.handlerAddService(w, r)
		assert.Equal(t, http.StatusCreated, w.Code)
	})
}

func TestHandlerAddOffer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock.NewMockRepository(ctrl)
	a := createMockApp(repo)

	t.Run("bad request", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/booking/offers", nil)
		w := httptest.NewRecorder()

		a.handlerAddOffer(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("bad parameters", func(t *testing.T) {
		offer := model.Offer{
			ServiceID: 1,
			CompanyID: 1,
			Price:     100,
		}
		reqBody, err := json.Marshal(offer)
		assert.NoError(t, err)

		r := httptest.NewRequest("POST", "/booking/offers", bytes.NewReader(reqBody))
		w := httptest.NewRecorder()
		a.handlerAddOffer(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("good request", func(t *testing.T) {
		offer := model.Offer{
			ServiceID:   1,
			CompanyID:   1,
			CompanyName: "My company",
			Price:       100,
			Datetime:    time.Now().Truncate(time.Second),
		}
		reqBody, err := json.Marshal(offer)
		assert.NoError(t, err)

		r := httptest.NewRequest("POST", "/booking/offers", bytes.NewReader(reqBody))
		w := httptest.NewRecorder()

		repo.EXPECT().OfferAdd(gomock.Any(), offer).Return(nil).Times(1)
		a.handlerAddOffer(w, r)
		assert.Equal(t, http.StatusCreated, w.Code)
	})
}

func TestHandlerDeleteOffer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock.NewMockRepository(ctrl)
	a := createMockApp(repo)

	payload := struct {
		CompanyId int64 `json:"company_id"`
	}{CompanyId: 1}

	reqBody, err := json.Marshal(payload)
	assert.NoError(t, err)

	t.Run("empty body", func(t *testing.T) {
		r := httptest.NewRequest("DELETE", "/booking/offers/1", nil)
		rctx := chi.NewRouteContext()
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		rctx.URLParams.Add("id", "1")
		w := httptest.NewRecorder()
		a.handlerDeleteOffer(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		r := httptest.NewRequest("DELETE", "/booking/offers/1", bytes.NewReader(reqBody))
		rctx := chi.NewRouteContext()
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		rctx.URLParams.Add("id", "1")
		w := httptest.NewRecorder()
		repo.EXPECT().OfferDelete(gomock.Any(), int64(1), int64(1)).Return(repository.ErrOfferNotFound).Times(1)
		a.handlerDeleteOffer(w, r)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("good request", func(t *testing.T) {
		r := httptest.NewRequest("DELETE", "/booking/offers/1", bytes.NewReader(reqBody))
		rctx := chi.NewRouteContext()
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		rctx.URLParams.Add("id", "1")
		w := httptest.NewRecorder()
		repo.EXPECT().OfferDelete(gomock.Any(), int64(1), int64(1)).Return(nil).Times(1)
		a.handlerDeleteOffer(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestHandlerChangeOfferStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock.NewMockRepository(ctrl)
	a := createMockApp(repo)

	payload := struct {
		Status model.OfferStatus `json:"status"`
	}{Status: model.OfferStatus("test")}

	reqBody, err := json.Marshal(payload)
	assert.NoError(t, err)

	t.Run("empty body", func(t *testing.T) {
		r := httptest.NewRequest("PUT", "/booking/offers/1/status", nil)
		rctx := chi.NewRouteContext()
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		rctx.URLParams.Add("id", "1")
		w := httptest.NewRecorder()
		a.handlerChangeOfferStatus(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		r := httptest.NewRequest("PUT", "/booking/offers/1/status", bytes.NewReader(reqBody))
		rctx := chi.NewRouteContext()
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		rctx.URLParams.Add("id", "1")
		w := httptest.NewRecorder()
		repo.EXPECT().OfferChangeStatus(gomock.Any(), int64(1), model.OfferStatus("test")).Return(repository.ErrOfferNotFound).Times(1)
		a.handlerChangeOfferStatus(w, r)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("good request", func(t *testing.T) {
		r := httptest.NewRequest("PUT", "/booking/offers/1/status", bytes.NewReader(reqBody))
		rctx := chi.NewRouteContext()
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		rctx.URLParams.Add("id", "1")
		w := httptest.NewRecorder()
		repo.EXPECT().OfferChangeStatus(gomock.Any(), int64(1), model.OfferStatus("test")).Return(nil).Times(1)
		a.handlerChangeOfferStatus(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestHandlerCancelOfferByCompany(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock.NewMockRepository(ctrl)
	a := createMockApp(repo)

	payload := struct {
		CompanyId int64  `json:"company_id"`
		Reason    string `json:"reason"`
	}{Reason: "test"}

	reqBody, err := json.Marshal(payload)
	assert.NoError(t, err)

	t.Run("empty body", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/booking/offers/1/cancel", nil)
		rctx := chi.NewRouteContext()
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		rctx.URLParams.Add("id", "1")
		r.Header.Set("X-User", "user_id1xxx")
		r.Header.Set("X-Email", "some_email@main.com")
		w := httptest.NewRecorder()
		a.handlerCancelOfferByCompany(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/booking/offers/1/cancel", bytes.NewReader(reqBody))
		rctx := chi.NewRouteContext()
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		rctx.URLParams.Add("id", "1")
		r.Header.Set("X-User", "user_id1xxx")
		r.Header.Set("X-Email", "some_email@main.com")
		w := httptest.NewRecorder()
		repo.EXPECT().OfferCancelByCompany(gomock.Any(), int64(1), payload.Reason, payload.CompanyId,
			r.Header.Get("X-User")).Return(repository.ErrOfferNotFound).Times(1)
		a.handlerCancelOfferByCompany(w, r)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("good request", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/booking/offers/1/cancel", bytes.NewReader(reqBody))
		rctx := chi.NewRouteContext()
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		rctx.URLParams.Add("id", "1")
		r.Header.Set("X-User", "user_id1xxx")
		r.Header.Set("X-Email", "some_email@main.com")
		w := httptest.NewRecorder()
		repo.EXPECT().OfferCancelByCompany(gomock.Any(), int64(1), payload.Reason,
			payload.CompanyId, r.Header.Get("X-User")).Return(nil).Times(1)
		a.handlerCancelOfferByCompany(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestHandlerCancelOfferByCustomer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock.NewMockRepository(ctrl)
	a := createMockApp(repo)

	payload := struct {
		Reason string `json:"reason"`
	}{Reason: "test"}

	reqBody, err := json.Marshal(payload)
	assert.NoError(t, err)

	t.Run("empty body", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/booking/offers/1/cancel/customer", nil)
		rctx := chi.NewRouteContext()
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		rctx.URLParams.Add("id", "1")
		r.Header.Set("X-User", "user_id1xxx")
		r.Header.Set("X-Email", "some@email.com")
		w := httptest.NewRecorder()
		a.handlerCancelOfferByCustomer(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/booking/offers/1/cancel/customer", bytes.NewReader(reqBody))
		rctx := chi.NewRouteContext()
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		rctx.URLParams.Add("id", "1")
		r.Header.Set("X-User", "user_id1xxx")
		r.Header.Set("X-Email", "some@email.com")
		w := httptest.NewRecorder()
		repo.EXPECT().OfferCancelByCustomer(gomock.Any(), int64(1), payload.Reason, r.Header.Get("X-User")).Return(repository.ErrOfferNotFound).Times(1)
		a.handlerCancelOfferByCustomer(w, r)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("good request", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/booking/offers/1/cancel/customer", bytes.NewReader(reqBody))
		rctx := chi.NewRouteContext()
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		rctx.URLParams.Add("id", "1")
		r.Header.Set("X-User", "user_id1xxx")
		r.Header.Set("X-Email", "some@email.com")
		w := httptest.NewRecorder()
		repo.EXPECT().OfferCancelByCustomer(gomock.Any(), int64(1), payload.Reason, r.Header.Get("X-User")).Return(nil).Times(1)
		a.handlerCancelOfferByCustomer(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestHandlerBookOffer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock.NewMockRepository(ctrl)
	a := createMockApp(repo)

	payload := struct {
		Reason string `json:"reason"`
	}{Reason: "test"}

	reqBody, err := json.Marshal(payload)
	assert.NoError(t, err)

	t.Run("not found", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/booking/offers/1/book", bytes.NewReader(reqBody))
		rctx := chi.NewRouteContext()
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		rctx.URLParams.Add("id", "1")
		r.Header.Set("X-User", "user_id1xxx")
		r.Header.Set("X-Email", "some@email.com")
		w := httptest.NewRecorder()
		repo.EXPECT().Book(gomock.Any(), int64(1), r.Header.Get("X-User")).Return(repository.ErrOfferNotFound).Times(1)
		a.handlerBookOffer(w, r)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("good request", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/booking/offers/1/book", bytes.NewReader(reqBody))
		rctx := chi.NewRouteContext()
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		rctx.URLParams.Add("id", "1")
		r.Header.Set("X-User", "user_id1xxx")
		r.Header.Set("X-Email", "some@email.com")
		w := httptest.NewRecorder()
		repo.EXPECT().Book(gomock.Any(), int64(1), r.Header.Get("X-User")).Return(nil).Times(1)
		a.handlerBookOffer(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestHandlerGetCompanyOffers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock.NewMockRepository(ctrl)
	a := createMockApp(repo)
	r := httptest.NewRequest("GET", "/booking/company/offers/123?page=1&limit=10", nil)
	rctx := chi.NewRouteContext()
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	rctx.URLParams.Add("id", "123")
	w := httptest.NewRecorder()
	repo.EXPECT().Services(gomock.Any()).Return([]model.Service{}, nil)
	repo.EXPECT().CompanyOffers(gomock.Any(), int64(123), 1, 10).Return([]model.Offer{
		{
			ID:          1,
			CompanyID:   123,
			ServiceID:   1,
			ServiceName: "Haircut",
			CompanyName: "My company",
		},
	}, nil)
	a.handlerGetCompanyOffers(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "My company")
}

func TestHandlerGetCustomerOffers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock.NewMockRepository(ctrl)
	a := createMockApp(repo)
	r := httptest.NewRequest("GET", "/booking/customer/offers?page=1&limit=10", nil)
	r.Header.Set("X-User", "user_id1xxx")
	r.Header.Set("X-Email", "some@email.com")
	w := httptest.NewRecorder()
	repo.EXPECT().Services(gomock.Any()).Return([]model.Service{}, nil)
	repo.EXPECT().CustomerOffers(gomock.Any(), "user_id1xxx", 1, 10).Return([]model.Offer{
		{
			ID:          1,
			CompanyID:   123,
			ServiceID:   1,
			ServiceName: "Haircut",
			CompanyName: "My company",
		},
	}, nil)
	a.handlerGetCustomerOffers(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "My company")
}

func TestHandlerSearchOffers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock.NewMockRepository(ctrl)
	a := createMockApp(repo)
	r := httptest.NewRequest("GET", "/booking/offers?service_id=1&from=2024-01-01T10:00:00%2B01:00&to=2024-01-01T11:00:00%2B01:00&page=1&limit=10", nil)
	w := httptest.NewRecorder()
	from := time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	repo.EXPECT().Services(gomock.Any()).Return([]model.Service{
		{
			ID:       1,
			ParentID: 0,
			Name:     "test",
		},
	}, nil)
	repo.EXPECT().OfferSearch(gomock.Any(), int64(1), from, to, 1, 10).Return([]model.Offer{
		{
			ID:          1,
			CompanyID:   123,
			ServiceID:   1,
			ServiceName: "Haircut",
			CompanyName: "My company",
		},
	}, nil)
	a.handlerSearchOffers(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "My company")
}
