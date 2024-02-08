package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/spf13/viper"
)

type (
	Client interface {
		GetWeather(city, state string) (Response, error)
	}
	client struct {
		BaseURL   string
		TimeoutMS int
		ApiKey    string
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

func NewClient(configs *viper.Viper) Client {
	baseURL := viper.GetString("application.clients.weather.base_url")
	timeout := viper.GetInt("application.clients.weather.timeout.ms")
	apiKey := os.Getenv(viper.GetString("application.clients.weather.api_key"))
	return &client{
		BaseURL:   baseURL,
		TimeoutMS: timeout,
		ApiKey:    apiKey,
	}
}

func (c *client) GetWeather(city, state string) (Response, error) {
	if c.ApiKey == "" {
		return Response{}, fmt.Errorf("weather: api key not found")
	}

	http.DefaultClient.Timeout = time.Duration(c.TimeoutMS) * time.Millisecond
	urlString := fmt.Sprintf("%s-%s", city, state)
	cityParse, _ := url.Parse(urlString)
	url := fmt.Sprintf("%s/current.json?key=%s&q=%s&aqi=no&lang=pt", c.BaseURL, c.ApiKey, cityParse)
	fmt.Println(url)
	response, err := http.Get(url)
	if err != nil {
		return Response{}, err
	}

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		fmt.Println(response.StatusCode, string(body))

		return Response{}, fmt.Errorf("weather: status code %d", response.StatusCode)
	}

	var resp Response
	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		return Response{}, err
	}

	return resp, nil
}
