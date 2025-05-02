package app

import (
	// "authservice/internal/adapter/kafka/consumers"
	"authservice/internal/config"
	v1 "authservice/internal/controller/http"
	"authservice/internal/di/setup"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Run(cfg config.Config) {

	deps := setup.MustContainer(cfg)
	defer setup.DeferContainer(deps)

	go deps.Bus.AuthConsumer.ConsumerAnswerTokens()

	app := fiber.New(fiber.Config{
		AppName: "Auth Service",
	})
	v1.NewRouter(app, deps)

	go func() {
		if err := app.Listen(":" + cfg.Server.Port); err != nil {
			deps.Logger.Error("Server error:" + err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	deps.Logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		deps.Logger.Fatal("Server forced to shutdown: " + err.Error())
	}

	deps.Logger.Info("Server gracefully stopped")
}
