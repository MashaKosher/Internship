package v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"adminservice/internal/adapter/kafka/producers"
	"adminservice/internal/usecase"
	"adminservice/internal/usecase/plan"
	"adminservice/internal/usecase/settings"
	"adminservice/internal/usecase/statistic"
	"adminservice/pkg"

	middleware "adminservice/internal/controller/http/middlewares"
	dailytask "adminservice/internal/usecase/daily_task"
)

type seasonRoutes struct {
	u usecase.Plan
}

type settingsRoutes struct {
	u usecase.Settings
}

type dailyTasksRoutes struct {
	u usecase.DailyTask
}

type statisticRoutes struct {
	u usecase.Statistic
}

func InitAdminRoutes(r *chi.Mux, planUseCase *plan.UseCase, settingsUseCase *settings.UseCase, dailyTaskUseCase *dailytask.UseCase, statisticUseCase *statistic.UseCase) {

	plans := &seasonRoutes{u: planUseCase}
	settings := &settingsRoutes{u: settingsUseCase}
	dailyTask := &dailyTasksRoutes{u: dailyTaskUseCase}
	statistic := &statisticRoutes{u: statisticUseCase}

	r.Route("/", func(r chi.Router) {
		r.Use(middleware.CheckToken)

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

// @Summary Детальное палнирование сезона
// @Description Обрабатывает запрос на планирование сезона и проверяет права пользователя
// @Tags Season Planing
// @Accept json
// @Produce json
// @Param season body entity.DetailSeasonJson true "Информация о сезоне"
// @Success 200 {string} string "User is admin"
// @Router /deatil-plan [post]
func (sr *seasonRoutes) planSeason(w http.ResponseWriter, r *http.Request) {
	season, err := sr.u.PlanSeasons(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(season)

}

// @Summary Update game settings
// @Description Update game configuration settings for authenticated user
// @Tags Settings
// @Accept  json
// @Produce  json
// @Param settings body entity.SettingsJson true "Game settings object"
// @Success 200 {object} entity.SettingsJson
// @Router			/settings [post]
func (gr *settingsRoutes) gameSettings(w http.ResponseWriter, r *http.Request) {
	// Checking Auth Responce

	settings, err := gr.u.UpdateSettings(w, r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	go producers.SendGameSettings(settings)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(settings)
}

// @Summary Create daily tasks
// @Description Create a new set of daily tasks
// @Tags DailyTasks
// @Accept  json
// @Produce  json
// @Param tasks body entity.DBDailyTasks true "Daily tasks object"
// @Success 201 {object} entity.DailyTasks
// @Router /daily-tasks [post]
func (dr *dailyTasksRoutes) createDailyTask(w http.ResponseWriter, r *http.Request) {
	dailyTask, err := dr.u.CreateDailyTask(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dailyTask)
}

// @Summary Delete daily tasks
// @Description Deelete todays task
// @Tags DailyTasks
// @Accept  json
// @Produce  json
// @Router /daily-tasks [delete]
func (dr *dailyTasksRoutes) deleteTodaysTask(w http.ResponseWriter, r *http.Request) {

	err := dr.u.DeleteDailyTask(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Task deleted")

}

// @Summary		Get season statistic
// @Description	Get season statistic
// @Tags			Statistic
// @Accept			json
// @Produce		json
// @Success		200
// @Router			/statistic/players [get]
func (sr *statisticRoutes) players(w http.ResponseWriter, r *http.Request) {
	// Checking Auth Responce
	if err := pkg.CheckToken(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte("There must be Players statistic"))
}

// @Summary		Get seasons statistic
// @Description	Get seasons statistic
// @Tags			Statistic
// @Accept			json
// @Produce		json
// @Success		200
// @Router			/statistic/seasons [get]
func (sr *statisticRoutes) seasons(w http.ResponseWriter, r *http.Request) {
	// Checking Auth Responce
	if err := pkg.CheckToken(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte("There must be Season statistic"))
}
