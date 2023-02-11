package entities

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	amount int
	from   uuid.UUID
	to     uuid.UUID
	when   time.Time
}
