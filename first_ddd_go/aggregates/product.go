package aggregates

import (
	"errors"

	"github.com/google/uuid"
	"github.com/johnnrails/ddd_go/first_ddd_go/entities"
)

type Product struct {
	item     *entities.Item
	price    float64
	quantity int
}

func NewProduct(name string, description string, price float64) (Product, error) {
	if name == "" || description == "" {
		return Product{}, errors.New("missing values")
	}

	return Product{
		item: &entities.Item{
			ID:          uuid.New(),
			Name:        name,
			Description: description,
		},
		price:    price,
		quantity: 0,
	}, nil
}

func (p Product) GetID() uuid.UUID {
	return p.item.GetID()
}

func (p Product) GetName() string {
	return p.item.GetName()
}

func (p Product) GetDescription() string {
	return p.item.GetDescription()
}

func (p Product) GetItem() *entities.Item {
	return p.item
}

func (p Product) GetPrice() float64 {
	return p.price
}

func (p Product) GetQuantity() int {
	return p.quantity
}
