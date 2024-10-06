package main

import (
	"log"

	"github.com/fauzancodes/sales-demo-api/app/config"
	"github.com/fauzancodes/sales-demo-api/app/routes"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
)

func main() {
	app := echo.New()

	config.Database()

	routes.IndexRoute(app)
	routes.AuthRoute(app)
	routes.ProductRoute(app)
	routes.CustomerRoute(app)
	routes.SaleRoute(app)

	port := config.LoadConfig().IndexPort

	log.Printf("Server: " + config.LoadConfig().BaseUrl + ":" + port)
	app.Logger.Fatal(app.Start(":" + port))
}
