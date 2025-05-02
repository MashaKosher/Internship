package season

import (
	repo "coreservice/internal/adapter/db/postgres"
	elasticRepo "coreservice/internal/adapter/elastic"
	redis "coreservice/internal/adapter/redis"
	"coreservice/internal/di"
	"coreservice/internal/entity"
	db "coreservice/internal/repository/sqlc/generated"
	"coreservice/pkg"
	"errors"
	"fmt"
)

type UseCase struct {
	repo                    repo.SeasonRepo
	logger                  di.LoggerType
	elasticSeasonStatusRepo elasticRepo.SeasonStatusRepo
	redis                   redis.LeaderBoardCache
}

func New(repo repo.SeasonRepo, logger di.LoggerType, elasticSeasonStatusRepo elasticRepo.SeasonStatusRepo, redis redis.LeaderBoardCache) *UseCase {
	return &UseCase{
		repo:                    repo,
		logger:                  logger,
		elasticSeasonStatusRepo: elasticSeasonStatusRepo,
		redis:                   redis,
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

	leaderboard, err := u.redis.GetSeasonLeaderBoard(id)
	if err != nil {
		return []entity.Leaderboard{}, errors.New("problems with redis: " + err.Error())
	}

	u.logger.Info("LeaderBoard from Redis: " + fmt.Sprint(leaderboard))

	if len(leaderboard) == 0 {
		leaderBoard, err := u.repo.GetSeasonLeaderBoard(int64(id))
		if err != nil {
			return []entity.Leaderboard{}, err
		}

		leaderboard = pkg.ConvertLeaderBoardDBListToJson(leaderBoard)

		u.logger.Info("LeaderBoard from DB: " + fmt.Sprint(leaderboard))

		err = u.redis.UpdateLeaderBoard(leaderboard, id)
		if err != nil {
			return []entity.Leaderboard{}, errors.New("some troubles with adding info to redis")
		}

	}

	return leaderboard, nil
}

func (u *UseCase) CurrentSeason() ([]db.Season, error) {
	ids, err := u.elasticSeasonStatusRepo.ActiveSeason()
	if err != nil {
		return []db.Season{}, err
	}

	u.logger.Info("Elastic found current season index: " + fmt.Sprint(ids))

	if len(ids) == 0 {
		return []db.Season{}, errors.New("there is no active season right now")
	}

	seasons, err := u.repo.GetSeasonsByIds(ids)
	if err != nil {
		return []db.Season{}, err
	}

	return seasons, nil
}

func (u *UseCase) PlannedSeason() ([]db.Season, error) {
	ids, err := u.elasticSeasonStatusRepo.PlannedSeasons()
	if err != nil {
		return []db.Season{}, err
	}
	u.logger.Info("Elastic found planned seasons indexes: " + fmt.Sprint(ids))

	if len(ids) == 0 {
		return []db.Season{}, errors.New("there is no planned season right now")
	}

	seasons, err := u.repo.GetSeasonsByIds(ids)
	if err != nil {
		return []db.Season{}, err
	}

	return seasons, nil
}
