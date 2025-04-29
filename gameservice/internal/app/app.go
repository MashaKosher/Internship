package app

import (
	"context"

	"gameservice/internal/adapter/kafka/consumers"
	redisRepo "gameservice/internal/adapter/redis/game_settings"
	v1 "gameservice/internal/controller"
	"gameservice/internal/di"
	"gameservice/internal/di/setup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
)

func Run(cfg di.ConfigType) {

	deps := setup.MustContainer(cfg)
	defer setup.DeferContainer(deps)

	go consumers.GameSettingsConsumer(redisRepo.New(deps.Cache, context.Background()), cfg, deps.Bus)

	// Graceful Shutdown
	e := echo.New()
	go func() {
		if err := e.Start(":" + cfg.Server.Port); err != nil && err != http.ErrServerClosed {
			deps.Logger.Fatal("Server error:" + err.Error())
		}
	}()

	v1.NewRouter(e, deps)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		deps.Logger.Fatal("Forced Shutdown:" + err.Error())
	}

	deps.Logger.Info("Server gracefully stopped")
}
