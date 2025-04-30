package season

import (
	"context"
	"coreservice/internal/entity"
	db "coreservice/internal/repository/sqlc/generated"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type SeasonRepo struct {
	Query *db.Queries
}

func New(queries *db.Queries) *SeasonRepo {
	if queries == nil {
		panic("queries is nil")
	}

	return &SeasonRepo{
		Query: queries,
	}
}

func (r *SeasonRepo) AddSeason(season entity.Season) error {
	layout := "2006-01-02 15:04:05 -0700 MST"

	startTime, err := time.Parse(layout, season.StartDate)
	if err != nil {
		return errors.New("ошибка парсинга даты начала сезона")
	}

	endTime, err := time.Parse(layout, season.EndDate)
	if err != nil {
		return errors.New("ошибка парсинга даты конца сезона")
	}

	start := pgtype.Timestamptz{
		Time:  startTime,
		Valid: true,
	}

	end := pgtype.Timestamptz{
		Time:  endTime,
		Valid: true,
	}

	fund := pgtype.Int4{
		Int32: int32(season.Fund),
		Valid: true, // Не забудьте установить Valid в true
	}

	err = r.Query.CreateSeason(context.Background(), db.CreateSeasonParams{ID: int64(season.ID), SeasonStart: start, SeasonEnd: end, SeasonFund: fund})
	if err != nil {
		return errors.New("error while adding season to DB" + err.Error())
	}

	return nil

}

func (r *SeasonRepo) GetSeasonById(id int64) (db.Season, error) {
	season, err := r.Query.GetSeason(context.Background(), id)
	return season, err
}

func (r *SeasonRepo) GetAllSeasons() ([]db.Season, error) {
	seasons, err := r.Query.GetAllSeasons(context.Background())
	return seasons, err
}

func (r *SeasonRepo) GetSeasonLeaderBoard(seasonID int64) ([]db.GetSeasonLeaderBoardRow, error) {
	leaderboard, err := r.Query.GetSeasonLeaderBoard(context.Background(), seasonID)
	return leaderboard, err
}

func (r *SeasonRepo) StartSeason(seasonID int) error {
	err := r.Query.StartSeason(context.Background(), int64(seasonID))
	return err
}

func (r *SeasonRepo) EndSeason(seasonID int) error {
	err := r.Query.EndSeason(context.Background(), int64(seasonID))
	return err
}

func (r *SeasonRepo) GetSeasonsByIds(seasonIDs []int32) ([]db.Season, error) {
	seasons, err := r.Query.GetSeasonsByID(context.Background(), seasonIDs)
	if err != nil {
		return []db.Season{}, err
	}
	return seasons, nil
}
