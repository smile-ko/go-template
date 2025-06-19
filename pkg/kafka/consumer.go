package kafkabus

import (
	"context"
	"github/smile-ko/go-template/config"
	"log"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

type MessageHandler func(ctx context.Context, message kafka.Message) error

type Consumer struct {
	conf   *config.Config
	notify chan error
	wg     sync.WaitGroup
}

func NewConsumer(conf *config.Config) *Consumer {
	return &Consumer{
		conf:   conf,
		notify: make(chan error, 1),
	}
}

func (c *Consumer) Notify() <-chan error {
	return c.notify
}

func (c *Consumer) Handler(ctx context.Context, topic string, handler MessageHandler) {
	c.wg.Add(1)
	defer c.wg.Done()

	conf := c.conf.Kafka

	maxWait, _ := time.ParseDuration(conf.Consumer.MaxWait)
	startOffset := kafka.FirstOffset
	if conf.Consumer.StartOffset == "last" {
		startOffset = kafka.LastOffset
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     conf.Brokers,
		Topic:       topic,
		GroupID:     conf.GroupID,
		MinBytes:    conf.Consumer.MinBytes,
		MaxBytes:    conf.Consumer.MaxBytes,
		MaxWait:     maxWait,
		StartOffset: startOffset,
	})
	defer func() {
		_ = reader.Close()
		log.Printf("Kafka [%s] stopped", topic)
	}()

	log.Printf("Kafka [%s] started", topic)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := reader.FetchMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return // shutdown triggered
				}
				log.Printf("Kafka [%s] FetchMessage error: %v", topic, err)
				c.notify <- err
				return
			}

			if err := handler(ctx, msg); err != nil {
				log.Printf("Kafka [%s] Handler error: %v", topic, err)
			}

			if err := reader.CommitMessages(ctx, msg); err != nil {
				log.Printf("Kafka [%s] Commit error: %v", topic, err)
			}
		}
	}
}

func (c *Consumer) Wait() {
	c.wg.Wait()
}
