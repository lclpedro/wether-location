package main

import (
	"github.com/SpyxBR/spyx-financial-control/internal/financial_control/views"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	app = views.NewAllHandlerViews(app)
	app.Listen(":8080")

}
