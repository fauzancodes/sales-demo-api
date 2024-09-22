package controllers

import (
	"net/http"

	"github.com/fauzancodes/sales-demo-api/app/dto"
	"github.com/fauzancodes/sales-demo-api/app/service"
	"github.com/fauzancodes/sales-demo-api/pkg/utils"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo/v4"
)

func Register(c echo.Context) error {
	var request dto.AuthRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(
			http.StatusUnprocessableEntity,
			dto.ErrorResponse{
				Status:  422,
				Message: "Invalid request body",
				Error:   err.Error(),
			},
		)
	}

	if err := request.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(
			http.StatusBadRequest,
			dto.ErrorResponse{
				Status:  400,
				Message: "Invalid request value",
				Error:   errVal.Error(),
			},
		)
	}

	param := utils.PopulatePaging(c, "")
	_, check, _ := service.GetUsers("", "", request.Email, param)
	if len(check) > 0 {
		return c.JSON(
			http.StatusBadRequest,
			dto.ErrorResponse{
				Status:  400,
				Message: "Email has been registered",
				Error:   "Email found",
			},
		)
	}

	result, err := service.CreateUser(dto.UserRequest{
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			dto.ErrorResponse{
				Status:  500,
				Message: "Failed to register",
				Error:   err,
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		dto.SuccessResponse{
			Status:  200,
			Message: "Success to register",
			Data:    result,
		},
	)
}
