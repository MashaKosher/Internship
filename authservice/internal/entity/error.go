package entity

import "errors"

type Error struct {
	Error string `json:"error"`
}

// DB custom erors
var (
	ErrTokenInCacheNotFound  = errors.New("token in cache not found")
	ErrUserIsNil             = errors.New("user is nil")
	ErrUserNameCannotBeEmpty = errors.New("username cannot be empty")
	ErrUserAlreadyExists     = errors.New("user with such username already exists")
	ErrInvalidUserID         = errors.New("user id cannot be less then 1")
	ErrUserNotFoundInDB      = errors.New("user not found in DB")
	ErrPasswordCannotBeEmpty = errors.New("password cannot be empty")
)

// Token Errors
var (
	ErrTokenExpired         = errors.New("token expired")
	ErrInvalidSigningMethod = errors.New("invalid signing method")
	ErrInvalidTokenType     = errors.New("invalid Token type")
	ErrInavlidID            = errors.New("invalid id value")
)
