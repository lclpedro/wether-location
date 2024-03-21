package recursiveapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/lclpedro/weather-location/pkg/requester"
	"github.com/spf13/viper"
)

type (
	ResponseAPI struct {
		City   string  `json:"city_name"`
		Temp_C float64 `json:"temp_C"`
		Temp_F float64 `json:"temp_F"`
		Temp_K float64 `json:"temp_K"`
	}

	Client interface {
		GetWeatherLocation(ctx context.Context, cep string) (ResponseAPI, error)
	}

	client struct {
		BaseURL   string
		Requester requester.Client
	}
)

const Timeout = "application.clients.recursiveapi.timeout.ms"

func NewClient(requester requester.Client) Client {
	baseURL := viper.GetString("application.clients.recursiveapi.base_url")

	return &client{
		BaseURL:   baseURL,
		Requester: requester,
	}
}

var (
	ErrNotFoundAddress = errors.New("can not find zipcode")
	ErrInvalidCep      = errors.New("invalid zipcode")
	ErrInternalError   = errors.New("internal error")
)

func (c *client) GetWeatherLocation(ctx context.Context, cep string) (ResponseAPI, error) {
	url := fmt.Sprintf("%s/weather/%s", c.BaseURL, cep)
	response, err := c.Requester.Get(url)

	if err != nil {
		return ResponseAPI{}, ErrInternalError
	}

	if response.StatusCode == http.StatusNotFound {
		return ResponseAPI{}, ErrNotFoundAddress
	}

	if response.StatusCode == http.StatusUnprocessableEntity {
		return ResponseAPI{}, ErrInvalidCep
	}

	var resp ResponseAPI
	if err = json.NewDecoder(response.Body).Decode(&resp); err != nil {
		log.Fatal(err)
		return ResponseAPI{}, err
	}

	return resp, nil
}
