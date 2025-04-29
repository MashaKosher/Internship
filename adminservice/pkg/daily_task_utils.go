package pkg

import (
	"adminservice/internal/entity"
)

func ParseDailyTaskToKafkaJSON(DBDailyTasks entity.DBDailyTasks) entity.DailyTasks {
	var DailyTasks entity.DailyTasks
	DailyTasks.GamesAmount = DBDailyTasks.GamesAmount
	DailyTasks.ReferalsAmount = DBDailyTasks.ReferalsAmount
	DailyTasks.TaskDate = DBDailyTasks.TaskDate.Format("2006-01-02")
	return DailyTasks
}
