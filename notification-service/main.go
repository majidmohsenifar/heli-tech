package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/majidmohsenifar/heli-tech/notification-service/config"
	"github.com/majidmohsenifar/heli-tech/notification-service/core"
	"github.com/majidmohsenifar/heli-tech/notification-service/handler/consumer"
	"github.com/majidmohsenifar/heli-tech/notification-service/logger"
	"github.com/majidmohsenifar/heli-tech/notification-service/service/notification"
)

func main() {
	viper := config.NewViper("./config/")
	logger := logger.NewLogger()
	kafkaURLs := viper.GetStringSlice("kafka.urls")
	notificationService := notification.NewService()
	kafkaReaderBuilder := core.NewKafkaReaderBuilder(kafkaURLs, "videoCreatedConsumer", 100000)
	kafkaReader := kafkaReaderBuilder.SetTopic(consumer.TopicTransactionCreated).Build()
	transactionCreatedConsumer := consumer.NewTransactionCreatedConsumer(
		kafkaReader,
		notificationService,
		logger)

	ctx := contextWithSigTerm(context.Background())
	forever := make(chan bool)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, os.Interrupt)

	transactionCreatedConsumer.Consume(ctx)
	go func() {
		for range sigChan {
			forever <- true
		}
	}()
	<-forever
}

func contextWithSigTerm(ctx context.Context) context.Context {
	inner, cancel := context.WithCancel(ctx)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		signal.Stop(sigChan)
		close(sigChan)
		log.Println("shutting down...")
		cancel()
	}()
	return inner
}
