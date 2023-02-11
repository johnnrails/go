package entities

import (
	"errors"

	"github.com/google/uuid"
)

type Customer struct {
	person       *Person
	products     []*Item
	transactions []Transaction
}

func NewCustomer(name string, age int) (Customer, error) {
	if name == "" {
		return Customer{}, errors.New("Invalid Name")
	}

	if age < 18 || age > 120 {
		return Customer{}, errors.New("Invalid Age")
	}

	return Customer{
		person: &Person{
			ID:   uuid.New(),
			Name: name,
			Age:  age,
		},
		products:     make([]*Item, 0),
		transactions: make([]Transaction, 0),
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
func (c *Customer) GetProducts() []*Item {
	return c.products
}
func (c *Customer) GetTransactions() []Transaction {
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

func (c *Customer) SetProducts(ps []*Item) {
	c.products = ps
}

func (c *Customer) SetTransactions(ts []Transaction) {
	c.transactions = ts
}
