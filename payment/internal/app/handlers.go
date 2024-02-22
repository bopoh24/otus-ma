package app

import (
	"encoding/json"
	"errors"
	"github.com/bopoh24/ma_1/payment/internal/repository"
	"github.com/bopoh24/ma_1/pkg/http/helper"
	"net/http"
)

func (a *App) handlerCreateAccount(w http.ResponseWriter, r *http.Request) {
	claims, err := helper.ExtractClaims(r)
	if err != nil {
		helper.ErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	err = a.service.CreateAccount(r.Context(), claims.Id)
	if err != nil {
		if errors.Is(err, repository.ErrAccountNotFound) {
			helper.ErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.JSONResponse(w, http.StatusCreated, nil)
}

func (a *App) handlerTopUp(w http.ResponseWriter, r *http.Request) {
	claims, err := helper.ExtractClaims(r)
	if err != nil {
		helper.ErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	payload := struct {
		Amount float32 `json:"amount"`
	}{}

	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	err = a.service.TopUp(r.Context(), claims.Id, payload.Amount)
	if err != nil {
		if errors.Is(err, repository.ErrAccountNotFound) {
			helper.ErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.JSONResponse(w, http.StatusOK, nil)
}

func (a *App) handlerBalance(w http.ResponseWriter, r *http.Request) {
	claims, err := helper.ExtractClaims(r)
	if err != nil {
		helper.ErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	balance, err := a.service.Balance(r.Context(), claims.Id)
	if err != nil {
		if errors.Is(err, repository.ErrAccountNotFound) {
			helper.ErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.JSONResponse(w, http.StatusOK, struct {
		Balance float32 `json:"balance"`
	}{Balance: balance})

}

func (a *App) handlerPaymentMake(w http.ResponseWriter, r *http.Request) {
	claims, err := helper.ExtractClaims(r)
	if err != nil {
		helper.ErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	payload := struct {
		OfferID int64   `json:"offer_id"`
		Amount  float32 `json:"amount"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	err = a.service.PaymentMake(r.Context(), payload.OfferID, claims.Id, payload.Amount)
	if err != nil {
		if errors.Is(err, repository.ErrAccountNotFound) {
			helper.ErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, repository.ErrInsufficientFunds) {
			helper.ErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.JSONResponse(w, http.StatusOK, nil)
}

func (a *App) handlerPaymentCancel(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		OrderID int64 `json:"offer_id"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	err = a.service.PaymentCancel(r.Context(), payload.OrderID)
	if err != nil {
		if errors.Is(err, repository.ErrPaymentNotFound) {
			helper.ErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.JSONResponse(w, http.StatusOK, nil)
}
