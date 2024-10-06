package routes

import (
	"github.com/fauzancodes/sales-demo-api/app/controllers"
	"github.com/fauzancodes/sales-demo-api/app/middlewares"
	"github.com/labstack/echo/v4"
)

func SaleRoute(app *echo.Echo) {
	sale := app.Group("/sale", middlewares.Auth, middlewares.StripHTMLMiddleware)
	{
		sale.POST("", controllers.CreateSale)
		sale.GET("", controllers.GetSales)
		sale.GET("/:id", controllers.GetSaleByID)
		sale.PATCH("/:id", controllers.UpdateSale)
		sale.DELETE("/:id", controllers.DeleteSale)
	}
}
