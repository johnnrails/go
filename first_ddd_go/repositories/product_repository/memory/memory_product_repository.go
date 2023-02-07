package memory

import (
	"sync"

	"github.com/google/uuid"
	"github.com/johnnrails/ddd_go/first_ddd_go/aggregates"
	"github.com/johnnrails/ddd_go/first_ddd_go/repositories/product_repository"
)

type MemoryProductRepository struct {
	products map[uuid.UUID]aggregates.Product
	sync.Mutex
}

func New() *MemoryProductRepository {
	return &MemoryProductRepository{
		products: make(map[uuid.UUID]aggregates.Product),
	}
}

func (mr *MemoryProductRepository) GetAll() ([]aggregates.Product, error) {
	var products []aggregates.Product
	for _, p := range mr.products {
		products = append(products, p)
	}
	return products, nil
}

func (mr *MemoryProductRepository) GetById(id uuid.UUID) (aggregates.Product, error) {
	if product, ok := mr.products[id]; ok {
		return product, nil
	}
	return aggregates.Product{}, product_repository.ErrProductNotFound
}

func (mr *MemoryProductRepository) Add(prod aggregates.Product) error {
	mr.Lock()
	defer mr.Unlock()
	if _, ok := mr.products[prod.GetID()]; ok {
		return product_repository.ErrProductAlreadyExist
	}
	mr.products[prod.GetID()] = prod
	return nil
}

func (mr *MemoryProductRepository) Update(prod aggregates.Product) error {
	mr.Lock()
	defer mr.Unlock()
	if _, ok := mr.products[prod.GetID()]; !ok {
		return product_repository.ErrProductNotFound
	}
	mr.products[prod.GetID()] = prod
	return nil
}

func (mr *MemoryProductRepository) Delete(id uuid.UUID) error {
	mr.Lock()
	defer mr.Unlock()
	if _, ok := mr.products[id]; !ok {
		return product_repository.ErrProductNotFound
	}
	delete(mr.products, id)
	return nil
}
