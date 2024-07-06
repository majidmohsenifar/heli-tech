package transaction

import (
	"context"
	"log/slog"

	transactionpb "github.com/majidmohsenifar/heli-tech/data-contracts/proto/transaction"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

const (
	TopicTransactionCreated = "transaction.created"
)

type TransactionEventManager interface {
	PublishTransactionCreatedEvent(ctx context.Context, params TransactionCreatedEventParams)
}

type transactionEventManager struct {
	kafkaWriter *kafka.Writer
	logger      *slog.Logger
}

type TransactionCreatedEventParams struct {
	UserID        int64
	TransactionID int64
	Kind          string
	Amount        float64
	Balance       float64
	CreatedAt     int64
}

func (em *transactionEventManager) PublishTransactionCreatedEvent(
	ctx context.Context,
	params TransactionCreatedEventParams,
) {
	data := transactionpb.TransactionCreatedEvent{
		TransactionID: params.TransactionID,
		UserID:        params.UserID,
		Amount:        params.Amount,
		Balance:       params.Balance,
		Kind:          params.Kind,
		CreatedAt:     params.CreatedAt,
	}
	buff, err := proto.Marshal(&data)
	if err != nil {
		em.logger.Error("failed to marshal proto:", err)
		return
	}
	err = em.kafkaWriter.WriteMessages(
		ctx,
		kafka.Message{
			Key:   nil,
			Value: buff,
			Topic: TopicTransactionCreated,
		},
	)
	if err != nil {
		em.logger.Error("failed to write messages:", err)
	}
}

func NewTransactionEventManager(
	kafkaWriter *kafka.Writer,
	logger *slog.Logger,
) TransactionEventManager {
	return &transactionEventManager{
		kafkaWriter: kafkaWriter,
		logger:      logger,
	}
}
