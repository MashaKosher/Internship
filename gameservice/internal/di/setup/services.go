package setup

import (
	"context"
	gameRepo "gameservice/internal/adapter/db/clickhouse/game"
	redisRepo "gameservice/internal/adapter/redis/game_settings"
	"gameservice/internal/di"
	"gameservice/internal/usecase/game"
)

func mustServices(db di.DBType, cache di.CacheType, logger di.LoggerType, bus di.Bus) di.Services {

	gameRepository := gameRepo.New(db)
	redisRepository := redisRepo.New(cache, context.Background())

	if err := gameRepository.CreateGameResultsTable(); err != nil {
		logger.Fatal("Problems with creating tables")
	}

	logger.Info("ClickHouse table created successfully")

	gameUseCase := game.New(
		gameRepository,
		redisRepository,
		logger,
		bus,
	)

	return di.Services{
		Game: gameUseCase,
	}

}
