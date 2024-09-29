package dto

import "github.com/google/uuid"

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error"`
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
