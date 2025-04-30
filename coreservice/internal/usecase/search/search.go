package search

import (
	"coreservice/internal/adapter/elastic"
	"coreservice/internal/di"
	"coreservice/internal/entity"
	"fmt"
)

type UseCase struct {
	logger di.LoggerType
}

func New(logger di.LoggerType) *UseCase {
	return &UseCase{
		logger: logger,
	}
}

func (u *UseCase) UserBuildSearchIndex() ([]entity.User, error) {
	users, err := elastic.AddingAllUsersToIndex()
	if err != nil {
		return []entity.User{}, err
	}
	u.logger.Info("Users from DB: " + fmt.Sprint(users))
	return users, nil
}

func (u *UseCase) SearchElasticByNameStrict(name string) ([]entity.User, error) {
	users, err := elastic.GetUserByName(name, elastic.Strict)
	if err != nil {

		return []entity.User{}, err
	}
	return users, nil
}

func (u *UseCase) SearchElasticByNameWildcard(name string) ([]entity.User, error) {
	users, err := elastic.GetUserByName(name, elastic.Wildcard)
	if err != nil {
		return []entity.User{}, err
	}
	return users, nil
}

func (u *UseCase) SearchElasticByNameFuzzy(name string) ([]entity.User, error) {
	users, err := elastic.GetUserByName(name, elastic.Fuzzy)
	if err != nil {
		return []entity.User{}, err
	}
	return users, nil
}
