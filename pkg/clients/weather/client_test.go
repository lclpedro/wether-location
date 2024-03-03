package weather_test

import (
	"bytes"
	"fmt"
	"github.com/lclpedro/weather-location/pkg/clients/weather"
	"github.com/lclpedro/weather-location/pkg/requester"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"os"
	"testing"
)

var mockRequester *requester.MockRequester

func setup() {
	mockRequester = new(requester.MockRequester)
	viper.Set("application.clients.weather.base_url", "http://api.weatherapi.com/v1")
	viper.Set("application.clients.weather.api_key", "WEATHER_API_KEY")
	_ = os.Setenv("WEATHER_API_KEY", "123")
}

const expectedResponse = `{"location":{"name":"San Paulo","region":"Sao Paulo","country":"Brazil","lat":-23.53,"lon":-46.62,"tz_id":"America/Sao_Paulo","localtime_epoch":1708044422,"localtime":"2024-02-15 21:47"},"current":{"last_updated_epoch":1708044300,"last_updated":"2024-02-15 21:45","temp_c":18.7,"temp_f":65.7,"is_day":0,"condition":{"text":"Patchy rain nearby","icon":"//cdn.weatherapi.com/weather/64x64/night/176.png","code":1063},"wind_mph":5.4,"wind_kph":8.6,"wind_degree":135,"wind_dir":"SE","pressure_mb":1018.0,"pressure_in":30.06,"precip_mm":0.05,"precip_in":0.0,"humidity":93,"cloud":99,"feelslike_c":18.7,"feelslike_f":65.7,"vis_km":10.0,"vis_miles":6.0,"uv":1.0,"gust_mph":6.5,"gust_kph":10.4}}`

func TestClient_GetWeather(t *testing.T) {
	t.Run("should return error api key weather", func(t *testing.T) {
		setup()
		_ = os.Setenv("WEATHER_API_KEY", "")
		client := weather.NewClient(mockRequester)
		response := io.NopCloser(bytes.NewBufferString(expectedResponse))

		mockRequester.On("Get",
			"http://api.weatherapi.com/v1/current.json?key=123&q=Sao Paulo-BR&aqi=no&lang=pt",
		).Return(
			&http.Response{
				StatusCode: 200,
				Body:       response,
			}, nil)

		_, err := client.GetWeather("Sao Paulo", "BR")
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "weather: api key not found")
	})

	t.Run("should return error when request fails", func(t *testing.T) {
		setup()
		client := weather.NewClient(mockRequester)
		mockRequester.On(
			"Get",
			"http://api.weatherapi.com/v1/current.json?key=123&q=Sao%20Paulo-BR&aqi=no&lang=pt").
			Return(&http.Response{}, fmt.Errorf("error"))
		_, err := client.GetWeather("Sao Paulo", "BR")
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "error")
	})

	t.Run("should return status code 400", func(t *testing.T) {
		setup()
		client := weather.NewClient(mockRequester)
		response := io.NopCloser(bytes.NewBufferString(expectedResponse))

		mockRequester.On("Get",
			"http://api.weatherapi.com/v1/current.json?key=123&q=Sao%20Paulo-BR&aqi=no&lang=pt",
		).Return(
			&http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       response,
			}, nil)

		_, err := client.GetWeather("Sao Paulo", "BR")
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "weather: status code 400")
	})

	t.Run("should return status code 200 - invalid payload", func(t *testing.T) {
		setup()
		client := weather.NewClient(mockRequester)
		response := io.NopCloser(bytes.NewBufferString("{message:Invalid payload}"))

		mockRequester.On("Get",
			"http://api.weatherapi.com/v1/current.json?key=123&q=Sao%20Paulo-BR&aqi=no&lang=pt",
		).Return(
			&http.Response{
				StatusCode: http.StatusOK,
				Body:       response,
			}, nil)

		_, err := client.GetWeather("Sao Paulo", "BR")
		assert.Error(t, err)
	})

	t.Run("should return status code 200", func(t *testing.T) {
		setup()
		client := weather.NewClient(mockRequester)
		response := io.NopCloser(bytes.NewBufferString(expectedResponse))

		mockRequester.On("Get",
			"http://api.weatherapi.com/v1/current.json?key=123&q=Sao%20Paulo-BR&aqi=no&lang=pt",
		).Return(
			&http.Response{
				StatusCode: http.StatusOK,
				Body:       response,
			}, nil)

		_, err := client.GetWeather("Sao Paulo", "BR")
		assert.NoError(t, err)
	})

}
