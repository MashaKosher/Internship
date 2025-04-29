package game

import (
	repo "gameservice/internal/adapter/db/clickhouse"
	redisRepo "gameservice/internal/adapter/redis"
	"gameservice/internal/di"
	"gameservice/internal/entity"
)

type UseCase struct {
	repo   repo.GameRepo
	redis  redisRepo.GameSettingsRepo
	logger di.LoggerType
	bus    di.Bus
}

func New(r repo.GameRepo, redis redisRepo.GameSettingsRepo, logger di.LoggerType, bus di.Bus) *UseCase {
	return &UseCase{
		repo:   r,
		redis:  redis,
		logger: logger,
		bus:    bus,
	}
}

func (u *UseCase) GetGameSettings() (entity.GameSettings, error) {
	return u.redis.GetGameSettings()
}

func (u *UseCase) SaveGame(result entity.GameResult) error {
	return u.repo.AddGame(result)
}

func (u *UseCase) GetPlayerStatistic(playerID int) (entity.PlayerStats, error) {
	return u.repo.GetPlayerStatistic(playerID)
}
