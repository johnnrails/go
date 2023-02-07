package main

import "testing"

func TestChecksValidat(t *testing.T) {
	p := &Product{
		Name:  "nics",
		Price: 1.00,
		SKU:   "abs-abc-ade",
	}
	validation := NewValidation()
	err := validation.Validate(p)
	if err != nil {
		t.Fatal(err)
	}
}
