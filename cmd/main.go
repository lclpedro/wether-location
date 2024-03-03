package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/lclpedro/weather-location/configs"
	"github.com/lclpedro/weather-location/internal/scaffold/services"
	"github.com/lclpedro/weather-location/internal/scaffold/views"
)

func init() {
	configs.InitConfigs()
}

func main() {

	app := fiber.New()

	allServices := services.NewAllServices()
	app = views.NewAllHandlerViews(app, allServices)

	if err := app.Listen(":8080"); err != nil {
		log.Fatal("Error to up service", err)
	}
}
