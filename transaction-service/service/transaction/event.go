package transaction

import (
	"context"
	"log/slog"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

const (
	TopicTransactionCreated = "transaction.Created"
)

type TransactionEventManager struct {
	kafkaWriter *kafka.Writer
	logger      *slog.Logger
}

type TransactionCreatedEventParams struct {
	UserID        int64
	TransactionID int64
	Type          string
	Amount        float64
	Balance       float64
	CreatedAt     int64
}

func (em *TransactionEventManager) PublishTransactionCreatedEvent(
	ctx context.Context,
	tx TransactionCreatedEventParams,
) {
	data := pb.ContentPurchasedEventParams{
		ContentID: tx.ContentID,
		UserID:    tx.UserID,
		Type:      tx.Type,
	}
	buff, err := proto.Marshal(&data)
	if err != nil {
		em.logger.ErrToSentry("failed to marshal proto:", err)
	}
	err := em.kafkaWriter.WriteMessages(
		ctx,
		kafka.Message{
			Key:   nil,
			Value: buff,
			Topic: TopicTransactionCreated,
		},
	)
	if err != nil {
		em.logger.ErrToSentry("failed to write messages:", err)
	}

}

func NewTransactionEventManager(
	kafkaWriter *kafka.Writer,
	logger *slog.Logger,
) *TransactionEventManager {
	return &TransactionEventManager{
		kafkaWriter: kafkaWriter,
		logger:      logger,
	}
}
