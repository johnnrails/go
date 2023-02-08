package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/thoas/go-funk"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedAt   string  `json:"-"`
	UpdatedAt   string  `json:"-"`
}

var products = []*Product{
	{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
		CreatedAt:   time.Now().UTC().String(),
		UpdatedAt:   time.Now().UTC().String(),
	},
	{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
		CreatedAt:   time.Now().UTC().String(),
		UpdatedAt:   time.Now().UTC().String(),
	},
}

func GetProducts() []*Product {
	return products
}

func FindProductByID(id int) (*Product, int, error) {
	product := &Product{}
	index := -1

	for i, p := range products {
		if p.ID == id {
			product = p
			index = i
		}
	}

	if index == -1 {
		return nil, index, errors.New("Product Not Found")
	}

	return product, index, nil
}

func Contains(p *Product) bool {
	return funk.Contains(products, p)
}

func DeleteProductByID(id int) error {
	product, i, err := FindProductByID(id)

	if err != nil {
		return err
	}

	last := products[len(products)-1]
	products[i] = last
	last = product

	products = funk.Initial(products).([]*Product)
	return nil
}

func AddProduct(p *Product) {
	p.ID = getNextID()
	products = append(products, p)
}

func UpdateProduct(id int, p *Product) error {
	_, pos, err := FindProduct(id)
	if err != nil {
		return err
	}
	products[pos] = p
	return nil
}

func FindProduct(id int) (*Product, int, error) {
	for i, p := range products {
		if p.ID == id {
			return p, i, nil
		}
	}
	return nil, -1, fmt.Errorf("Product Not Found")
}

func getNextID() int {
	return len(products)
}
