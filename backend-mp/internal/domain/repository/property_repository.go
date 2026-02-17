package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/taufiksoleh/backend-mp/internal/domain/entity"
	"gorm.io/gorm"
)

type PropertyFilter struct {
	Keyword  string
	Page     int
	PageSize int
}

type PaginatedProperties struct {
	Properties []entity.Property
	Total      int64
}

type PropertyRepository interface {
	FindAllWithFilter(ctx context.Context, filter PropertyFilter) (PaginatedProperties, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Property, error)
	FindByDocumentID(ctx context.Context, documentID string) (*entity.Property, error)
	Create(ctx context.Context, property *entity.Property) error
	Update(ctx context.Context, property *entity.Property) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindOrCreateFacility(ctx context.Context, name string) (*entity.Facility, error)
	WithTx(tx *gorm.DB) PropertyRepository
}
