package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/bopoh24/ma_1/company/internal/config"
	"github.com/bopoh24/ma_1/company/internal/model"
	"github.com/bopoh24/ma_1/company/internal/repository"
	"github.com/bopoh24/ma_1/pkg/sql/builder"
	_ "github.com/lib/pq"
)

type Repository struct {
	psql *sq.StatementBuilderType
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

func (r *Repository) CompanyCreate(ctx context.Context, company model.Company) error {
	q := r.psql.Insert("company").
		Columns("owner", "name", "description", "address", "phone", "active").
		Values(company.Owner, company.Name, company.Description, company.Address, company.Phone, company.Active)
	_, err := q.ExecContext(ctx)
	return err
}

func (r *Repository) CompanyUpdate(ctx context.Context, company model.Company) error {
	q := r.psql.Update("company").
		Set("name", company.Name).
		Set("description", company.Description).
		Set("address", company.Address).
		Set("phone", company.Phone).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": company.ID})
	_, err := q.ExecContext(ctx)
	return err

}

func (r *Repository) CompanyByID(ctx context.Context, id int64) (model.Company, error) {
	q := r.psql.Select("id", "owner", "logo", "name", "description", "address",
		"phone", "active", "created_at", "updated_at").
		From("company").
		Where(sq.Eq{"id": id})
	row := q.QueryRowContext(ctx)
	var company model.Company
	err := row.Scan(&company.ID, &company.Owner, &company.Logo, &company.Name, &company.Description, &company.Address,
		&company.Phone, &company.Active, &company.Created, &company.Updated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return company, repository.ErrCompanyNotFound
		}
		return model.Company{}, err
	}
	return company, nil
}

func (r *Repository) CompanyActivateDeactivate(ctx context.Context, id int64, active bool) error {
	q := r.psql.Update("company").
		Set("active", active).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": id})
	_, err := q.ExecContext(ctx)
	return err
}

func (r *Repository) CompanyUpdateLocation(ctx context.Context, id int64, lat float64, lng float64) error {
	q := r.psql.Update("company").
		Set("location", sq.Expr("point(?, ?)", lat, lng)).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": id})
	_, err := q.ExecContext(ctx)
	return err
}

func (r *Repository) CompanyUpdateLogo(ctx context.Context, id int64, logo string) error {
	q := r.psql.Update("company").
		Set("logo", logo).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": id})
	_, err := q.ExecContext(ctx)
	return err
}
