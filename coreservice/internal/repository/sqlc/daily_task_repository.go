package sqlc

import (
	"coreservice/internal/entity"
	"coreservice/internal/logger"
	db "coreservice/internal/repository/sqlc/generated"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func AddDailyTask(dailyTask entity.DailyTask) error {

	layout := "2006-01-02"
	taskTime, err := time.Parse(layout, dailyTask.TaskDate)
	if err != nil {
		return err
	}

	date := pgtype.Date{
		Time:  taskTime,
		Valid: true,
	}

	ref := pgtype.Int4{
		Int32: int32(dailyTask.ReferalsAmount),
		Valid: true,
	}

	win := pgtype.Int4{
		Int32: int32(dailyTask.GamesAmount),
		Valid: true,
	}

	if err := Query.AddDailyTask(Ctx, db.AddDailyTaskParams{TaskDate: date, ReferalsAmount: ref, WinsAmount: win}); err != nil {
		logger.Logger.Error("Error while adding user to DB")
		return err
	}
	return nil
}

func GetDailyTask() (db.DailyTask, error) {
	task, err := Query.GetDailyTask(Ctx, pgtype.Date{Time: time.Now(), Valid: true})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.DailyTask{}, errors.New("there is no tasks for today")
		}
		return db.DailyTask{}, err

	}
	return task, nil
}
