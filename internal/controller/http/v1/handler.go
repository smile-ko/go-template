package v1

import (
	"github.com/gofiber/fiber/v2"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// Hello godoc
// @Summary  test api
// @Tags     internal
// @Success  200
// @Failure  500
// @Router   /hello [get]
func (h *Handler) Hello(c *fiber.Ctx) error {
	name := c.Query("name", "world")
	return c.JSON(fiber.Map{
		"message": "Hello " + name,
	})
}
