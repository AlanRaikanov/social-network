package models

type SuccessResponse struct {
	Response string `json:"response"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
