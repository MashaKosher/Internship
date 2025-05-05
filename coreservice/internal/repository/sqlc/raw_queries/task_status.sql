
-- name: AddWin :one
INSERT INTO task_status_tbl (task_date, user_id, win)
VALUES ($1, $2, 1)
ON CONFLICT (task_date, user_id)
DO UPDATE SET win = task_status_tbl.win + 1
RETURNING win;


-- name: AddReferal :one
INSERT INTO task_status_tbl (task_date, user_id, referals)
VALUES ($1, $2, 1)
ON CONFLICT (task_date, user_id)
DO UPDATE SET referals = task_status_tbl.referals + 1
RETURNING referals;



-- name: CompleteWinTask :exec
UPDATE task_status_tbl 
SET win_status = TRUE 
WHERE task_date = $1 AND user_id = $2;


-- name: CompleteReferalsTask :exec
UPDATE task_status_tbl 
SET referals_status = TRUE 
WHERE task_date = $1 AND user_id = $2;


-- name: WinTaskStatus :one
SELECT win_status FROM task_status_tbl WHERE task_date = $1 AND user_id = $2;


-- name: ReferalsTaskStatus :one
SELECT referals_status FROM task_status_tbl WHERE task_date = $1 AND user_id = $2;