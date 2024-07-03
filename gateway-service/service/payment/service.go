package payment

import (
	"context"
	"log/slog"

	userpb "github.com/majidmohsenifar/heli-tech/data-contracts/proto/user"
)

type Service struct {
	//TODO: change this to payment
	paymentClient userpb.UserClient
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
	paymentClient userpb.UserClient,
	logger *slog.Logger,
) *Service {
	return &Service{
		paymentClient: paymentClient,
		logger:        logger,
	}
}
