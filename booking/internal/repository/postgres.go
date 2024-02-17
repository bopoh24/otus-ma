package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/bopoh24/ma_1/booking/internal/config"
	"github.com/bopoh24/ma_1/booking/pkg/model"
	"github.com/bopoh24/ma_1/pkg/sql/builder"
	_ "github.com/lib/pq"
	"time"
)

type Repository struct {
	psql *builder.Postgres
}

// New returns a new Repository struct
func New(dbConf config.Postgres) (*Repository, error) {
	psqlConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbConf.Host, dbConf.Port, dbConf.User, dbConf.Pass, dbConf.Database)
	psql, err := builder.NewPostgresBuilder(psqlConn)
	if err != nil {
		return nil, err
	}
	return &Repository{psql: psql}, nil

}

// Services returns a list of services
func (r *Repository) Services(ctx context.Context) ([]model.Service, error) {
	q := r.psql.Builder().Select("id", "parent_id", "name").From("service")
	rows, err := q.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var services []model.Service
	for rows.Next() {
		var s model.Service
		err = rows.Scan(&s.ID, &s.ParentID, &s.Name)
		if err != nil {
			return nil, err
		}
		services = append(services, s)
	}
	return services, nil
}

// ServiceAdd adds a new service
func (r *Repository) ServiceAdd(ctx context.Context, service model.Service) error {
	q := r.psql.Builder().Insert("service").
		Columns("parent_id", "name").
		Values(service.ParentID, service.Name)
	_, err := q.ExecContext(ctx)
	return err
}

// OfferAdd adds a new offer
func (r *Repository) OfferAdd(ctx context.Context, offer model.Offer) error {
	q := r.psql.Builder().Insert("offer").
		Columns("service_id", "company_id", "location", "datetime",
			"description", "price", "status", "created_by").
		Values(offer.ServiceID, offer.CompanyID, offer.Location, offer.Datetime,
			offer.Description, offer.Price, offer.Status, offer.CreatedBy)
	_, err := q.ExecContext(ctx)
	return err
}

// OfferDelete deletes an offer
func (r *Repository) OfferDelete(ctx context.Context, id int64, companyId int64) error {
	q := r.psql.Builder().Delete("offer").Where(sq.Eq{"id": id}, sq.Eq{"company_id": companyId})
	_, err := q.ExecContext(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrOfferNotFound
	}
	return err
}

// OfferChangeStatus changes the status of an offer
func (r *Repository) OfferChangeStatus(ctx context.Context, id int64, status model.OfferStatus) error {
	q := r.psql.Builder().Update("offer").
		Set("status", status).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": id})
	_, err := q.ExecContext(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrOfferNotFound
	}
	return err
}

// OfferCancelByCompany cancels an offer by company
func (r *Repository) OfferCancelByCompany(ctx context.Context, id int64, reason string, companyId int64, managerId string) error {
	q := r.psql.Builder().Update("offer").
		Set("status", model.OfferStatusCanceledByCompany).
		Set("cancel_reason", reason).
		Set("updated_by", managerId).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": id}, sq.Eq{"company_id": companyId})
	_, err := q.ExecContext(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrOfferNotFound
	}
	return err
}

// OfferCancelByCustomer cancels an offer by customer
func (r *Repository) OfferCancelByCustomer(ctx context.Context, id int64, reason string, customerId string) error {
	q := r.psql.Builder().Update("offer").
		Set("status", model.OfferStatusCanceledByCustomer).
		Set("cancel_reason", reason).
		Set("updated_by", customerId).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": id}, sq.Eq{"customer": customerId})
	_, err := q.ExecContext(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrOfferNotFound
	}
	return err
}

// OfferSearch searches for offers
func (r *Repository) OfferSearch(ctx context.Context, serviceId int64, from, to time.Time, page, limit int) ([]model.Offer, error) {
	q := r.psql.Builder().Select("id", "service_id", "company_id", "company_name", "location",
		"datetime", "description", "price", "status").
		From("offer").
		Where(
			sq.Eq{"status": model.OfferStatusOpen},
			sq.Eq{"service_id": serviceId},
			sq.GtOrEq{"datetime": from},
			sq.LtOrEq{"datetime": to}).
		Offset(uint64(page * limit)).Limit(uint64(limit))
	rows, err := q.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var offers []model.Offer
	for rows.Next() {
		var o model.Offer
		err = rows.Scan(&o.ID, &o.ServiceID, &o.CompanyID, &o.CompanyName, &o.Location,
			&o.Datetime, &o.Description, &o.Price, &o.Status)
		if err != nil {
			return nil, err
		}
	}
	return offers, nil
}

// Book books an offer
func (r *Repository) Book(ctx context.Context, offerId int64, customerId string) error {
	q := r.psql.Builder().Update("offer").
		Set("status", model.OfferStatusReserved).
		Set("customer", customerId).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": offerId}, sq.Eq{"status": model.OfferStatusOpen})
	_, err := q.ExecContext(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrOfferNotFound
	}
	return err

}

func (r *Repository) CompanyOffers(ctx context.Context, companyId int64, page, limit int) ([]model.Offer, error) {
	q := r.psql.Builder().Select("id", "service_id", "customer", "location", "datetime",
		"description", "price", "status", "cancel_reason", "created_by", "updated_by", "created_at", "updated_at").
		From("offer").
		Where(sq.Eq{"company_id": companyId}).
		Offset(uint64(page * limit)).Limit(uint64(limit))
	rows, err := q.QueryContext(ctx)
	var offers []model.Offer
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return offers, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var o model.Offer
		err = rows.Scan(&o.ID, &o.ServiceID, &o.Customer, &o.Location, &o.Datetime,
			&o.Description, &o.Price, &o.Status, &o.CancelReason, &o.CreatedBy, &o.UpdatedBy, &o.CreatedAt, &o.UpdatedAt)
		if err != nil {
			return nil, err
		}
		offers = append(offers, o)
	}
	return offers, nil
}

func (r *Repository) CustomerOffers(ctx context.Context, customerId string, page, limit int) ([]model.Offer, error) {
	q := r.psql.Builder().Select("id", "service_id", "company_id", "location", "datetime",
		"description", "price", "status", "cancel_reason").
		From("offer").
		Where(sq.Eq{"customer": customerId}).
		Offset(uint64(page * limit)).Limit(uint64(limit))
	rows, err := q.QueryContext(ctx)
	var offers []model.Offer
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return offers, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var o model.Offer
		err = rows.Scan(&o.ID, &o.ServiceID, &o.CompanyID, &o.Location, &o.Datetime,
			&o.Description, &o.Price, &o.Status, &o.CancelReason)
		if err != nil {
			return nil, err
		}
		offers = append(offers, o)
	}
	return offers, nil
}

// Close closes the repository
func (r *Repository) Close(ctx context.Context) error {
	return r.psql.Close()
}
