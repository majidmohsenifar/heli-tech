package test

import (
	"context"
	"testing"
	"time"

	transactionpb "github.com/majidmohsenifar/heli-tech/data-contracts/proto/transaction"
	"github.com/majidmohsenifar/heli-tech/notification-service/core"
	"github.com/majidmohsenifar/heli-tech/notification-service/handler/consumer"
	"github.com/majidmohsenifar/heli-tech/notification-service/service/notification"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func Test_TransactionCreatedConsumer(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	logger := getLogger()
	topic := consumer.TopicTransactionCreated
	viper := getViperConfig()
	kafkaURLs := viper.GetStringSlice("kafka.urls")
	kafkaClient := core.NewKafkaClient(kafkaURLs)
	_, err := kafkaClient.DeleteTopics(ctx, &kafka.DeleteTopicsRequest{
		Topics: []string{topic},
	})
	assert.Nil(err)
	res, err := kafkaClient.CreateTopics(ctx, &kafka.CreateTopicsRequest{
		Topics: []kafka.TopicConfig{
			{
				Topic:             topic,
				ReplicationFactor: 1,
				NumPartitions:     1,
			},
		},
	})
	assert.Nil(err)
	assert.Nil(res.Errors[topic])
	createdAt := time.Now().Add(-5 * time.Minute).Unix()
	msgValue := transactionpb.TransactionCreatedEvent{
		TransactionID: 1,
		UserID:        1,
		Amount:        100.0,
		Balance:       150,
		Kind:          "DEPOSIT",
		CreatedAt:     createdAt,
	}
	buff, err := proto.Marshal(&msgValue)
	assert.Nil(err)
	kafkaWriter := core.NewKafkaWriter(kafkaURLs, "")
	kafkaWriter.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Value: buff,
		Key:   nil,
	})
	kafkaReaderBuilder := core.NewKafkaReaderBuilder(kafkaURLs, "", 1000)
	kafkaReader := kafkaReaderBuilder.SetTopic(topic).Build()
	notificationService := notification.NewService()
	consumer := consumer.NewTransactionCreatedConsumer(kafkaReader, notificationService, logger)
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		//to be sure the data is consumed
		time.Sleep(2 * time.Second)
		cancel()
	}()
	consumer.Consume(ctx)

	notifs := notificationService.GetAllNotifications()
	assert.Len(notifs, 1)
	n1 := notifs[0]
	assert.Equal(n1.TransactionID, int64(1))
	assert.Equal(n1.UserID, int64(1))
	assert.Equal(n1.Amount, 100.0)
	assert.Equal(n1.Balance, 150.0)
	assert.Equal(n1.Kind, "DEPOSIT")
	assert.Equal(n1.CreatedAt, createdAt)
}
