package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/taufiksoleh/backend-mp/internal/domain/entity"
)

type Property struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	DocumentID  string         `gorm:"uniqueIndex;size:100"`
	Title       string         `gorm:"not null;size:255"`
	Description string         `gorm:"type:text"`
	BannerURL   string         `gorm:"type:text"`
	Price       float64        `gorm:"type:numeric(15,2);not null"`
	Terms       string         `gorm:"type:text"`
	Conditions  string         `gorm:"type:text"`
	Images      []PropertyImage `gorm:"foreignKey:PropertyID;constraint:OnDelete:CASCADE"`
	Facilities  []Facility     `gorm:"many2many:property_facilities"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Property) TableName() string { return "properties" }

func (m Property) ToEntity() entity.Property {
	images := make([]entity.PropertyImage, len(m.Images))
	for i, img := range m.Images {
		images[i] = img.ToEntity()
	}
	facilities := make([]entity.Facility, len(m.Facilities))
	for i, f := range m.Facilities {
		facilities[i] = f.ToEntity()
	}
	return entity.Property{
		ID:          m.ID,
		DocumentID:  m.DocumentID,
		Title:       m.Title,
		Description: m.Description,
		BannerURL:   m.BannerURL,
		Price:       m.Price,
		Terms:       m.Terms,
		Conditions:  m.Conditions,
		Images:      images,
		Facilities:  facilities,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func PropertyFromEntity(e entity.Property) Property {
	images := make([]PropertyImage, len(e.Images))
	for i, img := range e.Images {
		images[i] = PropertyImageFromEntity(img)
	}
	facilities := make([]Facility, len(e.Facilities))
	for i, f := range e.Facilities {
		facilities[i] = FacilityFromEntity(f)
	}
	return Property{
		ID:          e.ID,
		DocumentID:  e.DocumentID,
		Title:       e.Title,
		Description: e.Description,
		BannerURL:   e.BannerURL,
		Price:       e.Price,
		Terms:       e.Terms,
		Conditions:  e.Conditions,
		Images:      images,
		Facilities:  facilities,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

