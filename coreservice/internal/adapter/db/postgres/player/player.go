package player

import (
	"context"
	"coreservice/internal/entity"
	"coreservice/internal/logger"
	db "coreservice/internal/repository/sqlc/generated"
	"fmt"
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"
)

type PlayerRepo struct {
	Query *db.Queries
}

func New(queries *db.Queries) *PlayerRepo {
	if queries == nil {
		panic("queries is nil")
	}

	return &PlayerRepo{
		Query: queries,
	}
}

func (r *PlayerRepo) GetPlayerById(id int32) (db.User, bool) {
	player, err := r.Query.GetPlayer(context.Background(), id)
	return player, err == nil
}

func (r *PlayerRepo) AddPlayer(player entity.AuthAnswer) (db.User, error) {
	user, err := r.Query.AddPlayer(context.Background(), db.AddPlayerParams{ID: player.ID, Login: player.Login})
	if err != nil {
		logger.Logger.Error("Error while adding user to DB")
		return user, err
	}
	return user, nil
}

func (r *PlayerRepo) UpdateBalance(playerID int32, newBalance float64) (db.User, error) {

	var num pgtype.Numeric

	num.Int = big.NewInt(int64(newBalance * 100))
	num.Exp = -2
	num.Valid = true

	logger.Logger.Info(fmt.Sprint(num))

	player, err := r.Query.UpdateBalance(context.Background(), db.UpdateBalanceParams{ID: playerID, Column2: num})
	if err != nil {
		return db.User{}, err // Возвращаем ошибку, если она произошла
	}

	return player, nil // Возвращаем обновленного игрока
}

func (r *PlayerRepo) GetAllUsers() ([]db.User, error) {
	players, err := r.Query.GetAllPlayer(context.Background())
	if err != nil {
		return []db.User{}, err
	}
	return players, nil
}

func (r *PlayerRepo) GetUsersByIds(userIDs []int32) ([]db.User, error) {
	users, err := r.Query.GetUsersByIDs(context.Background(), userIDs)
	if err != nil {
		return []db.User{}, err
	}
	return users, nil
}
