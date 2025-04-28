package game

import (
	repo "gameservice/internal/adapter/db/clickhouse"
	redisRepo "gameservice/internal/adapter/redis"
	"gameservice/internal/entity"
)

type UseCase struct {
	repo  repo.GameRepo
	redis redisRepo.GameSettingsRepo
}

func New(r repo.GameRepo, redis redisRepo.GameSettingsRepo) *UseCase {
	return &UseCase{
		repo:  r,
		redis: redis,
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
