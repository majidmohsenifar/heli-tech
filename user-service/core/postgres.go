package core

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxInterface interface {
	Begin(context.Context) (pgx.Tx, error)
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	Close()
}

func NewDBClient(ctx context.Context, url string) (PgxInterface, error) {
	dbPool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}
	if err := dbPool.Ping(ctx); err != nil {
		return nil, err
	}
	return dbPool, nil
}
