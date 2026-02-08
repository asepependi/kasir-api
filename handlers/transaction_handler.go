package handlers

import (
	"encoding/json"
	"kasir-api/models"
	"kasir-api/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	service *services.TransactionService
}

func NewTransactionHandler(service *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

func (h *TransactionHandler) Checkout(c *gin.Context) {
	var req models.CheckoutRequest
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"error": "Invalid request body",
		})
		return
	}

	if len(req.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Items tidak boleh kosong",
		})
		return
	}

	for _, item := range req.Items {
		if item.ProductID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "product_id wajib diisi dan harus lebih dari 0",
			})
			return
		}
		if item.Quantity <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "quantity wajib diisi dan harus lebih dari 0",
			})
			return
		}
	}

	transaction, err := h.service.Checkout(req.Items, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

func (h *TransactionHandler) GetReport(c *gin.Context) {
	report, err := h.service.GetReport()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, report)
}

func (h *TransactionHandler) GetReportByDateRange(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "start_date dan end_date wajib diisi",
		})
		return
	}

	report, err := h.service.GetReportByDateRange(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, report)
}
