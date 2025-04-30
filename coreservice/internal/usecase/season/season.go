package season

import (
	repo "coreservice/internal/adapter/db/postgres"
	"coreservice/internal/adapter/elastic"
	"coreservice/internal/di"
	"coreservice/internal/entity"
	db "coreservice/internal/repository/sqlc/generated"
	"coreservice/pkg"
	"fmt"
)

type UseCase struct {
	repo   repo.SeasonRepo
	logger di.LoggerType
}

func New(repo repo.SeasonRepo, logger di.LoggerType) *UseCase {
	return &UseCase{
		repo:   repo,
		logger: logger,
	}
}

func (u *UseCase) SeasonInfo(id int) (db.Season, error) {
	season, err := u.repo.GetSeasonById(int64(id))
	if err != nil {
		return db.Season{}, err
	}
	return season, nil
}

func (u *UseCase) Seasons() ([]entity.SeasonListElement, error) {
	seasons, err := u.repo.GetAllSeasons()
	if err != nil {
		return []entity.SeasonListElement{}, err
	}

	return pkg.ConvertSeasonDBListToJson(seasons), nil
}

func (u *UseCase) SeasonLeaderBoard(id int) ([]entity.Leaderboard, error) {
	leaderBoard, err := u.repo.GetSeasonLeaderBoard(int64(id))
	if err != nil {
		return []entity.Leaderboard{}, err
	}
	return pkg.ConvertLeaderBoardDBListToJson(leaderBoard), nil
}

func (u *UseCase) CurrentSeason() ([]db.Season, error) {
	ids, err := elastic.SearchSeasonsByStatus(elastic.CurrentSeason)
	if err != nil {
		return []db.Season{}, err
	}

	u.logger.Info("Elastic found current season index: " + fmt.Sprint(ids))

	seasons, err := u.repo.GetSeasonsByIds(ids)
	if err != nil {
		return []db.Season{}, err
	}

	return seasons, nil
}

func (u *UseCase) PlannedSeason() ([]db.Season, error) {
	ids, err := elastic.SearchSeasonsByStatus(elastic.PlannedSeason)
	if err != nil {
		return []db.Season{}, err
	}
	u.logger.Info("Elastic found planned seasons indexes: " + fmt.Sprint(ids))

	seasons, err := u.repo.GetSeasonsByIds(ids)
	if err != nil {
		return []db.Season{}, err
	}

	return seasons, nil
}
