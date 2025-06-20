package kafkabus

import (
	"context"
	"github/smile-ko/go-template/config"
	"github/smile-ko/go-template/pkg/logger"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type MessageHandler func(ctx context.Context, message kafka.Message) error

type Consumer struct {
	conf   *option
	notify chan error
	wg     sync.WaitGroup
	logger logger.ILogger
}

type Producer struct {
	writer *kafka.Writer
	logger logger.ILogger
}

func NewConsumer(conf *config.Config, l logger.ILogger) *Consumer {
	return &Consumer{
		conf:   newOption(conf),
		notify: make(chan error, 1),
		logger: l,
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
		c.logger.Info("Kafka Consumer closed", zap.String("topic", topic))
	}()

	c.logger.Info("Kafka Consumer started", zap.String("topic", topic), zap.Strings("brokers", c.conf.Brokers))

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
				c.logger.Error("Kafka FetchMessage error", zap.String("topic", topic), zap.Error(err))
				c.notify <- err
				return
			}

			if err := handler(ctx, msg); err != nil {
				c.logger.Error("Kafka Handler error", zap.String("topic", topic), zap.Error(err))
			}

			if err := reader.CommitMessages(ctx, msg); err != nil {
				c.logger.Error("Kafka CommitMessages error", zap.String("topic", topic), zap.Error(err))
			}
		}
	}
}

func (c *Consumer) Wait() {
	c.wg.Wait()
}

func NewProducer(conf *config.Config, topic string, l logger.ILogger) *Producer {
	opt := newOption(conf)

	writer := &kafka.Writer{
		Addr:         kafka.TCP(opt.Brokers...),
		Topic:        topic,
		RequiredAcks: opt.RequiredAcks,
		MaxAttempts:  opt.MaxAttempts,
	}

	l.Info("Kafka Producer created", zap.String("topic", topic), zap.Strings("brokers", opt.Brokers))

	return &Producer{
		writer: writer,
		logger: l,
	}
}

func (p *Producer) Produce(ctx context.Context, key, value []byte) error {
	err := p.writer.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	})
	if err != nil {
		p.logger.Error("Kafka Produce error", zap.Error(err))
	}
	return err
}

func (p *Producer) Close() error {
	p.logger.Info("Kafka Producer closed")
	return p.writer.Close()
}
