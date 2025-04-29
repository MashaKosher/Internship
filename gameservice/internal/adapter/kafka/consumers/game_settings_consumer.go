package consumers

import (
	"fmt"
	"gameservice/internal/di"
	"gameservice/pkg"

	redisRepo "gameservice/internal/adapter/redis/game_settings"
)

func GameSettingsConsumer(redisRepo *redisRepo.GameSettingsRepo, cfg di.ConfigType, bus di.Bus) {
	// c, err := kafka.NewConsumer(&kafka.ConfigMap{
	// 	"bootstrap.servers": cfg.Kafka.Host + ":" + cfg.Kafka.Port,
	// 	"group.id":          "authService",
	// 	"auto.offset.reset": "earliest",
	// })
	// if err != nil {
	// 	logger.Fatal("Failed to create consumer: " + err.Error())
	// }
	// defer c.Close()

	bus.Logger.Info("Kafka Game Settings Consumer connected successfully")

	if err := bus.GameSettingsConsumer.SubscribeTopics([]string{cfg.Kafka.GameSettingsTopicRecieve}, nil); err != nil {
		bus.Logger.Fatal("Failed to subscribe to topics: " + err.Error())
	}

	for {
		msg, err := bus.GameSettingsConsumer.ReadMessage(-1)
		if err == nil {

			// Тут поменять и вообще все кафка утилс
			gameSettings, err := pkg.DeserializeGameSettings(msg.Value, bus.Logger)

			if err != nil {
				bus.Logger.Fatal("Error while consuming: " + err.Error())
			}

			bus.Logger.Info("Received message: " + string(msg.Value) + " from topic:" + fmt.Sprintln(msg.TopicPartition))

			if err := redisRepo.RefreshGameSettings(gameSettings); err != nil {
				bus.Logger.Fatal("Troubles with adding game setiings to Redis: " + err.Error())
			}

			bus.Logger.Info("Game Settings added successfuly: " + fmt.Sprint(gameSettings))

		} else {
			bus.Logger.Fatal("Error while consuming: " + err.Error())
		}
	}

}
