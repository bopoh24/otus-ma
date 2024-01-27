package pg

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/bopoh24/ma_1/customer/internal/config"
	"github.com/bopoh24/ma_1/customer/internal/model"
	"github.com/bopoh24/ma_1/customer/internal/repository"
	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

// New returns a new Repository struct
func New(dbConf config.Postgres) (*Repository, error) {

	psqlConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbConf.Host, dbConf.Port, dbConf.User, dbConf.Pass, dbConf.Database)
	db, err := sql.Open("postgres", psqlConn)
	if err != nil {
		return nil, err
	}

	return &Repository{db: db}, nil
}

// UserCreate creates a new user
func (r *Repository) UserCreate(user *model.User) error {
	q := `INSERT INTO users (external_id, username, first_name, last_name, email, phone) 
			VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := r.db.QueryRow(q, user.ExternalID, user.Username, user.FirstName, user.LastName, user.Email, user.Phone).
		Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}

// UserByID returns a user by id
func (r *Repository) UserByID(id int64) (*model.User, error) {
	row := r.db.QueryRow("SELECT username, first_name, last_name, email, phone FROM users WHERE id = $1", id)
	user := &model.User{}
	err := row.Scan(&user.Username, &user.FirstName, &user.LastName, &user.Email, &user.Phone)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

// UserByExternalID returns a user by external id
func (r *Repository) UserByExternalID(externalId string) (model.User, error) {
	row := r.db.QueryRow(
		`SELECT id, external_id, username, first_name, last_name, email, phone, description 
				FROM users WHERE external_id = $1`,
		externalId)
	user := model.User{}
	err := row.Scan(&user.ID, &user.ExternalID, &user.Username, &user.FirstName, &user.LastName, &user.Email,
		&user.Phone, &user.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, repository.ErrUserNotFound
		}
		return user, err
	}
	return user, nil
}

// UserUpdate updates a user
func (r *Repository) UserUpdate(user *model.User) error {
	q := `UPDATE users SET username = $1, first_name = $2, last_name = $3, email = $4, phone = $5 WHERE id = $6`
	res, err := r.db.Exec(q, user.Username, user.FirstName, user.LastName, user.Email, user.Phone, user.ID)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrUserNotFound
	}
	return nil
}

// UserUpdateByExternalID updates a user by external id
func (r *Repository) UserUpdateByExternalID(user *model.User) error {
	q := `UPDATE users SET username = $1, first_name = $2, last_name = $3, phone = $4, description = $5 
             WHERE external_id = $6`
	res, err := r.db.Exec(q, user.Username, user.FirstName, user.LastName, user.Phone, user.Description,
		user.ExternalID)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrUserNotFound
	}
	return nil
}

// UserDelete deletes a user by id
func (r *Repository) UserDelete(id int64) error {
	res, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrUserNotFound
	}
	return nil
}
