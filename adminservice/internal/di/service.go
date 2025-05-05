package di

import (
	"adminservice/internal/entity"
)

type Services struct {
	Plan      SesaonPlanService
	Settings  SettingsService
	DailyTask DailyTaskService
	Statistic StatisticService
}

type (
	SesaonPlanService interface {
		PlanSeasons(season entity.DetailSeasonJson) error
		Seasons() ([]entity.Season, error)
	}

	SettingsService interface {
		UpdateSettings(settings entity.SettingsJson) error
		GameSettings() (entity.SettingsJson, error)
	}

	DailyTaskService interface {
		CreateDailyTask(dailyTask entity.DBDailyTasks) (entity.DailyTasks, error)
		DeleteDailyTask() error
		GetDailyTask() (entity.DailyTasks, error)
	}

	StatisticService interface {
	}
)
