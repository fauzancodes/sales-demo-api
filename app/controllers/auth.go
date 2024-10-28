package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/fauzancodes/sales-demo-api/app/dto"
	"github.com/fauzancodes/sales-demo-api/app/pkg/bcrypt"
	webToken "github.com/fauzancodes/sales-demo-api/app/pkg/jwt"
	"github.com/fauzancodes/sales-demo-api/app/pkg/utils"
	"github.com/fauzancodes/sales-demo-api/app/repository"
	"github.com/fauzancodes/sales-demo-api/app/service"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func Register(c echo.Context) error {
	var request dto.RegisterRequest
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

	param := utils.PopulatePaging(c, "")
	_, check, _, _ := service.GetUsers("", "", request.Email, param, []string{})
	if len(check) > 0 {
		return c.JSON(
			http.StatusBadRequest,
			dto.Response{
				Status:  http.StatusBadRequest,
				Message: "Email has been registered",
				Error:   "",
			},
		)
	}

	result, statusCode, err := service.CreateUser(dto.UserRequest{
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		return c.JSON(
			statusCode,
			dto.Response{
				Status:  statusCode,
				Message: "Failed to register",
				Error:   err.Error(),
			},
		)
	}

	go service.SendEmailVerification(result, request.SuccessVerificationUrl, request.FailedVerificationUrl, utils.GetBaseUrl(c))

	return c.JSON(
		statusCode,
		dto.Response{
			Status:  statusCode,
			Message: "Success to register",
			Data:    result,
		},
	)
}

func Login(c echo.Context) error {
	var request dto.LoginRequest
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

	param := utils.PopulatePaging(c, "")
	_, user, statusCode, _ := service.GetUsers("", "", request.Email, param, []string{})
	if len(user) == 0 {
		return c.JSON(
			statusCode,
			dto.Response{
				Status:  statusCode,
				Message: "Email not found",
			},
		)
	}

	err := bcrypt.VerifyPassword(request.Password, user[0].Password)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			dto.Response{
				Status:  http.StatusBadRequest,
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
		statusCode,
		dto.Response{
			Status:  statusCode,
			Message: "Success to login",
			Data:    token,
		},
	)
}

func GetCurrentUser(c echo.Context) error {
	userID := c.Get("currentUser").(jwt.MapClaims)["id"].(string)
	log.Printf("Current user ID: %v", userID)

	data, statusCode, err := service.GetUserByID(userID, []string{
		"Products",
		"ProductCategories",
		"ProductStocks",
		"Customers",
		"Sale",
		"SaleDetails",
		"SaleDetails.Product",
		"Payment",
		"Payment.PaymentMethod",
	})
	if err != nil {
		return c.JSON(
			statusCode,
			dto.Response{
				Status:  statusCode,
				Message: "Data not found",
				Error:   err.Error(),
			},
		)
	}

	return c.JSON(
		statusCode,
		dto.Response{
			Status:  statusCode,
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
				Status:  http.StatusUnprocessableEntity,
				Message: "Invalid request body",
				Error:   err.Error(),
			},
		)
	}

	userID := c.Get("currentUser").(jwt.MapClaims)["id"].(string)
	log.Printf("Current user ID: %v", userID)

	data, statusCode, err := service.UpdateUser(userID, request)
	if err != nil {
		return c.JSON(
			statusCode,
			dto.Response{
				Status:  statusCode,
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

	statusCode, err := service.DeleteUser(userID)
	if err != nil {
		return c.JSON(
			statusCode,
			dto.Response{
				Status:  statusCode,
				Message: "Failed to delete data",
				Error:   err.Error(),
			},
		)
	}

	return c.JSON(
		statusCode,
		dto.Response{
			Status:  statusCode,
			Message: "Success to delete data",
		},
	)
}

func VerifyUser(c echo.Context) error {
	token := c.Param("token")

	data, successUrl, failedUrl, err := service.VerifyUser(token)
	if err != nil {
		if failedUrl != "" {
			return c.Redirect(http.StatusTemporaryRedirect, failedUrl)
		}
		return c.JSON(
			http.StatusNotFound,
			dto.Response{
				Status:  500,
				Message: "Failed to verify user",
				Error:   err.Error(),
			},
		)
	}

	if successUrl != "" {
		return c.Redirect(http.StatusTemporaryRedirect, successUrl)
	}
	return c.JSON(
		http.StatusOK,
		dto.Response{
			Status:  200,
			Message: "Success to verify user",
			Data:    data,
		},
	)
}

func ResendEmailVerification(c echo.Context) error {
	var request dto.ResendEmailVerification
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

	user, _, _, err := repository.GetUsers(dto.FindParameter{
		Filter: "deleted_at IS NULL AND email = '" + request.Email + "'",
	}, []string{})
	if err != nil || len(user) == 0 {
		return c.JSON(
			http.StatusNotFound,
			dto.Response{
				Status:  404,
				Message: "Failed to get user",
				Error:   err.Error(),
			},
		)
	}

	if user[0].IsVerified {
		return c.JSON(
			http.StatusBadRequest,
			dto.Response{
				Status:  http.StatusBadRequest,
				Message: "User has been verified",
			},
		)
	}

	go service.SendEmailVerification(user[0], request.SuccessVerificationUrl, request.FailedVerificationUrl, utils.GetBaseUrl(c))

	return c.JSON(
		http.StatusOK,
		dto.Response{
			Status:  200,
			Message: "Success to send email verification",
		},
	)
}

func SendForgotPasswordRequest(c echo.Context) error {
	var request dto.SendForgotPasswordRequest
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

	user, _, _, err := repository.GetUsers(dto.FindParameter{
		Filter: "deleted_at IS NULL AND email = '" + request.Email + "'",
	}, []string{})
	if err != nil || len(user) == 0 {
		return c.JSON(
			http.StatusNotFound,
			dto.Response{
				Status:  404,
				Message: "Failed to get user",
				Error:   err.Error(),
			},
		)
	}

	go service.SendResetPasswordRequest(user[0], request.RedirectUrl, utils.GetBaseUrl(c))

	return c.JSON(
		http.StatusOK,
		dto.Response{
			Status:  200,
			Message: "Success to send reset password request",
		},
	)
}

func SendResetPasswordRequestInstruction(c echo.Context) error {
	token := c.Param("token")

	return c.JSON(
		http.StatusOK,
		dto.Response{
			Status:  200,
			Message: "You should include a redirect_url field, so that the request will be forwarded to your url, then in that url create a page for the user to fill in their new password, then send the password from the user along with the token in the url param to the POST /auth/reset-password endpoint",
			Data:    token,
		},
	)
}

func ResetPassword(c echo.Context) error {
	var request dto.ResetPasswordRequest
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

	data, statusCode, err := service.ResetPassword(request)
	if err != nil {
		return c.JSON(
			statusCode,
			dto.Response{
				Status:  statusCode,
				Message: "Failed to reset password",
				Error:   err.Error(),
			},
		)
	}

	return c.JSON(
		statusCode,
		dto.Response{
			Status:  statusCode,
			Message: "Success to reset password",
			Data:    data,
		},
	)
}
