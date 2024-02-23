package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/bopoh24/ma_1/company/internal/config"
	"github.com/bopoh24/ma_1/company/internal/model"
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
	return &Repository{psql: psql}, nil

}

// CompanyCreate creates a new company profile
func (r *Repository) CompanyCreate(ctx context.Context, userId, email, firstName, lastName string, company model.Company) (int64, error) {
	tx, err := r.psql.Tx(ctx)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			err = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	q := r.psql.Builder().Insert("company").
		Columns("name", "description", "address", "phone", "email", "active").
		Values(company.Name, company.Description, company.Address,
			company.Phone, company.Email, company.Active).
		Suffix("RETURNING id").RunWith(tx)
	row := q.QueryRowContext(ctx)
	var companyId int64
	err = row.Scan(&companyId)
	if err != nil {
		return 0, err
	}

	manager := model.Manager{
		CompanyID: companyId,
		UserID:    userId,
		Email:     email,
		Role:      model.MangerRoleAdmin,
		Active:    true,
	}
	q = r.psql.Builder().Insert("company_manager").
		Columns("company_id", "user_id", "email", "first_name", "last_name", "role", "active").
		Values(manager.CompanyID, manager.UserID, manager.Email, firstName, lastName, manager.Role, manager.Active).RunWith(tx)
	_, err = q.ExecContext(ctx)
	return companyId, err
}

// CompanyUpdate updates a company profile
func (r *Repository) CompanyUpdate(ctx context.Context, company model.Company) error {
	q := r.psql.Builder().Update("company").
		Set("name", company.Name).
		Set("description", company.Description).
		Set("address", company.Address).
		Set("phone", company.Phone).
		Set("email", company.Email).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": company.ID})
	_, err := q.ExecContext(ctx)
	return err

}

// CompanyByID returns a company profile by its ID
func (r *Repository) CompanyByID(ctx context.Context, id int64) (model.Company, error) {
	q := r.psql.Builder().Select("id", "logo", "name", "description", "address",
		"phone", "email", "active", "created_at", "updated_at").
		From("company").
		Where(sq.Eq{"id": id})
	row := q.QueryRowContext(ctx)
	var company model.Company
	err := row.Scan(&company.ID, &company.Logo, &company.Name, &company.Description, &company.Address,
		&company.Phone, &company.Email, &company.Active, &company.Created, &company.Updated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return company, ErrCompanyNotFound
		}
		return model.Company{}, err
	}
	return company, nil
}

// CompanyActivateDeactivate activates or deactivates a company
func (r *Repository) CompanyActivateDeactivate(ctx context.Context, id int64, active bool) error {
	q := r.psql.Builder().Update("company").
		Set("active", active).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": id})
	_, err := q.ExecContext(ctx)
	return err
}

// CompanyUpdateLocation updates a company location
func (r *Repository) CompanyUpdateLocation(ctx context.Context, id int64, lat float64, lng float64) error {
	q := r.psql.Builder().Update("company").
		Set("location", sq.Expr("point(?, ?)", lat, lng)).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": id})
	_, err := q.ExecContext(ctx)
	return err
}

// CompanyUpdateLogo updates a company logo
func (r *Repository) CompanyUpdateLogo(ctx context.Context, id int64, logo string) error {
	q := r.psql.Builder().Update("company").
		Set("logo", logo).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": id})
	_, err := q.ExecContext(ctx)
	return err
}

// CompanyManagers returns a list of company managers
func (r *Repository) CompanyManagers(ctx context.Context, companyID int64) ([]model.Manager, error) {
	q := r.psql.Builder().Select("id", "company_id", "user_id", "email", "role", "active", "created_at", "updated_at").
		From("company_manager").
		Where(sq.Eq{"company_id": companyID})
	rows, err := q.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var managers []model.Manager
	for rows.Next() {
		var manager model.Manager
		err = rows.Scan(&manager.ID, &manager.CompanyID, &manager.UserID,
			&manager.Email, &manager.Role, &manager.Active, &manager.Created, &manager.Updated)
		if err != nil {
			return nil, err
		}
		managers = append(managers, manager)
	}
	return managers, nil
}

