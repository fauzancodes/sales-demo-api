package routes

import (
	"github.com/fauzancodes/sales-demo-api/app/controllers"
	"github.com/fauzancodes/sales-demo-api/app/middlewares"
	"github.com/labstack/echo/v4"
)

func CustomerRoute(app *echo.Echo) {
	customer := app.Group("/customer", middlewares.CheckAPIKey, middlewares.Auth)
	{
		customer.POST("", controllers.CreateCustomer)
		customer.GET("", controllers.GetCustomers)
		customer.GET("/:id", controllers.GetCustomerByID)
		customer.PATCH("/:id", controllers.UpdateCustomer)
		customer.DELETE("/:id", controllers.DeleteCustomer)

		importData := customer.Group("/import")
		{
			importData.GET("/template", controllers.GetCustomerImportTemplate)
			importData.POST("", controllers.ImportCustomer)
		}

		customer.GET("/export", controllers.ExportCustomer)
	}
}
