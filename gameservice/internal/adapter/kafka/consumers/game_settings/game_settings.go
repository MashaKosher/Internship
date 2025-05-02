package gamesettings

import (
	"fmt"
	"gameservice/internal/di"
	utils "gameservice/pkg/kafka_utils"

	redisRepo "gameservice/internal/adapter/redis/game_settings"
)

type GameSettingsConsumer struct {
	consumer  di.KafkaConsumer
	logger    di.LoggerType
	cfg       di.ConfigType
	cacheRepo *redisRepo.GameSettingsRepo
}

func New(cfg di.ConfigType, logger di.LoggerType, consumer di.KafkaConsumer, cacheRepo *redisRepo.GameSettingsRepo) *GameSettingsConsumer {
	return &GameSettingsConsumer{
		consumer:  consumer,
		logger:    logger,
		cfg:       cfg,
		cacheRepo: cacheRepo,
	}
}

func (c *GameSettingsConsumer) Close() {
	c.consumer.Close()
}

// gamesettinggsconsumer
func (c *GameSettingsConsumer) ConsumeGameSettings() {

	c.logger.Info("Kafka Game Settings Consumer connected successfully")

	if err := c.consumer.SubscribeTopics([]string{c.cfg.Kafka.GameSettingsTopicRecieve}, nil); err != nil {
		c.logger.Fatal("Failed to subscribe to topics: " + err.Error())
	}

	for {
		msg, err := c.consumer.ReadMessage(-1)
		if err == nil {

			// Тут поменять и вообще все кафка утилс
			gameSettings, err := utils.DeserializeGameSettings(msg.Value, c.logger)

			if err != nil {
				c.logger.Fatal("Error while consuming: " + err.Error())
			}

			c.logger.Info("Received message: " + string(msg.Value) + " from topic:" + fmt.Sprintln(msg.TopicPartition))

			if err := c.cacheRepo.RefreshGameSettings(gameSettings); err != nil {
				c.logger.Fatal("Troubles with adding game setiings to Redis: " + err.Error())
			}

			c.logger.Info("Game Settings added successfuly: " + fmt.Sprint(gameSettings))

		} else {
			c.logger.Fatal("Error while consuming: " + err.Error())
		}
	}

}
