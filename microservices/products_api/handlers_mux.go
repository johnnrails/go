package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
)

type ProductHandlerMux struct {
	l *log.Logger
}

func NewProductHandlerMux(l *log.Logger) *ProductHandlerMux {
	return &ProductHandlerMux{l}
}

func (h *ProductHandlerMux) GetIDFromPath(path string) (int, error) {
	reg := regexp.MustCompile("/([0-9]+)")
	g := reg.FindAllStringSubmatch(path, -1)

	moreThanOneId := len(g) != 1
	if moreThanOneId {
		return -1, errors.New("More than one ID Found.")
	}

	moreThanOneCaptureGroup := len(g[0]) != 2
	if moreThanOneCaptureGroup {
		return -1, errors.New("More than one capture group found.")
	}

	idString := g[0][1]
	id, err := strconv.Atoi(idString)

	if err != nil {
		return -1, errors.New("Unable to parse ID to Int.")
	}

	return id, nil
}

func (h *ProductHandlerMux) GetProducts(w http.ResponseWriter, r *http.Request) {
	h.l.Println("Handle GET Products")
	products := GetProducts()
	if err := ToJSON(products, w); err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (h *ProductHandlerMux) AddProduct(w http.ResponseWriter, r *http.Request) {
	h.l.Println("Handle POST Product")
	product := r.Context().Value(KeyProduct{}).(Product)
	AddProduct(&product)
}

func (h *ProductHandlerMux) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	h.l.Println("Handle PUT Product")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Unable to convert id", http.StatusBadRequest)
		return
	}

	product := r.Context().Value(KeyProduct{}).(Product)
	err = UpdateProduct(id, &product)
	if err != nil {
		http.Error(w, "Product Not Found", http.StatusNotFound)
		return
	}
}

func (h ProductHandlerMux) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		product := Product{}
		if err := FromJSON(product, r.Body); err != nil {
			h.l.Println("[ERROR] deserializing product", err)
			http.Error(w, "Error reading product", http.StatusBadRequest)
			return
		}
		validation := NewValidation()
		if err := validation.Validate(product); err != nil {
			h.l.Println("[ERROR] Validating product", err)
			http.Error(w, fmt.Sprintf("Error validating product: %s", err), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, product)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
