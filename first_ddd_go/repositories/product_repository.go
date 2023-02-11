package repositories

import (
	"errors"

	"github.com/google/uuid"
	"github.com/johnnrails/ddd_go/first_ddd_go/entities"
)

var (
	ErrProductNotFound     = errors.New("product not found")
	ErrProductAlreadyExist = errors.New("product already existt")
)

type ProductRepository interface {
	GetAll() ([]entities.Product, error)
	GetById(id uuid.UUID) (entities.Product, error)
	Add(product entities.Product) error
	Update(product entities.Product) error
	Delete(id uuid.UUID) error
}
