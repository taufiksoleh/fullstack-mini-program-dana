package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/taufiksoleh/backend-mp/internal/domain/entity"
)

type Facility struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name      string    `gorm:"uniqueIndex;size:100;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Facility) TableName() string { return "facilities" }

func (m Facility) ToEntity() entity.Facility {
	return entity.Facility{
		ID:        m.ID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func FacilityFromEntity(e entity.Facility) Facility {
	return Facility{
		ID:        e.ID,
		Name:      e.Name,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
