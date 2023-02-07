package entities

import (
	"github.com/google/uuid"
)

type Person struct {
	ID   uuid.UUID
	Name string
	Age  int
}
