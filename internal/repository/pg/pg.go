package pg

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/bopoh24/ma_1/internal/config"
	"github.com/bopoh24/ma_1/internal/model"
	"github.com/bopoh24/ma_1/internal/repository"
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
	q := `INSERT INTO users (username, first_name, last_name, email, phone) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(q, user.Username, user.FirstName, user.LastName, user.Email, user.Phone)
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
