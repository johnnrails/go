package internal

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	var input CreateProductInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/products")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	repository := NewProductRepositoryMySQL(db)
	usecase := NewCreateProductUsecase(repository)

	err = usecase.Execute(input)
	db.Close()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func ListAllProductsHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/products")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	repository := NewProductRepositoryMySQL(db)
	usecase := NewListProductsUsecase(repository)

	products, err := usecase.Execute()
	db.Close()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}
