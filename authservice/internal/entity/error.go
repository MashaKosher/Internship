package entity

import "errors"

type Error struct {
	Error string `json:"error"`
}

var ErrTokenInCacheNotFound = errors.New("token in cache not found")
