package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/bopoh24/ma_1/customer/internal/config"
	"github.com/bopoh24/ma_1/customer/internal/model"
	"github.com/bopoh24/ma_1/customer/internal/repository"
	_ "github.com/lib/pq"
)

type Repository struct {
	db   *sql.DB
	psql sq.StatementBuilderType
}

// New returns a new Repository struct
func New(dbConf config.Postgres) (*Repository, error) {
	psqlConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbConf.Host, dbConf.Port, dbConf.User, dbConf.Pass, dbConf.Database)
	db, err := sql.Open("postgres", psqlConn)
	if err != nil {
		return nil, err
	}
	dbCache := sq.NewStmtCache(db)
	if err != nil {
		return nil, err
	}
	return &Repository{
		psql: sq.StatementBuilder.RunWith(dbCache).PlaceholderFormat(sq.Dollar),
		db:   db,
	}, nil

}

// CustomerCreate creates a new customer profile
func (r *Repository) CustomerCreate(ctx context.Context, customer model.Customer) error {
	query := r.psql.Insert("customer").
		Columns("id", "email", "first_name", "last_name").Values(
		customer.ID, customer.Email, customer.FirstName, customer.LastName)
	_, err := query.ExecContext(ctx)
	return err
}

// CustomerUpdate updates a customer profile
func (r *Repository) CustomerUpdate(ctx context.Context, customer model.Customer) error {
	query := r.psql.Update("customer").
		Set("first_name", customer.FirstName).
		Set("last_name", customer.LastName).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": customer.ID})
	_, err := query.ExecContext(ctx)
	return err
}

// CustomerByID returns a customer profile by id
func (r *Repository) CustomerByID(ctx context.Context, id string) (model.Customer, error) {
	query := r.psql.Select("id", "email", "first_name", "last_name", "phone").
		From("customer").Where(sq.Eq{"id": id})
	row := query.QueryRowContext(ctx)
	customer := model.Customer{}
	err := row.Scan(&customer.ID, &customer.Email, &customer.FirstName, &customer.LastName, &customer.Phone)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return customer, repository.ErrCustomerNotFound
		}
	}
	return customer, err
}

// CustomerUpdatePhone updates a customer phone
func (r *Repository) CustomerUpdatePhone(ctx context.Context, id string, phone string) error {
	query := r.psql.Update("customer").
		Set("phone", phone).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": id})
	_, err := query.ExecContext(ctx)
	return err
}

// CustomerUpdateLocation updates a customer location
func (r *Repository) CustomerUpdateLocation(ctx context.Context, id string, lat float64, lng float64) error {
	query := r.psql.Update("customer").
		Set("location", sq.Expr("point(?, ?)", lat, lng)).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": id})
	_, err := query.ExecContext(ctx)
	return err
}
