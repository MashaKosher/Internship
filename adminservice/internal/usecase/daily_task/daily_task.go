package dailytask

import (
	repo "adminservice/internal/adapter/db/sql"
	"adminservice/internal/di"
	"adminservice/internal/entity"
	"fmt"
)

type UseCase struct {
	repo   repo.DailyTaskRepo
	logger di.LoggerType
	cfg    di.ConfigType
	bus    di.Bus
}

func New(r repo.DailyTaskRepo, logger di.LoggerType, cfg di.ConfigType, bus di.Bus) *UseCase {
	return &UseCase{
		repo:   r,
		logger: logger,
		cfg:    cfg,
		bus:    bus,
	}
}

func (u *UseCase) CreateDailyTask(dailyTask entity.DBDailyTasks) (entity.DailyTasks, error) {

	if err := u.repo.AddDailyTask(dailyTask); err != nil {
		return entity.DailyTasks{}, err
	}

	u.logger.Info("Task added to DB successfully: " + fmt.Sprint(dailyTask))

	dailyTaskOut := dailyTask.ToDTO()

	go u.bus.DailyTaskProducer.SendDailyTask(dailyTaskOut)

	return dailyTaskOut, nil
}

func (u *UseCase) DeleteDailyTask() error {
	if err := u.repo.DeleteTodaysTask(); err != nil {
		return err
	}
	return nil
}
