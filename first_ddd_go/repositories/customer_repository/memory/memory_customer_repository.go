package memory

import (
	"sync"

	"github.com/google/uuid"
	"github.com/johnnrails/ddd_go/first_ddd_go/aggregates"
	"github.com/johnnrails/ddd_go/first_ddd_go/repositories/customer_repository"
)

type MemoryRepository struct {
	customers []aggregates.Customer
	sync.Mutex
}

func New() *MemoryRepository {
	return &MemoryRepository{
		customers: make([]aggregates.Customer, 0),
	}
}

func (mr *MemoryRepository) Get(id uuid.UUID) (aggregates.Customer, error) {
	for _, c := range mr.customers {
		if id.ID() == c.GetID().ID() {
			return c, nil
		}
	}
	return aggregates.Customer{}, customer_repository.ErrCustomerNotFound
}

func (mr *MemoryRepository) Add(c aggregates.Customer) error {
	if mr.customers == nil {
		mr.Lock()
		mr.customers = make([]aggregates.Customer, 0)
		mr.Unlock()
	}

	if _, err := mr.Get(c.GetID()); err == nil {
		return customer_repository.ErrFailedToAddCustomer
	}

	mr.Lock()
	mr.customers = append(mr.customers, c)
	mr.Unlock()
	return nil
}

func (mr *MemoryRepository) Update(c aggregates.Customer) error {
	cust, err := mr.Get(c.GetID())

	if err != nil {
		return customer_repository.ErrFailedToUpdateCustomer
	}

	mr.Lock()
	cust.SetName(c.GetName())
	cust.SetAge(c.GetAge())
	cust.SetProducts(c.GetProducts())
	cust.SetTransactions(c.GetTransactions())
	mr.Unlock()
	return nil
}
