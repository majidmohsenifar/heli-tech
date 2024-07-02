package client

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDBClient(ctx context.Context, url string) (*pgxpool.Pool, error) {
	dbPool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}
	if err := dbPool.Ping(ctx); err != nil {
		return nil, err
	}
	return dbPool, nil
}
