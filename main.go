// package main
package handler

import (
	"log"
	"net/http"

	"github.com/fauzancodes/sales-demo-api/app/config"
	"github.com/fauzancodes/sales-demo-api/app/routes"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
)

func Main(w http.ResponseWriter, r *http.Request) {
	app := Start()

	app.ServeHTTP(w, r)
}

func main() {
	app := Start()

	port := config.LoadConfig().IndexPort

	log.Printf("Server: " + config.LoadConfig().BaseUrl + ":" + port)
	app.Logger.Fatal(app.Start(":" + port))
}

func Start() *echo.Echo {
	app := echo.New()

	config.Database()

	routes.IndexRoute(app)
	routes.AuthRoute(app)
	routes.ProductRoute(app)
	routes.CustomerRoute(app)
	routes.SaleRoute(app)
	routes.PaymentGatewayRoute(app)

	return app
}
