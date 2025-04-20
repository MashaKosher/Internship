-- name: CreateSeason :exec
INSERT INTO seasons (id, season_start, season_end, season_fund, season_status)
VALUES ($1, $2, $3, $4, 'planned');


-- name: GetAllSeasons :many
SELECT * FROM seasons;

-- name: GetSeason :one
SELECT * FROM seasons 
WHERE id = $1 LIMIT 1;


-- name: GetSeasonsByID :many
SELECT * FROM seasons 
WHERE id IN (
    SELECT unnest(sqlc.arg('season_ids')::int[])
);



-- name: GetSeasonLeaderBoard :many
WITH cte AS (
    SELECT * FROM leaderboard
    WHERE season_id = $1
)
SELECT * FROM cte ORDER BY win DESC;


-- name: StartSeason :exec
UPDATE seasons SET season_status = 'current' WHERE id = $1;

-- name: EndSeason :exec
UPDATE seasons SET season_status = 'ended' WHERE id = $1;