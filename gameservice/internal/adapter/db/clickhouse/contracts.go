package clickhouse

import "gameservice/internal/entity"

type (
	GameRepo interface {
		AddGame(result entity.GameResult) error
		GetPlayerStatistic(playerID int) (entity.PlayerStats, error)
		CreateGameResultsTable() error
	}
)
