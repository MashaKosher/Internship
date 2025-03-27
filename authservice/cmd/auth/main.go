package main

import (
	_ "authservice/docs"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"

	"authservice/internal/config"
	"authservice/internal/db"
	"authservice/internal/handler"
	"authservice/internal/keys"
	"authservice/internal/logger"
)

// @title						Auth service
// @version					1.0
// @description				Auth server API
// @host						localhost:8080
// @BasePath					/
// @securityDefinitions.basic	BasicAuth
// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {

	config.LoadEnvs()

	fmt.Println("File Log Name: " + config.Cfg.LogFileName)

	// Creating Log File
	logFile, err := os.OpenFile(config.Cfg.LogFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic("Log error: " + err.Error())
	}
	defer logFile.Close()

	logger.CreateLogger()
	defer logger.Logger.Sync()
	defer logger.Logger.Info("Program end")

	keys.ReadRSAKeys()

	app := fiber.New(fiber.Config{
		AppName: "Auth Service",
	})

	db.ConncetDB()

	handler.Handlers(app)

	// app.Listen(":" + string(config.Cfg.Server.Port))
	app.Listen(":8080")

}

// package main

// import (
// 	"os"
// 	"authservice/internal/config"
// 	"authservice/internal/keys"
// 	"authservice/internal/logger"
// )

// func main() {
// 	// Loading env vars from .env
// 	config.LoadEnvs()
// 	// Creating Log File
// 	logFile, err := os.OpenFile(config.Cfg.LogFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
// 	if err != nil {
// 		panic("Log error: " + err.Error())
// 	}
// 	defer logFile.Close()

// 	logger.CreateLogger()
// 	defer logger.Logger.Sync()
// 	defer logger.Logger.Info("Program end")
// 	keys.ReadRSAKeys()
// }
