package http

import (
	routes "adminservice/internal/controller/http/v1"
	"adminservice/internal/di"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

// Очень много передается параметров, в твоем сервисе же есть общие депенденси, которые ты скорее всего
// будешь использовать. Тогда почему например не сделать какой-нибудь di контейнер, и передавать в него
// все эти зависимости?
//
// Например:
// ```
//
//	type UseCases struct {
//		PlanUseCase UseCasePlan
//		SettingsUseCase UseCaseSettings
//		DailyTaskUseCase UseCaseDailyTask
//		StatisticUseCase UseCaseStatistic
//	}
//
// ```
// Плюс общая рекомендация: чтобы избегать циклов зависимостей, лучше использовать интерфейсы, а не конкретные реализации.
// И потом передавать этот контейнер куда нужно.
func NewRouter(r *chi.Mux, deps di.Container) {
	// Middlewares
	middleWares(r)

	// Swagger route initialize
	routes.InitSwaggerRoutes(r)

	// Auth routes initialize
	routes.InitAdminRoutes(r, deps)
}

func middleWares(r *chi.Mux) {
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
}
