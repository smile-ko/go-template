package main

import (
	"fmt"
	"github/smile-ko/go-template/config"
	"github/smile-ko/go-template/internal/interfaces/http"
	"github/smile-ko/go-template/pkg/httpserver"
	"github/smile-ko/go-template/pkg/logger"
	"github/smile-ko/go-template/pkg/postgres"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	l := logger.New(cfg.Log.Level)

	// Init Postgres
	pg := postgres.NewOrGetSingleton(cfg)
	defer pg.Close()

	// Init HTTP server
	httpServer := httpserver.New(
		httpserver.Port(cfg.HTTP.Port),
		httpserver.Prefork(cfg.HTTP.UsePreforkMode),
	)

	// Register routes (router.go)
	http.NewRouter(httpServer.App, cfg, pg, l)

	// Start HTTP Server
	httpServer.Start()

	// Handle graceful shutdown
	waitForShutdown(httpServer, l)
}

func waitForShutdown(httpServer *httpserver.Server, l logger.Interface) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	var err error
	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
