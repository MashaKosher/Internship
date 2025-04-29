package dailytask

import (
	"context"
	"coreservice/internal/entity"
	"coreservice/internal/logger"
	db "coreservice/internal/repository/sqlc/generated"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type DailyTaskRepo struct {
	Query *db.Queries
}

func New(queries *db.Queries) *DailyTaskRepo {
	if queries == nil {
		panic("queries is nil")
	}

	return &DailyTaskRepo{
		Query: queries,
	}
}

func (r *DailyTaskRepo) AddDailyTask(dailyTask entity.DailyTask) error {

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

	if err := r.Query.AddDailyTask(context.Background(), db.AddDailyTaskParams{TaskDate: date, ReferalsAmount: ref, WinsAmount: win}); err != nil {
		logger.Logger.Error("Error while adding user to DB")
		return err
	}
	return nil
}

func (r *DailyTaskRepo) GetDailyTask() (db.DailyTask, error) {
	task, err := r.Query.GetDailyTask(context.Background(), pgtype.Date{Time: time.Now(), Valid: true})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.DailyTask{}, errors.New("there is no tasks for today")
		}
		return db.DailyTask{}, err

	}
	return task, nil
}
