package redis

import "coreservice/internal/entity"

type (
	LeaderBoardCache interface {
		UpdateLeaderBoard(leaderBoard []entity.Leaderboard, seasonID int) error
		GetSeasonLeaderBoard(seasonID int) ([]entity.Leaderboard, error)
	}
)
