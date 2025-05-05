package usersignup

import (
	"coreservice/internal/di"
	"coreservice/internal/entity"
	"encoding/json"
	"fmt"

	sqlcUtils "coreservice/pkg/sqlc_utils"

	dailyTaskRepo "coreservice/internal/adapter/db/postgres/daily_task"
	userRepo "coreservice/internal/adapter/db/postgres/user"
)

type UserSignUpConsumer struct {
	consumer      di.KafkaConsumer
	logger        di.LoggerType
	cfg           di.ConfigType
	dailyTaskRepo *dailyTaskRepo.DailyTaskRepo
	userRepo      *userRepo.UserRepo
}

func New(cfg di.ConfigType, logger di.LoggerType, consumer di.KafkaConsumer, dailyTaskRepo *dailyTaskRepo.DailyTaskRepo, userRepo *userRepo.UserRepo) *UserSignUpConsumer {
	return &UserSignUpConsumer{
		consumer:      consumer,
		logger:        logger,
		cfg:           cfg,
		dailyTaskRepo: dailyTaskRepo,
		userRepo:      userRepo,
	}
}

func (c *UserSignUpConsumer) Close() {
	c.consumer.Close()
}

func (c *UserSignUpConsumer) ReceiveSignedUpUser() {
	var err error
	var signedUpUser entity.SignedUpUser

	c.logger.Info("RecieveSeasonInfo is working")
	err = c.consumer.Subscribe(c.cfg.Kafka.UserSignupRecieve, nil)
	if err != nil {
		c.logger.Fatal("Failed to assign partition:" + err.Error())
		panic(err)
	}

	for {
		msg, err := c.consumer.ReadMessage(-1)

		if err == nil {

			c.logger.Info("Received message: " + string(msg.Value) + " from topic:" + msg.TopicPartition.String() + " with offset " + msg.TopicPartition.Offset.String())

			err := json.Unmarshal(msg.Value, &signedUpUser)
			c.logger.Info("Request recieved: " + fmt.Sprintln(signedUpUser))
			if err != nil {
				c.logger.Error("Error while consuming: " + err.Error())
			}

			c.logger.Info(fmt.Sprint(signedUpUser))

			_, err = c.userRepo.AddPlayer(signedUpUser.ToAuthAnswer())
			if err != nil {
				c.logger.Error(err.Error())
			}

			fmt.Println()

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
			c.logger.Info(fmt.Sprint(sqlcUtils.NumericToFloat64(dailyTask.ReferalsReward)))
			c.logger.Info(fmt.Sprint(sqlcUtils.NumberToNumeric(sqlcUtils.NumericToFloat64(dailyTask.ReferalsReward))))

			refStatus, err := c.dailyTaskRepo.ReferalsTaskStatus(int(signedUpUser.ReferalID), dailyTask)
			if err != nil {
				c.logger.Error(err.Error())
			}

			if refStatus {
				c.logger.Info("User with ID: " + fmt.Sprint(signedUpUser.ReferalID) + " done referal task")
				return
			}

			c.logger.Info("Referal Task Status: " + fmt.Sprint(refStatus) + " for user with ID: " + fmt.Sprint(signedUpUser.ReferalID))

			referalAmount, err := c.dailyTaskRepo.AddReferal(int(signedUpUser.ReferalID), dailyTask)
			if err != nil {
				c.logger.Error(err.Error())
			}

			if referalAmount == int(dailyTask.ReferalsAmount.Int32) {
				c.dailyTaskRepo.CompleteReferalsTask(int(signedUpUser.ReferalID), dailyTask)
				c.logger.Info("User with ID: " + fmt.Sprint(signedUpUser.ReferalID) + " done referal task")
				// c.userRepo.

				dbUser, _ := c.userRepo.GetPlayerById(int32(signedUpUser.ReferalID))

				c.userRepo.UpdateBalance(int32(signedUpUser.ReferalID), sqlcUtils.NumericToFloat64(dailyTask.ReferalsReward)+sqlcUtils.NumericToFloat64(dbUser.Balance))
			}

		} else {
			c.logger.Error("Error while consuming: " + err.Error())
		}
	}
}
