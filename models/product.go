package models

import "time"

type Product struct {
	ID						int					`json:"id"`
	CategoryID		int					`json:"category_id"`
	CategoryName	string			`json:"category_name"`
	Name					string			`json:"name"`
	Price					int					`json:"price"`
	Stock					int					`json:"stock"`
	CreatedAt			*time.Time	`json:"created_at"`
}
