package user

import (
	repo "coreservice/internal/adapter/db/postgres"
	"coreservice/internal/di"
	"coreservice/internal/entity"
	"coreservice/pkg"
	sqlcutils "coreservice/pkg/sqlc_utils"
)

type UseCase struct {
	repo   repo.UserRepo
	logger di.LoggerType
}

func New(repo repo.UserRepo, logger di.LoggerType) *UseCase {
	return &UseCase{
		repo:   repo,
		logger: logger,
	}
}

func (u *UseCase) UserInfo(data any) (entity.User, error) {
	user, err := pkg.ConvertAnyToDBUser(data)
	if err != nil {
		return entity.User{}, err
	}

	return pkg.GetUserInfo(&user), nil
}

func (u *UseCase) MakeDeposit(data any, deposit entity.Balance) (entity.Response, error) {

	player, err := pkg.ConvertAnyToDBUser(data)
	if err != nil {
		return entity.Response{}, err
	}

	// var currentBalance float64
	// if player.Balance.Valid { // Chcecking if balance is NOT NULL

	// 	currentBalance = sqlcutils.NumericToFloat64(player.Balance)

	// 	// left, _ := player.Balance.Int.Float64()
	// 	// currentBalance = left / math.Pow(float64(10), -float64(player.Balance.Exp))
	// }

	//
	currentBalance := sqlcutils.NumericToFloat64(player.Balance)

	player, err = u.repo.UpdateBalance(player.ID, deposit.Balance+currentBalance)
	if err != nil {
		return entity.Response{}, err
	}

	return entity.Response{Message: player.Balance.Int.String()}, nil
}
