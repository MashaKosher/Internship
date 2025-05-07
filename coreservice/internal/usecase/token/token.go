package token

import "coreservice/internal/entity"

// import (
// 	"coreservice/internal/entity"
// 	"coreservice/pkg"
// 	"log"
// )

type UseCase struct {
}

func New() *UseCase {
	return &UseCase{}
}

func (u *UseCase) VerifyToken(authAnswer, message any) (entity.TypeResponse, error) {

	// log.Println("We are in Use Case !!!!")
	// user, err := pkg.ConvertAnyToDBUser(authAnswer)
	// if err != nil {
	// 	return entity.TypeResponse{}, err
	// }
	// return entity.TypeResponse{User: pkg.GetUserInfo(&user), Message: message.(string)}, nil

	return entity.TypeResponse{}, nil
}
