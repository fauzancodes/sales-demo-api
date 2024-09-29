package routes

import (
	"github.com/fauzancodes/sales-demo-api/app/controllers"
	"github.com/fauzancodes/sales-demo-api/app/middlewares"
	"github.com/labstack/echo/v4"
)

func CustomerRoute(app *echo.Echo) {
	customer := app.Group("/customer", middlewares.Auth, middlewares.StripHTMLMiddleware)
	{
		customer.POST("", controllers.CreateCustomer)
		customer.GET("", controllers.GetCustomers)
		customer.GET("/:id", controllers.GetCustomerByID)
		customer.PATCH("/:id", controllers.UpdateCustomer)
		customer.DELETE("/:id", controllers.DeleteCustomer)
	}
}
