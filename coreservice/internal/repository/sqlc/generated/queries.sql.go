// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: queries.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const addPlayer = `-- name: AddPlayer :one
INSERT INTO users (
  id, 
  login,
  balance,
  win_rate
) VALUES (
  $1, $2, $3, $4
)
RETURNING id, login, balance, win_rate
`

type AddPlayerParams struct {
	ID      int32
	Login   string
	Balance pgtype.Numeric
	WinRate pgtype.Numeric
}

func (q *Queries) AddPlayer(ctx context.Context, arg AddPlayerParams) (User, error) {
	row := q.db.QueryRow(ctx, addPlayer,
		arg.ID,
		arg.Login,
		arg.Balance,
		arg.WinRate,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Login,
		&i.Balance,
		&i.WinRate,
	)
	return i, err
}

const getAllPlayer = `-- name: GetAllPlayer :many
SELECT id, login, balance, win_rate FROM users
`

func (q *Queries) GetAllPlayer(ctx context.Context) ([]User, error) {
	rows, err := q.db.Query(ctx, getAllPlayer)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Login,
			&i.Balance,
			&i.WinRate,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPlayer = `-- name: GetPlayer :one
SELECT id, login, balance, win_rate FROM users
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetPlayer(ctx context.Context, id int32) (User, error) {
	row := q.db.QueryRow(ctx, getPlayer, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Login,
		&i.Balance,
		&i.WinRate,
	)
	return i, err
}

const getUsersByIDs = `-- name: GetUsersByIDs :many
SELECT id, login, balance, win_rate FROM users
WHERE id IN (
    SELECT unnest($1::int[])
)
`

func (q *Queries) GetUsersByIDs(ctx context.Context, userIds []int32) ([]User, error) {
	rows, err := q.db.Query(ctx, getUsersByIDs, userIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Login,
			&i.Balance,
			&i.WinRate,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateBalance = `-- name: UpdateBalance :one
UPDATE users
  set balance = $2::numeric
WHERE id = $1
RETURNING id, login, balance, win_rate
`

type UpdateBalanceParams struct {
	ID      int32
	Column2 pgtype.Numeric
}

func (q *Queries) UpdateBalance(ctx context.Context, arg UpdateBalanceParams) (User, error) {
	row := q.db.QueryRow(ctx, updateBalance, arg.ID, arg.Column2)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Login,
		&i.Balance,
		&i.WinRate,
	)
	return i, err
}
