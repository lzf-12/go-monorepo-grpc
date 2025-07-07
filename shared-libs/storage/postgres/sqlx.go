package storage

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Sqlxdb struct {
	*sqlx.DB
}

const pgDriver = "postgres"

func NewSqlx(db *sql.DB) *sqlx.DB {
	return sqlx.NewDb(db, pgDriver)
}
