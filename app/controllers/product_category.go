package controllers

import (
	"log"
	"net/http"

	"github.com/fauzancodes/sales-demo-api/app/dto"
	"github.com/fauzancodes/sales-demo-api/app/service"
	"github.com/fauzancodes/sales-demo-api/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func CreateProductCategory(c echo.Context) error {
	userID := c.Get("currentUser").(jwt.MapClaims)["id"].(string)
	log.Printf("Current user ID: %v", userID)

	var request dto.ProductCategoryRequest
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

	result, err := service.CreateProductCategory(userID, request)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			dto.Response{
				Status:  500,
				Message: "Failed to create",
				Error:   err.Error(),
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		dto.Response{
			Status:  200,
			Message: "Success to create",
			Data:    result,
		},
	)
}

func GetProductCategories(c echo.Context) error {
	userID := c.Get("currentUser").(jwt.MapClaims)["id"].(string)
	log.Printf("Current user ID: %v", userID)

	name := c.QueryParam("name")

	param := utils.PopulatePaging(c, "status")
	data, _, err := service.GetProductCategories(name, userID, param)
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

func GetProductCategoryByID(c echo.Context) error {
	id := c.Param("id")

	data, err := service.GetProductCategoryByID(id)
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

	return c.JSON(
		http.StatusOK,
		dto.Response{
			Status:  200,
			Message: "Success to get data",
			Data:    data,
		},
	)
}

func UpdateProductCategory(c echo.Context) error {
	id := c.Param("id")

	var request dto.ProductCategoryRequest
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

	data, err := service.UpdateProductCategory(id, request)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			dto.Response{
				Status:  500,
				Message: "Failed to update data",
				Error:   err.Error(),
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		dto.Response{
			Status:  200,
			Message: "Success to update data",
			Data:    data,
		},
	)
}

func DeleteProductCategory(c echo.Context) error {
	id := c.Param("id")

	err := service.DeleteProductCategory(id)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			dto.Response{
				Status:  500,
				Message: "Failed to delete data",
				Error:   err.Error(),
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		dto.Response{
			Status:  200,
			Message: "Success to delete data",
		},
	)
}
