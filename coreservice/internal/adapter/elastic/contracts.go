package elastic

import (
	"coreservice/internal/entity"
)

type (
	UserNameRepo interface {
		AddingAllUsersToIndex() ([]entity.User, error)
		AddUserToIndex(user entity.User, userId int) error
		GetUserByNameStrict(name string) ([]entity.User, error)
		GetUserByNameWildcard(name string) ([]entity.User, error)
		GetUserByNameFuzzy(name string) ([]entity.User, error)
	}

	SeasonStatusRepo interface {
		AddSeasonToIndex(seasonID int) error
		StartSeason(seasonID int) error
		EndSeason(seasonID int) error
		ActiveSeason() ([]int32, error)
		PlannedSeasons() ([]int32, error)
	}
)
