package weatherlocation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	recursiveapi "github.com/lclpedro/weather-location/pkg/clients/recursive-api"
	"github.com/lclpedro/weather-location/pkg/clients/viacep"
	weatherClient "github.com/lclpedro/weather-location/pkg/clients/weather"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/gofiber/fiber/v2"
	weatherlocation "github.com/lclpedro/weather-location/internal/scaffold/services/weather_location"
	"github.com/lclpedro/weather-location/pkg/requester"
	"github.com/spf13/viper"
)

type view struct {
	WeatherLocationService weatherlocation.Service
	Configs                *viper.Viper
	Tracer                 trace.Tracer
}

type View interface {
	WeatherLocationHandler(*fiber.Ctx) error
	WeatherLocationPostHandler(*fiber.Ctx) error
}

func NewView(tracer trace.Tracer, weatherLocationService weatherlocation.Service) View {
	return &view{
		WeatherLocationService: weatherLocationService,
		Configs:                viper.GetViper(),
		Tracer:                 tracer,
	}
}

func (v *view) WeatherLocationHandler(c *fiber.Ctx) error {
	carrier := propagation.HeaderCarrier(c.GetReqHeaders())
	ctx := context.Background()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	ctx, spam := v.Tracer.Start(ctx, "weather-location-api")
	defer spam.End()

	requesterViaCep := requester.NewRequester(ctx, v.Tracer, v.Configs.GetInt(viacep.ViaCEPTimeout))
	requesterWeather := requester.NewRequester(ctx, v.Tracer, v.Configs.GetInt(weatherClient.WeatherAPITimeout))

	v.WeatherLocationService.SetClients(
		viacep.NewClient(requesterViaCep),
		weatherClient.NewClient(requesterWeather),
	)

	cep := c.Params("cep")
	weather, err := v.WeatherLocationService.GetWeatherLocation(ctx, cep)
	if errors.Is(err, viacep.ErrInvalidCep) {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": err.Error()})
	}

	if errors.Is(err, viacep.ErrNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": err.Error()})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(weather)
}

type location struct {
	Cep string `json:"cep"`
}

func (v *view) WeatherLocationPostHandler(c *fiber.Ctx) error {
	carrier := propagation.HeaderCarrier(c.GetReqHeaders())
	ctx := context.Background()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	ctx, spam := v.Tracer.Start(ctx, "weather-location-recursive-api")
	defer spam.End()

	var location location
	err := json.Unmarshal(c.Body(), &location)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": err.Error()})
	}

	if location.Cep == "" || len(location.Cep) != 8 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": "Invalid CEP"})
	}

	_requester := requester.NewRequester(ctx, v.Tracer, v.Configs.GetInt(recursiveapi.Timeout))
	client := recursiveapi.NewClient(_requester)

	response, err := client.GetWeatherLocation(ctx, location.Cep)
	fmt.Println(response)
	fmt.Println(err)

	if errors.Is(err, recursiveapi.ErrInvalidCep) {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": err.Error()})
	}

	if errors.Is(err, recursiveapi.ErrNotFoundAddress) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": err.Error()})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{
		"city":   response.City,
		"temp_C": response.Temp_C,
		"temp_F": response.Temp_F,
		"temp_K": response.Temp_K,
	})
}
