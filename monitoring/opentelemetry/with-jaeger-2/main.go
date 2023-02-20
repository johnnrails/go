package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/johnnrails/ddd_go/monitoring/opentelemetry/with-jaeger-2/storage"
	"github.com/johnnrails/ddd_go/monitoring/opentelemetry/with-jaeger-2/trace"
	"github.com/johnnrails/ddd_go/monitoring/opentelemetry/with-jaeger-2/users"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	ctx := context.Background()

	provider, err := trace.NewProvider(trace.ProviderConfig{
		JaegerEndpoint: "http://localhost:14268/api/traces",
		ServiceName:    "client",
		ServiceVersion: "1.0.0",
		Environment:    "dev",
		Disabled:       false,
	})
	if err != nil {
		log.Fatalln(err)
	}
	defer provider.Close(ctx)

	db, err := sql.Open("mysql", "root:@tcp(:3307)/jaeger")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	user := users.NewController(storage.NewUserStorage(db))
	router := http.DefaultServeMux
	router.HandleFunc("/users", trace.HTTPHandlerFunc(user.Execute, "users_create_handler").ServeHTTP)

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalln(err)
	}
}