package health

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lclpedro/scaffold-golang-fiber/internal/scaffold/services/health"
	"net/http"
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
	err := v.healthService.Ping(c.Context())
	if err != nil {
		c.Status(http.StatusBadRequest)
		return c.JSON(fiber.Map{"message": "Application with error running..", "error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Application running..", "error": nil})
}
