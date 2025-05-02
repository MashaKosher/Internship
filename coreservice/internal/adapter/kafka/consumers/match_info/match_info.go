package matchinfo

import (
	"coreservice/internal/di"
	"coreservice/internal/entity"
	utils "coreservice/pkg/kafka_utils"
	"fmt"
	"strconv"

	leaderboardRepo "coreservice/internal/adapter/db/postgres/leaderboard"
	userRepo "coreservice/internal/adapter/db/postgres/user"
	seasonStatusElasticRepo "coreservice/internal/adapter/elastic/seasons"
)

type MatchInfoConsumer struct {
	consumer        di.KafkaConsumer
	logger          di.LoggerType
	cfg             di.ConfigType
	userRepo        *userRepo.UserRepo
	leaderboardRepo *leaderboardRepo.LeaderboardRepo
	elastic         *seasonStatusElasticRepo.SeasonStatusRepo
}

func New(
	cfg di.ConfigType,
	logger di.LoggerType,
	consumer di.KafkaConsumer,
	userRepo *userRepo.UserRepo,
	leaderboardRepo *leaderboardRepo.LeaderboardRepo,
	elastic *seasonStatusElasticRepo.SeasonStatusRepo,
) *MatchInfoConsumer {
	return &MatchInfoConsumer{
		consumer:        consumer,
		logger:          logger,
		cfg:             cfg,
		userRepo:        userRepo,
		leaderboardRepo: leaderboardRepo,
		elastic:         elastic,
	}
}

func (c *MatchInfoConsumer) Close() {
	c.consumer.Close()
}

func (c *MatchInfoConsumer) RecieveMatchInfo() {
	var err error
	var season entity.Match

	c.logger.Info("Recieve Match Info is working")
	err = c.consumer.Subscribe(c.cfg.Kafka.MatchTopicRecieve, nil)
	if err != nil {
		c.logger.Fatal("Failed to assign partition:" + err.Error())
		panic(err)
	}

	for {
		msg, err := c.consumer.ReadMessage(-1)

		if err == nil {

			c.logger.Info("Received message: " + string(msg.Value) + " from topic:" + msg.TopicPartition.String() + " with offset " + msg.TopicPartition.Offset.String())

			answer, err := utils.DeserializeMatchAnswer(msg.Value, season, c.logger)
			if err != nil {
				c.logger.Error("Error while consuming: " + err.Error())
			}

			season, err := c.elastic.ActiveSeason()
			if err != nil {
				c.logger.Fatal("Error while checkig active season: " + err.Error())
			}

			if len(season) == 0 {
				c.logger.Info("There is no active season right now")
			} else {
				c.logger.Info("winner id: " + strconv.Itoa(answer.Winner) + " Season ID: " + strconv.Itoa(int(season[0])))
				err := c.leaderboardRepo.UpdateSeasonLeaderboard(int(season[0]), answer.Winner)
				if err != nil {
					c.logger.Fatal("Pizda while updating leaderboard: " + err.Error())
				}
				c.logger.Info("Leaderboard updated successfully!!!")

			}

			if _, err := c.userRepo.UpdateBalance(int32(answer.Winner), float64(answer.WinAmount)); err != nil {
				c.logger.Fatal("Pizda upfdating winner balance: " + err.Error())
			}
			c.logger.Info("Winner balance updated successfully")

			if _, err := c.userRepo.UpdateBalance(int32(answer.Loser), float64(-answer.LoseAmount)); err != nil {
				c.logger.Fatal("Pizda upfdating loser balance: " + err.Error())
			}
			c.logger.Info("Loser balance updated successfully")

			c.logger.Info(fmt.Sprint(season))

			c.logger.Info("Match recievd; " + fmt.Sprint(answer))

		} else {
			c.logger.Error("Error while consuming: " + err.Error())
		}
	}
}
