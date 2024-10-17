package service

import (
	"errors"
	"fmt"
	"log"
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

	htmlString := `
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<style>
				body {
					font-family: Arial, sans-serif;
					background-color: #f4f4f4;
					margin: 0;
					padding: 0;
				}
				.container {
					width: 100%;
					max-width: 600px;
					margin: 50px auto;
					background-color: #ffffff;
					border-radius: 8px;
					box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
					overflow: hidden;
				}
				.header {
					background-color: #006BFF;
					color: #ffffff;
					padding: 20px;
					text-align: center;
				}
				.header h1 {
					margin: 0;
					font-size: 24px;
				}
				.content {
					padding: 30px;
					text-align: center;
				}
				.content h2 {
					color: #333333;
					font-size: 22px;
				}
				.content p {
					color: #666666;
					font-size: 16px;
				}
				.content a {
					display: inline-block;
					margin-top: 20px;
					padding: 12px 24px;
					font-size: 16px;
					color: #006BFF;
					text-decoration: none;
					border-radius: 4px;
				}
				.content a.btn {
					color: #f4f4f4;
					background-color: #006BFF;
					font-weight: bold;
				}
				.link-alternative {
					margin-top: 20px;
					font-size: 14px;
					color: #999999;
				}
				.link-alternative a {
					color: #006BFF;
					text-decoration: underline;
					word-break: break-all;
				}
				.footer {
					background-color: #f4f4f4;
					padding: 20px;
					text-align: center;
					font-size: 14px;
					color: #999999;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h1>Verify Your Email Address</h1>
				</div>
				<div class="content">
					<h2>Hello, ` + fill.Name + `!</h2>
					<p>Thank you for signing up. Please click the button below to verify your email address:</p>
					<a class="btn" href="` + fill.VerificationUrl + `" target="_blank">Verify Email</a>
					<p>If you did not create an account, please ignore this email.</p>
					<div class="link-alternative">
						<p>If the button above does not work, please click or copy and paste the following link into your browser:</p>
						<a href="` + fill.VerificationUrl + `" target="_blank">` + fill.VerificationUrl + `</a>
					</div>
				</div>
				<div class="footer">
					&copy; 2024 <a href="` + fill.AppUrl + `" target="_blank">Sale Demo API</a>. All rights reserved.
				</div>
			</div>
		</body>
		</html>
	`

	smtp.SendEmail("email-verification", htmlString, "", user.Email, "Email Verification", "", fill)
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

func SendResetPasswordRequest(user models.SDAUser, redirectUrl, appUrl string) {
	claims := jwt.MapClaims{}
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	claims["redirectUrl"] = redirectUrl
	token, err := webToken.GenerateToken(&claims)
	if err != nil {
		log.Println("Failed to generate jwt token:", err.Error())
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

	htmlString := `
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<style>
				body {
					font-family: Arial, sans-serif;
					background-color: #f4f4f4;
					margin: 0;
					padding: 0;
				}
				.container {
					width: 100%;
					max-width: 600px;
					margin: 50px auto;
					background-color: #ffffff;
					border-radius: 8px;
					box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
					overflow: hidden;
				}
				.header {
					background-color: #FF6B6B;
					color: #ffffff;
					padding: 20px;
					text-align: center;
				}
				.header h1 {
					margin: 0;
					font-size: 24px;
				}
				.content {
					padding: 30px;
					text-align: center;
				}
				.content h2 {
					color: #333333;
					font-size: 22px;
				}
				.content p {
					color: #666666;
					font-size: 16px;
				}
				.content a {
					display: inline-block;
					margin-top: 20px;
					padding: 12px 24px;
					font-size: 16px;
					color: #ffffff;
					text-decoration: none;
					border-radius: 4px;
				}
				.content a.btn {
					background-color: #FF6B6B;
					font-weight: bold;
				}
				.link-alternative {
					margin-top: 20px;
					font-size: 14px;
					color: #999999;
				}
				.link-alternative a {
					color: #FF6B6B;
					text-decoration: underline;
					word-break: break-all;
				}
				.footer {
					background-color: #f4f4f4;
					padding: 20px;
					text-align: center;
					font-size: 14px;
					color: #999999;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h1>Reset Your Password</h1>
				</div>
				<div class="content">
					<h2>Hello, ` + fill.Name + `!</h2>
					<p>We received a request to reset your password. Click the button below to reset it:</p>
					<a class="btn" href="` + fill.ResetPasswordUrl + `" target="_blank">Reset Password</a>
					<p>If you did not request a password reset, you can safely ignore this email.</p>
					<div class="link-alternative">
						<p>If the button above does not work, please click or copy and paste the following link into your browser:</p>
						<a href="` + fill.ResetPasswordUrl + `" target="_blank">` + fill.ResetPasswordUrl + `</a>
					</div>
				</div>
				<div class="footer">
					&copy; 2024 <a href="` + fill.AppUrl + `" target="_blank">Sale Demo API</a>. All rights reserved.
				</div>
			</div>
		</body>
		</html>
	`

	smtp.SendEmail("reset-password", htmlString, "", user.Email, "Reset Your Password", "", fill)
}

func ResetPassword(request dto.ResetPasswordRequest) (user models.SDAUser, err error) {
	claims, err := webToken.DecodeToken(request.Token)
	if err != nil {
		return
	}

	userID := claims["id"].(string)
	user, err = GetUserByID(userID, []string{})
	if err != nil {
		return
	}

	user.Password = bcrypt.HashPassword(request.NewPassword)
	user, err = repository.UpdateUser(user)
	if err != nil {
		return
	}

	return
}
