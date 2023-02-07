package aggregates

import (
	"errors"

	"github.com/google/uuid"
	"github.com/johnnrails/ddd_go/first_ddd_go/entities"
	"github.com/johnnrails/ddd_go/first_ddd_go/valueobjects"
)

type Customer struct {
	person       *entities.Person
	products     []*entities.Item
	transactions []valueobjects.Transaction
}

func NewCustomer(name string, age int) (Customer, error) {
	if name == "" {
		return Customer{}, errors.New("Invalid Name")
	}

	if age < 18 || age > 120 {
		return Customer{}, errors.New("Invalid Age")
	}

	return Customer{
		person: &entities.Person{
			ID:   uuid.New(),
			Name: name,
			Age:  age,
		},
		products:     make([]*entities.Item, 0),
		transactions: make([]valueobjects.Transaction, 0),
	}, nil
}

func (c *Customer) GetID() uuid.UUID {
	return c.person.ID
}
func (c *Customer) GetName() string {
	return c.person.Name
}
func (c *Customer) GetAge() int {
	return c.person.Age
}
func (c *Customer) GetProducts() []*entities.Item {
	return c.products
}
func (c *Customer) GetTransactions() []valueobjects.Transaction {
	return c.transactions
}

func (c *Customer) SetName(name string) {
	c.person.Name = name
}

func (c *Customer) SetID(id uuid.UUID) {
	c.person.ID = id
}

func (c *Customer) SetAge(age int) {
	c.person.Age = age
}

func (c *Customer) SetProducts(ps []*entities.Item) {
	c.products = ps
}

func (c *Customer) SetTransactions(ts []valueobjects.Transaction) {
	c.transactions = ts
}
