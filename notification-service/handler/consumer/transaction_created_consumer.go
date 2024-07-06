package consumer

import (
	"context"
	"log/slog"

	transactionpb "github.com/majidmohsenifar/heli-tech/data-contracts/proto/transaction"
	"github.com/majidmohsenifar/heli-tech/notification-service/service/notification"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

const (
	TopicTransactionCreated = "transaction.created"
)

type TransactionCreatedConsumer struct {
	kafkaReader         *kafka.Reader
	notificationService *notification.Service
	logger              *slog.Logger
}

func (c *TransactionCreatedConsumer) Consume(ctx context.Context) {
	for {
		msg, err := c.kafkaReader.FetchMessage(ctx)
		if err != nil {
			c.logger.Error("failed to fetch message", err)
			break
		}
		var params transactionpb.TransactionCreatedEvent
		err = proto.Unmarshal(msg.Value, &params)
		if err != nil {
			c.logger.Error("failed to unmarshal message", err)
			break
		}
		err = c.notificationService.SendNotification(ctx, notification.SendNotificationParams{
			UserID:        0,
			TransactionID: 0,
			Amount:        0,
		})
		if err == nil {
			err = c.kafkaReader.CommitMessages(ctx, msg)
			if err != nil {
				c.logger.Error("failed to commit message", err)
				break
			}
		} else {
			c.logger.Error("failed to decrease video comment count", err)
			break
		}
	}
}

func NewTransactionCreatedConsumer(
	kafkaReader *kafka.Reader,
	notificationService *notification.Service,
	logger *slog.Logger,
) *TransactionCreatedConsumer {
	return &TransactionCreatedConsumer{
		kafkaReader:         kafkaReader,
		notificationService: notificationService,
		logger:              logger,
	}
}
