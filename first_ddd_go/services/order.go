package services

import (
	"github.com/johnnrails/ddd_go/first_ddd_go/aggregates"
	"github.com/johnnrails/ddd_go/first_ddd_go/repositories/customer_repository"
	"github.com/johnnrails/ddd_go/first_ddd_go/repositories/product_repository"
)

type OrderService struct {
	customerRepository customer_repository.CustomerRepository
	productRepository  product_repository.ProductRepository
}

func NewOrderService(cr customer_repository.CustomerRepository, pr product_repository.ProductRepository) *OrderService {
	return &OrderService{
		customerRepository: cr,
		productRepository:  pr,
	}
}

func (o *OrderService) CalculatePrice(products []aggregates.Product) (float64, error) {
	var price float64
	for _, p := range products {
		price += p.GetPrice()
	}
	return price, nil
}
