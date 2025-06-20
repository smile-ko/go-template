package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		App     App
		HTTP    HTTP
		Log     Log
		PG      PG
		GRPC    GRPC
		Metrics Metrics
		Swagger Swagger
		Kafka   Kafka
	}

	App struct {
		Name    string `env:"APP_NAME,required"`
		Version string `env:"APP_VERSION,required"`
		EnvName string `env:"ENV_NAME,required"`
	}

	HTTP struct {
		Port           string `env:"HTTP_PORT,required"`
		UsePreforkMode bool   `env:"HTTP_USE_PREFORK_MODE" envDefault:"false"`
	}

	Log struct {
		Level         string `env:"LOG_LEVEL,required"`
		FileLogName   string `env:"LOG_FILE_LOG_NAME"`
		MaxSize       int    `env:"MAX_SIZE"`
		MaxBackups    int    `env:"LOG_MAX_BACKUPS"`
		MaxAge        int    `env:"LOG_MAX_AGE"`
		Compress      bool   `env:"LOG_COMPRESS"`
		ConsoleOutput bool   `env:"LOG_CONSOLE_OUTPUT" default:"true"`
		UseJSON       bool   `env:"LOG_USE_JSON" default:"false"`
	}

	PG struct {
		PoolMax int    `env:"PG_POOL_MAX,required"`
		URL     string `env:"PG_URL,required"`
	}

	GRPC struct {
		Port string `env:"GRPC_PORT,required"`
	}

	Metrics struct {
		Enabled bool `env:"METRICS_ENABLED" envDefault:"true"`
	}

	Swagger struct {
		Enabled bool `env:"SWAGGER_ENABLED" envDefault:"false"`
	}

	Kafka struct {
		Brokers  []string `env:"KAFKA_BROKERS,required"`
		GroupID  string   `env:"KAFKA_GROUP_ID,required"`
		Producer Producer
		Consumer Consumer
	}

	Producer struct {
		RequiredAcks int    `env:"KAFKA_PRODUCER_REQUIRED_ACKS" envDefault:"-1"`
		Async        bool   `env:"KAFKA_PRODUCER_ASYNC" envDefault:"false"`
		Compression  string `env:"KAFKA_PRODUCER_COMPRESSION" envDefault:"none"`
	}

	Consumer struct {
		MinBytes    int    `env:"KAFKA_CONSUMER_MIN_BYTES" envDefault:"10000"`
		MaxBytes    int    `env:"KAFKA_CONSUMER_MAX_BYTES" envDefault:"10000000"`
		MaxWait     string `env:"KAFKA_CONSUMER_MAX_WAIT" envDefault:"1s"`
		StartOffset string `env:"KAFKA_CONSUMER_START_OFFSET" envDefault:"latest"`
	}
)

func NewConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
