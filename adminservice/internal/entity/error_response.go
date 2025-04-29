package entity

type Response struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"message"`
}
