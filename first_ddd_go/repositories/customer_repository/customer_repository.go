package customer_repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/johnnrails/ddd_go/first_ddd_go/aggregates"
)

var (
	ErrCustomerNotFound       = errors.New("customer not found")
	ErrFailedToAddCustomer    = errors.New("failed to add customer")
	ErrFailedToUpdateCustomer = errors.New("failed to update customer")
)

type CustomerRepository interface {
	Get(uuid.UUID) (aggregates.Customer, error)
	Add(aggregates.Customer) error
	Update(aggregates.Customer) error
}
