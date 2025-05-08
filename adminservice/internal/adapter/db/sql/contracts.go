package sql

import "adminservice/internal/entity"

type (
	PlanRepo interface {
		AddNewSeason(season *entity.Season) error
		Seasons() ([]entity.Season, error)
	}

	SettingsRepo interface {
		UpdateSettings(newSettings *entity.GameSettings) error
		GameSettings() (entity.GameSettings, error)
	}

	DailyTaskRepo interface {
		AddDailyTask(task *entity.DBDailyTasks) error
		DeleteTodaysTask() error
		GetDailyTask() (entity.DBDailyTasks, error)
	}

	StatisticRepo interface {
	}
)
