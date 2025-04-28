package game

import (
	"context"
	"database/sql"
	"fmt"
	"gameservice/internal/entity"
)

// Переписать
type GameRepo struct {
	*sql.DB
}

// Переписать
func New(db *sql.DB) *GameRepo {
	return &GameRepo{db}
}

func (r *GameRepo) AddGame(result entity.GameResult) error {
	ctx := context.Background()

	query := `
	INSERT INTO game_results (
		game_time,
		winner,
		loser,
		win_amount,
		lose_amount,
		winner_result,
		loser_result
	) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.DB.ExecContext(ctx, query,
		result.GameTime,
		result.Winner,
		result.Loser,
		result.WinAmount,
		result.LoseAmount,
		result.WinnerResult,
		result.LoserResult,
	)
	return err
}

func (r *GameRepo) GetPlayerStatistic(playerID int) (entity.PlayerStats, error) {
	ctx := context.Background()
	query := `
	SELECT 
		? AS player_id,
		count() AS total_games,
		sum(winner = ?) AS total_wins,
		sum(loser = ?) AS total_losses
	FROM (
		SELECT winner, loser FROM game_results
		WHERE winner = ? OR loser = ?
	)
	`

	var playerStats entity.PlayerStats
	err := r.DB.QueryRowContext(ctx, query,
		playerID,
		playerID,
		playerID,
		playerID,
		playerID,
	).Scan(
		&playerStats.PlayerID,
		&playerStats.TotalGames,
		&playerStats.TotalWins,
		&playerStats.TotalLosses,
	)

	if err != nil {
		return playerStats, fmt.Errorf("failed to get player stats: %w", err)
	}

	return playerStats, nil
}

func (r *GameRepo) CreateGameResultsTable() error {
	ctx := context.Background()

	query := `
	CREATE TABLE IF NOT EXISTS game_results (
		game_time      DateTime('UTC')  DEFAULT now(),
		winner         Int,
		loser          Int,
		win_amount     Int,
		lose_amount    Int,
		winner_result  Int,
		loser_result   Int,
		
		-- Дополнительные настройки движка
		INDEX winner_idx winner TYPE bloom_filter GRANULARITY 3,
		INDEX loser_idx loser TYPE bloom_filter GRANULARITY 3
	)
	ENGINE = MergeTree()
	PARTITION BY toYYYYMM(game_time)
	ORDER BY (game_time, winner, loser)
	SETTINGS index_granularity = 8192
	`

	_, err := r.DB.ExecContext(ctx, query)
	return err
}
