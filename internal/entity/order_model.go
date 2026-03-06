package entity

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Items     []string
	Status    string
	CreatedAt time.Time
}
