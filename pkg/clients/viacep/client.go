package viacep

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/lclpedro/weather-location/pkg/requester"
	"github.com/spf13/viper"
)

type (
	Client interface {
		GetAddress(ctx context.Context, cep string) (Response, error)
	}
	client struct {
		BaseURL   string
		Requester requester.Client
	}

	Response struct {
		Cep         string `json:"cep"`
		Logradouro  string `json:"logradouro"`
		Complemento string `json:"complemento"`
		Bairro      string `json:"bairro"`
		Localidade  string `json:"localidade"`
		Uf          string `json:"uf"`
		Err         bool   `json:"erro"`
	}
)

const ViaCEPTimeout = "application.clients.viacep.timeout.ms"

func NewClient(requester requester.Client) Client {
	baseURL := viper.GetString("application.clients.viacep.base_url")

	return &client{
		BaseURL:   baseURL,
		Requester: requester,
	}
}

var (
	ErrNotFound   = errors.New("viacep: can not find zipcode")
	ErrInvalidCep = errors.New("viacep: invalid zipcode")
)

func (c *client) cepIsValid(cep string) bool {
	return len(cep) == 8
}

func (c *client) GetAddress(ctx context.Context, cep string) (Response, error) {
	if !c.cepIsValid(cep) {
		fmt.Println("viacep: invalid zipcode")
		return Response{}, ErrInvalidCep
	}

	response, err := c.Requester.Get(fmt.Sprintf("%s/%s/json", c.BaseURL, cep))
	if err != nil {
		return Response{}, err
	}

	if response.StatusCode != http.StatusOK {
		return Response{}, fmt.Errorf("viacep: %w, status_code: %d", ErrNotFound, response.StatusCode)
	}

	var resp Response
	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		return Response{}, err
	}

	if resp.Err {
		return Response{}, ErrNotFound
	}

	return resp, nil
}
