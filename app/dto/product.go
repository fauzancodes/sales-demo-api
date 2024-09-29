package dto

import (
	"github.com/fauzancodes/sales-demo-api/app/models"
	validation "github.com/go-ozzo/ozzo-validation"
)

type ProductRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Status      bool     `json:"status"`
	Image       []string `json:"image"`
	Price       float64  `json:"price"`
	CategoryID  string   `json:"category_id"`
}

func (request ProductRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Name, validation.Required),
	)
}

type ProductResponse struct {
	models.CustomGormModel
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Status      bool                 `json:"status"`
	Image       []string             `json:"image"`
	Price       float64              `json:"price"`
	Stock       int                  `json:"stock"`
	Category    GlobalIDNameResponse `json:"category"`
}
