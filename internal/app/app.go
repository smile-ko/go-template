package app

import (
	"context"
	"fmt"
	"github/smile-ko/go-template/config"
	"github/smile-ko/go-template/internal/controller/grpc"
	"github/smile-ko/go-template/internal/controller/http"
	grpcserver "github/smile-ko/go-template/pkg/grpc"
	"github/smile-ko/go-template/pkg/httpserver"
	kafkabus "github/smile-ko/go-template/pkg/kafka"
	"github/smile-ko/go-template/pkg/logger"
	"github/smile-ko/go-template/pkg/postgres"
	"os"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Init Postgres
	pg := postgres.NewOrGetSingleton(cfg)
	defer pg.Close()

	// Init HTTP server
	httpServer := httpserver.New(
		httpserver.Port(cfg.HTTP.Port),
		httpserver.Prefork(cfg.HTTP.UsePreforkMode),
	)

	// Register routes
	http.NewRouter(httpServer.App, cfg, pg, l)

	// gRPC Server
	grpcServer := grpcserver.New(grpcserver.Port(cfg.GRPC.Port))
	grpc.RegisterGRPCServices(grpcServer.App, pg, l)

	// Kafka consumer
	kafkaConsumer := kafkabus.NewConsumer(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	kafkaTopic := map[string]kafkabus.MessageHandler{
		"demo-topic": func(ctx context.Context, msg kafka.Message) error {
			fmt.Printf("Kafka received: %s", string(msg.Value))
			return nil
		},
	}

	for topic, handler := range kafkaTopic {
		go kafkaConsumer.Handler(ctx, topic, handler)
	}

	// Start HTTP Server
	httpServer.Start()
	grpcServer.Start()

	// Handle graceful shutdown
	waitForShutdown(cancel, kafkaConsumer, httpServer, grpcServer, l)
}

func waitForShutdown(cancel context.CancelFunc, kafkaConsumer *kafkabus.Consumer, httpServer *httpserver.Server, grpcServer *grpcserver.Server, l logger.Interface) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	var err error
	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
		cancel()
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
		cancel()
	case err = <-grpcServer.Notify():
		l.Error(fmt.Errorf("app - Run - grpcServer.Notify: %w", err))
		cancel()
	case err = <-kafkaConsumer.Notify():
		l.Error(fmt.Errorf("app - Run - kafkaConsumer.Notify: %w", err))
		cancel()
	}

	_ = httpServer.Shutdown()
	_ = grpcServer.Shutdown()

	kafkaConsumer.Wait()

	l.Info("Shutdown complete")
}
