package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/fauzancodes/sales-demo-api/app/dto"
	"github.com/fauzancodes/sales-demo-api/app/pkg/utils"
	"github.com/fauzancodes/sales-demo-api/app/service"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func CreateProductStock(c echo.Context) error {
	userID := c.Get("currentUser").(jwt.MapClaims)["id"].(string)
	log.Printf("Current user ID: %v", userID)

	var request dto.ProductStockRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(
			http.StatusUnprocessableEntity,
			dto.Response{
				Status:  http.StatusUnprocessableEntity,
				Message: "Invalid request body",
				Error:   err.Error(),
			},
		)
	}

	if err := request.Validate(); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			dto.Response{
				Status:  http.StatusBadRequest,
				Message: "Invalid request value",
				Error:   err.Error(),
			},
		)
	}

	result, statusCode, err := service.CreateProductStock(userID, request)
	if err != nil {
		return c.JSON(
			statusCode,
			dto.Response{
				Status:  statusCode,
				Message: "Failed to add stock",
				Error:   err.Error(),
			},
		)
	}

	return c.JSON(
		statusCode,
		dto.Response{
			Status:  statusCode,
			Message: "Success to add stock",
			Data:    result,
		},
	)
}

func GetProductStocks(c echo.Context) error {
	withUser, _ := strconv.ParseBool(c.QueryParam("with_user"))
	withProduct, _ := strconv.ParseBool(c.QueryParam("with_product"))

	userID := c.Get("currentUser").(jwt.MapClaims)["id"].(string)
	log.Printf("Current user ID: %v", userID)

	productID := c.QueryParam("product_id")

	var preloadFields []string
	if withUser {
		preloadFields = append(preloadFields, "User")
	}
	if withProduct {
		preloadFields = append(preloadFields, "Product")
	}

	param := utils.PopulatePaging(c, "status")
	data, _, statusCode, err := service.GetProductStocks(productID, userID, param, preloadFields)
	if err != nil {
		return c.JSON(
			statusCode,
			dto.Response{
				Status:  statusCode,
				Message: "Failed to get data",
				Error:   err.Error(),
			},
		)
	}

	return c.JSON(statusCode, data)
}
