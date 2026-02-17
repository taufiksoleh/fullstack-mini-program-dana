package entity

import (
	"time"

	"github.com/google/uuid"
)

type PropertyImage struct {
	ID         uuid.UUID
	PropertyID uuid.UUID
	URL        string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
