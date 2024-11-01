package routes

import (
	"github.com/fauzancodes/sales-demo-api/app/controllers"
	"github.com/fauzancodes/sales-demo-api/app/middlewares"
	"github.com/labstack/echo/v4"
)

func AuthRoute(app *echo.Echo) {
	auth := app.Group("/auth", middlewares.CheckAPIKey)
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
		auth.GET("/user", controllers.GetCurrentUser, middlewares.Auth)
		auth.PATCH("/update-profile", controllers.UpdateProfile, middlewares.Auth)
		auth.DELETE("/remove-account", controllers.RemoveAccount, middlewares.Auth)

		emailVerfication := auth.Group("/email-verification")
		{
			emailVerfication.GET("/:token", controllers.VerifyUser)
			emailVerfication.POST("/resend", controllers.ResendEmailVerification)
		}

		resetPassword := auth.Group("/reset-password")
		{
			resetPassword.POST("/send", controllers.SendForgotPasswordRequest)
			resetPassword.GET("/instruction/:token", controllers.SendResetPasswordRequestInstruction)
			resetPassword.POST("", controllers.ResetPassword)
		}
	}
}
