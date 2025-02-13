// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: transaction.sql

package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createTransaction = `-- name: CreateTransaction :one
INSERT INTO transactions (
    user_id,
    kind,
    amount,
    created_at
) VALUES (
  $1, $2, $3, now()
) RETURNING id, user_id, kind, amount, created_at
`

type CreateTransactionParams struct {
	UserID int64
	Kind   Kind
	Amount pgtype.Numeric
}

func (q *Queries) CreateTransaction(ctx context.Context, db DBTX, arg CreateTransactionParams) (Transaction, error) {
	row := db.QueryRow(ctx, createTransaction, arg.UserID, arg.Kind, arg.Amount)
	var i Transaction
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Kind,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const getTransactionByID = `-- name: GetTransactionByID :one
SELECT id, user_id, kind, amount, created_at from transactions WHERE id = $1
`

func (q *Queries) GetTransactionByID(ctx context.Context, db DBTX, id int64) (Transaction, error) {
	row := db.QueryRow(ctx, getTransactionByID, id)
	var i Transaction
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Kind,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const getUserTransactionsByPagination = `-- name: GetUserTransactionsByPagination :many
SELECT id, user_id, kind, amount, created_at FROM transactions 
WHERE user_id = $1 
ORDER BY created_at DESC
OFFSET $2
LIMIT $3
`

type GetUserTransactionsByPaginationParams struct {
	UserID int64
	Offset int32
	Limit  int32
}

func (q *Queries) GetUserTransactionsByPagination(ctx context.Context, db DBTX, arg GetUserTransactionsByPaginationParams) ([]Transaction, error) {
	rows, err := db.Query(ctx, getUserTransactionsByPagination, arg.UserID, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Transaction
	for rows.Next() {
		var i Transaction
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Kind,
			&i.Amount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
