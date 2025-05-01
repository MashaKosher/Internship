package search

import (
	elasticRepo "coreservice/internal/adapter/elastic"
	"coreservice/internal/di"
	"coreservice/internal/entity"
	"fmt"
)

type UseCase struct {
	logger              di.LoggerType
	elasticUserNameRepo elasticRepo.UserNameRepo
}

func New(logger di.LoggerType, elasticUserNameRepo elasticRepo.UserNameRepo) *UseCase {
	return &UseCase{
		logger:              logger,
		elasticUserNameRepo: elasticUserNameRepo,
	}
}

func (u *UseCase) UserBuildSearchIndex() ([]entity.User, error) {
	users, err := u.elasticUserNameRepo.AddingAllUsersToIndex()
	if err != nil {
		return []entity.User{}, err
	}
	u.logger.Info("Users from DB: " + fmt.Sprint(users))
	return users, nil
}

func (u *UseCase) SearchElasticByNameStrict(name string) ([]entity.User, error) {
	users, err := u.elasticUserNameRepo.GetUserByNameStrict(name)
	if err != nil {

		return []entity.User{}, err
	}
	return users, nil
}

func (u *UseCase) SearchElasticByNameWildcard(name string) ([]entity.User, error) {
	users, err := u.elasticUserNameRepo.GetUserByNameWildcard(name)
	if err != nil {
		return []entity.User{}, err
	}
	return users, nil
}

func (u *UseCase) SearchElasticByNameFuzzy(name string) ([]entity.User, error) {
	users, err := u.elasticUserNameRepo.GetUserByNameWildcard(name)
	if err != nil {
		return []entity.User{}, err
	}
	return users, nil
}
