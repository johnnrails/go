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
