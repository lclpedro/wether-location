package views

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lclpedro/scaffold-golang-fiber/internal/scaffold/services"
	"github.com/lclpedro/scaffold-golang-fiber/internal/scaffold/views/health"
)

type AllViews struct {
	HealthView health.HealthView
}

func newAllViews(services *services.AllServices) *AllViews {
	return &AllViews{
		HealthView: health.NewHealthView(services.HealthService),
	}
}

func NewAllHandlerViews(app *fiber.App, services *services.AllServices) *fiber.App {
	views := newAllViews(services)
	app.Get("/health", views.HealthView.HealthHandler)
	return app
}
