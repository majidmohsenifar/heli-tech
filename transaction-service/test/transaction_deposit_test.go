package test

import (
	"context"
	"net"
	"testing"

	transactionpb "github.com/majidmohsenifar/heli-tech/data-contracts/proto/transaction"
	"github.com/majidmohsenifar/heli-tech/transaction-service/core"
	"github.com/majidmohsenifar/heli-tech/transaction-service/handler/transactiongrpc"
	"github.com/majidmohsenifar/heli-tech/transaction-service/helper"
	"github.com/majidmohsenifar/heli-tech/transaction-service/repository"
	"github.com/majidmohsenifar/heli-tech/transaction-service/service/transaction"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func TestTransaction_Deposit_InvalidInputs(t *testing.T) {
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
	req := transactionpb.DepositRequest{
		UserID: 0,
		Amount: 0,
	}
	res, err := client.Deposit(ctx, &req)
	assert.Nil(res)
	e, ok := status.FromError(err)
	assert.True(ok)
	assert.Equal(e.Code(), codes.Code(400))
	assert.Equal(e.Message(), "userID is required")

	//zero amount
	req = transactionpb.DepositRequest{
		UserID: 1,
		Amount: 0,
	}
	res, err = client.Deposit(ctx, &req)
	assert.Nil(res)
	e, ok = status.FromError(err)
	assert.True(ok)
	assert.Equal(e.Code(), codes.Code(400))
	assert.Equal(e.Message(), "amount should be positive")

	//negative amount
	req = transactionpb.DepositRequest{
		UserID: 1,
		Amount: -20,
	}
	res, err = client.Deposit(ctx, &req)
	assert.Nil(res)
	e, ok = status.FromError(err)
	assert.True(ok)
	assert.Equal(e.Code(), codes.Code(400))
	assert.Equal(e.Message(), "amount should be positive")
}

func TestTransaction_Deposit_Successful(t *testing.T) {
	//user balance is 150, and want to withdraw 100
	assert := assert.New(t)
	ctx := context.Background()
	db := getDB()
	err := truncateDB()
	assert.Nil(err)
	repo := repository.New()
	redisClient := getRedis()
	redisLocker := core.NewRedisLocker(redisClient)

	topic := transaction.TopicTransactionCreated
	kafkaURLs := getViperConfig().GetStringSlice("kafka.urls")
	kafkaClient := core.NewKafkaClient(kafkaURLs)
	_, err = kafkaClient.DeleteTopics(ctx, &kafka.DeleteTopicsRequest{
		Topics: []string{topic},
	})
	assert.Nil(err)
	creationTopicRes, err := kafkaClient.CreateTopics(ctx, &kafka.CreateTopicsRequest{
		Topics: []kafka.TopicConfig{
			{
				Topic:             topic,
				NumPartitions:     1,
				ReplicationFactor: 1,
			},
		},
	})
	assert.Nil(creationTopicRes.Errors[topic])
	assert.Nil(err)
	kafkaWriter := core.NewKafkaWriter(kafkaURLs, "")

	transactionEventManager := transaction.NewTransactionEventManager(kafkaWriter, getLogger())
	transactionService := transaction.NewService(db, repo, redisLocker, getLogger(), transactionEventManager)
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
	req := transactionpb.DepositRequest{
		UserID: 1,
		Amount: 100,
	}
	res, err := client.Deposit(ctx, &req)
	assert.Nil(err)
	assert.Greater(res.Id, int64(0))
	assert.Equal(res.Amount, 100.0)
	assert.Equal(res.NewBalance, 100.0)
	assert.Greater(res.CreatedAt, int64(0))

	createdTx, err := repo.GetTransactionByID(ctx, db, res.Id)
	assert.Nil(err)
	assert.Equal(createdTx.ID, res.Id)
	amount, err := helper.PGNumericToFloat64(createdTx.Amount)
	assert.Nil(err)
	assert.Equal(amount, 100.0)
	assert.Equal(createdTx.UserID, int64(1))

	//checking kafka
	kafkaReaderBuilder := core.NewKafkaReaderBuilder(kafkaURLs, "", 1000)
	kafkaReader := kafkaReaderBuilder.SetTopic(topic).Build()
	msg, err := kafkaReader.FetchMessage(ctx)
	assert.Nil(err)
	assert.NotNil(msg)
	assert.Equal(msg.Topic, topic)
	var event transactionpb.TransactionCreatedEvent
	err = proto.Unmarshal(msg.Value, &event)
	assert.Nil(err)
	assert.Equal(event.UserID, int64(1))
	assert.Equal(event.UserID, int64(1))

	assert.Equal(event.TransactionID, createdTx.ID)
	assert.Equal(event.UserID, int64(1))
	assert.Equal(event.Amount, 100.0)
	assert.Equal(event.Balance, 100.0)
	assert.Equal(event.Kind, string(repository.KindDEPOSIT))
	assert.Greater(event.CreatedAt, int64(0))
}
