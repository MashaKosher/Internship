package user

import (
	"context"
	"coreservice/internal/entity"
	db "coreservice/internal/repository/sqlc/generated"

	// "math/big"

	// "github.com/jackc/pgx/v5/pgtype"

	utils "coreservice/pkg/sqlc_utils"
)

type UserRepo struct {
	Query *db.Queries
}

func New(queries *db.Queries) *UserRepo {
	if queries == nil {
		panic("queries is nil")
	}

	return &UserRepo{
		Query: queries,
	}
}

func (r *UserRepo) GetPlayerById(id int32) (db.User, bool) {
	player, err := r.Query.GetPlayer(context.Background(), id)
	return player, err == nil
}

func (r *UserRepo) AddPlayer(player entity.AuthAnswer) (db.User, error) {
	user, err := r.Query.AddPlayer(context.Background(), db.AddPlayerParams{ID: player.ID, Login: player.Login})
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *UserRepo) UpdateBalance(playerID int32, newBalance float64) (db.User, error) {

	// var num pgtype.Numeric

	// num.Int = big.NewInt(int64(newBalance * 100))
	// num.Exp = -2
	// num.Valid = true

	num, err := utils.NumberToNumeric(newBalance)
	if err != nil {
		return db.User{}, err
	}

	player, err := r.Query.UpdateBalance(context.Background(), db.UpdateBalanceParams{ID: playerID, Column2: num})
	if err != nil {
		return db.User{}, err
	}

	return player, nil
}

func (r *UserRepo) GetAllUsers() ([]db.User, error) {
	players, err := r.Query.GetAllPlayer(context.Background())
	if err != nil {
		return []db.User{}, err
	}
	return players, nil
}

func (r *UserRepo) GetUsersByIds(userIDs []int32) ([]db.User, error) {
	users, err := r.Query.GetUsersByIDs(context.Background(), userIDs)
	if err != nil {
		return []db.User{}, err
	}
	return users, nil
}
