package test

import (
	"context"
	"net"
	"testing"

	transactionpb "github.com/majidmohsenifar/heli-tech/data-contracts/proto/transaction"
	"github.com/majidmohsenifar/heli-tech/transaction-service/handler/transactiongrpc"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func TestTransaction_Withdraw_InvalidInputs(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	transactionGrpcServer := transactiongrpc.NewServer(nil)

	l, err := net.Listen("tcp", "127.0.0.1:0")
	assert.Nil(err)
	defer l.Close()
	googleGrpcServer := grpc.NewServer()

	transactionpb.RegisterTransactionServer(googleGrpcServer, transactionGrpcServer)
	userConn, err := grpc.NewClient(
		l.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.Nil(err)
	go googleGrpcServer.Serve(l)
	client := transactionpb.NewTransactionClient(userConn)

	req := transactionpb.WithdrawRequest{
		UserID: 0,
		Amount: 0,
	}
	res, err := client.Withdraw(ctx, &req)
	assert.Nil(res)
	e, ok := status.FromError(err)
	assert.True(ok)
	assert.Equal(e.Code(), codes.Code(400))
	assert.Equal(e.Message(), "email is empty")
}
