package entity

import (
	"time"

	"github.com/google/uuid"
)

type Property struct {
	ID          uuid.UUID
	DocumentID  string
	Title       string
	Description string
	BannerURL   string
	Price       float64
	Terms       string
	Conditions  string
	Images      []PropertyImage
	Facilities  []Facility
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

