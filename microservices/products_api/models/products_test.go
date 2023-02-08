package models

import (
	"testing"

	"github.com/johnnrails/ddd_go/microservices/products_api/helpers"
	"github.com/stretchr/testify/assert"
)

func TestProductMissingNameAndSKUReturnError(t *testing.T) {
	p := Product{
		Price: 1.22,
	}
	v := helpers.NewValidation()
	err := v.Validate(p)
	assert.Len(t, err, 2)
}
func TestProductMissingPriceReturnError(t *testing.T) {
	p := Product{
		Name:  "abc",
		SKU:   "abc-abd-ded",
		Price: -1,
	}

	v := helpers.NewValidation()
	err := v.Validate(p)
	assert.Len(t, err, 1)
}

func TestProductsFunkContains(t *testing.T) {
	p := products[1]
	contains := Contains(p)
	assert.Equal(t, true, contains)
}

func TestProductsFunkFindProductByID(t *testing.T) {
	p := products[1]
	pd, _ := FindProductByID(p.ID)
	assert.Equal(t, p, pd)
}

func TestProductsFunkDeleteProduct(t *testing.T) {
	p := products[1]
	DeleteProductByID(p.ID)
	assert.Equal(t, 1, len(products))
}
