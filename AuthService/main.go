package main

import (
	_ "authservice/docs" // импортируйте generated swagger docs
	"os"

	config "authservice/config"

	routes "authservice/routes"

	"github.com/gofiber/fiber/v2"
	// fiber-swagger middleware
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
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err) // Обработка ошибок при создании файла
	}
	defer logFile.Close() // Закрытие файла после завершения работы
	config.CreateLogger()
	defer config.Logger.Sync() // Отложенная синхронизация логов
	defer config.Logger.Info("Program end")

	app := fiber.New(fiber.Config{
		AppName: "Auth Service",
	})

	config.ConncetDB()

	// Навешиваем пути
	routes.Handlers(app)

	// app.Get("/", hello)
	// app.Get("/hi", Hi)

	app.Listen(":8080")

}
