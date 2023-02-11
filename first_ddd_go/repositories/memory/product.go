package memory

import (
	"sync"

	"github.com/google/uuid"
	"github.com/johnnrails/ddd_go/first_ddd_go/entities"
	"github.com/johnnrails/ddd_go/first_ddd_go/repositories"
)

type MemoryProductRepository struct {
	products map[uuid.UUID]entities.Product
	sync.Mutex
}

func NewProductRepository() *MemoryProductRepository {
	return &MemoryProductRepository{
		products: make(map[uuid.UUID]entities.Product),
	}
}

func (mr *MemoryProductRepository) GetAll() ([]entities.Product, error) {
	var products []entities.Product
	for _, p := range mr.products {
		products = append(products, p)
	}
	return products, nil
}

func (mr *MemoryProductRepository) GetById(id uuid.UUID) (entities.Product, error) {
	if product, ok := mr.products[id]; ok {
		return product, nil
	}
	return entities.Product{}, repositories.ErrProductNotFound
}

func (mr *MemoryProductRepository) Add(prod entities.Product) error {
	mr.Lock()
	defer mr.Unlock()
	if _, ok := mr.products[prod.GetID()]; ok {
		return repositories.ErrProductAlreadyExist
	}
	mr.products[prod.GetID()] = prod
	return nil
}

func (mr *MemoryProductRepository) Update(prod entities.Product) error {
	mr.Lock()
	defer mr.Unlock()
	if _, ok := mr.products[prod.GetID()]; !ok {
		return repositories.ErrProductNotFound
	}
	mr.products[prod.GetID()] = prod
	return nil
}

func (mr *MemoryProductRepository) Delete(id uuid.UUID) error {
	mr.Lock()
	defer mr.Unlock()
	if _, ok := mr.products[id]; !ok {
		return repositories.ErrProductNotFound
	}
	delete(mr.products, id)
	return nil
}
