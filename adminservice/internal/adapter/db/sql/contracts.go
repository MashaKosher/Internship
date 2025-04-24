package sql

import "adminservice/internal/entity"

type (
	PlanRepo interface {
		AddNewSeason(season *entity.Season) error
		FindSeasonCross(season *entity.Season) error
	}

	SettingsRepo interface {
		UpdateSettings(newSettings entity.GameSettings) error
	}

	DailyTaskRepo interface {
		AddDailyTask(task entity.DBDailyTasks) error
		DeleteTodaysTask() error
	}

	StatisticRepo interface {
	}
)
