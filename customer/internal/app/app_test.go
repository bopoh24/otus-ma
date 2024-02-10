package app

import (
	"bytes"
	"encoding/json"
	"github.com/bopoh24/ma_1/customer/internal/model"
	"github.com/bopoh24/ma_1/customer/internal/repository"
	"github.com/bopoh24/ma_1/customer/internal/service"
	mock "github.com/bopoh24/ma_1/customer/mocks"
	"github.com/bopoh24/ma_1/pkg/verifier/phone"
	verifierMock "github.com/bopoh24/ma_1/pkg/verifier/phone/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func createMockApp(repo service.Repository, verifier phone.Verifier) *App {
	return &App{
		service: service.New(repo, verifier),
	}
}

func TestHandlerProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := httptest.NewRequest("GET", "/customer/profile", nil)
	userID := "123"
	r.Header.Set("X-User", userID)

	t.Run("OK", func(t *testing.T) {

		repo := mock.NewMockRepository(ctrl)
		app := createMockApp(repo, nil)
		repo.EXPECT().CustomerByID(gomock.Any(), userID).Return(model.Customer{
			ID:        userID,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "test@email.com",
			Phone:     "123456789",
		}, nil)
		w := httptest.NewRecorder()
		app.handlerProfile(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "John")
	})
}

func TestHandlerProfileUpdate(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userID := "123"

	t.Run("Bad request", func(t *testing.T) {
		r := httptest.NewRequest("PUT", "/customer/profile", nil)
		r.Header.Set("X-User", userID)

		repo := mock.NewMockRepository(ctrl)
		app := createMockApp(repo, nil)
		w := httptest.NewRecorder()
		app.handlerProfileUpdate(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)

	})

	t.Run("exists", func(t *testing.T) {
		body := model.Customer{
			FirstName: "John",
			LastName:  "Doe",
		}
		data, err := json.Marshal(body)
		assert.NoError(t, err)

		r := httptest.NewRequest("PUT", "/customer/profile", bytes.NewBuffer(data))
		r.Header.Set("X-User", userID)

		repo := mock.NewMockRepository(ctrl)
		app := createMockApp(repo, nil)
		body.ID = userID
		repo.EXPECT().CustomerUpdate(gomock.Any(), body).Return(nil)
		repo.EXPECT().CustomerByID(gomock.Any(), userID).Return(body, nil)
		w := httptest.NewRecorder()
		app.handlerProfileUpdate(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "John")
	})

	t.Run("not exists", func(t *testing.T) {
		body := model.Customer{
			FirstName: "John",
			LastName:  "Doe",
		}
		data, err := json.Marshal(body)
		assert.NoError(t, err)

		r := httptest.NewRequest("PUT", "/customer/profile", bytes.NewBuffer(data))
		r.Header.Set("X-User", userID)

		repo := mock.NewMockRepository(ctrl)
		app := createMockApp(repo, nil)
		body.ID = userID
		repo.EXPECT().CustomerUpdate(gomock.Any(), body).Return(repository.ErrCustomerNotFound)
		repo.EXPECT().CustomerCreate(gomock.Any(), body).Return(nil)
		repo.EXPECT().CustomerByID(gomock.Any(), userID).Return(body, nil)
		w := httptest.NewRecorder()
		app.handlerProfileUpdate(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "John")
	})
}

func TestHandlerRequestPhoneVerification(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userID := "123"
	phone := "123456789"

	t.Run("Bad request", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/customer/phone/verify", nil)
		r.Header.Set("X-User", userID)

		repo := mock.NewMockRepository(ctrl)
		app := createMockApp(repo, nil)
		w := httptest.NewRecorder()
		app.handlerRequestPhoneVerification(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)

	})

	t.Run("OK", func(t *testing.T) {
		body := model.Customer{
			Phone: phone,
		}
		data, err := json.Marshal(body)
		assert.NoError(t, err)

		r := httptest.NewRequest("POST", "/customer/phone/verify", bytes.NewBuffer(data))

		phoneVerifier := verifierMock.NewMockVerifier(ctrl)
		app := createMockApp(nil, phoneVerifier)
		w := httptest.NewRecorder()
		phoneVerifier.EXPECT().Send(gomock.Any(), phone).Return(nil)
		app.handlerRequestPhoneVerification(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestHandlerVerifyPhone(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userID := "123"
	phone := "123456789"
	code := "1234"

	t.Run("Bad request", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/customer/phone/verify/check", nil)
		r.Header.Set("X-User", userID)

		repo := mock.NewMockRepository(ctrl)
		app := createMockApp(repo, nil)
		w := httptest.NewRecorder()
		app.handlerVerifyPhone(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)

	})

	t.Run("OK", func(t *testing.T) {
		body := struct {
			Phone string `json:"phone"`
			Code  string `json:"code"`
		}{
			Phone: phone,
			Code:  code,
		}
		data, err := json.Marshal(body)
		assert.NoError(t, err)

		r := httptest.NewRequest("POST", "/customer/phone/verify/check", bytes.NewBuffer(data))
		r.Header.Set("X-User", userID)

		repo := mock.NewMockRepository(ctrl)
		phoneVerifier := verifierMock.NewMockVerifier(ctrl)
		app := createMockApp(repo, phoneVerifier)
		w := httptest.NewRecorder()
		phoneVerifier.EXPECT().Check(gomock.Any(), phone, code).Return(nil)
		repo.EXPECT().CustomerUpdatePhone(gomock.Any(), userID, phone).Return(nil)
		app.handlerVerifyPhone(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
