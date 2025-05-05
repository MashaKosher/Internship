package entity

import "errors"

type Error struct {
	Error string `json:"error"`
}

var ErrNoDailyTask = errors.New("there is no tasks for today")
