package app

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-chi/chi/v5"
	"github.com/johnnrails/ddd_go/imersao12-full-cycle/internal"
)

func main() {

	r := chi.NewRouter()
	r.Post("/products", internal.CreateProductHandler)
	r.Get("/products", internal.ListAllProductsHandler)

	go http.ListenAndServe(":8080", r)

	msgChan := make(chan *kafka.Message)
	go internal.Consume([]string{"product"}, "localhost:9094", msgChan)

	db, err := sql.Open("mysql", "...")
	if err != nil {
		log.Fatal("")
	}
	defer db.Close()

	for msg := range msgChan {
		input := internal.CreateProductInput{}
		err := json.Unmarshal(msg.Value, &input)

		if err != nil {
		}

		repository := internal.NewProductRepositoryMySQL(db)
		usecase := internal.NewCreateProductUsecase(repository)

		err = usecase.Execute(input)
	}
}
