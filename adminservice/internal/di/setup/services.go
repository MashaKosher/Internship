package setup

import (
	"adminservice/internal/di"
	dailyTask "adminservice/internal/usecase/daily_task"
	"adminservice/internal/usecase/plan"
	"adminservice/internal/usecase/settings"
	"adminservice/internal/usecase/statistic"

	dailyTaskRepo "adminservice/internal/adapter/db/sql/daily_task"
	planRepo "adminservice/internal/adapter/db/sql/plan"
	settingsRepo "adminservice/internal/adapter/db/sql/settings"
	statisticRepo "adminservice/internal/adapter/db/sql/statistic"
)

func mustServices(db di.DBType, logger di.LoggerType, cfg di.ConfigType, bus di.Bus) di.Services {

	planUseCase := plan.New(
		planRepo.New(db),
		logger,
		cfg,
		bus,
	)

	settingsUseCase := settings.New(
		settingsRepo.New(db),
		logger,
		cfg,
		bus,
	)

	dailyTasksUseCase := dailyTask.New(
		dailyTaskRepo.New(db),
		logger,
		cfg,
		bus,
	)

	statisticUseCase := statistic.New(
		statisticRepo.New(db),
		logger,
		cfg,
		bus,
	)

	return di.Services{
		Plan:      planUseCase,
		Settings:  settingsUseCase,
		DailyTask: dailyTasksUseCase,
		Statistic: statisticUseCase,
	}

}
