package main

import (
	"github/smile-ko/go-template/config"
	"github/smile-ko/go-template/internal/app"
	"log"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(cfg)
}
