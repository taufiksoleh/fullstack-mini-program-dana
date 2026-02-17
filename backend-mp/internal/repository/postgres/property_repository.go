package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/taufiksoleh/backend-mp/internal/domain/apperror"
	"github.com/taufiksoleh/backend-mp/internal/domain/entity"
	domainrepo "github.com/taufiksoleh/backend-mp/internal/domain/repository"
	pgmodel "github.com/taufiksoleh/backend-mp/internal/repository/postgres/model"
	"github.com/taufiksoleh/backend-mp/pkg/logger"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type propertyRepository struct {
	db *gorm.DB
}

func NewPropertyRepository(db *gorm.DB) domainrepo.PropertyRepository {
	return &propertyRepository{db: db}
}

func (r *propertyRepository) WithTx(tx *gorm.DB) domainrepo.PropertyRepository {
	return &propertyRepository{db: tx}
}

func (r *propertyRepository) FindAllWithFilter(ctx context.Context, filter domainrepo.PropertyFilter) (result domainrepo.PaginatedProperties, err error) {
	defer logger.Track(ctx, "repo.FindAllWithFilter", "keyword", filter.Keyword, "page", filter.Page)(&err)

	var models []pgmodel.Property
	q := r.db.WithContext(ctx).Model(&pgmodel.Property{})
	if filter.Keyword != "" {
		like := "%" + filter.Keyword + "%"
		q = q.Where("title ILIKE ?", like)
	}

	if err = q.Count(&result.Total).Error; err != nil {
		err = fmt.Errorf(errFindAllCount, err)
		return
	}

	offset := (filter.Page - 1) * filter.PageSize
	if err = q.Preload("Images").Preload("Facilities").
		Limit(filter.PageSize).Offset(offset).
		Find(&models).Error; err != nil {
		err = fmt.Errorf(errFindAll, err)
		return
	}

	result.Properties = make([]entity.Property, len(models))
	for i, m := range models {
		result.Properties[i] = m.ToEntity()
	}
	return
}

func (r *propertyRepository) FindByID(ctx context.Context, id uuid.UUID) (p *entity.Property, err error) {
	defer logger.Track(ctx, "repo.FindByID", "id", id)(&err)

	var m pgmodel.Property
	if err = r.db.WithContext(ctx).Preload("Images").Preload("Facilities").
		First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = apperror.ErrNotFound
			return
		}
		err = fmt.Errorf(errFindByID, err)
		return
	}
	e := m.ToEntity()
	p = &e
	return
}

func (r *propertyRepository) FindByDocumentID(ctx context.Context, documentID string) (p *entity.Property, err error) {
	defer logger.Track(ctx, "repo.FindByDocumentID", "document_id", documentID)(&err)

	var m pgmodel.Property
	if err = r.db.WithContext(ctx).Preload("Images").Preload("Facilities").
		First(&m, "document_id = ?", documentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = apperror.ErrNotFound
			return
		}
		err = fmt.Errorf(errFindByDocumentID, err)
		return
	}
	e := m.ToEntity()
	p = &e
	return
}

func (r *propertyRepository) Create(ctx context.Context, property *entity.Property) (err error) {
	defer logger.Track(ctx, "repo.Create", "id", property.ID)(&err)

	m := pgmodel.PropertyFromEntity(*property)
	if err = r.db.WithContext(ctx).Create(&m).Error; err != nil {
		err = fmt.Errorf(errCreate, err)
		return
	}
	*property = m.ToEntity()
	return
}

func (r *propertyRepository) Update(ctx context.Context, property *entity.Property) (err error) {
	defer logger.Track(ctx, "repo.Update", "id", property.ID)(&err)

	m := pgmodel.PropertyFromEntity(*property)
	if err = r.db.WithContext(ctx).
		Session(&gorm.Session{FullSaveAssociations: true}).
		Save(&m).Error; err != nil {
		err = fmt.Errorf(errUpdate, err)
		return
	}
	*property = m.ToEntity()
	return
}

func (r *propertyRepository) Delete(ctx context.Context, id uuid.UUID) (err error) {
	defer logger.Track(ctx, "repo.Delete", "id", id)(&err)

	if err = r.db.WithContext(ctx).Delete(&pgmodel.Property{}, "id = ?", id).Error; err != nil {
		err = fmt.Errorf(errDelete, err)
	}
	return
}

func (r *propertyRepository) FindOrCreateFacility(ctx context.Context, name string) (f *entity.Facility, err error) {
	defer logger.Track(ctx, "repo.FindOrCreateFacility", "name", name)(&err)

	var m pgmodel.Facility
	result := r.db.WithContext(ctx).
		Where(pgmodel.Facility{Name: name}).
		Clauses(clause.OnConflict{DoNothing: true}).
		FirstOrCreate(&m)
	if result.Error != nil {
		err = fmt.Errorf(errFindOrCreate, result.Error)
		return
	}
	e := m.ToEntity()
	f = &e
	return
}
