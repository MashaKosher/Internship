package pkg

import (
	"coreservice/internal/entity"
	db "coreservice/internal/repository/sqlc/generated"
)

func ConvertSeasonDBListToJson(seasons []db.Season) []entity.SeasonListElement {
	result := make([]entity.SeasonListElement, 0, len(seasons))
	for _, season := range seasons {
		result = append(result, convertSeasonDBElemToJson(&season))
	}
	return result
}

func convertSeasonDBElemToJson(season *db.Season) entity.SeasonListElement {
	return entity.SeasonListElement{ID: uint(season.ID), Satatus: season.SeasonStatus.String}
}

func ConvertLeaderBoardDBListToJson(leaderbord []db.GetSeasonLeaderBoardRow) []entity.Leaderboard {
	result := make([]entity.Leaderboard, 0, len(leaderbord))
	for _, record := range leaderbord {
		result = append(result, convertLeaderBoardRecordDBElemToJson(&record))
	}

	return result
}

func convertLeaderBoardRecordDBElemToJson(record *db.GetSeasonLeaderBoardRow) entity.Leaderboard {
	return entity.Leaderboard{UserID: record.UserID, Win: record.Win}
}
