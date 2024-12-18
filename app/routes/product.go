package routes

import (
	"github.com/fauzancodes/sales-demo-api/app/controllers"
	"github.com/fauzancodes/sales-demo-api/app/middlewares"
	"github.com/labstack/echo/v4"
)

func ProductRoute(app *echo.Echo) {
	product := app.Group("/product", middlewares.CheckAPIKey, middlewares.Auth)
	{
		category := product.Group("/category")
		{
			category.POST("", controllers.CreateProductCategory)
			category.GET("", controllers.GetProductCategories)
			category.GET("/:id", controllers.GetProductCategoryByID)
			category.PATCH("/:id", controllers.UpdateProductCategory)
			category.DELETE("/:id", controllers.DeleteProductCategory)

			importData := category.Group("/import")
			{
				importData.GET("/template", controllers.GetProductCategoryImportTemplate)
				importData.POST("", controllers.ImportProductCategory)
			}

			category.GET("/export", controllers.ExportProductCategory)
		}

		product.POST("/upload-image", controllers.UploadProductPicture)
		product.POST("", controllers.CreateProduct)
		product.GET("", controllers.GetProducts)
		product.GET("/:id", controllers.GetProductByID)
		product.PATCH("/:id", controllers.UpdateProduct)
		product.DELETE("/:id", controllers.DeleteProduct)

		importData := product.Group("/import")
		{
			importData.GET("/template", controllers.GetProductImportTemplate)
			importData.POST("", controllers.ImportProduct)
		}

		stock := product.Group("/stock")
		{
			stock.POST("", controllers.CreateProductStock)
			stock.GET("", controllers.GetProductStocks)
		}

		product.GET("/export", controllers.ExportProduct)
	}
}
