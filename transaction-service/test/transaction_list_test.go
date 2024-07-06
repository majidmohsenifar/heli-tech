package test

import (
	"context"
	"math/big"
	"net"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	transactionpb "github.com/majidmohsenifar/heli-tech/data-contracts/proto/transaction"
	"github.com/majidmohsenifar/heli-tech/transaction-service/handler/transactiongrpc"
	"github.com/majidmohsenifar/heli-tech/transaction-service/repository"
	"github.com/majidmohsenifar/heli-tech/transaction-service/service/transaction"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestTransaction_GetUserTransactions_Successful(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	db := getDB()
	err := truncateDB()
	assert.Nil(err)
	repo := repository.New()
	transactionService := transaction.NewService(db, repo, nil, nil, nil)
	transactionGrpcServer := transactiongrpc.NewServer(transactionService)
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

	tx1, err := repo.CreateTransaction(ctx, db, repository.CreateTransactionParams{
		UserID: 1,
		Kind:   repository.KindDEPOSIT,
		Amount: pgtype.Numeric{Int: big.NewInt(100), Valid: true},
	})
	assert.Nil(err)

	tx2, err := repo.CreateTransaction(ctx, db, repository.CreateTransactionParams{
		UserID: 1,
		Kind:   repository.KindWITHDRAW,
		Amount: pgtype.Numeric{Int: big.NewInt(50), Valid: true},
	})
	assert.Nil(err)

	client := transactionpb.NewTransactionClient(userConn)
	//first we get page 0 with pageSize 1
	req := transactionpb.GetTransactionsRequest{
		Page:     0,
		PageSize: 1,
		UserID:   1,
	}
	result, err := client.GetTransactions(ctx, &req)
	assert.Nil(err)
	assert.Equal(len(result.Transactions), 1)
	//the order is desc so the first one is tx2
	tx := result.Transactions[0]
	assert.Equal(tx.ID, tx2.ID)
	assert.Equal(tx.Amount, 50.0)
	assert.Equal(tx.Kind, "WITHDRAW")
	assert.Equal(tx.CreatedAt, tx2.CreatedAt.Time.Unix())

	//first we get page 1 with pageSize 1
	req = transactionpb.GetTransactionsRequest{
		Page:     1,
		PageSize: 1,
		UserID:   1,
	}
	result, err = client.GetTransactions(ctx, &req)
	assert.Nil(err)
	assert.Equal(len(result.Transactions), 1)
	//the order is desc so this one is tx1
	tx = result.Transactions[0]
	assert.Equal(tx.ID, tx1.ID)
	assert.Equal(tx.Amount, 100.0)
	assert.Equal(tx.Kind, "DEPOSIT")
	assert.Equal(tx.CreatedAt, tx1.CreatedAt.Time.Unix())
}
