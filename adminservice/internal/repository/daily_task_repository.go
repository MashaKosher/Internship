package repository

import (
	"adminservice/internal/db"
	"adminservice/internal/entity"
	"time"
)

func AddDailyTask(task entity.DBDailyTasks) error {
	if err := db.DB.Create(&task).Error; err != nil {
		return err
	}
	return nil
}

func DeleteTodaysTask() error {

	var dailyTask entity.DBDailyTasks

	if err := db.DB.Where("task_date = ?", time.Now().Format("2006-01-02")).First(&dailyTask).Error; err != nil {
		return err
	}

	if err := db.DB.Delete(&dailyTask).Error; err != nil {
		return err
	}

	return nil
}
