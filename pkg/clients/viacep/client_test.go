package viacep_test

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/lclpedro/weather-location/pkg/clients/viacep"
	"github.com/lclpedro/weather-location/pkg/requester"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

var mockRequester *requester.MockRequester

func setup() {
	mockRequester = new(requester.MockRequester)
	viper.Set("application.clients.viacep.base_url", "http://viacep.com.br/ws")
}

func TestClient_GetAddress(t *testing.T) {
	expectedResponse := `{"cep":"12345678","logradouro":"Rua XPTO","complemento":"","bairro":"Vila","localidade":"São Paulo","uf":"SP","ibge":"000","gia":"000","ddd":"11","siafi":"00"}`
	t.Run("should return invalid cep", func(t *testing.T) {
		setup()
		client := viacep.NewClient(mockRequester)
		_, err := client.GetAddress("12345")
		assert.Error(t, err)
		assert.Equal(t, viacep.ErrInvalidCep, err)
	})

	t.Run("should return error in request", func(t *testing.T) {
		setup()
		client := viacep.NewClient(mockRequester)
		mockRequester.On("Get", "http://viacep.com.br/ws/12345678/json").Return(&http.Response{}, fmt.Errorf("error"))
		_, err := client.GetAddress("12345678")
		assert.Error(t, err)
	})

	t.Run("should return status code 400", func(t *testing.T) {
		setup()
		client := viacep.NewClient(mockRequester)
		response := io.NopCloser(bytes.NewBufferString(expectedResponse))
		mockRequester.On("Get", "http://viacep.com.br/ws/12345678/json").Return(&http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       response,
		}, nil)
		_, err := client.GetAddress("12345678")
		assert.Error(t, err)
		assert.Equal(t, true, errors.Is(err, viacep.ErrNotFound))
	})

	t.Run("should return status code 200 invalid json", func(t *testing.T) {
		setup()
		client := viacep.NewClient(mockRequester)
		response := io.NopCloser(bytes.NewBufferString("{message:Invalid payload}"))
		mockRequester.On("Get", "http://viacep.com.br/ws/12345678/json").Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       response,
		}, nil)
		_, err := client.GetAddress("12345678")
		assert.Error(t, err)
	})

	t.Run("should return status code 200 return error", func(t *testing.T) {
		setup()
		client := viacep.NewClient(mockRequester)
		response := io.NopCloser(bytes.NewBufferString(`{"erro":true}`))
		mockRequester.On("Get", "http://viacep.com.br/ws/12345678/json").Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       response,
		}, nil)
		_, err := client.GetAddress("12345678")
		assert.Error(t, err)
	})

	t.Run("should return status code 200", func(t *testing.T) {
		setup()
		client := viacep.NewClient(mockRequester)
		response := io.NopCloser(bytes.NewBufferString(expectedResponse))
		mockRequester.On("Get", "http://viacep.com.br/ws/12345678/json").Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       response,
		}, nil)
		_, err := client.GetAddress("12345678")
		assert.NoError(t, err)
	})
}
