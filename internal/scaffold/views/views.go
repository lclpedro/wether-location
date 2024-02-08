package views

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lclpedro/weather-location/internal/scaffold/services"
	"github.com/lclpedro/weather-location/internal/scaffold/views/health"
	weatherlocation "github.com/lclpedro/weather-location/internal/scaffold/views/weather_location"
)

type AllViews struct {
	HealthView  health.View
	WeatherView weatherlocation.View
}

func newAllViews(services *services.AllServices) *AllViews {
	return &AllViews{
		HealthView:  health.NewHealthView(services.HealthService),
		WeatherView: weatherlocation.NewView(services.WeatherLocationService),
	}
}

func NewAllHandlerViews(app *fiber.App, services *services.AllServices) *fiber.App {
	views := newAllViews(services)
	app.Get("/health", views.HealthView.HealthHandler)
	app.Get("/weather/:cep", views.WeatherView.WeatherLocationHandler)
	return app
}
