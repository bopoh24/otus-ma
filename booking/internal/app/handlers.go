package app

import (
	"encoding/json"
	"errors"
	"github.com/bopoh24/ma_1/booking/internal/model"
	"github.com/bopoh24/ma_1/booking/internal/repository"
	"github.com/bopoh24/ma_1/pkg/http/helper"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"time"
)

func (a *App) handlerGetServices(w http.ResponseWriter, r *http.Request) {
	services, err := a.service.Services(r.Context())
	if err != nil {
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.JSONResponse(w, http.StatusOK, services)
}

func (a *App) handlerAddService(w http.ResponseWriter, r *http.Request) {
	var service model.Service
	err := json.NewDecoder(r.Body).Decode(&service)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	err = a.service.ServiceAdd(r.Context(), service)
	if err != nil {
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.JSONResponse(w, http.StatusCreated, nil)
}

func (a *App) handlerAddOffer(w http.ResponseWriter, r *http.Request) {
	var offer model.Offer
	err := json.NewDecoder(r.Body).Decode(&offer)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if offer.ServiceID == 0 || offer.Datetime.IsZero() || offer.Price == 0 ||
		offer.CompanyID == 0 || offer.CompanyName == "" {
		helper.ErrorResponse(w, http.StatusBadRequest, "invalid offer")
		return
	}

	err = a.service.OfferAdd(r.Context(), offer)
	if err != nil {
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.JSONResponse(w, http.StatusCreated, nil)
}

func (a *App) handlerDeleteOffer(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	payload := struct {
		CompanyId int64 `json:"company_id"`
	}{}
	if err = json.NewDecoder(r.Body).Decode(&payload); err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	err = a.service.OfferDelete(r.Context(), id, payload.CompanyId)
	if err != nil {
		if errors.Is(err, repository.ErrOfferNotFound) {
			helper.ErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.JSONResponse(w, http.StatusOK, nil)
}

func (a *App) handlerChangeOfferStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	payload := struct {
		Status model.OrderStatus `json:"status"`
	}{}
	if err = json.NewDecoder(r.Body).Decode(&payload); err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	err = a.service.OfferChangeStatus(r.Context(), id, payload.Status)
	if err != nil {
		if errors.Is(err, repository.ErrOfferNotFound) {
			helper.ErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.JSONResponse(w, http.StatusOK, nil)
}

func (a *App) handlerCancelOfferByCompany(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	claims, err := helper.ExtractClaims(r)
	if err != nil {
		helper.ErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	payload := struct {
		CompanyId int64  `json:"company_id"`
		Reason    string `json:"reason"`
	}{}
	if err = json.NewDecoder(r.Body).Decode(&payload); err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	err = a.service.OfferCancelByCompany(r.Context(), id, payload.Reason, payload.CompanyId, claims.Id)
	if err != nil {
		if errors.Is(err, repository.ErrOfferNotFound) {
			helper.ErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.JSONResponse(w, http.StatusOK, nil)
}

func (a *App) handlerCancelOfferByCustomer(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	claims, err := helper.ExtractClaims(r)
	if err != nil {
		helper.ErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	payload := struct {
		Reason string `json:"reason"`
	}{}
	if err = json.NewDecoder(r.Body).Decode(&payload); err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	err = a.service.OfferCancelByCustomer(r.Context(), id, payload.Reason, claims.Id)
	if err != nil {
		if errors.Is(err, repository.ErrOfferNotFound) {
			helper.ErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.JSONResponse(w, http.StatusOK, nil)
}

func (a *App) handlerSearchOffers(w http.ResponseWriter, r *http.Request) {
	serviceIdStr := r.URL.Query().Get("service_id")
	serviceId, err := strconv.ParseInt(serviceIdStr, 10, 64)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	fromStr := r.URL.Query().Get("from")
	from, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	toStr := r.URL.Query().Get("to")
	to, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	offers, err := a.service.OfferSearch(r.Context(), serviceId, from.UTC(), to.UTC(), page, limit)
	if err != nil {
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.JSONResponse(w, http.StatusOK, offers)
}

func (a *App) handlerBookOffer(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	claims, err := helper.ExtractClaims(r)
	if err != nil {
		helper.ErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	err = a.service.Book(r.Context(), id, claims.Id)
	if err != nil {
		if errors.Is(err, repository.ErrOfferNotFound) {
			helper.ErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.JSONResponse(w, http.StatusOK, nil)
}

func (a *App) handlerGetCompanyOffers(w http.ResponseWriter, r *http.Request) {
	companyIdStr := chi.URLParam(r, "companyId")
	companyId, err := strconv.ParseInt(companyIdStr, 10, 64)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	offers, err := a.service.CompanyOffers(r.Context(), companyId, page, limit)
	if err != nil {
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.JSONResponse(w, http.StatusOK, offers)
}

func (a *App) handlerGetCustomerOffers(w http.ResponseWriter, r *http.Request) {
	claims, err := helper.ExtractClaims(r)
	if err != nil {
		helper.ErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	offers, err := a.service.CustomerOffers(r.Context(), claims.Id, page, limit)
	if err != nil {
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.JSONResponse(w, http.StatusOK, offers)
}
