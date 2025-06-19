package v2

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(r fiber.Router) {
	h := NewHandler()

	r.Get("/healthz", h.HealthCheck)
	r.Get("/hello", h.Hello)
}
