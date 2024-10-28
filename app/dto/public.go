package dto

import "github.com/google/uuid"

type Response struct {
	Status  int         `json:"status,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type FindParameter struct {
	BaseFilter string
	Filter     string
	Order      string
	Limit      int
	Offset     int
}

type GlobalIDNameResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
