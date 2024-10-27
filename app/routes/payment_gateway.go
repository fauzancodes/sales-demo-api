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
			midtrans.POST("/charge/core", controllers.MidtransChargeCore, middlewares.Auth)
			midtrans.POST("/charge/snap", controllers.MidtransChargeSnap, middlewares.Auth)
			midtrans.POST("/notification", controllers.MidtransNotification)
		}

		ipaymu := paymentGateway.Group("/ipaymu")
		{
			ipaymu.GET("/available-payment-method", controllers.GetIPaymuPaymentMethods, middlewares.Auth)
			ipaymu.POST("/charge", controllers.IPaymuCharge, middlewares.Auth)
			ipaymu.POST("/notification", controllers.IPaymuNotification)
		}

		xendit := paymentGateway.Group("/xendit")
		{
			xendit.GET("/available-payment-method", controllers.GetXenditPaymentMethods, middlewares.Auth)
			xendit.POST("/charge/payment", controllers.XenditChargePayment, middlewares.Auth)
			xendit.POST("/charge/invoice", controllers.XenditChargeInvoice, middlewares.Auth)
			xendit.POST("/notification", controllers.XenditNotification)
		}
	}
}
