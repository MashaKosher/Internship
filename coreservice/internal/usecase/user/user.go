package user

import (
	repo "coreservice/internal/adapter/db/postgres"
	"coreservice/internal/di"
	"coreservice/internal/entity"
	"coreservice/pkg"
	sqlcutils "coreservice/pkg/sqlc_utils"
	"fmt"
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

	user, err := pkg.ConvertAnyToDBUser(data)
	if err != nil {
		return entity.Response{}, err
	}

	currentBalance := sqlcutils.NumericToFloat64(user.Balance)

	user, err = u.repo.UpdateBalance(user.ID, deposit.Balance+currentBalance)
	if err != nil {
		return entity.Response{}, err
	}

	return entity.Response{Message: fmt.Sprint(sqlcutils.NumericToFloat64(user.Balance))}, nil
}
