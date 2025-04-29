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
	// Создаем Use Case
	planUseCase := plan.New(
		planRepo.New(db), // создаем конкретный репозиторий и передаем в конретный Use Case
		logger,
		cfg,
		bus,
	)

	// Создаем Use Case
	settingsUseCase := settings.New(
		settingsRepo.New(db), // создаем конкретный репозиторий и передаем в конретный Use Case
		logger,
		cfg,
		bus,
	)

	// Создаем Use Case
	dailyTasksUseCase := dailyTask.New(
		dailyTaskRepo.New(db), // создаем конкретный репозиторий и передаем в конретный Use Case
		logger,
		cfg,
		bus,
	)

	// Создаем Use Case
	statisticUseCase := statistic.New(
		statisticRepo.New(db), // создаем конкретный репозиторий и передаем в конретный Use Case
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
