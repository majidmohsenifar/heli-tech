package core

import (
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaReaderBuilder struct {
	urls     []string
	topic    string
	groupID  string
	maxBytes int
}

func NewKafkaReaderBuilder(
	urls []string,
	groupID string,
	maxBytes int,
) *KafkaReaderBuilder {
	return &KafkaReaderBuilder{
		urls:     urls,
		groupID:  groupID,
		maxBytes: maxBytes,
	}
}

func (b *KafkaReaderBuilder) SetTopic(topic string) *KafkaReaderBuilder {
	b.topic = topic
	return b
}

func (b *KafkaReaderBuilder) Build() *kafka.Reader {
	return NewKafkaReader(b.urls, b.topic, b.groupID, b.maxBytes)
}

func NewKafkaReader(urls []string, topic, groupID string, MaxBytes int) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  urls,
		GroupID:  groupID,
		Topic:    topic,
		MaxBytes: MaxBytes, // 10MB
	})
}

func NewKafkaWriter(urls []string, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:                   kafka.TCP(urls...),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
		Async:                  true,
	}
}

func NewKafkaClient(urls []string) *kafka.Client {
	return &kafka.Client{
		Addr:    kafka.TCP(urls...),
		Timeout: 10 * time.Second,
	}
}

