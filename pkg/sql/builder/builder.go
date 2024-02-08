package builder

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
)

// NewPostgresBuilder is a wrapper around squirrel.StatementBuilderType
func NewPostgresBuilder(psqlConn string) (*sq.StatementBuilderType, error) {
	db, err := sql.Open("postgres", psqlConn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	dbCache := sq.NewStmtCache(db)
	if err != nil {
		return nil, err
	}
	builder := sq.StatementBuilder.RunWith(dbCache).PlaceholderFormat(sq.Dollar)
	return &builder, nil
}
