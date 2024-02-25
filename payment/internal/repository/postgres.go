package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/bopoh24/ma_1/payment/internal/config"
	"github.com/bopoh24/ma_1/payment/internal/model"
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

func (r *Repository) CreateAccount(ctx context.Context, customerID string) error {
	q := r.psql.Builder().Insert("account").Columns("customer", "balance").Values(customerID, 0)
	_, err := q.ExecContext(ctx)
	return err
}

func (r *Repository) TopUp(ctx context.Context, customerID string, amount float32) error {
	tx, err := r.psql.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			err = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()
	q := r.psql.Builder().Update("account").
		Set("balance", sq.Expr("balance + ?", amount)).
		Where(sq.Eq{"customer": customerID}).RunWith(tx)
	_, err = q.ExecContext(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrAccountNotFound
		}
	}
	qt := r.psql.Builder().Insert("transaction").Columns("type", "customer", "amount").
		Values(model.TransactionTypeTopUp, customerID, amount).RunWith(tx)
	_, err = qt.ExecContext(ctx)
	return err
}

func (r *Repository) Balance(ctx context.Context, customerID string) (float32, error) {
	q := r.psql.Builder().Select("balance").From("account").Where(sq.Eq{"customer": customerID})
	row := q.QueryRowContext(ctx)
	var balance float32
	err := row.Scan(&balance)
	return balance, err
}

func (r *Repository) PaymentMake(ctx context.Context, offerId int64, customerID string, amount float32) error {
	balance, err := r.Balance(ctx, customerID)
	if err != nil {
		return err
	}
	balance -= amount
	if balance < 0 {
		return ErrInsufficientFunds
	}

	tx, err := r.psql.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			err = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	q := r.psql.Builder().Update("account").
		Set("balance", balance).
		Where(sq.Eq{"customer": customerID}).RunWith(tx)
	_, err = q.ExecContext(ctx)

	qt := r.psql.Builder().Insert("transaction").Columns("type", "customer", "offer_id", "amount").
		Values(model.TransactionTypePayment, customerID, offerId, amount).RunWith(tx)
	_, err = qt.ExecContext(ctx)
	return err
}

func (r *Repository) PaymentCancel(ctx context.Context, offerId int64) error {
	tx, err := r.psql.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			err = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()
	q := r.psql.Builder().Select("customer", "amount").From("transaction").Where(sq.Eq{"offer_id": offerId})
	row := q.QueryRowContext(ctx)
	var customerID string
	var amount float32
	err = row.Scan(&customerID, &amount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrPaymentNotFound
		}
		return err
	}

	qa := r.psql.Builder().Update("account").
		Set("balance", sq.Expr("balance + ?", amount)).
		Where(sq.Eq{"customer": customerID}).RunWith(tx)
	_, err = qa.ExecContext(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrAccountNotFound
		}
	}
	qt := r.psql.Builder().Insert("transaction").Columns("type", "customer", "offer_id", "amount").
		Values(model.TransactionTypeRefund, customerID, offerId, amount).RunWith(tx)
	_, err = qt.ExecContext(ctx)
	return err
}

// Close closes the repository
func (r *Repository) Close(ctx context.Context) error {
	return r.psql.Close()
}
