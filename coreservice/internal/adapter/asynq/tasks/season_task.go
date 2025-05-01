package tasks

import (
	"context"
	"coreservice/internal/di"
	"encoding/json"
	"fmt"
	"time"

	seasonRepo "coreservice/internal/adapter/db/postgres/season"
	seasonStatusElasticRepo "coreservice/internal/adapter/elastic/seasons"

	"github.com/hibiken/asynq"
)

type actionType string

const (
	CurrentSeason actionType = "Current"
	CancelSeason  actionType = "Cancled"
)

type SeasonPayload struct {
	SeasonID   int        `json:"season-id"`
	SeasonTime time.Time  `json:"season-time"`
	ActionType actionType `json:"action-type"`
}

func NewSeasonTask(seasonId int, seasonTime time.Time, action actionType) (*asynq.Task, error) {

	if action != CurrentSeason && action != CancelSeason {
		return nil, fmt.Errorf("invalid action type, it can onlly be %s or %s", CurrentSeason, CancelSeason)
	}

	payload, err := json.Marshal(SeasonPayload{
		SeasonID: seasonId, SeasonTime: seasonTime, ActionType: action,
	})

	if err != nil {
		return nil, err
	}

	task := asynq.NewTask(TypeSeason, payload)
	return task, nil
}

func SeasonTaskHadler(logger di.LoggerType, db di.DBType, ESClient di.ESClient, Index di.ElasticIndex) func(ctx context.Context, t *asynq.Task) error {
	return func(ctx context.Context, t *asynq.Task) error {
		seasonRepo := seasonRepo.New(db)
		elastic := seasonStatusElasticRepo.New(ESClient, Index, logger)

		logger.Info("Season task produce....")
		var season SeasonPayload
		if err := json.Unmarshal(t.Payload(), &season); err != nil {
			logger.Info(err.Error())
			return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
		}

		logger.Info("Serialize successfully")

		now := time.Now()

		if season.ActionType == CancelSeason && season.SeasonTime.Before(now) {
			logger.Info(fmt.Sprintf("ВНИМАНИЕ: Сезон просрочен (End time: %v, Current: %v)", season.SeasonTime, now))
		}

		if season.ActionType == CurrentSeason {
			seasonRepo.StartSeason(season.SeasonID)
			logger.Info(fmt.Sprintf("Season %d started\n", season.SeasonID))
			elastic.StartSeason(season.SeasonID)
		} else {
			time.Sleep(time.Second)
			seasonRepo.EndSeason(season.SeasonID)
			logger.Info(fmt.Sprintf("Season %d ended\n", season.SeasonID))
			elastic.EndSeason(season.SeasonID)
		}

		return nil
	}
}

// func SeasonTaskHadler(ctx context.Context, t *asynq.Task, logger di.LoggerType) error {

// 	seasonRepo := seasonRepo.New(db)
// 	elastic := seasonStatusElasticRepo.New(ESClient, Index, bus.Logger)

// 	fmt.Println("Season task produce....")
// 	var season SeasonPayload
// 	if err := json.Unmarshal(t.Payload(), &season); err != nil {
// 		logger.Info(err.Error())
// 		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
// 	}

// 	logger.Info("Serialize successfully")

// 	now := time.Now()

// 	fmt.Println(now)

// 	if season.ActionType == CancelSeason && season.SeasonTime.Before(now) {
// 		log.Printf("ВНИМАНИЕ: Сезон просрочен (End time: %v, Current: %v)", season.SeasonTime, now)
// 	}
// 	// else if season.ActionType == CancelSeason && now.Before(season.SeasonTime) {
// 	// 	return fmt.Errorf("task too early (now: %v, start: %v)", now, season.SeasonTime)
// 	// }

// 	if season.ActionType == CurrentSeason {
// 		repo.StartSeason(season.SeasonID)
// 		fmt.Printf("Season %d started\n", season.SeasonID)
// 		elastic.UpdateSeasonInIndex(season.SeasonID, elastic.CurrentSeason)
// 	} else {
// 		time.Sleep(time.Second)
// 		repo.EndSeason(season.SeasonID)
// 		fmt.Printf("Season %d ended\n", season.SeasonID)
// 		elastic.UpdateSeasonInIndex(season.SeasonID, elastic.EndedSeason)
// 	}

// 	return nil
// }
