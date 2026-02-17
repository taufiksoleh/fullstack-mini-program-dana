package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/taufiksoleh/backend-mp/internal/domain/entity"
)

type PropertyImage struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	PropertyID uuid.UUID `gorm:"type:uuid;not null;index"`
	URL        string    `gorm:"type:text;not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (PropertyImage) TableName() string { return "property_images" }

func (m PropertyImage) ToEntity() entity.PropertyImage {
	return entity.PropertyImage{
		ID:         m.ID,
		PropertyID: m.PropertyID,
		URL:        m.URL,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

func PropertyImageFromEntity(e entity.PropertyImage) PropertyImage {
	return PropertyImage{
		ID:         e.ID,
		PropertyID: e.PropertyID,
		URL:        e.URL,
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
	}
}
