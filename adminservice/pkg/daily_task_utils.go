package pkg

import (
	"adminservice/internal/entity"
	"encoding/json"
	"io"
)

func ParseDailyTaskBody(body io.ReadCloser, DBDailyTasks *entity.DBDailyTasks) error {
	if err := json.NewDecoder(body).Decode(DBDailyTasks); err != nil {
		return err
	}
	return nil
}

func ParseDailyTaskToKafkaJSON(DBDailyTasks entity.DBDailyTasks) entity.DailyTasks {
	var DailyTasks entity.DailyTasks
	DailyTasks.GamesAmount = DBDailyTasks.GamesAmount
	DailyTasks.ReferalsAmount = DBDailyTasks.ReferalsAmount
	DailyTasks.TaskDate = DBDailyTasks.TaskDate.Format("2006-01-02")

	return DailyTasks
}
