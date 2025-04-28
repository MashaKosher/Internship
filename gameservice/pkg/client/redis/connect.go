package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func InitRedis() (*redis.Client, context.Context) {
	ctx := context.Background()
	redisDB := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // адрес Redis
		Password: "",               // пароль, если есть
		DB:       0,                // номер базы данных
	})

	// Проверка подключения
	pong, err := redisDB.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to Redis:", pong)

	return redisDB, ctx
}

// func StoreMessage(key string, message models.MessageData) error {
// 	jsonData, err := json.Marshal(message)
// 	if err != nil {
// 		return err
// 	}
// 	return rdb.Set(ctx, key, jsonData, 0).Err()
// }

// func GetMessage(key string) (models.MessageData, error) {
// 	var message models.MessageData
// 	val, err := rdb.Get(ctx, key).Result()
// 	if err != nil {
// 		return message, err
// 	}
// 	err = json.Unmarshal([]byte(val), &message)
// 	return message, err
// }
