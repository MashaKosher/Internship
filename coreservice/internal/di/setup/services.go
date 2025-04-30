package setup

import (
	"coreservice/internal/di"
	dailyTask "coreservice/internal/usecase/daily_task"
	"coreservice/internal/usecase/search"
	"coreservice/internal/usecase/season"
	"coreservice/internal/usecase/token"
	"coreservice/internal/usecase/user"

	dailyTaskRepo "coreservice/internal/adapter/db/postgres/daily_task"
	seasonRepo "coreservice/internal/adapter/db/postgres/season"
	userRepo "coreservice/internal/adapter/db/postgres/user"
)

func mustServices(db di.DBType, logger di.LoggerType, cfg di.ConfigType) di.Services {

	// Создаем Use Case
	dailyTasksUseCase := dailyTask.New(
		dailyTaskRepo.New(db),
	)

	tokenUseCase := token.New()

	userUseCase := user.New(
		userRepo.New(db),
		logger,
	)

	seasonUseCase := season.New(
		seasonRepo.New(db),
		logger,
	)

	searchUseCase := search.New(
		logger,
	)

	return di.Services{
		DailyTask: dailyTasksUseCase,
		Token:     tokenUseCase,
		User:      userUseCase,
		Season:    seasonUseCase,
		Search:    searchUseCase,
	}

}
