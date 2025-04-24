package v1

import (
	"github.com/go-chi/chi/v5"

	controllers "adminservice/internal/service"

	middleware "adminservice/internal/controller/http/middlewares"
)

func InitAdminRoutes(r *chi.Mux) {
	r.Route("/", func(r chi.Router) {
		r.Use(middleware.CheckToken)
		r.Post("/deatil-plan", controllers.PlanSeason)
		r.Post("/settings", controllers.GameSettings)
		r.Route("/statistic", func(r chi.Router) {
			r.Get("/players", controllers.Players)
			r.Get("/seasons", controllers.Seasons)
		})

		r.Post("/daily-tasks", controllers.CreateDailyTasks)
		r.Delete("/daily-tasks", controllers.DeleteTodaysTask)
	})

}
