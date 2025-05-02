package gamesettings

import (
	"context"
	"encoding/json"
	"gameservice/internal/entity"

	"github.com/redis/go-redis/v9"
)

const gameSettingsKey string = "gameSettings"

type GameSettingsRepo struct {
	*redis.Client
	context.Context
}

func New(redisDB *redis.Client, ctx context.Context) *GameSettingsRepo {
	return &GameSettingsRepo{redisDB, ctx}
}

func (r *GameSettingsRepo) RefreshGameSettings(gameSettings entity.GameSettings) error {
	jsonData, err := json.Marshal(gameSettings)
	if err != nil {
		return err
	}

	if err := r.Client.Set(r.Context, gameSettingsKey, jsonData, 0).Err(); err != nil {
		return err
	}
	return nil
}

func (r *GameSettingsRepo) GetGameSettings() (entity.GameSettings, error) {
	var settings entity.GameSettings

	jsonData, err := r.Client.Get(r.Context, gameSettingsKey).Bytes()
	if err == redis.Nil {
		return settings, nil
	} else if err != nil {
		return settings, err
	}

	if err := json.Unmarshal(jsonData, &settings); err != nil {
		return settings, err
	}

	return settings, nil
}
