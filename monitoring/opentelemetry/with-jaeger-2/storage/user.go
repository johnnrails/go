package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/johnnrails/ddd_go/monitoring/opentelemetry/with-jaeger-2/trace"
)

type UserStorer interface {
	Insert(ctx context.Context, user User) error
}

type User struct {
	ID   int
	Name string
}

type UserStorage struct {
	database *sql.DB
}

func NewUserStorage(dtb *sql.DB) UserStorage {
	return UserStorage{
		dtb,
	}
}

func (u UserStorage) Insert(ctx context.Context, user User) error {
	ctx, span := trace.NewSpan(ctx, "UserStorage.Insert", nil)
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	if _, err := u.database.ExecContext(ctx, `INSERT INTO users (name) VALUES (?)`, user.Name); err != nil {
		log.Println(err)
		return fmt.Errorf("insert: failed to execute query")
	}

	return nil
}