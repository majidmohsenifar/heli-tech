// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package repository

import (
	"context"
)

type Querier interface {
	CreateTransaction(ctx context.Context, db DBTX, arg CreateTransactionParams) (Transaction, error)
	CreateUserBalanceOrDecreaseAmount(ctx context.Context, db DBTX, arg CreateUserBalanceOrDecreaseAmountParams) (UserBalance, error)
	CreateUserBalanceOrIncreaseAmount(ctx context.Context, db DBTX, arg CreateUserBalanceOrIncreaseAmountParams) (UserBalance, error)
}

var _ Querier = (*Queries)(nil)
