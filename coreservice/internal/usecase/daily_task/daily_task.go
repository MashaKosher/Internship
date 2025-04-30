package dailytask

import (
	repo "coreservice/internal/adapter/db/postgres"
	db "coreservice/internal/repository/sqlc/generated"
)

type UseCase struct {
	repo repo.DailyTaskRepo
}

func New(r repo.DailyTaskRepo) *UseCase {
	return &UseCase{
		repo: r,
	}
}

func (u *UseCase) DailyTask() (db.DailyTask, error) {
	task, err := u.repo.GetDailyTask()
	if err != nil {
		return db.DailyTask{}, err
	}
	return task, nil
}
