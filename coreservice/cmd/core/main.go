package main

import (
	"coreservice/internal/app"
	"coreservice/internal/config"

	_ "coreservice/docs"
)

// @title Example API
// @version 1.0
// @description This is a sample API for demonstrating Swagger with Gin.
// @host localhost:8006
// @BasePath /
func main() {
	cfg := config.MustParseConfig()
	app.Run(cfg)
}
