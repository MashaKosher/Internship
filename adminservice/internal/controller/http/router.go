package http

import (
	routes "adminservice/internal/controller/http/v1"
	dailytask "adminservice/internal/usecase/daily_task"
	"adminservice/internal/usecase/plan"
	"adminservice/internal/usecase/settings"
	"adminservice/internal/usecase/statistic"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func NewRouter(r *chi.Mux, planUseCase *plan.UseCase, settingsUseCase *settings.UseCase, dailyTaskUseCase *dailytask.UseCase, statisticUseCase *statistic.UseCase) {
	// Middlewares
	middleWares(r)

	// Swagger route initialize
	routes.InitSwaggerRoutes(r)

	// Auth routes initialize
	routes.InitAdminRoutes(r, planUseCase, settingsUseCase, dailyTaskUseCase, statisticUseCase)
}

func middleWares(r *chi.Mux) {
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
}
