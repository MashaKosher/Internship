package statistic

import (
	repo "adminservice/internal/adapter/db/sql"
)

type UseCase struct {
	repo repo.StatisticRepo
}

func New(r repo.StatisticRepo) *UseCase {
	return &UseCase{
		repo: r,
	}
}
