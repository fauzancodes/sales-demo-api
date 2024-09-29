package routes

import (
	"github.com/fauzancodes/sales-demo-api/app/controllers"
	"github.com/fauzancodes/sales-demo-api/app/middlewares"
	"github.com/labstack/echo/v4"
)

func RouteInit(app *echo.Echo) {
	app.Static("/assets", "assets")
	app.Static("/docs", "docs")

	app.GET("/", controllers.Index)

	api := app.Group("/v1", middlewares.StripHTMLMiddleware)

	auth := api.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
		auth.GET("/user", controllers.GetCurrentUser, middlewares.Auth)
		auth.PATCH("/update-profile", controllers.UpdateProfile, middlewares.Auth)
		auth.DELETE("/remove-account", controllers.RemoveAccount, middlewares.Auth)
	}

	product := api.Group("/product", middlewares.Auth)
	{
		category := product.Group("/category")
		{
			category.POST("", controllers.CreateProductCategory)
			category.GET("", controllers.GetProductCategories)
			category.GET("/:id", controllers.GetProductCategoryByID)
			category.PATCH("/:id", controllers.UpdateProductCategory)
			category.DELETE("/:id", controllers.DeleteProductCategory)
		}

		product.POST("/upload-image", controllers.UploadFile)
		product.POST("", controllers.CreateProduct)
		product.GET("", controllers.GetProducts)
		product.GET("/:id", controllers.GetProductByID)
		product.PATCH("/:id", controllers.UpdateProduct)
		product.DELETE("/:id", controllers.DeleteProduct)

		stock := product.Group("/stock")
		{
			stock.POST("", controllers.CreateProductStock)
			stock.GET("", controllers.GetProductStocks)
		}
	}

	customer := api.Group("/customer", middlewares.Auth)
	{
		customer.POST("", controllers.CreateCustomer)
		customer.GET("", controllers.GetCustomers)
		customer.GET("/:id", controllers.GetCustomerByID)
		customer.PATCH("/:id", controllers.UpdateCustomer)
		customer.DELETE("/:id", controllers.DeleteCustomer)
	}
}
