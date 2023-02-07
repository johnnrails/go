package product_repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/johnnrails/ddd_go/first_ddd_go/aggregates"
)

var (
	ErrProductNotFound     = errors.New("product not found")
	ErrProductAlreadyExist = errors.New("product already existt")
)

type ProductRepository interface {
	GetAll() ([]aggregates.Product, error)
	GetById(id uuid.UUID) (aggregates.Product, error)
	Add(product aggregates.Product) error
	Update(product aggregates.Product) error
	Delete(id uuid.UUID) error
}
