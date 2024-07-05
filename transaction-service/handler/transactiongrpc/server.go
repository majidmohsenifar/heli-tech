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

func (s *server) CreateTransaction(
	ctx context.Context,
	req *transactionpb.CreateTransactionRequest,
) (*transactionpb.CreateTransactionResponse, error) {
	panic("not implemented") // TODO: Implement
}

func NewServer(
	transactionService *transaction.Service,
) transactionpb.TransactionServer {
	return &server{
		transactionService: transactionService,
	}
}
