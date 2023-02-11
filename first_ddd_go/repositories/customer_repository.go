package repositories

import (
	"errors"

	"github.com/google/uuid"
	"github.com/johnnrails/ddd_go/first_ddd_go/entities"
)

var (
	ErrCustomerNotFound       = errors.New("customer not found")
	ErrFailedToAddCustomer    = errors.New("failed to add customer")
	ErrFailedToUpdateCustomer = errors.New("failed to update customer")
)

type CustomerRepository interface {
	Get(uuid.UUID) (entities.Customer, error)
	Add(entities.Customer) error
	Update(entities.Customer) error
}
