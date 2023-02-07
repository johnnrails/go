package sqlcapp

import (
	"context"
	"database/sql"
	"log"

	"github.com/johnnrails/ddd_go/imersao12-full-cycle/internal/db/sqlite"
)

var ddl string

func run() error {
	ctx := context.Background()
	db, err := sql.Open("sqlite3", ":memory:")

	if err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return err
	}

	queries := sqlite.New(db)

	authors, err := queries.ListAuthors(ctx)
	if err != nil {
		return err
	}
	log.Println(authors)

	err = queries.CreateAuthor(ctx, sqlite.CreateAuthorParams{
		Name: "Brian kerninghan",
		Bio:  sql.NullString{String: "Co-author of the C Programming Language"},
	})

	if err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
