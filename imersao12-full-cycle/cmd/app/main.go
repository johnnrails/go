package app

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-chi/chi/v5"
	"github.com/johnnrails/ddd_go/imersao12-full-cycle/internal/infra/akafka"
	"github.com/johnnrails/ddd_go/imersao12-full-cycle/internal/infra/repositories/mysql"
	"github.com/johnnrails/ddd_go/imersao12-full-cycle/internal/infra/web"
	"github.com/johnnrails/ddd_go/imersao12-full-cycle/internal/usecases"
)

func main() {

	r := chi.NewRouter()
	r.Post("/products", web.CreateProductHandler)
	r.Get("/products", web.ListAllProductsHandler)

	go http.ListenAndServe(":8080", r)

	msgChan := make(chan *kafka.Message)
	go akafka.Consume([]string{"product"}, "localhost:9094", msgChan)

	db, err := sql.Open("mysql", "...")
	if err != nil {
		log.Fatal("")
	}
	defer db.Close()

	for msg := range msgChan {
		input := usecases.CreateProductInput{}
		err := json.Unmarshal(msg.Value, &input)

		if err != nil {
		}

		repository := mysql.NewProductRepositoryMySQL(db)
		usecase := usecases.NewCreateProductUsecase(repository)

		err = usecase.Execute(input)
	}
}
