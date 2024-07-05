package transactiongrpc

import (
	"context"

	transactionpb "github.com/majidmohsenifar/heli-tech/data-contracts/proto/transaction"
	"github.com/majidmohsenifar/heli-tech/transaction-service/service/transaction"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	transactionService *transaction.Service
	transactionpb.UnimplementedTransactionServer
}

func (s *server) Withdraw(
	ctx context.Context,
	req *transactionpb.WithdrawRequest,
) (*transactionpb.WithdrawResponse, error) {
	if req.UserID < 1 {
		return nil, status.Error(codes.Code(400), "userID is required")
	}
	if req.Amount <= 0 {
		return nil, status.Error(codes.Code(400), "amount should be positive")
	}
	result, err := s.transactionService.Withdraw(ctx, transaction.WithdrawParams{
		UserID: req.UserID,
		Amount: req.Amount,
	})
	if err == transaction.ErrOngoingRequest {
		return nil, status.Error(codes.Code(422), err.Error())
	}
	if err == transaction.ErrInsufficientBalance {
		return nil, status.Error(codes.Code(422), err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Code(500), "something went wrong")
	}
	return &transactionpb.WithdrawResponse{
		CreatedAt:  result.CreatedAt,
		Amount:     result.Amount,
		NewBalance: result.NewBalance,
	}, nil
}

func (s *server) Deposit(
	ctx context.Context,
	req *transactionpb.DepositRequest,
) (*transactionpb.DepositResponse, error) {
	if req.UserID < 1 {
		return nil, status.Error(codes.Code(400), "userID is required")
	}
	if req.Amount <= 0 {
		return nil, status.Error(codes.Code(400), "amount should be positive")
	}
	result, err := s.transactionService.Deposit(ctx, transaction.DepositParams{
		UserID: req.UserID,
		Amount: req.Amount,
	})
	if err == transaction.ErrOngoingRequest {
		return nil, status.Error(codes.Code(422), err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Code(500), "something went wrong")
	}
	return &transactionpb.DepositResponse{
		CreatedAt:  result.CreatedAt,
		Amount:     result.Amount,
		NewBalance: result.NewBalance,
	}, nil
}

func NewServer(
	transactionService *transaction.Service,
) transactionpb.TransactionServer {
	return &server{
		transactionService: transactionService,
	}
}
