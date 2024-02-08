package services

import (
	"github.com/lclpedro/weather-location/internal/scaffold/services/health"
	weatherlocation "github.com/lclpedro/weather-location/internal/scaffold/services/weather_location"
)

type AllServices struct {
	HealthService          health.Service
	WeatherLocationService weatherlocation.Service
}

func NewAllServices() *AllServices {
	return &AllServices{
		HealthService:          health.NewHealthService(),
		WeatherLocationService: weatherlocation.NewService(),
	}
}
