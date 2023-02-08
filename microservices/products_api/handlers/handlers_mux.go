package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/johnnrails/ddd_go/microservices/products_api/helpers"
	"github.com/johnnrails/ddd_go/microservices/products_api/models"
)

type ProductHandlerMux struct {
	l *log.Logger
}

func NewProductHandlerMux(l *log.Logger) *ProductHandlerMux {
	return &ProductHandlerMux{l}
}

func (h *ProductHandlerMux) GetProducts(w http.ResponseWriter, r *http.Request) {
	h.l.Println("Handle GET Products")
	products := models.GetProducts()
	if err := helpers.ToJSON(products, w); err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (h *ProductHandlerMux) AddProduct(w http.ResponseWriter, r *http.Request) {
	h.l.Println("Handle POST Product")
	product := r.Context().Value(KeyProduct{}).(*models.Product)
	models.AddProduct(product)
}

func (h *ProductHandlerMux) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	h.l.Println("Handle PUT Product")

	id := GetProductIDFromRequest(r)
	product := r.Context().Value(KeyProduct{}).(models.Product)
	if err := models.UpdateProduct(id, &product); err != nil {
		http.Error(w, "Product Not Found", http.StatusNotFound)
		return
	}
}

func GetProductIDFromRequest(r *http.Request) int {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic(err)
	}
	return id
}

func (h ProductHandlerMux) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		product := &models.Product{}

		if err := helpers.FromJSON(product, r.Body); err != nil {
			h.l.Println("[ERROR] deserializing product", err)
			w.WriteHeader(http.StatusBadRequest)
			helpers.ToJSON(err.Error(), w)
			return
		}

		validation := helpers.NewValidation()
		errs := validation.Validate(product)

		if len(errs) != 0 {
			h.l.Println("[ERROR] validating product", errs)
			w.WriteHeader(http.StatusUnprocessableEntity)
			helpers.ToJSON(errs.Errors()[0], w)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, product)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
