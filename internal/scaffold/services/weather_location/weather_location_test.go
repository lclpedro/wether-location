package weatherlocation_test

import (
	"errors"
	weatherLocation "github.com/lclpedro/weather-location/internal/scaffold/services/weather_location"
	"github.com/lclpedro/weather-location/pkg/clients/viacep"
	"github.com/lclpedro/weather-location/pkg/clients/weather"
	"github.com/stretchr/testify/suite"
	"testing"
)

type suiteService struct {
	suite.Suite

	viaCEPMock  *viacep.MockViaCEP
	weatherMock *weather.MockWeather
}

func TestRunSuite(t *testing.T) {
	suite.Run(t, new(suiteService))
}

func (s *suiteService) SetupSubTest() {
	s.viaCEPMock = new(viacep.MockViaCEP)
	s.weatherMock = new(weather.MockWeather)
}

func (s *suiteService) TestGetWeatherLocation() {
	s.Run("should return error in get address", func() {
		s.viaCEPMock.On("GetAddress", "12345678").Return(viacep.Response{}, errors.New("error"))

		service := weatherLocation.NewService()
		service.SetClients(
			s.viaCEPMock,
			s.weatherMock,
		)

		_, err := service.GetWeatherLocation("12345678")
		s.Error(viacep.ErrNotFound, err)
	})

	s.Run("should return error in get weather by address", func() {
		s.viaCEPMock.On("GetAddress", "12345678").Return(viacep.Response{
			Cep:         "12345678",
			Logradouro:  "Rua Mock",
			Complemento: "Mock",
			Bairro:      "Mock",
			Localidade:  "São Paulo",
			Uf:          "SP",
			Err:         false,
		}, nil)

		s.weatherMock.On("GetWeather", "São Paulo", "SP").Return(weather.Response{}, errors.New("error"))

		service := weatherLocation.NewService()
		service.SetClients(
			s.viaCEPMock,
			s.weatherMock,
		)

		_, err := service.GetWeatherLocation("12345678")
		s.Error(viacep.ErrNotFound, err)
	})
}
