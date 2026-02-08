package routes

import (
	"database/sql"
	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"

	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine, db *sql.DB) {
	// Category
	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	category := handlers.NewCategoryHandler(categoryService)
	// Products
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	product := handlers.NewProductHandler(productService)
	// Transaction
	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)
	transaction := handlers.NewTransactionHandler(transactionService)

	r.GET("/", func(c *gin.Context){
		c.JSON(200, gin.H{
			"status":  "OK",
			"message": "Kasir API is running",
		})
	})

	api := r.Group("/api/")
	{
		categoryGroup := api.Group("/category")
		categoryGroup.GET("/", category.GetAll)
		categoryGroup.POST("/", category.Create)
		categoryGroup.GET("/:id", category.GetByID)
		categoryGroup.PUT("/:id", category.Update)
		categoryGroup.DELETE("/:id", category.Delete)

		productGroup := api.Group("/product")
		productGroup.GET("/", product.GetAll)
		productGroup.POST("/", product.Create)
		productGroup.GET("/:id", product.GetByID)
		productGroup.PUT("/:id", product.Update)
		productGroup.DELETE("/:id", product.Delete)

		api.POST("checkout", transaction.Checkout)
		api.GET("/report/hari-ini", transaction.GetReport)
		api.GET("/report", transaction.GetReportByDateRange)
	}
}
