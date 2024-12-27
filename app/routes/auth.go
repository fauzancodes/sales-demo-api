package routes

import (
	"github.com/fauzancodes/sales-demo-api/app/controllers"
	"github.com/fauzancodes/sales-demo-api/app/middlewares"
	"github.com/labstack/echo/v4"
)

func AuthRoute(app *echo.Echo) {
	auth := app.Group("/auth")
	{
		auth.POST("/register", controllers.Register, middlewares.CheckAPIKey)
		auth.POST("/login", controllers.Login, middlewares.CheckAPIKey)
		auth.GET("/user", controllers.GetCurrentUser, middlewares.CheckAPIKey, middlewares.Auth)
		auth.PATCH("/update-profile", controllers.UpdateProfile, middlewares.CheckAPIKey, middlewares.Auth)
		auth.DELETE("/remove-account", controllers.RemoveAccount, middlewares.CheckAPIKey, middlewares.Auth)

		emailVerfication := auth.Group("/email-verification")
		{
			emailVerfication.GET("/:token", controllers.VerifyUser)
			emailVerfication.POST("/resend", controllers.ResendEmailVerification, middlewares.CheckAPIKey)
		}

		resetPassword := auth.Group("/reset-password")
		{
			resetPassword.POST("/send", controllers.SendForgotPasswordRequest, middlewares.CheckAPIKey)
			resetPassword.GET("/instruction/:token", controllers.SendResetPasswordRequestInstruction)
			resetPassword.POST("", controllers.ResetPassword, middlewares.CheckAPIKey)
		}
	}
}
