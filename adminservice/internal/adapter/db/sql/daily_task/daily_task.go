package dailytask

import (
	"adminservice/internal/entity"
	"errors"
	"time"

	"gorm.io/gorm"
)

type DailyTaskRepo struct {
	*gorm.DB
}

func New(db *gorm.DB) *DailyTaskRepo {
	return &DailyTaskRepo{db}
}

func (r DailyTaskRepo) AddDailyTask(task entity.DBDailyTasks) error {
	if err := r.DB.Create(&task).Error; err != nil {
		return err
	}
	return nil
}

func (r DailyTaskRepo) DeleteTodaysTask() error {

	var dailyTask entity.DBDailyTasks

	if err := r.DB.Where("task_date = ?", time.Now().Format("2006-01-02")).First(&dailyTask).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.ErrRecordNotFound
		}
		return err
	}

	if err := r.DB.Delete(&dailyTask).Error; err != nil {
		return err
	}

	return nil
}

func (r DailyTaskRepo) GetDailyTask() (entity.DBDailyTasks, error) {
	var dailyTask entity.DBDailyTasks

	if err := r.DB.Where("task_date = ?", time.Now().Format("2006-01-02")).First(&dailyTask).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.DBDailyTasks{}, entity.ErrRecordNotFound
		}
		return entity.DBDailyTasks{}, err // Это уже внутренняя ошибка
	}

	return dailyTask, nil
}
