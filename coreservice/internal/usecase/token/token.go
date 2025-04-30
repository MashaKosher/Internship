package token

import (
	"coreservice/internal/entity"
	"coreservice/pkg"
)

type UseCase struct {
}

func New() *UseCase {
	return &UseCase{}
}

func (u *UseCase) VerifyToken(message, data any) (entity.TypeResponse, error) {
	user, err := pkg.ConvertAnyToDBUser(data)
	if err != nil {
		return entity.TypeResponse{}, err
	}
	return entity.TypeResponse{User: pkg.GetUserInfo(&user), Message: message.(string)}, nil
}
