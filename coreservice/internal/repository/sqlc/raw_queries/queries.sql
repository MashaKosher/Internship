-- name: GetPlayer :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;


-- name: GetAllPlayer :many
SELECT * FROM users;


-- name: AddPlayer :one
INSERT INTO users (
  id, 
  login,
  balance,
  win_rate
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;



-- name: UpdateBalance :one
UPDATE users
  set balance = $2::numeric
WHERE id = $1
RETURNING *;


-- name: GetUsersByIDs :many
SELECT * FROM users
WHERE id IN (
    SELECT unnest(sqlc.arg('user_ids')::int[])
);



