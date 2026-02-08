package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
	"time"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	tx, err:= repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	
	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	for _, item := range items {
		var productPrice, stock int
		var productName string

		err := tx.QueryRow("SELECT name, price, stock FROM products WHERE id = $1", item.ProductID).Scan(&productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Product Id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		subTotal := productPrice * item.Quantity
		totalAmount += subTotal

		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID: item.ProductID,
			ProductName: productName,
			Quantity: item.Quantity,
			Subtotal: subTotal,
		})
	}

	var transactionID int
	var createdAt time.Time
	err = tx.QueryRow(
		"INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id, created_at",
		totalAmount,
	).Scan(&transactionID, &createdAt)
	if err != nil {
		return nil, err
	}

	insertDetailQuery := "INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4) RETURNING id"

	for i := range details {
		detail := &details[i]
		detail.TransactionID = transactionID

		err = tx.QueryRow(
			insertDetailQuery,
			detail.TransactionID,
			detail.ProductID,
			detail.Quantity,
			detail.Subtotal,
		).Scan(&detail.ID)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID: transactionID,
		TotalAmount: totalAmount,
		CreatedAt: createdAt,
		Details: details,
	}, nil
}

func (repo *TransactionRepository) GetReport() (*models.TransactionReport, error) {
	var totalRevenue int
	var totalTransaksi int

	err := repo.db.QueryRow(`
		SELECT COALESCE(SUM(total_amount), 0), COALESCE(COUNT(*), 0) 
		FROM transactions
		WHERE DATE(created_at) = CURRENT_DATE`).Scan(&totalRevenue, &totalTransaksi)
	if err != nil {
		return nil, err
	}

	var productName string
	var qtyTerjual int
	err = repo.db.QueryRow(`
		SELECT p.name, COALESCE(SUM(td.quantity), 0) AS total_qty
		FROM transaction_details td
		JOIN products p ON p.id = td.product_id
		JOIN transactions t ON t.id = td.transaction_id
		WHERE DATE(t.created_at) = CURRENT_DATE
		GROUP BY p.id, p.name
		ORDER BY total_qty DESC
		LIMIT 1
	`).Scan(&productName, &qtyTerjual)
	if err != nil {
		if err == sql.ErrNoRows {
			productName = ""
			qtyTerjual = 0
		} else {
			return nil, err
		}
	}

	return &models.TransactionReport{
		TotalRevenue:   totalRevenue,
		TotalTransaksi: totalTransaksi,
		ProdukTerlaris: models.BestSellProduct{
			Nama:       productName,
			QtyTerjual: qtyTerjual,
		},
	}, nil
}

func (repo *TransactionRepository) GetReportByDateRange(startDate, endDate string) (*models.TransactionReport, error) {
	var totalRevenue int
	var totalTransaksi int

	err := repo.db.QueryRow(`
		SELECT COALESCE(SUM(total_amount), 0), COALESCE(COUNT(*), 0)
		FROM transactions
		WHERE created_at >= $1 AND created_at <= $2`, 
		startDate, endDate).Scan(&totalRevenue, &totalTransaksi)
	if err != nil {
		return nil, err
	}

	var productName string
	var qtyTerjual int
	err = repo.db.QueryRow(`
		SELECT p.name, COALESCE(SUM(td.quantity), 0) AS total_qty
		FROM transaction_details td
		JOIN products p ON p.id = td.product_id
		JOIN transactions t ON t.id = td.transaction_id
		WHERE t.created_at >= $1 AND t.created_at <= $2
		GROUP BY p.id, p.name
		ORDER BY total_qty DESC
		LIMIT 1
	`, startDate, endDate).Scan(&productName, &qtyTerjual)
	if err != nil {
		if err == sql.ErrNoRows {
			productName = ""
			qtyTerjual = 0
		} else {
			return nil, err
		}
	}

	return &models.TransactionReport{
		TotalRevenue:   totalRevenue,
		TotalTransaksi: totalTransaksi,
		ProdukTerlaris: models.BestSellProduct{
			Nama:       productName,
			QtyTerjual: qtyTerjual,
		},
	}, nil
}
