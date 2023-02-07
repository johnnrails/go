package entities

import "github.com/google/uuid"

type Item struct {
	ID          uuid.UUID
	Name        string
	Description string
}

func (i Item) GetID() uuid.UUID {
	return i.ID
}

func (i Item) GetName() string {
	return i.Name
}

func (i Item) GetDescription() string {
	return i.Description
}
