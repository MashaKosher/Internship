package v1

import (
	"github.com/go-chi/chi/v5"

	"adminservice/internal/di"

	middleware "adminservice/internal/controller/http/middlewares"
)

func InitAdminRoutes(r *chi.Mux, deps di.Container) {

	plans := initPlanRoutes(deps)
	settings := initSettingsRoutes(deps)
	dailyTask := initDailyTaskRoutes(deps)
	statistic := initStatisticRoutes(deps)

	r.Route("/", func(r chi.Router) {
		r.Use(middleware.CheckToken(deps.Config, deps.Bus))

		// Plan
		r.Post("/deatil-plan", plans.planSeason)

		// Settings
		r.Post("/settings", settings.gameSettings)

		// Statistic
		r.Route("/statistic", func(r chi.Router) {
			r.Get("/players", statistic.players)
			r.Get("/seasons", statistic.seasons)
		})

		// Daily Task
		r.Post("/daily-tasks", dailyTask.createDailyTask)
		r.Delete("/daily-tasks", dailyTask.deleteTodaysTask)
	})

}
