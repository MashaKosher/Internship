package entity

import "errors"

type Response struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"message"`
}

var ErrRecordNotFound = errors.New("record not found in DB")
