package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type TransactionService struct {
	transactionRepo *repositories.TransactionRepository
}

func NewTransactionService(transactionRepo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{transactionRepo: transactionRepo}
}

func (s *TransactionService) Checkout(items []models.CheckoutItem, useLock bool) (*models.Transaction, error) {
	return s.transactionRepo.CreateTransaction(items)
}

func (s *TransactionService) GetReport() (*models.TransactionReport, error) {
	return s.transactionRepo.GetReport()
}

func (s *TransactionService) GetReportByDateRange(startDate, endDate string) (*models.TransactionReport, error) {
	return s.transactionRepo.GetReportByDateRange(startDate, endDate)
}
