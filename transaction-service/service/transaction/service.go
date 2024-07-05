package transaction

import (
	"context"
	"log/slog"

	"github.com/majidmohsenifar/heli-tech/transaction-service/core"
	"github.com/majidmohsenifar/heli-tech/transaction-service/repository"
)

var ()

type Service struct {
	db     core.PgxInterface
	repo   repository.Querier
	logger *slog.Logger
}

type WithdrawParams struct {
	UserID int64
	Amount float64
}

type TransactionDetail struct {
}

type DepositParams struct {
	UserID int64
	Amount float64
}

func (s *Service) Withdraw(ctx context.Context, params WithdrawParams) (TransactionDetail, error) {
	//TODO: we should use lock here
	panic("here we go")
}

func (s *Service) Deposit(ctx context.Context, params DepositParams) (TransactionDetail, error) {
	panic("here we go")
}

func NewService(
	db core.PgxInterface,
	repo repository.Querier,
	logger *slog.Logger,
) *Service {
	return &Service{
		db:     db,
		repo:   repo,
		logger: logger,
	}
}
