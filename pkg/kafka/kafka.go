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
	conf   *option
	notify chan error
	wg     sync.WaitGroup
}

type Producer struct {
	writer *kafka.Writer
}

func NewConsumer(conf *config.Config) *Consumer {
	return &Consumer{
		conf:   newOption(conf),
		notify: make(chan error, 1),
	}
}

func (c *Consumer) Notify() <-chan error {
	return c.notify
}

func (c *Consumer) Handler(ctx context.Context, topic string, handler MessageHandler) {
	c.wg.Add(1)
	defer c.wg.Done()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     c.conf.Brokers,
		Topic:       topic,
		GroupID:     c.conf.GroupID,
		MinBytes:    c.conf.MinBytes,
		MaxBytes:    c.conf.MaxBytes,
		MaxWait:     c.conf.MaxWait,
		StartOffset: c.conf.StartOffset,
	})
	defer func() {
		_ = reader.Close()
		log.Printf("Kafka Consumer [%s] stopped", topic)
	}()

	log.Printf("Kafka Consumer [%s] started", topic)

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

func NewProducer(conf *config.Config, topic string) *Producer {
	opt := newOption(conf)

	return &Producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(opt.Brokers...),
			Topic:        topic,
			RequiredAcks: opt.RequiredAcks,
			MaxAttempts:  opt.MaxAttempts,
		},
	}
}

func (p *Producer) Produce(ctx context.Context, key, value []byte) error {
	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	})
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
