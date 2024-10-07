package service

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fauzancodes/sales-demo-api/app/config"
	"github.com/fauzancodes/sales-demo-api/app/dto"
	"github.com/fauzancodes/sales-demo-api/app/models"
	webToken "github.com/fauzancodes/sales-demo-api/app/pkg/jwt"
	"github.com/fauzancodes/sales-demo-api/app/pkg/utils"
	"github.com/fauzancodes/sales-demo-api/app/repository"
	"github.com/golang-jwt/jwt/v5"
)

func SendEmailVerification(user models.SDAUser, successUrl, failedUrl string) {
	var appUrl string
	if config.LoadConfig().BaseUrl == "http://localhost" {
		appUrl = fmt.Sprintf("%v:%v", config.LoadConfig().BaseUrl, config.LoadConfig().IndexPort)
	} else {
		appUrl = config.LoadConfig().BaseUrl
	}

	claims := jwt.MapClaims{}
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	claims["successUrl"] = successUrl
	claims["failedUrl"] = failedUrl
	token, err := webToken.GenerateToken(&claims)
	if err != nil {
		log.Println("Failed to generate jwt token:", err.Error())
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

	utils.SendEmail("email-verification", user.Email, "Email Verification", "", fill)
}

func VerifyUser(token string) (user models.SDAUser, successUrl, failedUrl string, err error) {
	if token == "" {
		err = errors.New("no jwt token provided")
		return
	}

	claims, err := webToken.DecodeToken(token)
	if err != nil {
		return
	}

	successUrl = claims["successUrl"].(string)
	failedUrl = claims["failedUrl"].(string)

	userID := claims["id"].(string)
	user, err = GetUserByID(userID, []string{})
	if err != nil {
		return
	}

	user.IsVerified = true
	user, err = repository.UpdateUser(user)
	if err != nil {
		return
	}

	return
}
