package kafkabus

import (
	"github/smile-ko/go-template/config"
	"time"

	"github.com/segmentio/kafka-go"
)

type option struct {
	Brokers      []string
	GroupID      string
	MinBytes     int
	MaxBytes     int
	MaxWait      time.Duration
	StartOffset  int64
	RequiredAcks kafka.RequiredAcks
	MaxAttempts  int
}

func newOption(conf *config.Config) *option {
	maxWait, _ := time.ParseDuration(conf.Kafka.Consumer.MaxWait)

	startOffset := kafka.FirstOffset
	if conf.Kafka.Consumer.StartOffset == "last" {
		startOffset = kafka.LastOffset
	}

	requiredAcks := kafka.RequireAll
	if conf.Kafka.Producer.RequiredAcks == 0 {
		requiredAcks = kafka.RequireNone
	}

	return &option{
		Brokers:      conf.Kafka.Brokers,
		GroupID:      conf.Kafka.GroupID,
		MinBytes:     conf.Kafka.Consumer.MinBytes,
		MaxBytes:     conf.Kafka.Consumer.MaxBytes,
		MaxWait:      maxWait,
		StartOffset:  startOffset,
		RequiredAcks: requiredAcks,
	}
}
