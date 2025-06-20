package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func NotFound() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Not Found",
			"message": fmt.Sprintf("Route %s %s not found", c.Method(), c.OriginalURL()),
		})
	}
}
