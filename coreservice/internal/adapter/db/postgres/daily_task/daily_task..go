package dailytask

import (
	"context"
	"coreservice/internal/entity"
	db "coreservice/internal/repository/sqlc/generated"
	utils "coreservice/pkg/sqlc_utils"
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

	// layout := "2006-01-02"
	// taskTime, err := time.Parse(layout, dailyTask.TaskDate)
	// if err != nil {
	// 	return err
	// }

	// date := pgtype.Date{
	// 	Time:  taskTime,
	// 	Valid: true,
	// }

	date, err := utils.StringToDate(dailyTask.TaskDate)
	if err != nil {
		return err
	}

	// ref := pgtype.Int4{
	// 	Int32: int32(dailyTask.ReferalsAmount),
	// 	Valid: true,
	// }

	ref := utils.IntToInt4(dailyTask.ReferalsAmount)

	// win := pgtype.Int4{
	// 	Int32: int32(dailyTask.GamesAmount),
	// 	Valid: true,
	// }

	win := utils.IntToInt4(dailyTask.GamesAmount)

	referalsReward, err := utils.NumberToNumeric(dailyTask.ReferalsTaskReward)
	if err != nil {
		return err
	}

	winReward, err := utils.NumberToNumeric(dailyTask.GameTaskReward)
	if err != nil {
		return err
	}

	if err := r.Query.AddDailyTask(
		context.Background(),
		db.AddDailyTaskParams{
			TaskDate:       date,
			ReferalsAmount: ref,
			ReferalsReward: referalsReward,
			WinsAmount:     win,
			WinReward:      winReward}); err != nil {
		return errors.New("error while adding season to DB" + err.Error())
	}
	return nil
}

func (r *DailyTaskRepo) GetDailyTask() (db.DailyTask, error) {
	task, err := r.Query.GetDailyTask(context.Background(), pgtype.Date{Time: time.Now(), Valid: true})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.DailyTask{}, entity.ErrNoDailyTask
		}
		return db.DailyTask{}, err

	}
	return task, nil
}

//

func (r *DailyTaskRepo) AddWin(userID int, dailyTask db.DailyTask) (int, error) {

	// date, err := utils.StringToDate(dailyTask.TaskDate)
	// if err != nil {
	// 	return -1, err
	// }

	winAmount, err := r.Query.AddWin(context.Background(), db.AddWinParams{TaskDate: dailyTask.TaskDate, UserID: int32(userID)})
	if err != nil {
		return -1, err
	}

	return utils.Int4ToInt(winAmount), nil
}

func (r *DailyTaskRepo) AddReferal(userID int, dailyTask db.DailyTask) (int, error) {

	// date, err := utils.StringToDate(dailyTask.TaskDate)
	// if err != nil {
	// 	return -1, err
	// }

	winAmount, err := r.Query.AddReferal(context.Background(), db.AddReferalParams{TaskDate: dailyTask.TaskDate, UserID: int32(userID)})
	if err != nil {
		return -1, err
	}

	return utils.Int4ToInt(winAmount), nil
}

func (r *DailyTaskRepo) CompleteWinTask(userID int, dailyTask db.DailyTask) error {
	// date, err := utils.StringToDate(dailyTask.TaskDate)
	// if err != nil {
	// 	return err
	// }

	err := r.Query.CompleteWinTask(context.Background(), db.CompleteWinTaskParams{TaskDate: dailyTask.TaskDate, UserID: int32(userID)})
	if err != nil {
		return err
	}
	return nil
}

func (r *DailyTaskRepo) CompleteReferalsTask(userID int, dailyTask db.DailyTask) error {
	// date, err := utils.StringToDate(dailyTask.TaskDate)
	// if err != nil {
	// 	return err
	// }

	err := r.Query.CompleteReferalsTask(context.Background(), db.CompleteReferalsTaskParams{TaskDate: dailyTask.TaskDate, UserID: int32(userID)})
	if err != nil {
		return err
	}
	return nil
}

func (r *DailyTaskRepo) WinTaskStatus(userID int, dailyTask db.DailyTask) (bool, error) {

	// date, err := utils.StringToDate(dailyTask.TaskDate)
	// if err != nil {
	// 	return false, err
	// }

	status, err := r.Query.WinTaskStatus(context.Background(), db.WinTaskStatusParams{TaskDate: dailyTask.TaskDate, UserID: int32(userID)})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return status, nil
}

func (r *DailyTaskRepo) ReferalsTaskStatus(userID int, dailyTask db.DailyTask) (bool, error) {

	// date, err := utils.StringToDate(dailyTask.TaskDate)
	// if err != nil {
	// 	return false, err
	// }

	status, err := r.Query.ReferalsTaskStatus(context.Background(), db.ReferalsTaskStatusParams{TaskDate: dailyTask.TaskDate, UserID: int32(userID)})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return status, nil
}
