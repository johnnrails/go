package internal

import (
	"database/sql"
)

type ProductRepositoryMySQL struct {
	DB *sql.DB
}

func NewProductRepositoryMySQL(db *sql.DB) *ProductRepositoryMySQL {
	return &ProductRepositoryMySQL{db}
}

func (r *ProductRepositoryMySQL) Create(product *Product) error {
	_, err := r.DB.Exec("Insert into products (id, name, price) values(?,?,?)", &product.ID, &product.Name, &product.Price)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductRepositoryMySQL) FindAll() ([]*Product, error) {
	rows, err := r.DB.Query("select id, name, price from products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price); err != nil {
			return nil, err
		}
		products = append(products, &product)
	}
	return products, nil
}
