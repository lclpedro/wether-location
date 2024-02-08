package viacep

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

type (
	Client interface {
		GetAddress(cep string) (Response, error)
	}
	client struct {
		BaseURL   string
		TimeoutMS int
	}

	Response struct {
		Cep         string `json:"cep"`
		Logradouro  string `json:"logradouro"`
		Complemento string `json:"complemento"`
		Bairro      string `json:"bairro"`
		Localidade  string `json:"localidade"`
		Uf          string `json:"uf"`
	}
)

func NewClient(configs *viper.Viper) Client {
	baseURL := viper.GetString("application.clients.viacep.base_url")
	timeout := viper.GetInt("application.clients.viacep.timeout.ms")
	return &client{
		BaseURL:   baseURL,
		TimeoutMS: timeout,
	}
}

var (
	ErrNotFound   = errors.New("viacep: can not find zipcode")
	ErrInvalidCep = errors.New("viacep: invalid zipcode")
)

func (c *client) cepIsValid(cep string) bool {
	return len(cep) == 8
}

func (c *client) GetAddress(cep string) (Response, error) {
	if !c.cepIsValid(cep) {
		fmt.Println("viacep: invalid zipcode")
		return Response{}, ErrInvalidCep
	}
	http.DefaultClient.Timeout = time.Duration(c.TimeoutMS) * time.Millisecond
	response, err := http.Get(fmt.Sprintf("%s/%s/json", c.BaseURL, cep))
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

	return resp, nil
}
