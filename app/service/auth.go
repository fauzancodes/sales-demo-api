package service

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/fauzancodes/sales-demo-api/app/dto"
	"github.com/fauzancodes/sales-demo-api/app/models"
	"github.com/fauzancodes/sales-demo-api/app/pkg/bcrypt"
	webToken "github.com/fauzancodes/sales-demo-api/app/pkg/jwt"
	"github.com/fauzancodes/sales-demo-api/app/pkg/smtp"
	"github.com/fauzancodes/sales-demo-api/app/repository"
	"github.com/golang-jwt/jwt/v5"
)

func SendEmailVerification(user models.SDAUser, successUrl, failedUrl, appUrl string) {
	claims := jwt.MapClaims{}
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	claims["successUrl"] = successUrl
	claims["failedUrl"] = failedUrl
	token, err := webToken.GenerateToken(&claims)
	if err != nil {
		log.Println(err.Error())
		return
	}

	verificationUrl := fmt.Sprintf("%v/auth/email-verification/%v", appUrl, token)

	name := fmt.Sprintf("%v %v", user.FirstName, user.LastName)
	if strings.ReplaceAll(name, " ", "") == "" {
		name = user.Email
	}

	fill := dto.EmailVerfication{
		Name:            name,
		AppUrl:          appUrl,
		VerificationUrl: verificationUrl,
	}

	smtp.SendEmail("email-verification", "", user.Email, "Email Verification", "", fill)
}

func VerifyUser(token string) (user models.SDAUser, successUrl, failedUrl string, err error) {
	if token == "" {
		err = errors.New("no jwt token provided")
		return
	}

	claims, err := webToken.DecodeToken(token)
	if err != nil {
		log.Println(err.Error())
		return
	}

	successUrl = claims["successUrl"].(string)
	failedUrl = claims["failedUrl"].(string)

	userID := claims["id"].(string)
	user, _, err = GetUserByID(userID, []string{})
	if err != nil {
		log.Println(err.Error())
		return
	}

	user.IsVerified = true
	user, err = repository.UpdateUser(user)
	if err != nil {
		log.Println("Failed to update user: " + err.Error())
		return
	}

	return
}

func SendResetPasswordRequest(user models.SDAUser, redirectUrl, appUrl string) {
	claims := jwt.MapClaims{}
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	claims["redirectUrl"] = redirectUrl
	token, err := webToken.GenerateToken(&claims)
	if err != nil {
		log.Println(err.Error())
		return
	}

	var resetPasswordUrl string
	if redirectUrl != "" {
		resetPasswordUrl = fmt.Sprintf("%v/%v", redirectUrl, token)
	} else {
		resetPasswordUrl = fmt.Sprintf("%v/auth/reset-password/instruction/%v", appUrl, token)
	}

	name := fmt.Sprintf("%v %v", user.FirstName, user.LastName)
	if strings.ReplaceAll(name, " ", "") == "" {
		name = user.Email
	}

	fill := dto.ResetPassword{
		Name:             name,
		AppUrl:           appUrl,
		ResetPasswordUrl: resetPasswordUrl,
	}

	smtp.SendEmail("reset-password", "", user.Email, "Reset Your Password", "", fill)
}

func ResetPassword(request dto.ResetPasswordRequest) (user models.SDAUser, statusCode int, err error) {
	claims, err := webToken.DecodeToken(request.Token)
	if err != nil {
		statusCode = http.StatusInternalServerError
		return
	}

	userID := claims["id"].(string)
	user, statusCode, err = GetUserByID(userID, []string{})
	if err != nil {
		return
	}

	user.Password = bcrypt.HashPassword(request.NewPassword)
	user, err = repository.UpdateUser(user)
	if err != nil {
		log.Println("Failed to update user: ", err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	return
}
