package setup

import (
	"context"
	"coreservice/internal/di"

	"github.com/redis/go-redis/v9"
)

func mustCache(cfg di.ConfigType, logger di.LoggerType) di.CacheType {
	ctx := context.Background()
	r := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port, // "localhost:6379", // адрес Redis
		Password: cfg.Redis.Password,                    //"" // пароль, если есть
		DB:       cfg.Redis.DB,                          // 0 // номер базы данных
	})

	// Проверка подключения
	_, err := r.Ping(ctx).Result()
	if err != nil {
		logger.Fatal(err.Error())
	}

	logger.Info("Connected to Redis successfully")

	return r
}
