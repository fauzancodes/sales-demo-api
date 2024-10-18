package routes

import (
	"github.com/fauzancodes/sales-demo-api/app/controllers"
	"github.com/fauzancodes/sales-demo-api/app/middlewares"
	"github.com/labstack/echo/v4"
)

func PaymentGatewayRoute(app *echo.Echo) {
	paymentGateway := app.Group("/payment-gateway", middlewares.StripHTMLMiddleware)
	{
		midtrans := paymentGateway.Group("/midtrans")
		{
			midtrans.GET("/available-payment-method", controllers.GetMidtransPaymentMethods, middlewares.Auth)
			midtrans.POST("/charge", controllers.MidtransCharge, middlewares.Auth)
			midtrans.POST("/notification", controllers.MidtransNotification)
		}

		ipaymu := paymentGateway.Group("/ipaymu")
		{
			ipaymu.GET("/available-payment-method", controllers.GetIPaymuPaymentMethods, middlewares.Auth)
		}
	}
}
