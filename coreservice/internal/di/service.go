package di

import (
	"coreservice/internal/entity"
	db "coreservice/internal/repository/sqlc/generated"
)

type Services struct {
	DailyTask DailyTaskService
	Token     TokenService
	User      UserService
	Season    SeasonService
	Search    SearchService
}

type (
	DailyTaskService interface {
		DailyTask() (db.DailyTask, error)
	}

	TokenService interface {
		VerifyToken(message, data any) (entity.TypeResponse, error)
	}

	UserService interface {
		UserInfo(data any) (entity.User, error)
		MakeDeposit(data any, deposit entity.Balance) (entity.Response, error)
	}

	SeasonService interface {
		SeasonInfo(id int) (db.Season, error)
		Seasons() ([]entity.SeasonListElement, error)
		SeasonLeaderBoard(id int) ([]entity.Leaderboard, error)
		CurrentSeason() ([]db.Season, error)
		PlannedSeason() ([]db.Season, error)
	}

	SearchService interface {
		UserBuildSearchIndex() ([]entity.User, error)
		SearchElasticByNameStrict(name string) ([]entity.User, error)
		SearchElasticByNameWildcard(name string) ([]entity.User, error)
		SearchElasticByNameFuzzy(name string) ([]entity.User, error)
	}
)
