package v1

import (
	"github/smile-ko/go-template/pkg/logger"
	"github/smile-ko/go-template/pkg/postgres"

	"github.com/gofiber/fiber/v2"
)

func NewV1Router(r fiber.Router, pg *postgres.Postgres, l logger.ILogger) {
	h := NewHandler()

	r.Get("/healthz", h.HealthCheck)
	r.Get("/hello", h.Hello)
}
