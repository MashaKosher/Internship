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
		r.Use(middleware.CheckToken(deps.Config, deps.Bus, deps.Logger))

		// Plan
		r.Post("/deatil-plan", plans.planSeason)
		r.Get("/seasons", plans.seasons)

		// Settings
		r.Get("/settings", settings.gameSettings)
		r.Put("/settings", settings.updatGameSettings)

		// Statistic
		r.Route("/statistic", func(r chi.Router) {
			r.Get("/players", statistic.players)
			r.Get("/seasons", statistic.seasons)
		})

		// Daily Task
		r.Get("/daily-tasks", dailyTask.dailyTask)
		r.Post("/daily-tasks", dailyTask.createDailyTask)
		r.Delete("/daily-tasks", dailyTask.deleteDailyTask)
	})

}
