package main

import (
	"adminservice/internal/config"
	"adminservice/internal/di/setup"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"adminservice/internal/app"

	_ "adminservice/docs"
)

// @title						Admin service
// @version						1.0
// @description					Admin server API
// @host						localhost:8004
// @BasePath					/
// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	cfg := config.MustParseConfig()

	deps := setup.MustContainer(cfg)
	defer setup.DeferContainer(deps)

	server := &http.Server{Addr: cfg.Server.Host + ":" + cfg.Server.Port, Handler: app.Run(deps)}

	deps.Logger.Info("Program started")
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				deps.Logger.Info("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Println(err)
		}
		serverStopCtx()
	}()

	// Run the server
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		deps.Logger.Error(err.Error())
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}
