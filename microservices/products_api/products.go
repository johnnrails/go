package main

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	SKU         string  `json:"sku"`
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

func (p *Product) FromJSON(r io.Reader) error {
	return json.NewDecoder(r).Decode(p)
}

type Products []*Product

func (p *Products) ToJSON(w io.Writer) error {
	return json.NewEncoder(w).Encode(p)
}

func GetProducts() Products {
	return products
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