// ManagerByID returns a manager by its ID
func (r *Repository) ManagerByID(ctx context.Context, id int64) (model.Manager, error) {
	q := r.psql.Builder().Select("id", "company_id", "user_id", "email", "role", "active", "created_at", "updated_at").
		From("company_manager").
		Where(sq.Eq{"id": id})
	row := q.QueryRowContext(ctx)
	var manager model.Manager
	err := row.Scan(&manager.ID, &manager.CompanyID, &manager.UserID,
		&manager.Email, &manager.Role, &manager.Active, &manager.Created, &manager.Updated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return manager, ErrManagerNotFound
		}
		return model.Manager{}, err
	}
	return manager, nil
}

// ManagerByUserID returns a manager by its user ID
func (r *Repository) ManagerByUserID(ctx context.Context, userId string) (model.Manager, error) {
	q := r.psql.Builder().Select("id", "company_id", "user_id", "email", "role", "active", "created_at", "updated_at").
		From("company_manager").
		Where(sq.Eq{"user_id": userId})
	row := q.QueryRowContext(ctx)
	var manager model.Manager
	err := row.Scan(&manager.ID, &manager.CompanyID, &manager.UserID,
		&manager.Email, &manager.Role, &manager.Active, &manager.Created, &manager.Updated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return manager, ErrManagerNotFound
		}
		return model.Manager{}, err
	}
	return manager, nil

}

// ManagerByEmail returns a manager by its email
func (r *Repository) ManagerByEmail(ctx context.Context, email string) (model.Manager, error) {
	q := r.psql.Builder().Select("id", "company_id", "user_id", "email", "role", "active", "created_at", "updated_at").
		From("company_manager").
		Where(sq.Eq{"email": email})
	row := q.QueryRowContext(ctx)
	var manager model.Manager
	err := row.Scan(&manager.ID, &manager.CompanyID, &manager.UserID,
		&manager.Email, &manager.Role, &manager.Active, &manager.Created, &manager.Updated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return manager, ErrManagerNotFound
		}
		return model.Manager{}, err
	}
	return manager, nil
}

// ManagerActivateDeactivate activates or deactivates a manager
func (r *Repository) ManagerActivateDeactivate(ctx context.Context, id int64, active bool) error {
	q := r.psql.Builder().Update("company_manager").
		Set("active", active).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": id})
	_, err := q.ExecContext(ctx)
	return err
}

// ManagerSetRole sets a manager role
func (r *Repository) ManagerSetRole(ctx context.Context, id int64, role model.MangerRole) error {
	q := r.psql.Builder().Update("company_manager").
		Set("role", role).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": id})
	_, err := q.ExecContext(ctx)
	return err
}

// ManagerDelete deletes a manager
func (r *Repository) ManagerDelete(ctx context.Context, id int64) error {
	q := r.psql.Builder().Delete("company_manager").
		Where(sq.Eq{"id": id})
	_, err := q.ExecContext(ctx)
	return err
}

// ManagerInvite invites a manager to a company
func (r *Repository) ManagerInvite(ctx context.Context, companyID int64, email string, role model.MangerRole) error {
	//TODO implement me
	panic("implement me")
}

// ManagerRole returns a manager role
func (r *Repository) ManagerRole(ctx context.Context, companyId int64, userId string) (model.MangerRole, error) {
	q := r.psql.Builder().Select("role").From("company_manager").
		Where(sq.Eq{"company_id": companyId, "user_id": userId})
	row := q.QueryRowContext(ctx)
	var role model.MangerRole
	err := row.Scan(&role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return role, ErrManagerNotFound
		}
		return "", err
	}
	return role, nil
}

func (r *Repository) MyCompanies(ctx context.Context, userId string) ([]model.Company, error) {
	q := r.psql.Builder().Select("c.id", "c.logo", "c.name", "c.description", "c.address",
		"c.phone", "c.email", "c.active", "c.created_at", "c.updated_at").
		From("company c").
		Join("company_manager m ON c.id = m.company_id").
		Where(sq.Eq{"m.user_id": userId})

	companies := make([]model.Company, 0)
	rows, err := q.QueryContext(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return companies, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var company model.Company
		err = rows.Scan(&company.ID, &company.Logo, &company.Name, &company.Description, &company.Address,
			&company.Phone, &company.Email, &company.Active, &company.Created, &company.Updated)
		if err != nil {
			return nil, err
		}
		companies = append(companies, company)
	}
	return companies, nil
}

// Close closes the database connection
func (r *Repository) Close(ctx context.Context) error {
	return r.psql.Close()
}
