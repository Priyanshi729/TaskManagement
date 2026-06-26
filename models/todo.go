package models

type ClientError struct {
	MessageToUser string `json:"messageToUser"`
	Err           string `json:"error"`
	StatusCode    int    `json:"statusCode"`
}

type TodoRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type Todo struct {
	ID          string `json:"id" db:"id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	IsCompleted bool   `json:"done" db:"is_completed"`
}
