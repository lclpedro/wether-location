package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lclpedro/weather-location/configs"
	"github.com/lclpedro/weather-location/internal/scaffold/services"
	"github.com/lclpedro/weather-location/internal/scaffold/views"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	configs.InitConfigs()
	app := fiber.New()

	allServices := services.NewAllServices()
	app = views.NewAllHandlerViews(app, allServices)

	app.Listen(":8080")
}
