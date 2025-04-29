package usecase

import (
	"adminservice/internal/entity"
	"net/http"
)

// Почему в бизнес-логике лежат параметры http хендлера?
type (
	Plan interface {
		PlanSeasons(w http.ResponseWriter, r *http.Request) (entity.Season, error)
	}

	Settings interface {
		UpdateSettings(w http.ResponseWriter, r *http.Request) (entity.SettingsJson, error)
	}

	DailyTask interface {
		CreateDailyTask(w http.ResponseWriter, r *http.Request) (entity.DailyTasks, error)
		DeleteDailyTask(w http.ResponseWriter, r *http.Request) error
	}

	Statistic interface {
	}
)
