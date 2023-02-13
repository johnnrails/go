package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestInsertGorm(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	db.Create(&Product{Code: "D42", Price: 100})

	var product Product
	db.First(&product, Product{
		Price: 100,
	})

	assert.Equal(t, int(product.Price), 100)

	db.Delete(&product, product.Price)
}

func TestInsertASliceOfProducts(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	var products = []*Product{
		&Product{Code: "A1", Price: 1},
		&Product{Code: "B2", Price: 2},
		&Product{Code: "D3", Price: 3},
	}

	result := db.Create(&products)
	assert.Equal(t, 3, int(result.RowsAffected))
}

func TestProductsWithCategory(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	product := &Product{Code: "A1", Price: 1, Categories: []*Category{
		&Category{
			Name: "product-category-1",
		},
		&Category{
			Name: "product-category-2",
		},
	}}

	result := db.Create(&product)

	category := &Category{}
	db.Find(&category, Category{
		Name: "product-category-2",
	})
	fmt.Println(category)

	var product2 = &Product{}
	db.Find(&product2, "id = ?", product.ID)
	db.Model(&product2).Association("Categories").Find(&product2.Categories)

	fmt.Println(product2)
	fmt.Println(product2.Categories)

	db.Delete(&product2, product.Price)

	assert.NoError(t, result.Error)
	// assert.Equal(t, 2, len(product2.Categories))
	// assert.Equal(t, "product-category-2", product2.Categories[0].Name)
}
