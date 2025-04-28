package redis

import "gameservice/internal/entity"

type (
	GameSettingsRepo interface {
		RefreshGameSettings(gameSettings entity.GameSettings) error
		GetGameSettings() (entity.GameSettings, error)
	}
)
