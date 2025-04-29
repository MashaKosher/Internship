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
	}

	SettingsService interface {
		UpdateSettings(settings entity.SettingsJson) error
	}

	DailyTaskService interface {
		CreateDailyTask(dailyTask entity.DBDailyTasks) (entity.DailyTasks, error)
		DeleteDailyTask() error
	}

	StatisticService interface {
	}
)
