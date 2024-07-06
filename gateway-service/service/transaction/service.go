package transaction

import (
	"context"
	"log/slog"

	transactionpb "github.com/majidmohsenifar/heli-tech/data-contracts/proto/transaction"
)

type Service struct {
	transactionClient transactionpb.TransactionClient
	logger            *slog.Logger
}

type WithdrawParams struct {
	UserID int64   `json:"-"`
	Amount float64 `json:"amount" validate:"required,gt=0"`
}

type TransactionDetail struct {
	ID         int64   `json:"id"`
	CreatedAt  int64   `json:"createdAt"`
	Amount     float64 `json:"amount"`
	NewBalance float64 `json:"newBalance"`
}

type DepositParams struct {
	UserID int64   `json:"-"`
	Amount float64 `json:"amount" validate:"required,gt=0"`
}

type GetUserTransactionsParams struct {
	UserID   int64  `json:"-"`
	Page     uint32 `form:"page"`
	PageSize uint32 `form:"pageSize"`
}

type Transaction struct {
	ID        int64
	Amount    float64
	Kind      string
	CreatedAt int64
}

func (s *Service) Withdraw(
	ctx context.Context,
	params WithdrawParams,
) (TransactionDetail, error) {
	res, err := s.transactionClient.Withdraw(ctx, &transactionpb.WithdrawRequest{
		UserID: params.UserID,
		Amount: params.Amount,
	})
	if err != nil {
		return TransactionDetail{}, err
	}
	return TransactionDetail{
		ID:         res.Id,
		CreatedAt:  res.CreatedAt,
		Amount:     res.Amount,
		NewBalance: res.NewBalance,
	}, nil
}

func (s *Service) Deposit(
	ctx context.Context,
	params DepositParams,
) (TransactionDetail, error) {
	res, err := s.transactionClient.Deposit(ctx, &transactionpb.DepositRequest{
		UserID: params.UserID,
		Amount: params.Amount,
	})
	if err != nil {
		return TransactionDetail{}, err
	}
	return TransactionDetail{
		ID:         res.Id,
		CreatedAt:  res.CreatedAt,
		Amount:     res.Amount,
		NewBalance: res.NewBalance,
	}, nil
}

func (s *Service) GetUserTransactions(
	ctx context.Context,
	params GetUserTransactionsParams,
) ([]Transaction, error) {
	res, err := s.transactionClient.GetTransactions(ctx, &transactionpb.GetTransactionsRequest{
		Page:     params.Page,
		PageSize: params.PageSize,
		UserID:   params.UserID,
	})
	if err != nil {
		return nil, err
	}
	txs := make([]Transaction, len(res.Transactions))
	for i, tx := range res.Transactions {
		txs[i] = Transaction{
			ID:        tx.ID,
			Amount:    tx.Amount,
			Kind:      tx.Kind,
			CreatedAt: tx.CreatedAt,
		}
	}
	return txs, nil
}

func NewService(
	transactionClient transactionpb.TransactionClient,
	logger *slog.Logger,
) *Service {
	return &Service{
		transactionClient: transactionClient,
		logger:            logger,
	}
}
