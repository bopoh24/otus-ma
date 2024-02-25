package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/bopoh24/ma_1/customer/internal/config"
	"github.com/bopoh24/ma_1/customer/internal/model"
	"github.com/bopoh24/ma_1/pkg/sql/builder"
	_ "github.com/lib/pq"
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
	return &Repository{
		psql: psql,
	}, nil

}

// CustomerCreate creates a new customer profile
func (r *Repository) CustomerCreate(ctx context.Context, customer model.Customer) error {
	query := r.psql.Builder().Insert("customer").
		Columns("id", "email", "first_name", "last_name").Values(
		customer.ID, customer.Email, customer.FirstName, customer.LastName)
	_, err := query.ExecContext(ctx)
	return err
}

// CustomerUpdate updates a customer profile
func (r *Repository) CustomerUpdate(ctx context.Context, customer model.Customer) error {
	query := r.psql.Builder().Update("customer").
		Set("first_name", customer.FirstName).
		Set("last_name", customer.LastName).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": customer.ID})
	_, err := query.ExecContext(ctx)
	return err
}

// CustomerByID returns a customer profile by id
func (r *Repository) CustomerByID(ctx context.Context, id string) (model.Customer, error) {
	query := r.psql.Builder().Select("id", "email", "first_name", "last_name", "phone", "created_at", "updated_at").
		From("customer").Where(sq.Eq{"id": id})
	row := query.QueryRowContext(ctx)
	customer := model.Customer{}
	err := row.Scan(&customer.ID, &customer.Email, &customer.FirstName, &customer.LastName, &customer.Phone,
		&customer.Created, &customer.Updated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return customer, ErrCustomerNotFound
		}
	}
	return customer, err
}

// CustomerUpdatePhone updates a customer phone
func (r *Repository) CustomerUpdatePhone(ctx context.Context, id string, phone string) error {
	query := r.psql.Builder().Update("customer").
		Set("phone", phone).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": id})
	_, err := query.ExecContext(ctx)
	return err
}

// CustomerUpdateLocation updates a customer location
func (r *Repository) CustomerUpdateLocation(ctx context.Context, id string, lat float64, lng float64) error {
	query := r.psql.Builder().Update("customer").
		Set("location", sq.Expr("point(?, ?)", lat, lng)).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": id})
	_, err := query.ExecContext(ctx)
	return err
}

// Close closes the repository
func (r *Repository) Close(ctx context.Context) error {
	return r.psql.Close()
}
