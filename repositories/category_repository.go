package repositories

import (
	"database/sql"
	"kasir-api/models"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (repo *CategoryRepository) GetAll() ([]models.Category, error) {
	query := "SELECT id, name, description, created_at FROM categories"
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]models.Category, 0)
	for rows.Next() {
		var c models.Category
		err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}

func (repo *CategoryRepository) Create(category *models.Category) error {
	query := "INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id, created_at"
	err := repo.db.QueryRow(query, category.Name, category.Description).Scan(&category.ID, &category.CreatedAt)
	return err
}

func (repo *CategoryRepository) GetByID(id string) (*models.Category, error) {
	query := "SELECT id, name, description, created_at FROM categories WHERE id = $1"
	row := repo.db.QueryRow(query, id)
	var c models.Category
	err := row.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (repo *CategoryRepository) Update(id string, category *models.Category) error {
	query := "UPDATE categories SET name = $1, description = $2 WHERE id = $3"
	_, err := repo.db.Exec(query, category.Name, category.Description, id)
	return err
}

func (repo *CategoryRepository) Delete(id string) error {
	query := "DELETE FROM categories WHERE id = $1"
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
