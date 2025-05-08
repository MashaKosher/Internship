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

func (r *DailyTaskRepo) GetDailyTask() (entity.DBDailyTasks, error) {
	var dailyTask entity.DBDailyTasks
	if err := r.DB.Where("task_date = ?", time.Now().Format("2006-01-02")).First(&dailyTask).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.DBDailyTasks{}, entity.ErrRecordNotFound
		}
		return entity.DBDailyTasks{}, err
	}
	return dailyTask, nil
}

func (r *DailyTaskRepo) AddDailyTask(dailyTask *entity.DBDailyTasks) error {
	if dailyTask == nil {
		return entity.ErrDailytaskIsNil
	}
	var count int64
	if err := r.DB.Model(&entity.DBDailyTasks{}).Where("task_date = ?", dailyTask.TaskDate).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return entity.ErrDailyTaskAlreadyExists
	}

	return r.DB.Create(dailyTask).Error
}

func (r *DailyTaskRepo) DeleteTodaysTask() error {

	result := r.DB.Where("task_date = ?", time.Now().Format("2006-01-02")).Delete(&entity.DBDailyTasks{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return entity.ErrRecordNotFound
	}

	return nil
}
