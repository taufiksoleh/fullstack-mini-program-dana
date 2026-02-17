package entity

import (
	"time"

	"github.com/google/uuid"
)

type Facility struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
