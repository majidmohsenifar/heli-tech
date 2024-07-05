package test

import (
	"context"
	"math/big"
	"net"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	transactionpb "github.com/majidmohsenifar/heli-tech/data-contracts/proto/transaction"
	"github.com/majidmohsenifar/heli-tech/transaction-service/core"
	"github.com/majidmohsenifar/heli-tech/transaction-service/handler/transactiongrpc"
	"github.com/majidmohsenifar/heli-tech/transaction-service/repository"
	"github.com/majidmohsenifar/heli-tech/transaction-service/service/transaction"

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

	//empty userID
	req := transactionpb.WithdrawRequest{
		UserID: 0,
		Amount: 0,
	}
	res, err := client.Withdraw(ctx, &req)
	assert.Nil(res)
	e, ok := status.FromError(err)
	assert.True(ok)
	assert.Equal(e.Code(), codes.Code(400))
	assert.Equal(e.Message(), "userID is required")

	//zero amount
	req = transactionpb.WithdrawRequest{
		UserID: 1,
		Amount: 0,
	}
	res, err = client.Withdraw(ctx, &req)
	assert.Nil(res)
	e, ok = status.FromError(err)
	assert.True(ok)
	assert.Equal(e.Code(), codes.Code(400))
	assert.Equal(e.Message(), "amount should be positive")

	//negative amount
	req = transactionpb.WithdrawRequest{
		UserID: 1,
		Amount: -20,
	}
	res, err = client.Withdraw(ctx, &req)
	assert.Nil(res)
	e, ok = status.FromError(err)
	assert.True(ok)
	assert.Equal(e.Code(), codes.Code(400))
	assert.Equal(e.Message(), "amount should be positive")
}

func TestTransaction_Withdraw_NoBalanceInDBForUser(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	db := getDB()
	err := truncateDB()
	assert.Nil(err)
	repo := repository.New()
	redisClient := getRedis()
	redisLocker := core.NewRedisLocker(redisClient)
	transactionService := transaction.NewService(db, repo, redisLocker, getLogger(), nil)
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
	client := transactionpb.NewTransactionClient(userConn)

	req := transactionpb.WithdrawRequest{
		UserID: 1,
		Amount: 100,
	}
	res, err := client.Withdraw(ctx, &req)
	assert.Nil(res)
	e, ok := status.FromError(err)
	assert.True(ok)
	assert.Equal(e.Code(), codes.Code(422))
	assert.Equal(e.Message(), "the balance is not enough for withdraw")
}

func TestTransaction_Withdraw_BalanceExistInDB_ButIsNotEnough(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	db := getDB()
	err := truncateDB()
	assert.Nil(err)
	repo := repository.New()
	redisClient := getRedis()
	redisLocker := core.NewRedisLocker(redisClient)
	transactionService := transaction.NewService(db, repo, redisLocker, getLogger(), nil)
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
	_, err = repo.CreateUserBalanceOrIncreaseAmount(ctx, db, repository.CreateUserBalanceOrIncreaseAmountParams{
		UserID: 1,
		Amount: pgtype.Numeric{Int: big.NewInt(50), Valid: true},
	})
	assert.Nil(err)

	client := transactionpb.NewTransactionClient(userConn)
	req := transactionpb.WithdrawRequest{
		UserID: 1,
		Amount: 100,
	}
	res, err := client.Withdraw(ctx, &req)
	assert.Nil(res)
	e, ok := status.FromError(err)
	assert.True(ok)
	assert.Equal(e.Code(), codes.Code(422))
	assert.Equal(e.Message(), "the balance is not enough for withdraw")
}

func TestTransaction_Withdraw_Successful(t *testing.T) {
	//user balance is 150, and want to withdraw 100
	assert := assert.New(t)
	ctx := context.Background()
	db := getDB()
	err := truncateDB()
	assert.Nil(err)
	repo := repository.New()
	redisClient := getRedis()
	redisLocker := core.NewRedisLocker(redisClient)
	transactionEventManager:= transaction.NewTransactionEventManager(nil, getLogger())

	transactionService := transaction.NewService(db, repo, redisLocker, getLogger(), nil)
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
	_, err = repo.CreateUserBalanceOrIncreaseAmount(ctx, db, repository.CreateUserBalanceOrIncreaseAmountParams{
		UserID: 1,
		Amount: pgtype.Numeric{Int: big.NewInt(150), Valid: true},
	})
	assert.Nil(err)

	client := transactionpb.NewTransactionClient(userConn)
	req := transactionpb.WithdrawRequest{
		UserID: 1,
		Amount: 100,
	}
	res, err := client.Withdraw(ctx, &req)
	assert.Nil(err)
	assert.Equal(res.Amount, 100.0)
	assert.Equal(res.NewBalance, 50.0)
	assert.Greater(res.CreatedAt, int64(0))
}
