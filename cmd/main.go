package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lclpedro/scaffold-golang-fiber/internal/scaffold/views"
)

func main() {
	app := fiber.New()
	app = views.NewAllHandlerViews(app)
	app.Listen(":8080")
}
