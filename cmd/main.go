package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/lclpedro/scaffold-golang-fiber/configs"
	"github.com/lclpedro/scaffold-golang-fiber/internal/scaffold/repositories"
	"github.com/lclpedro/scaffold-golang-fiber/internal/scaffold/services"
	"github.com/lclpedro/scaffold-golang-fiber/internal/scaffold/views"
	"github.com/lclpedro/scaffold-golang-fiber/pkg/mysql"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	configs.InitConfigs()
	app := fiber.New()

	databaseConfig := mysql.GetDatabaseConfiguration()
	read := mysql.InitMySQLConnection(databaseConfig[mysql.ReadOperation], mysql.ReadOperation)
	write := mysql.InitMySQLConnection(databaseConfig[mysql.WriteOperation], mysql.WriteOperation)
	connMysql, err := mysql.NewConnection(write, read)
	checkError(err)
	uowInstance := mysql.NewUnitOfWork(connMysql)
	repositories.RegistryRepositories(uowInstance, connMysql)
	allServices := services.NewAllServices(uowInstance)
	app = views.NewAllHandlerViews(app, allServices)

	app.Listen(":8080")
}
