package models

import (
	"errors"

	"github.com/thoas/go-funk"
)

func Contains(p *Product) bool {
	return funk.Contains(products, p)
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
