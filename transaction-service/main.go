package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"

	transactionpb "github.com/majidmohsenifar/heli-tech/data-contracts/proto/transaction"
	"github.com/majidmohsenifar/heli-tech/transaction-service/config"
	"github.com/majidmohsenifar/heli-tech/transaction-service/core"
	"github.com/majidmohsenifar/heli-tech/transaction-service/handler/transactiongrpc"
	"github.com/majidmohsenifar/heli-tech/transaction-service/logger"
	"github.com/majidmohsenifar/heli-tech/transaction-service/repository"
	"github.com/majidmohsenifar/heli-tech/transaction-service/service/transaction"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	grpcPort = "50051"
)

func main() {
	ctx := context.Background()
	viper := config.NewViper("./config/")
	logger := logger.NewLogger()
	dbClient, err := core.NewDBClient(ctx, viper.GetString("db.dsn"))
	if err != nil {
		logger.Error("failed to initiate a db client", err)
		os.Exit(1)
	}
	defer dbClient.Close()
	repo := repository.New()
	redisClient, err := core.NewRedisClient(viper.GetString("redis.dsn"))
	if err != nil {
		logger.Error("failed to initiate a redis client", err)
		os.Exit(1)
	}
	defer redisClient.Close()

	kafkaURLs := viper.GetStringSlice("kafka.urls")
	kafkaWriter := core.NewKafkaWriter(kafkaURLs, "")
	transactionEventManager := transaction.NewTransactionEventManager(kafkaWriter, logger)
	redisLocker := core.NewRedisLocker(redisClient)
	transactionService := transaction.NewService(
		dbClient,
		repo,
		redisLocker,
		logger,
		transactionEventManager,
	)
	grpcPanicRecoveryHandler := func(p any) error {
		err := errors.New("recovered from panic")
		tempErr, ok := p.(error)
		if ok {
			err = tempErr
		} else {
			panicStr, ok := p.(string)
			if ok {
				err = errors.New(panicStr)
			}
		}
		logger.Error("recovered from panic", err)
		return status.Errorf(codes.Internal, "%s", "something went wrong")
	}
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
		),
		grpc.ChainStreamInterceptor(
			recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
		),
	)
	transactionGrpcServer := transactiongrpc.NewServer(
		transactionService,
	)
	transactionpb.RegisterTransactionServer(grpcServer, transactionGrpcServer)
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", grpcPort))
	if err != nil {
		logger.Error("can not listen to grpcPort", err)
		os.Exit(1)
	}
	err = grpcServer.Serve(l)
	if err != nil {
		logger.Error("can not serv grpc", err)
		os.Exit(1)
	}
}
