package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/lclpedro/scaffold-golang-fiber/internal/scaffold/repositories"
	"github.com/lclpedro/scaffold-golang-fiber/internal/scaffold/services"
	"github.com/lclpedro/scaffold-golang-fiber/internal/scaffold/views"
	uow "github.com/lclpedro/scaffold-golang-fiber/pkg/unity_of_work"
)

func main() {
	app := fiber.New()

	// TODO: Adding mysql connection factory
	dbConnection, err := sqlx.Open("mysql", "root:root@tcp(localhost:3306)/scaffold")
	if err != nil {
		panic(err)
	}

	uowInstance := uow.NewUnityOfWork(dbConnection)
	repositories.RegistryRepositories(uowInstance, dbConnection)
	allServices := services.NewAllServices(uowInstance)
	app = views.NewAllHandlerViews(app, allServices)

	app.Listen(":8080")
}
