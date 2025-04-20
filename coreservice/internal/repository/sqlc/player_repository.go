package sqlc

import (
	"coreservice/internal/entity"
	"coreservice/internal/logger"
	db "coreservice/internal/repository/sqlc/generated"
	"fmt"
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"
)

func GetPlayerById(id int32) (db.User, bool) {
	player, err := Query.GetPlayer(Ctx, id)
	return player, err == nil
}

func AddPlayer(player entity.AuthAnswer) (db.User, error) {
	user, err := Query.AddPlayer(Ctx, db.AddPlayerParams{ID: player.ID, Login: player.Login})
	if err != nil {
		logger.Logger.Error("Error while adding user to DB")
		return user, err
	}
	return user, nil
}

func UpdateBalance(playerID int32, newBalance float64) (db.User, error) {

	var num pgtype.Numeric

	num.Int = big.NewInt(int64(newBalance * 100))
	num.Exp = -2
	num.Valid = true

	logger.Logger.Info(fmt.Sprint(num))

	player, err := Query.UpdateBalance(Ctx, db.UpdateBalanceParams{ID: playerID, Column2: num})
	if err != nil {
		return db.User{}, err // Возвращаем ошибку, если она произошла
	}

	return player, nil // Возвращаем обновленного игрока
}

func GetAllUsers() ([]db.User, error) {
	players, err := Query.GetAllPlayer(Ctx)
	if err != nil {
		return []db.User{}, err
	}
	return players, nil
}

func GetUsersByIds(userIDs []int32) ([]db.User, error) {
	users, err := Query.GetUsersByIDs(Ctx, userIDs)
	if err != nil {
		return []db.User{}, err
	}
	return users, nil
}
