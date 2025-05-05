package matchinfo

import (
	"coreservice/internal/di"
	"coreservice/internal/entity"
	utils "coreservice/pkg/kafka_utils"
	sqlcUtils "coreservice/pkg/sqlc_utils"
	"fmt"
	"strconv"

	dailyTaskRepo "coreservice/internal/adapter/db/postgres/daily_task"
	leaderboardRepo "coreservice/internal/adapter/db/postgres/leaderboard"
	userRepo "coreservice/internal/adapter/db/postgres/user"
	elasticRepo "coreservice/internal/adapter/elastic"
)

type MatchInfoConsumer struct {
	consumer        di.KafkaConsumer
	logger          di.LoggerType
	cfg             di.ConfigType
	userRepo        *userRepo.UserRepo
	dailyTaskRepo   *dailyTaskRepo.DailyTaskRepo
	leaderboardRepo *leaderboardRepo.LeaderboardRepo
	elastic         elasticRepo.SeasonStatusRepo
}

func New(
	cfg di.ConfigType,
	logger di.LoggerType,
	consumer di.KafkaConsumer,
	userRepo *userRepo.UserRepo,
	dailyTaskRepo *dailyTaskRepo.DailyTaskRepo,
	leaderboardRepo *leaderboardRepo.LeaderboardRepo,
	elastic elasticRepo.SeasonStatusRepo,
) *MatchInfoConsumer {
	return &MatchInfoConsumer{
		consumer:        consumer,
		logger:          logger,
		cfg:             cfg,
		userRepo:        userRepo,
		dailyTaskRepo:   dailyTaskRepo,
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

			winner, _ := c.userRepo.GetPlayerById(int32(answer.Winner))

			if _, err := c.userRepo.UpdateBalance(int32(answer.Winner), float64(answer.WinAmount+sqlcUtils.NumericToFloat64(winner.Balance))); err != nil {
				c.logger.Fatal("Pizda upfdating winner balance: " + err.Error())
			}
			c.logger.Info("Winner balance updated successfully")

			loser, _ := c.userRepo.GetPlayerById(int32(answer.Loser))

			if _, err := c.userRepo.UpdateBalance(int32(answer.Loser), float64(sqlcUtils.NumericToFloat64(loser.Balance)-answer.LoseAmount)); err != nil {
				c.logger.Fatal("Pizda upfdating loser balance: " + err.Error())
			}
			c.logger.Info("Loser balance updated successfully")

			c.logger.Info(fmt.Sprint(season))

			c.logger.Info("Match recievd; " + fmt.Sprint(answer))

			// /////////////////////
			dailyTask, err := c.dailyTaskRepo.GetDailyTask()
			if err != nil {
				if err == entity.ErrNoDailyTask {
					c.logger.Info("There is no daily task for today")
					return
				}
				c.logger.Error(err.Error())
				return
			}
			c.logger.Info("There is a task for today: " + fmt.Sprint(dailyTask))
			c.logger.Info(fmt.Sprint(sqlcUtils.NumericToFloat64(dailyTask.WinReward)))
			c.logger.Info(fmt.Sprint(sqlcUtils.NumberToNumeric(sqlcUtils.NumericToFloat64(dailyTask.WinReward))))

			winStatus, err := c.dailyTaskRepo.WinTaskStatus(int(winner.ID), dailyTask)
			if err != nil {
				c.logger.Error(err.Error())
			}

			if winStatus {
				c.logger.Info("User with ID: " + fmt.Sprint(winner.ID) + " done referal task")
				return
			}

			c.logger.Info("Referal Task Status: " + fmt.Sprint(winStatus) + " for user with ID: " + fmt.Sprint(winner.ID))

			winsAmount, err := c.dailyTaskRepo.AddWin(int(winner.ID), dailyTask)
			if err != nil {
				c.logger.Error(err.Error())
			}

			if winsAmount == int(dailyTask.WinsAmount.Int32) {
				c.dailyTaskRepo.CompleteWinTask(int(winner.ID), dailyTask)
				c.logger.Info("User with ID: " + fmt.Sprint(winner.ID) + " done referal task")

				c.userRepo.UpdateBalance(int32(winner.ID), sqlcUtils.NumericToFloat64(dailyTask.ReferalsReward)+sqlcUtils.NumericToFloat64(winner.Balance))
			}

		} else {
			c.logger.Error("Error while consuming: " + err.Error())
		}
	}
}
