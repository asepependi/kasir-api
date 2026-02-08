package handlers

import (
	"database/sql"
	"encoding/json"
	"kasir-api/models"
	"kasir-api/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service *services.ProductService
}

func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) GetAll(c *gin.Context) {
	searchQuery := c.Query("name")
	products, err := h.service.GetAll(searchQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"error": "Internal server error",
		})
		return
	}

	c.JSON(200, products)
}

func (h *ProductHandler) Create(c *gin.Context) {
	var newProduct models.Product
	err := json.NewDecoder(c.Request.Body).Decode(&newProduct)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	if newProduct.CategoryID == 0 || newProduct.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Kategori produk wajib diisi",
		})
		return
	}

	if newProduct.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Nama produk wajib diisi",
		})
		return
	}

	if newProduct.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Harga produk wajib diisi dan harus lebih dari 0",
		})
		return
	}

	if newProduct.Stock <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Stok produk wajib diisi dan harus lebih dari 0",
		})
		return
	}

	newData, err := h.service.Create(&newProduct)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H {
		"data": gin.H{
			"id":          newData.ID,
			"category_id": newData.CategoryID,
			"name":        newData.Name,
			"price":       newData.Price,
			"stock":       newData.Stock,
			"created_at":  newData.CreatedAt,
		},
		"message": "Berhasil disimpan",
	})
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product ID",
		})
		return
	}

	product, err := h.service.GetByID(idInt)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Product not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H {
			"error": "Internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":          product.ID,
		"category_id": product.CategoryID,
		"name":        product.Name,
		"price":       product.Price,
		"stock":       product.Stock,
		"created_at":  product.CreatedAt,
	})
}

func (h *ProductHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var updateProduct models.Product
	err := json.NewDecoder(c.Request.Body).Decode(&updateProduct)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product ID",
		})
		return
	}

	if updateProduct.CategoryID == 0 || updateProduct.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Kategori produk wajib diisi",
		})
		return
	}

	if updateProduct.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Nama produk wajib diisi",
		})
		return
	}

	if updateProduct.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Harga produk wajib diisi dan harus lebih dari 0",
		})
		return
	}

	if updateProduct.Stock <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Stok produk wajib diisi dan harus lebih dari 0",
		})
		return
	}

	updated, err := h.service.Update(idInt, &updateProduct)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":          updated.ID,
			"category_id": updated.CategoryID,
			"name":        updated.Name,
			"price":       updated.Price,
			"stock":       updated.Stock,
			"created_at":  updated.CreatedAt,
		},
		"message": "Berhasil diupdate",
	})
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product ID",
		})
		return
	}

	if err := h.service.Delete(idInt); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Product not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil dihapus",
	})
}
