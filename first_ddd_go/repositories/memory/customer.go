package memory

import (
	"sync"

	"github.com/google/uuid"
	"github.com/johnnrails/ddd_go/first_ddd_go/entities"
	"github.com/johnnrails/ddd_go/first_ddd_go/repositories"
)

type MemoryCustomerRepository struct {
	customers []entities.Customer
	sync.Mutex
}

func NewMemoryCostumerRepository() *MemoryCustomerRepository {
	return &MemoryCustomerRepository{
		customers: make([]entities.Customer, 0),
	}
}

func (mr *MemoryCustomerRepository) Get(id uuid.UUID) (entities.Customer, error) {
	for _, c := range mr.customers {
		if id.ID() == c.GetID().ID() {
			return c, nil
		}
	}
	return entities.Customer{}, repositories.ErrCustomerNotFound
}

func (mr *MemoryCustomerRepository) Add(c entities.Customer) error {
	if mr.customers == nil {
		mr.Lock()
		mr.customers = make([]entities.Customer, 0)
		mr.Unlock()
	}

	if _, err := mr.Get(c.GetID()); err == nil {
		return repositories.ErrFailedToAddCustomer
	}

	mr.Lock()
	mr.customers = append(mr.customers, c)
	mr.Unlock()
	return nil
}

func (mr *MemoryCustomerRepository) Update(c entities.Customer) error {
	cust, err := mr.Get(c.GetID())

	if err != nil {
		return repositories.ErrFailedToUpdateCustomer
	}

	mr.Lock()
	cust.SetName(c.GetName())
	cust.SetAge(c.GetAge())
	cust.SetProducts(c.GetProducts())
	cust.SetTransactions(c.GetTransactions())
	mr.Unlock()
	return nil
}
