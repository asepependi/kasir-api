package routes

import (
	"database/sql"
	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"

	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine, db *sql.DB) {
	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	category := handlers.NewCategoryHandler(categoryService)

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
	}
}