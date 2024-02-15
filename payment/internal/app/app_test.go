package app

import (
	"bytes"
	"encoding/json"
	"github.com/bopoh24/ma_1/payment/internal/repository"
	"github.com/bopoh24/ma_1/payment/internal/service"
	mock "github.com/bopoh24/ma_1/payment/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func createMockApp(repo service.Repository) *App {
	return &App{
		service: service.New(repo),
	}
}

func TestHandlerCreateAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock.NewMockRepository(ctrl)
	a := createMockApp(repo)

	r := httptest.NewRequest("POST", "/payment/create-account", nil)
	w := httptest.NewRecorder()
	r.Header.Set("X-User", "user_id1xxx")
	r.Header.Set("X-Email", "some@emeil.com")
	repo.EXPECT().CreateAccount(gomock.Any(), "user_id1xxx").Return(nil)
	a.handlerCreateAccount(w, r)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestHandlerTopUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock.NewMockRepository(ctrl)
	a := createMockApp(repo)

	payload := struct {
		Amount float32 `json:"amount"`
	}{Amount: 10.5}

	reqBody, err := json.Marshal(payload)
	assert.NoError(t, err)

	t.Run("Error", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/payment/top-up", bytes.NewReader(reqBody))
		w := httptest.NewRecorder()
		r.Header.Set("X-User", "user_id1xxx")
		r.Header.Set("X-Email", "some@emeil.com")
		repo.EXPECT().TopUp(gomock.Any(), "user_id1xxx", float32(10.5)).Return(repository.ErrAccountNotFound)
		a.handlerTopUp(w, r)
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "account not found")
	})

	t.Run("Success", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/payment/top-up", bytes.NewReader(reqBody))
		w := httptest.NewRecorder()
		r.Header.Set("X-User", "user_id1xxx")
		r.Header.Set("X-Email", "some@emeil.com")
		repo.EXPECT().TopUp(gomock.Any(), "user_id1xxx", float32(10.5)).Return(nil)
		a.handlerTopUp(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestHandlerBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock.NewMockRepository(ctrl)
	a := createMockApp(repo)

	t.Run("Error", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/payment/balance", nil)
		w := httptest.NewRecorder()
		r.Header.Set("X-User", "user_id1xxx")
		r.Header.Set("X-Email", "some@emeil.com")
		repo.EXPECT().Balance(gomock.Any(), "user_id1xxx").Return(float32(0), repository.ErrAccountNotFound)
		a.handlerBalance(w, r)
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "account not found")
	})

	t.Run("Success", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/payment/balance", nil)
		w := httptest.NewRecorder()
		r.Header.Set("X-User", "user_id1xxx")
		r.Header.Set("X-Email", "some@emeil.com")
		repo.EXPECT().Balance(gomock.Any(), "user_id1xxx").Return(float32(10.5), nil)
		a.handlerBalance(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "10.5")
	})
}

func TestHandlerPaymentCancel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock.NewMockRepository(ctrl)
	a := createMockApp(repo)

	payload := struct {
		OrderId int64 `json:"order_id"`
	}{OrderId: 1}

	reqBody, err := json.Marshal(payload)
	assert.NoError(t, err)

	t.Run("Error", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/payment/cancel", bytes.NewReader(reqBody))
		r.Header.Set("X-User", "user_id1xxx")
		r.Header.Set("X-Email", "some@emeil.com")
		w := httptest.NewRecorder()
		repo.EXPECT().PaymentCancel(gomock.Any(), int64(1)).Return(repository.ErrPaymentNotFound)
		a.handlerPaymentCancel(w, r)
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "payment not found")
	})

	t.Run("Success", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/payment/cancel", bytes.NewReader(reqBody))
		r.Header.Set("X-User", "user_id1xxx")
		r.Header.Set("X-Email", "some@emeil.com")
		w := httptest.NewRecorder()
		repo.EXPECT().PaymentCancel(gomock.Any(), int64(1)).Return(nil)
		a.handlerPaymentCancel(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestHandlerPaymentMake(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock.NewMockRepository(ctrl)
	a := createMockApp(repo)

	payload := struct {
		OrderID int64   `json:"order_id"`
		Amount  float32 `json:"amount"`
	}{OrderID: 1, Amount: 10.5}

	reqBody, err := json.Marshal(payload)
	assert.NoError(t, err)

	t.Run("Not found", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/payment/make", bytes.NewReader(reqBody))
		r.Header.Set("X-User", "user_id1xxx")
		r.Header.Set("X-Email", "some@emeil.com")
		w := httptest.NewRecorder()
		repo.EXPECT().PaymentMake(gomock.Any(), int64(1), "user_id1xxx", float32(10.5)).Return(repository.ErrAccountNotFound)
		a.handlerPaymentMake(w, r)
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "account not found")
	})
	t.Run("Insufficient funds", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/payment/make", bytes.NewReader(reqBody))
		r.Header.Set("X-User", "user_id1xxx")
		r.Header.Set("X-Email", "some@emeil.com")
		w := httptest.NewRecorder()
		repo.EXPECT().PaymentMake(gomock.Any(), int64(1), "user_id1xxx", float32(10.5)).Return(repository.ErrInsufficientFunds)
		a.handlerPaymentMake(w, r)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
		assert.Contains(t, w.Body.String(), "insufficient funds")
	})
	t.Run("Success", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/payment/make", bytes.NewReader(reqBody))
		r.Header.Set("X-User", "user_id1xxx")
		r.Header.Set("X-Email", "some@emeil.com")
		w := httptest.NewRecorder()
		repo.EXPECT().PaymentMake(gomock.Any(), int64(1), "user_id1xxx", float32(10.5)).Return(nil)
		a.handlerPaymentMake(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
