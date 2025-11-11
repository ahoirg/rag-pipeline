package models

type AskRequest struct {
	Query string `json:"query" validate:"required,min=3"`
}
