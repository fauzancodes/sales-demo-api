package dto

import "github.com/google/uuid"

type Response struct {
	Status  int         `json:"status,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type FindParameter struct {
	BaseFilter       string
	BaseFilterValues []any
	Filter           string
	FilterValues     []any
	Order            string
	Limit            int
	Offset           int
}

type GlobalIDNameResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
