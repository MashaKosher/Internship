package postgres

import (
	"coreservice/internal/entity"
	db "coreservice/internal/repository/sqlc/generated"
)

type (
	DailyTaskRepo interface {
		GetDailyTask() (db.DailyTask, error)
		AddDailyTask(dailyTask entity.DailyTask) error
	}

	PlayerRepo interface {
		GetPlayerById(id int32) (db.User, bool)
		AddPlayer(player entity.AuthAnswer) (db.User, error)
		UpdateBalance(playerID int32, newBalance float64) (db.User, error)
		GetAllUsers() ([]db.User, error)
		GetUsersByIds(userIDs []int32) ([]db.User, error)
	}

	SeasonRepo interface {
		AddSeason(season entity.Season) error
		GetSeasonById(id int64) (db.Season, error)
		GetAllSeasons() ([]db.Season, error)
		GetSeasonLeaderBoard(seasonID int64) ([]db.GetSeasonLeaderBoardRow, error)
		StartSeason(seasonID int) error
		EndSeason(seasonID int) error
		GetSeasonsByIds(seasonIDs []int32) ([]db.Season, error)
	}
)
