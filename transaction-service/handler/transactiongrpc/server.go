package transactiongrpc

import (
	"context"

	transactionpb "github.com/majidmohsenifar/heli-tech/data-contracts/proto/transaction"
	"github.com/majidmohsenifar/heli-tech/transaction-service/service/transaction"
)

type server struct {
	transactionService *transaction.Service
	transactionpb.UnimplementedTransactionServer
}

func (s *server) Withdraw(
	ctx context.Context,
	req *transactionpb.WithdrawRequest,
) (*transactionpb.WithdrawResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s *server) Deposit(
	ctx context.Context,
	req *transactionpb.DepositRequest,
) (*transactionpb.DepositResponse, error) {
	panic("not implemented") // TODO: Implement
	//TODO: validate req params

	//result, err := s.transactionService.Deposit(transaction.DepositParams{
	//UserID: req.UserID,
	//Amount: float64(req.Amount),
	//})
}

func NewServer(
	transactionService *transaction.Service,
) transactionpb.TransactionServer {
	return &server{
		transactionService: transactionService,
	}
}
