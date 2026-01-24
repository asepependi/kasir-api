package main

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Categories struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
}

var category = []Categories{
	{ID: 1, Name: "Makanan", Description: "Kategori Makanan" },
	{ID: 2, Name: "Minuman", Description: "Kategori Minuman" },
	{ID: 3, Name: "Dessert", Description: "Kategori Dessert" },
	{ID: 4, Name: "Snack", Description: "Kategori Snack" },
	{ID: 5, Name: "Lain-lain", Description: "Kategori Lain-lain" },
}

func main() {
	router := gin.Default()

	if err := router.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		log.Fatal("Failed to set trusted proxies:", err)
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"message": "Kasir API is running",
		})
	})

	// List Category
	router.GET("/api/category", func(c *gin.Context) {
		c.JSON(200, category)
	})

	// Create Category
	router.POST("/api/category", func(c *gin.Context) {
		var newCategory Categories
		err := json.NewDecoder(c.Request.Body).Decode(&newCategory)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "Invalid request",
			})
			return
		}

		newCategory.ID = len(category) + 1
		category = append(category, newCategory)

		c.JSON(201, newCategory)
	})

	// Get Category
	router.GET("/api/:id/category", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "Invalid category ID",
			})
			return
		}

		for _, p := range category {
			if p.ID == id {
				c.JSON(200, p)
				return
			}
		}

		c.JSON(404, gin.H {
			"error": "Kategori belum tersedia",
		})
	})

	// Update Category
	router.PUT("/api/:id/category", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(400, gin.H {
				"error": "Invalid category ID",
			})
			return
		}

		var updateCategory Categories
		err = json.NewDecoder(c.Request.Body).Decode(&updateCategory)
		if err != nil {
			c.JSON(400, gin.H {
				"error": "Invalid request",
			})
			return
		}

		for i := range category {
			if category[i].ID == id {
				updateCategory.ID = id
				category[i] = updateCategory
				c.JSON(200, updateCategory)
				return
			}
		}

		c.JSON(404, gin.H {
			"error": "Category belum tersedia",
		})
	})

	// Delete Category
	router.DELETE("/api/:id/category", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(400, gin.H {
				"error": "Invalid category ID",
			})
			return
		}

		for i, p := range category {
			if p.ID == id {
				category = append(category[:i], category[i+1:]...)
				c.JSON(200, gin.H {
					"data": category,
					"message": "Sukses delete",
				})

				return
			}
		}

		c.JSON(404, gin.H {
			"error": "Category belum tersedia",
		})
	})

	const port = "8090"

	log.Printf("Kasir API server is running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
