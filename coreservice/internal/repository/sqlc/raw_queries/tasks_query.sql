-- name: GetDailyTask :one
SELECT * FROM daily_tasks 
WHERE task_date = $1 LIMIT 1;



-- name: AddDailyTask :exec
INSERT INTO daily_tasks (task_date, referals_amount, wins_amount)
VALUES ($1, $2, $3);