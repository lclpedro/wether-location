package services

import (
	"github.com/lclpedro/weather-location/internal/scaffold/services/health"
	weatherLocation "github.com/lclpedro/weather-location/internal/scaffold/services/weather_location"
)

type AllServices struct {
	HealthService          health.Service
	WeatherLocationService weatherLocation.Service
}

func NewAllServices() *AllServices {
	return &AllServices{
		HealthService:          health.NewHealthService(),
		WeatherLocationService: weatherLocation.NewService(),
	}
}
