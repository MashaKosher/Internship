package setup

import (
	"context"
	leaderBoardCache "coreservice/internal/adapter/redis/leaderBoard"
	"coreservice/internal/di"
	dailyTask "coreservice/internal/usecase/daily_task"
	"coreservice/internal/usecase/search"
	"coreservice/internal/usecase/season"
	"coreservice/internal/usecase/token"
	"coreservice/internal/usecase/user"

	dailyTaskRepo "coreservice/internal/adapter/db/postgres/daily_task"
	seasonRepo "coreservice/internal/adapter/db/postgres/season"
	userRepo "coreservice/internal/adapter/db/postgres/user"
	seasonStatusElasticRepo "coreservice/internal/adapter/elastic/seasons"
	userNameElasticRepo "coreservice/internal/adapter/elastic/user"
)

func mustServices(db di.DBType, logger di.LoggerType, elastic di.ElasticType, redis di.CacheType) di.Services {

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
		seasonStatusElasticRepo.New(elastic.ESClient, elastic.SeasonSearchIndex, logger),
		leaderBoardCache.New(redis, context.Background()),
	)

	searchUseCase := search.New(
		logger,
		userNameElasticRepo.New(elastic.ESClient, elastic.UserSearchIndex, logger, userRepo.New(db)),
	)

	return di.Services{
		DailyTask: dailyTasksUseCase,
		Token:     tokenUseCase,
		User:      userUseCase,
		Season:    seasonUseCase,
		Search:    searchUseCase,
	}

}
