package main

import (
	"context"
	"database/sql"
	"log"
	"strconv"

	_ "embed"

	"github.com/google/uuid"
	"github.com/johnnrails/ddd_go/working-with-db/sqlc"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var ddl string

func main() {
	ctx := context.Background()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatalln(err)
	}

	if _, err := db.ExecContext(ctx, ddl); err != nil {
		log.Fatalln(err)
	}

	queries := sqlc.New(db)

	for i := 0; i < 10; i++ {
		idx := strconv.FormatInt(int64(i), 10)
		userID := uuid.New().String()
		err = queries.CreateUser(ctx, sqlc.CreateUserParams{
			ID:    userID,
			Name:  "john" + idx,
			Email: "joaocarfe9@gmail.com" + idx,
			Bio: sql.NullString{
				String: "Bio",
				Valid:  true,
			},
			Password: "jhsfsdlfhsdlf" + idx,
		})
		if err != nil {
			log.Println("CREATE USER ERROR::::::", err)
		}
		price, err := strconv.Atoi(idx)
		if err != nil {
			log.Println("STRCONV ERROR::::::", err)
		}
		productID := uuid.New().String()
		err = queries.CreateProduct(ctx, sqlc.CreateProductParams{
			ID:     productID,
			Name:   "name" + idx,
			Code:   "code" + idx,
			Price:  int64(price),
			UserID: userID,
		})
		if err != nil {
			log.Println("CREATE PRODUCT ERROR::::::", err)
		}
	}

	users, err := queries.ListUsers(ctx)
	products, err := queries.ListProducts(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(users)
	log.Println(products)
}
