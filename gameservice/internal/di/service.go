package di

import "gameservice/internal/entity"

type Services struct {
	Game GameService
}

type (
	GameService interface {
		GetGameSettings() (entity.GameSettings, error)
		SaveGame(result entity.GameResult) error
		GetPlayerStatistic(playerID int) (entity.PlayerStats, error)
	}
)
