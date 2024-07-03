package payment

import (
	"context"
	"log/slog"

	paymentpb "github.com/majidmohsenifar/heli-tech/data-contracts/proto/payment"
)

type Service struct {
	paymentClient paymentpb.PaymentClient
	logger        *slog.Logger
}

type WithdrawParams struct {
}

type WithdrawResponse struct {
}

type DepositParams struct {
}

type DepostiResponse struct {
}

func (s *Service) Withdraw(
	ctx context.Context,
	params WithdrawParams,
) error {
	panic("handle this later")
}

func (s *Service) Deposit(
	ctx context.Context,
	params DepositParams,
) error {
	panic("handle this later")
}

func NewService(
	paymentClient paymentpb.PaymentClient,
	logger *slog.Logger,
) *Service {
	return &Service{
		paymentClient: paymentClient,
		logger:        logger,
	}
}
