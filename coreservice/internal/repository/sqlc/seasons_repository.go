package sqlc

import (
	"coreservice/internal/entity"
	"coreservice/internal/logger"
	"time"

	db "coreservice/internal/repository/sqlc/generated"

	"github.com/jackc/pgx/v5/pgtype"
)

func AddSeason(season entity.Season) error {

	// dateStr := "2026-12-10 01:00:00 +0000 UTC"

	// Определяем формат строки (layout)
	layout := "2006-01-02 15:04:05 -0700 MST"

	startTime, err := time.Parse(layout, season.StartDate)
	if err != nil {
		logger.Logger.Error("Ошибка парсинга:" + err.Error())
		return err
	}

	endTime, err := time.Parse(layout, season.EndDate)
	if err != nil {
		logger.Logger.Error("Ошибка парсинга:" + err.Error())
		return err
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

	err = Query.CreateSeason(Ctx, db.CreateSeasonParams{ID: int64(season.ID), SeasonStart: start, SeasonEnd: end, SeasonFund: fund})
	if err != nil {
		logger.Logger.Error("Error while adding season to DB" + err.Error())
	}

	return nil

}

func GetSeasonById(id int64) (db.Season, error) {
	season, err := Query.GetSeason(Ctx, id)
	return season, err
}

func GetAllSeasons() ([]db.Season, error) {
	seasons, err := Query.GetAllSeasons(Ctx)
	return seasons, err
}

func GetSeasonLeaderBoard(seasonID int64) ([]db.GetSeasonLeaderBoardRow, error) {
	leaderboard, err := Query.GetSeasonLeaderBoard(Ctx, seasonID)
	return leaderboard, err
}

func StartSeason(seasonID int) error {
	err := Query.StartSeason(Ctx, int64(seasonID))
	return err
}

func EndSeason(seasonID int) error {
	err := Query.EndSeason(Ctx, int64(seasonID))
	return err
}

func GetSeasonsByIds(seasonIDs []int32) ([]db.Season, error) {
	seasons, err := Query.GetSeasonsByID(Ctx, seasonIDs)
	if err != nil {
		return []db.Season{}, err
	}
	return seasons, nil
}
