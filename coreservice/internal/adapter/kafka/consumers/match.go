package consumers

import (
	"coreservice/internal/di"
	"coreservice/internal/entity"
	"coreservice/pkg"
	"fmt"
	"strconv"

	// seasonRepo "coreservice/internal/adapter/db/postgres/season"
	leaderboardRepo "coreservice/internal/adapter/db/postgres/leaderboard"
	userRepo "coreservice/internal/adapter/db/postgres/user"
	seasonStatusElasticRepo "coreservice/internal/adapter/elastic/seasons"
)

func RecieveMatchInfo(cfg di.ConfigType, bus di.Bus, db di.DBType, ESClient di.ESClient, Index di.ElasticIndex) {

	// seasonRepo := seasonRepo.New(db)
	elastic := seasonStatusElasticRepo.New(ESClient, Index, bus.Logger)
	leaderboardRepo := leaderboardRepo.New(db)
	userRepo := userRepo.New(db)

	var err error
	var season entity.Match

	bus.Logger.Info("Recieve Match Info is working")
	err = bus.GameConsumer.Subscribe(cfg.Kafka.MatchTopicRecieve, nil)
	if err != nil {
		bus.Logger.Fatal("Failed to assign partition:" + err.Error())
		panic(err)
	}

	for {
		msg, err := bus.GameConsumer.ReadMessage(-1)

		if err == nil {

			bus.Logger.Info("Received message: " + string(msg.Value) + " from topic:" + msg.TopicPartition.String() + " with offset " + msg.TopicPartition.Offset.String())

			answer, err := pkg.DeserializeMatchAnswer(msg.Value, season, bus.Logger)
			if err != nil {
				bus.Logger.Error("Error while consuming: " + err.Error())
			}

			season, err := elastic.ActiveSeason()
			if err != nil {
				bus.Logger.Fatal("Error while checkig active season: " + err.Error())
			}

			if len(season) == 0 {
				bus.Logger.Info("There is no active season right now")
			} else {
				bus.Logger.Info("winner id: " + strconv.Itoa(answer.Winner) + " Season ID: " + strconv.Itoa(int(season[0])))
				err := leaderboardRepo.UpdateSeasonLeaderboard(int(season[0]), answer.Winner)
				if err != nil {
					bus.Logger.Fatal("Pizda while updating leaderboard: " + err.Error())
				}
				bus.Logger.Info("Leaderboard updated successfully!!!")

			}

			if _, err := userRepo.UpdateBalance(int32(answer.Winner), float64(answer.WinAmount)); err != nil {
				bus.Logger.Fatal("Pizda upfdating winner balance: " + err.Error())
			}
			bus.Logger.Info("Winner balance updated successfully")

			if _, err := userRepo.UpdateBalance(int32(answer.Loser), float64(-answer.LoseAmount)); err != nil {
				bus.Logger.Fatal("Pizda upfdating loser balance: " + err.Error())
			}
			bus.Logger.Info("Loser balance updated successfully")

			bus.Logger.Info(fmt.Sprint(season))

			bus.Logger.Info("Match recievd; " + fmt.Sprint(answer))

		} else {
			bus.Logger.Error("Error while consuming: " + err.Error())
		}
	}
}
