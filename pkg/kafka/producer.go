package kafkabus

import (
	"context"
	"github/smile-ko/go-template/config"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(conf *config.Config, topic string) *Producer {
	requiredAcks := kafka.RequireAll
	if conf.Kafka.Producer.RequiredAcks == 0 {
		requiredAcks = kafka.RequireNone
	}

	return &Producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(conf.Kafka.Brokers...),
			Topic:        topic,
			RequiredAcks: requiredAcks,
			MaxAttempts:  conf.Kafka.Producer.RequiredAcks,
		},
	}
}

// Produce sends a Kafka message with key and value
func (p *Producer) Produce(ctx context.Context, key []byte, value []byte) error {
	msg := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}

	return p.writer.WriteMessages(ctx, msg)
}

// Close the producer when done
func (p *Producer) Close() error {
	return p.writer.Close()
}
