package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/lclpedro/weather-location/pkg/requester"
	"github.com/spf13/viper"
)

type (
	Client interface {
		GetWeather(city, state string) (Response, error)
	}
	client struct {
		BaseURL   string
		ApiKey    string
		Requester requester.Client
	}

	Location struct {
		CityName  string `json:"name"`
		Region    string `json:"region"`
		Country   string `json:"country"`
		Localtime string `json:"localtime"`
	}

	Current struct {
		LastUpdated string  `json:"last_updated"`
		Temperature float64 `json:"temp_c"`
	}

	Response struct {
		Location Location `json:"location"`
		Current  Current  `json:"current"`
	}
)

const WeatherAPITimeout = "application.clients.weather.timeout.ms"

func NewClient(requester requester.Client) Client {
	baseURL := viper.GetString("application.clients.weather.base_url")
	apiKey := os.Getenv(viper.GetString("application.clients.weather.api_key"))
	return &client{
		BaseURL:   baseURL,
		ApiKey:    apiKey,
		Requester: requester,
	}
}

func (c *client) GetWeather(city, state string) (Response, error) {
	if c.ApiKey == "" {
		return Response{}, fmt.Errorf("weather: api key not found")
	}
	urlString := fmt.Sprintf("%s-%s", city, state)
	cityParse, _ := url.Parse(urlString)
	url := fmt.Sprintf("%s/current.json?key=%s&q=%s&aqi=no&lang=pt", c.BaseURL, c.ApiKey, cityParse)
	response, err := c.Requester.Get(url)
	if err != nil {
		return Response{}, err
	}

	if response.StatusCode != http.StatusOK {
		return Response{}, fmt.Errorf("weather: status code %d", response.StatusCode)
	}

	var resp Response
	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		return Response{}, err
	}

	return resp, nil
}
