package views

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lclpedro/weather-location/internal/scaffold/services"
	"github.com/lclpedro/weather-location/internal/scaffold/views/health"
	weatherlocation "github.com/lclpedro/weather-location/internal/scaffold/views/weather_location"
	"go.opentelemetry.io/otel/trace"
)

type AllViews struct {
	HealthView  health.View
	WeatherView weatherlocation.View
}

func newAllViews(services *services.AllServices, tracer trace.Tracer) *AllViews {
	return &AllViews{
		HealthView:  health.NewHealthView(services.HealthService),
		WeatherView: weatherlocation.NewView(tracer, services.WeatherLocationService),
	}
}

func NewAllHandlerViews(app *fiber.App, tracer trace.Tracer, services *services.AllServices) *fiber.App {
	views := newAllViews(services, tracer)
	app.Get("/health", views.HealthView.HealthHandler)
	app.Get("/weather/:cep", views.WeatherView.WeatherLocationHandler)
	return app
}
