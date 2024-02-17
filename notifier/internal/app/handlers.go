package app

import (
	"encoding/json"
	offerModel "github.com/bopoh24/ma_1/booking/pkg/model"
	"github.com/bopoh24/ma_1/notifier/internal/model"
	"net/http"
)

func (a *App) handlerSend(w http.ResponseWriter, r *http.Request) {
	var notification model.BookingNotification
	if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	status := notification.Offer.Status
	switch status {
	case offerModel.OfferStatusFailed:
		if err := a.service.BookingFailed(notification); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case offerModel.OfferStatusPaid:
		if err := a.service.BookingPaid(notification); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case offerModel.OfferStatusSubmitted:
		if err := a.service.BookingSubmitted(notification); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case offerModel.OfferStatusCanceledByCustomer:
		if err := a.service.BookingCancelledByCustomer(notification); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case offerModel.OfferStatusCanceledByCompany:
		if err := a.service.BookingCancelledByCompany(notification); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case offerModel.OfferStatusCompleted:
		if err := a.service.BookingCompleted(notification); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
