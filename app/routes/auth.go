package routes

import (
	"github.com/fauzancodes/sales-demo-api/app/controllers"
	"github.com/fauzancodes/sales-demo-api/app/middlewares"
	"github.com/labstack/echo/v4"
)

func AuthRoute(app *echo.Echo) {
	auth := app.Group("/auth", middlewares.StripHTMLMiddleware)
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
		auth.GET("/user", controllers.GetCurrentUser, middlewares.Auth)
		auth.PATCH("/update-profile", controllers.UpdateProfile, middlewares.Auth)
		auth.DELETE("/remove-account", controllers.RemoveAccount, middlewares.Auth)
	}
}
