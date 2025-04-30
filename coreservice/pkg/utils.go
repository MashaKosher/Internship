package pkg

import (
	"coreservice/internal/entity"
	db "coreservice/internal/repository/sqlc/generated"
	"errors"
	"math"

	"github.com/jackc/pgx/v5/pgtype"
)

func GetUserInfo(user *db.User) entity.User {

	var jsonUser entity.User

	jsonUser.Login = user.Login
	jsonUser.Balance = ValidateNumeric(user.Balance)
	jsonUser.WinRate = ValidateNumeric(user.WinRate)

	return jsonUser

}

func ConvertDBUserSliceToUser(users []db.User) []entity.User {
	result := make([]entity.User, 0, len(users))
	for _, user := range users {
		result = append(result, GetUserInfo(&user))
	}
	return result
}

func ValidateNumeric(entity pgtype.Numeric) float64 {
	var res float64
	if !entity.Valid { // if field is NULL in DB
		return res
	}
	res, _ = entity.Int.Float64()

	return res / math.Pow(float64(10), -float64(entity.Exp))
}

func ConvertAnyToDBUser(entity any) (db.User, error) {
	user, ok := entity.(db.User)
	if !ok {
		return db.User{}, errors.New("value cannot be converted to db.User")
	}

	return user, nil
}
