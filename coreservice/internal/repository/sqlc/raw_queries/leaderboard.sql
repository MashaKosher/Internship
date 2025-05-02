-- name: UpdateSeasonLeaderBoard :exec
INSERT INTO leaderboard (season_id, user_id, win)
VALUES ($1, $2, 1)
ON CONFLICT (season_id, user_id) 
DO UPDATE SET win = leaderboard.win + 1
RETURNING *;