package main

import (
	"log"

	"github.com/fauzancodes/sales-demo-api/app/config"
	"github.com/fauzancodes/sales-demo-api/app/middlewares"
	"github.com/fauzancodes/sales-demo-api/app/routes"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
)

func main() {
	app := Init()

	port := config.LoadConfig().IndexPort

	log.Printf("Server: " + config.LoadConfig().BaseUrl + ":" + port)
	app.Logger.Fatal(app.Start(":" + port))
}

func Init() *echo.Echo {
	app := echo.New()

	app.Use(middlewares.Cors())
	app.Use(middlewares.Gzip())
	app.Use(middlewares.Logger())
	app.Use(middlewares.Secure())
	app.Use(middlewares.Recover())

	config.Database()

	routes.IndexRoute(app)
	routes.AuthRoute(app)
	routes.ProductRoute(app)
	routes.CustomerRoute(app)
	routes.SaleRoute(app)
	routes.PaymentGatewayRoute(app)

	return app
}
