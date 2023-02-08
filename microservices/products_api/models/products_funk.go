package models

import (
	"errors"

	"github.com/thoas/go-funk"
)

func Contains(p *Product) bool {
	return funk.Contains(products, p)
}

func FindProductByID(id int) (*Product, error) {
	r := funk.Find(products, func(p *Product) bool {
		return p.ID == id
	}).(*Product)
	if r == nil {
		return nil, errors.New("Product Not Found")
	}
	return r, nil
}

func DeleteProduct(id int) {
	products = funk.Filter(products, func(p *Product) bool {
		return p.ID != id
	}).([]*Product)
}
