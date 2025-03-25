package main

import (
	_ "authservice/docs"
	"os"

	config "authservice/config"
	routes "authservice/routes"

	"github.com/gofiber/fiber/v2"
)

// @title           Auth service
// @version         1.0
// @description     Auth server API

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

func main() {

	// Loading env vars from .env
	config.LoadEnvs()

	// Creating Log File
	logFile, err := os.OpenFile(config.Envs.LogFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	// Configure Logger
	config.CreateLogger()
	defer config.Logger.Sync() // Syncing all logs at the end of program
	defer config.Logger.Info("Program end")

	// Reading RSA keys
	config.ReadKeys()

	// Serving App
	app := fiber.New(fiber.Config{
		AppName: "Auth Service",
	})

	// Connect to DB
	config.ConncetDB()

	// Collect routes
	routes.Handlers(app)

	app.Listen(":8080")

}
