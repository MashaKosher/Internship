-- name: GetDailyTask :one
SELECT * FROM daily_tasks 
WHERE task_date = $1 LIMIT 1;



-- name: AddDailyTask :exec
INSERT INTO daily_tasks (task_date, referals_amount, wins_amount, referals_reward, win_reward)
VALUES ($1, $2, $3, $4, $5);