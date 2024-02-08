package weatherlocation

import (
	"github.com/lclpedro/weather-location/internal/scaffold/domains"
	"github.com/lclpedro/weather-location/pkg/viacep"
	"github.com/lclpedro/weather-location/pkg/weather"
)

type Service interface {
	GetWeatherLocation(cep string) (Output, error)
	SetClients(viacepClient viacep.Client, weatherClient weather.Client)
}

type service struct {
	weatherClient weather.Client
	viacepClient  viacep.Client
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

func NewService() Service {
	return &service{}
}

func (s *service) SetClients(viacepClient viacep.Client, weatherClient weather.Client) {
	s.viacepClient = viacepClient
	s.weatherClient = weatherClient
}

func (s *service) GetWeatherLocation(cep string) (Output, error) {
	address, err := s.viacepClient.GetAddress(cep)
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
