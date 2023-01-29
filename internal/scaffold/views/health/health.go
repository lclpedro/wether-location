package health

import (
	"github.com/gofiber/fiber/v2"
)

type IHealthView interface {
	HealthHandler(c *fiber.Ctx) error
}
type healthView struct{}

func NewHealthView() *healthView {
	return &healthView{}
}

func (v healthView) HealthHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Application running.."})
}
