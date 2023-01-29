package views

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lclpedro/scaffold-golang-fiber/internal/scaffold/views/health"
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
