package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/fauzancodes/sales-demo-api/app/dto"
	"github.com/fauzancodes/sales-demo-api/app/pkg/bcrypt"
	webToken "github.com/fauzancodes/sales-demo-api/app/pkg/jwt"
	"github.com/fauzancodes/sales-demo-api/app/pkg/utils"
	"github.com/fauzancodes/sales-demo-api/app/service"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func Register(c echo.Context) error {
	var request dto.AuthRequest
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

	param := utils.PopulatePaging(c, "")
	_, check, _ := service.GetUsers("", "", request.Email, param)
	if len(check) > 0 {
		return c.JSON(
			http.StatusBadRequest,
			dto.Response{
				Status:  400,
				Message: "Email has been registered",
				Error:   "",
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
			dto.Response{
				Status:  500,
				Message: "Failed to register",
				Error:   err.Error(),
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		dto.Response{
			Status:  200,
			Message: "Success to register",
			Data:    result,
		},
	)
}

func Login(c echo.Context) error {
	var request dto.AuthRequest
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

	param := utils.PopulatePaging(c, "")
	_, user, err := service.GetUsers("", "", request.Email, param)
	if len(user) == 0 {
		return c.JSON(
			http.StatusNotFound,
			dto.Response{
				Status:  404,
				Message: "Email not found",
				Error:   err.Error(),
			},
		)
	}

	err = bcrypt.VerifyPassword(request.Password, user[0].Password)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			dto.Response{
				Status:  400,
				Message: "Failed to verify password",
				Error:   err.Error(),
			},
		)
	}

	claims := jwt.MapClaims{}
	claims["id"] = user[0].ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	token, err := webToken.GenerateToken(&claims)
	if err != nil {
		return c.JSON(
			http.StatusUnauthorized,
			dto.Response{
				Status:  401,
				Message: "Failed to generate jwt token",
				Error:   err.Error(),
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		dto.Response{
			Status:  200,
			Message: "Success to login",
			Data:    token,
		},
	)
}

func GetCurrentUser(c echo.Context) error {
	userID := c.Get("currentUser").(jwt.MapClaims)["id"].(string)
	log.Printf("Current user ID: %v", userID)

	data, err := service.GetUserByID(userID)
	if err != nil {
		return c.JSON(
			http.StatusNotFound,
			dto.Response{
				Status:  404,
				Message: "Data not found",
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

func UpdateProfile(c echo.Context) error {
	var request dto.UserRequest
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

	userID := c.Get("currentUser").(jwt.MapClaims)["id"].(string)
	log.Printf("Current user ID: %v", userID)

	data, err := service.UpdateUser(userID, request)
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

func RemoveAccount(c echo.Context) error {
	userID := c.Get("currentUser").(jwt.MapClaims)["id"].(string)
	log.Printf("Current user ID: %v", userID)

	err := service.DeleteUser(userID)
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
