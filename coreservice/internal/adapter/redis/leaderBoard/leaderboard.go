package leaderboard

import (
	"context"
	"coreservice/internal/entity"
	"encoding/json"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const cacheTime = 1

type LeaderboardCache struct {
	*redis.Client
	context.Context
}

func New(redisDB *redis.Client, ctx context.Context) *LeaderboardCache {
	return &LeaderboardCache{redisDB, ctx}
}

func (r *LeaderboardCache) UpdateLeaderBoard(leaderBoard []entity.Leaderboard, seasonID int) error {
	jsonData, err := json.Marshal(leaderBoard)
	if err != nil {
		return err
	}

	if err := r.Client.Set(r.Context, "season"+strconv.Itoa(seasonID), jsonData, time.Minute*cacheTime).Err(); err != nil {
		return err
	}
	return nil
}

func (r *LeaderboardCache) GetSeasonLeaderBoard(seasonID int) ([]entity.Leaderboard, error) {
	var leaderBoard []entity.Leaderboard

	jsonData, err := r.Client.Get(r.Context, "season"+strconv.Itoa(seasonID)).Bytes()
	if err == redis.Nil {
		// Если ключ не существует, возвращаем пустую структуру
		return leaderBoard, nil
	} else if err != nil {
		return leaderBoard, err
	}

	if err := json.Unmarshal(jsonData, &leaderBoard); err != nil {
		return leaderBoard, err
	}

	return leaderBoard, nil
}
