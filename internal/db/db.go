package db

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type DB interface {
	// QueryRow(query string, args ...interface{}) *sql.Row
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type postgresDB struct {
	*sqlx.DB
}

func (p *postgresDB) QueryRow(query string, args ...interface{}) *sql.Row {
	return p.DB.QueryRow(query, args...)
}

func (p *postgresDB) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return p.DB.SelectContext(ctx, dest, query, args...)
}

func (p *postgresDB) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return p.DB.GetContext(ctx, dest, query, args...)
}

func (p *postgresDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return p.DB.ExecContext(ctx, query, args...)
}
