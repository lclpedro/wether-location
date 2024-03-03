package services

import (
	"github.com/lclpedro/weather-location/internal/scaffold/services/health"
	weatherLocation "github.com/lclpedro/weather-location/internal/scaffold/services/weather_location"
	"go.opentelemetry.io/otel/trace"
)

type AllServices struct {
	HealthService          health.Service
	WeatherLocationService weatherLocation.Service
}

func NewAllServices(trace trace.Tracer) *AllServices {
	return &AllServices{
		HealthService:          health.NewHealthService(),
		WeatherLocationService: weatherLocation.NewService(trace),
	}
}
