package game

import (
	repo "gameservice/internal/adapter/db/mongo"
)

type UseCase struct {
	repo repo.GameRepo
}

func New(r repo.GameRepo) *UseCase {
	return &UseCase{
		repo: r,
	}
}
