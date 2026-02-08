package repositories

import (
	"database/sql"
	"kasir-api/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (repo *ProductRepository) GetAll(nameFilter string) ([]models.Product, error) {
	query := `
		SELECT 
			p.id,
			p.category_id,
			p.name, 
			p.price, 
			p.stock, 
			p.created_at,
			c.name AS category_name
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
	`

	var args []interface{}
	if nameFilter != "" {
		query += " WHERE p.name ILIKE $1"
		args = append(args, "%"+nameFilter+"%")
	}

	rows, err := repo.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		err := rows.Scan(
			&p.ID,
			&p.CategoryID,
			&p.Name,
			&p.Price,
			&p.Stock,
			&p.CreatedAt,
			&p.CategoryName,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (repo *ProductRepository) Create(product *models.Product) error {
	query := "INSERT INTO products (category_id, name, price, stock) VALUES ($1, $2, $3, $4) RETURNING id, created_at"
	err := repo.db.QueryRow(query, product.CategoryID, product.Name, product.Price, product.Stock).Scan(&product.ID, &product.CreatedAt)
	return err
}

func (repo *ProductRepository) GetByID(id string) (*models.Product, error) {
	query := `
		SELECT 
			p.id, 
			p.category_id, 
			p.name, 
			p.price, 
			p.stock, 
			p.created_at
		FROM products p
		WHERE p.id = $1
	`
	row := repo.db.QueryRow(query, id)
	var p models.Product
	err := row.Scan(
		&p.ID,
		&p.CategoryID,
		&p.Name,
		&p.Price,
		&p.Stock,
		&p.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (repo *ProductRepository) Update(id string, product *models.Product) error {
	query := "UPDATE products SET category_id = $1, name = $2, price = $3, stock = $4 WHERE id = $5"
	_, err := repo.db.Exec(query, product.CategoryID, product.Name, product.Price, product.Stock, id)
	return err
}

func (repo *ProductRepository) Delete(id string) error {
	query := "DELETE FROM products WHERE id = $1"
	res, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
