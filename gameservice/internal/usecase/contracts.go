package usecase

import "gameservice/internal/entity"

type (
	Game interface {
		GetGameSettings() (entity.GameSettings, error)
		SaveGame(result entity.GameResult) error
		GetPlayerStatistic(playerID int) (entity.PlayerStats, error)
	}
)
