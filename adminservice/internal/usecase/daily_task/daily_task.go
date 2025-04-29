package dailytask

import (
	repo "adminservice/internal/adapter/db/sql"
	"adminservice/internal/adapter/kafka/producers"
	"adminservice/internal/di"
	"adminservice/internal/entity"
	"fmt"
)

type UseCase struct {
	// Пиши интерфейсы по месту использования, а не реализации.
	// Интерфейсы - контракт, заключаемый между вызывающим и вызываемым кодом.
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

func (uc *UseCase) CreateDailyTask(dailyTask entity.DBDailyTasks) (entity.DailyTasks, error) {

	if err := uc.repo.AddDailyTask(dailyTask); err != nil {
		return entity.DailyTasks{}, err
	}

	uc.logger.Info("Task added to DB successfully: " + fmt.Sprint(dailyTask))

	dailyTaskOut := dailyTask.ToDTO()

	go producers.SendDailyTask(dailyTaskOut, uc.cfg, uc.bus)

	return dailyTaskOut, nil
}

func (uc *UseCase) DeleteDailyTask() error {
	if err := uc.repo.DeleteTodaysTask(); err != nil {
		return err
	}
	return nil
}
