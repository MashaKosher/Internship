package entity

import "errors"

type Response struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"message"`
}

// Custom DB Errors
var (
	// General
	ErrRecordNotFound = errors.New("record not found in DB")

	// Daily Task
	ErrDailytaskIsNil         = errors.New("daily task entity is nil")
	ErrDailyTaskAlreadyExists = errors.New("there is already task for today")

	//Plan Season
	ErrSeasonIsNil        = errors.New("season entity is nil")
	ErrSeasonsAreCrossing = errors.New("seasons are crossing")

	//Game Settings
	ErrGameSettingsIsNil = errors.New("game settings entity is nil")
)
