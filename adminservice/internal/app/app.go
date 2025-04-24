package app

import (
	v1 "adminservice/internal/controller/http"
	db "adminservice/pkg/client/sql"
	"net/http"

	dailyTask "adminservice/internal/usecase/daily_task"
	"adminservice/internal/usecase/plan"
	"adminservice/internal/usecase/settings"
	"adminservice/internal/usecase/statistic"

	dailyTaskRepo "adminservice/internal/adapter/db/sql/daily_task"
	planRepo "adminservice/internal/adapter/db/sql/plan"
	settingsRepo "adminservice/internal/adapter/db/sql/settings"
	statisticRepo "adminservice/internal/adapter/db/sql/statistic"

	"github.com/go-chi/chi/v5"
)

func Run() http.Handler {
	r := chi.NewRouter()

	db := db.ConncetDB()

	// Создаем Use Case
	planUseCase := plan.New(
		planRepo.New(db), // создаем конкретный репозиторий и передаем в конретный Use Case
	)

	// Создаем Use Case
	settingsUseCase := settings.New(
		settingsRepo.New(db), // создаем конкретный репозиторий и передаем в конретный Use Case
	)

	// Создаем Use Case
	dailyTasksUseCase := dailyTask.New(
		dailyTaskRepo.New(db), // создаем конкретный репозиторий и передаем в конретный Use Case
	)

	// Создаем Use Case
	statisticUseCase := statistic.New(
		statisticRepo.New(db), // создаем конкретный репозиторий и передаем в конретный Use Case
	)

	v1.NewRouter(r, planUseCase, settingsUseCase, dailyTasksUseCase, statisticUseCase)

	return r
}
