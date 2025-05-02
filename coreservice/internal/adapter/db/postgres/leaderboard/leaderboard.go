package leaderboard

import (
	"context"
	db "coreservice/internal/repository/sqlc/generated"
)

type LeaderboardRepo struct {
	Query *db.Queries
}

func New(queries *db.Queries) *LeaderboardRepo {
	if queries == nil {
		panic("queries is nil")
	}

	return &LeaderboardRepo{
		Query: queries,
	}
}

func (r *LeaderboardRepo) UpdateSeasonLeaderboard(seasonID, playerID int) error {

	return r.Query.UpdateSeasonLeaderBoard(context.Background(), db.UpdateSeasonLeaderBoardParams{
		SeasonID: int64(seasonID),
		UserID:   int32(playerID),
	})
}
