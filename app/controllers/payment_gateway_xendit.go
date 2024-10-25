package controllers

import (
	"log"
	"net/http"

	"github.com/fauzancodes/sales-demo-api/app/dto"
	"github.com/fauzancodes/sales-demo-api/app/pkg/utils"
	"github.com/fauzancodes/sales-demo-api/app/service"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func GetXenditPaymentMethods(c echo.Context) error {
	code := c.QueryParam("code")

	param := utils.PopulatePaging(c, "")
	data, _, err := service.GetXenditPaymentMethods(code, param)
	if err != nil {
		return c.JSON(
			http.StatusNotFound,
			dto.Response{
				Status:  404,
				Message: "Failed to get data",
				Error:   err.Error(),
			},
		)
	}

	return c.JSON(http.StatusOK, data)
}

func XenditCharge(c echo.Context) error {
	userID := c.Get("currentUser").(jwt.MapClaims)["id"].(string)
	log.Printf("Current user ID: %v", userID)

	var request dto.XenditRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(
			http.StatusUnprocessableEntity,
			dto.Response{
				Status:  422,
				Message: "Invalid request body",
				Error:   err.Error(),
			},
		)
	}

	if err := request.Validate(); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			dto.Response{
				Status:  400,
				Message: "Invalid request value",
				Error:   err.Error(),
			},
		)
	}

	response, err := service.XenditCharge(userID, utils.GetBaseUrl(c), request)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			dto.Response{
				Status:  500,
				Message: "Failed to charge payment to xendit",
				Error:   err.Error(),
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		dto.Response{
			Status:  200,
			Message: "Success to charge payment to midtrans",
			Data:    response,
		},
	)
}

func XenditNotification(c echo.Context) error {
	callbackToken := c.Request().Header.Get("x-callback-token")

	var request dto.XenditNotificationRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(
			http.StatusUnprocessableEntity,
			dto.Response{
				Status:  422,
				Message: "Invalid request body",
				Error:   err.Error(),
			},
		)
	}

	err := service.XenditHandleNotification(request, callbackToken)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			dto.Response{
				Status:  500,
				Message: "Failed to handle xendit notification",
				Error:   err.Error(),
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		dto.Response{
			Status:  200,
			Message: "Success to process xendit notification",
			Data:    c.Request().Body,
		},
	)
}
