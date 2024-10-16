package routes

import (
	"github.com/fauzancodes/sales-demo-api/app/controllers"
	"github.com/fauzancodes/sales-demo-api/app/middlewares"
	"github.com/labstack/echo/v4"
)

func PaymentGatewayRoute(app *echo.Echo) {
	paymentGateway := app.Group("/payment-gateway", middlewares.Auth, middlewares.StripHTMLMiddleware)
	{
		midtrans := paymentGateway.Group("/midtrans")
		{
			midtrans.GET("/available-payment-method", controllers.GetMidtransPaymentMethods)
			midtrans.POST("/charge", controllers.MidtransCharge)
			midtrans.POST("/notification", controllers.MidtransNotification)
		}
	}
}
