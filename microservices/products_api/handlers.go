package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type ProductHandler struct {
	l *log.Logger
}

func NewProductHandler(l *log.Logger) *ProductHandler {
	return &ProductHandler{l}
}

func (h *ProductHandler) GetIDFromPath(path string) (int, error) {
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

func (h *ProductHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.GetProducts(w, r)
		return
	}

	if r.Method == http.MethodPost {
		h.AddProduct(w, r)
		return
	}

	if r.Method == http.MethodPut {
		h.l.Println("PUT", r.URL.Path)
		id, err := h.GetIDFromPath(r.URL.Path)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		h.UpdateProduct(id, w, r)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	h.l.Println("Handle GET Products")
	products := GetProducts()
	if err := ToJSON(products, w); err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (h *ProductHandler) AddProduct(w http.ResponseWriter, r *http.Request) {
	h.l.Println("Handle POST Product")
	product := &Product{}
	if err := FromJSON(products, r.Body); err != nil {
		http.Error(w, "Unable to marshal json", http.StatusBadRequest)
	}
	AddProduct(product)
}

func (h *ProductHandler) UpdateProduct(id int, w http.ResponseWriter, r *http.Request) {
	h.l.Println("Handle POST Product")
	product := &Product{}
	err := FromJSON(product, r.Body)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusBadRequest)
	}
	err = UpdateProduct(id, product)
	if err != nil {
		http.Error(w, "Product Not Found", http.StatusNotFound)
		return
	}
}

type KeyProduct struct{}

func (h ProductHandler) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		product := Product{}
		if err := FromJSON(product, r.Body); err != nil {
			h.l.Println("[ERROR] deserializing product", err)
			http.Error(w, "Error reading product", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, product)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
