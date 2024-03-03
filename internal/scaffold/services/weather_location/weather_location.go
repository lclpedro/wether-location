package weatherlocation

import (
	"context"
	"github.com/lclpedro/weather-location/internal/scaffold/domains"
	viaCepClient "github.com/lclpedro/weather-location/pkg/clients/viacep"
	weatherAPIClient "github.com/lclpedro/weather-location/pkg/clients/weather"
	"go.opentelemetry.io/otel/trace"
)

type Service interface {
	GetWeatherLocation(ctx context.Context, cep string) (Output, error)
	SetClients(viaCepClient viaCepClient.Client, weatherClient weatherAPIClient.Client)
}

type service struct {
	weatherClient weatherAPIClient.Client
	viaCepClient  viaCepClient.Client
	trace         trace.Tracer
}

type Output struct {
	TempC      float64 `json:"temp_C"`
	TempK      float64 `json:"temp_K"`
	TempF      float64 `json:"temp_F"`
	CityName   string  `json:"city_name"`
	Region     string  `json:"region"`
	Country    string  `json:"country"`
	LastUpdate string  `json:"last_update"`
}

func NewService(trace trace.Tracer) Service {
	return &service{
		trace: trace,
	}
}

func (s *service) SetClients(viacepClient viaCepClient.Client, weatherClient weatherAPIClient.Client) {
	s.viaCepClient = viacepClient
	s.weatherClient = weatherClient
}

func (s *service) GetWeatherLocation(ctx context.Context, cep string) (Output, error) {

	address, err := s.viaCepClient.GetAddress(ctx, cep)
	if err != nil {
		return Output{}, err
	}

	weather, err := s.weatherClient.GetWeather(address.Localidade, address.Uf)
	if err != nil {
		return Output{}, err
	}

	weatherDomain := domains.NewWeather(
		weather.Location.CityName,
		weather.Location.Region,
		weather.Location.Country,
		weather.Current.Temperature,
		weather.Current.LastUpdated,
	)

	return Output{
		TempC:      weatherDomain.GetCelcius(),
		TempK:      weatherDomain.GetKelvin(),
		TempF:      weatherDomain.GetFahrenheit(),
		CityName:   weatherDomain.City,
		Region:     weatherDomain.State,
		Country:    weatherDomain.Country,
		LastUpdate: weatherDomain.LastUpdated,
	}, nil
}
