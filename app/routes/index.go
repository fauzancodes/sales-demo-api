package routes

import (
	"github.com/fauzancodes/sales-demo-api/app/controllers"
	"github.com/fauzancodes/sales-demo-api/app/middlewares"
	"github.com/labstack/echo/v4"
)

func Route(app *echo.Echo) {
	app.Static("/assets", "assets")
	app.Static("/docs", "docs")

	app.GET("/", controllers.Index, middlewares.StripHTMLMiddleware)

	AuthRoute(app)
	CustomerRoute(app)
	ProductRoute(app)
	SaleRoute(app)
	PaymentGatewayRoute(app)
}
