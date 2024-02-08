package health

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lclpedro/weather-location/internal/scaffold/services/health"
)

type View interface {
	HealthHandler(c *fiber.Ctx) error
}
type healthView struct {
	healthService health.Service
}

func NewHealthView(healthService health.Service) View {
	return &healthView{
		healthService: healthService,
	}
}

func (v healthView) HealthHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Application running..", "context": v.healthService.Ping(c.Context())})
}
