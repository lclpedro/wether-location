package views

import (
	"github.com/SpyxBR/spyx-financial-control/internal/financial_control/views/health"
	"github.com/gofiber/fiber/v2"
)

type AllViews struct {
	HealthView health.IHealthView
}

func newAllViews() *AllViews {
	return &AllViews{
		HealthView: health.NewHealthView(),
	}
}

func NewAllHandlerViews(app *fiber.App) *fiber.App {
	views := newAllViews()
	app.Get("/health", views.HealthView.HealthHandler)
	return app
}
