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

type CreateTransactionParams struct {
	UserID int64
	Amount float64
	Kind   repository.Kind
}

func (s *Service) CreateTransaction(ctx context.Context, params CreateTransactionParams) error {
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
