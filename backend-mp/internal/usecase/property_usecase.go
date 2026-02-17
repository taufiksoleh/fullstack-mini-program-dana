package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/taufiksoleh/backend-mp/internal/domain/apperror"
	"github.com/taufiksoleh/backend-mp/internal/domain/entity"
	"github.com/taufiksoleh/backend-mp/internal/domain/repository"
	"github.com/taufiksoleh/backend-mp/pkg/logger"
	"gorm.io/gorm"
)

type PropertyRequest struct {
	DocumentID    string
	Title         string
	Description   string
	BannerURL     string
	Price         float64
	Terms         string
	Conditions    string
	ImageURLs     []string
	FacilityNames []string
}

type PropertyUsecase interface {
	GetAll(ctx context.Context, filter repository.PropertyFilter) (repository.PaginatedProperties, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Property, error)
	GetByDocumentID(ctx context.Context, documentID string) (*entity.Property, error)
	Create(ctx context.Context, req PropertyRequest) (*entity.Property, error)
	Update(ctx context.Context, id uuid.UUID, req PropertyRequest) (*entity.Property, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type propertyUsecase struct {
	repo repository.PropertyRepository
	db   *gorm.DB
}

func NewPropertyUsecase(repo repository.PropertyRepository, db *gorm.DB) PropertyUsecase {
	return &propertyUsecase{repo: repo, db: db}
}

func validatePropertyRequest(req PropertyRequest) error {
	if strings.TrimSpace(req.Title) == "" {
		return fmt.Errorf("%w: %s", apperror.ErrValidation, errTitleRequired)
	}
	if req.Price <= 0 {
		return fmt.Errorf("%w: %s", apperror.ErrValidation, errPriceInvalid)
	}
	return nil
}

func (uc *propertyUsecase) GetAll(ctx context.Context, filter repository.PropertyFilter) (result repository.PaginatedProperties, err error) {
	defer logger.Track(ctx, "usecase.GetAll", "keyword", filter.Keyword, "page", filter.Page)(&err)

	result, err = uc.repo.FindAllWithFilter(ctx, filter)
	return
}

func (uc *propertyUsecase) GetByID(ctx context.Context, id uuid.UUID) (p *entity.Property, err error) {
	defer logger.Track(ctx, "usecase.GetByID", "id", id)(&err)

	p, err = uc.repo.FindByID(ctx, id)
	if err != nil {
		err = fmt.Errorf(errGetByID, err)
	}
	return
}

func (uc *propertyUsecase) GetByDocumentID(ctx context.Context, documentID string) (p *entity.Property, err error) {
	defer logger.Track(ctx, "usecase.GetByDocumentID", "document_id", documentID)(&err)

	p, err = uc.repo.FindByDocumentID(ctx, documentID)
	if err != nil {
		err = fmt.Errorf(errGetByDocumentID, err)
	}
	return
}

func (uc *propertyUsecase) Create(ctx context.Context, req PropertyRequest) (property *entity.Property, err error) {
	defer logger.Track(ctx, "usecase.Create", "title", req.Title)(&err)

	if err = validatePropertyRequest(req); err != nil {
		return
	}

	txErr := uc.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := uc.repo.WithTx(tx)

		facilities, fErr := uc.resolveFacilitiesWithRepo(ctx, txRepo, req.FacilityNames)
		if fErr != nil {
			return fErr
		}

		images := make([]entity.PropertyImage, len(req.ImageURLs))
		for i, url := range req.ImageURLs {
			images[i] = entity.PropertyImage{ID: uuid.New(), URL: url}
		}

		property = &entity.Property{
			ID:          uuid.New(),
			DocumentID:  req.DocumentID,
			Title:       req.Title,
			Description: req.Description,
			BannerURL:   req.BannerURL,
			Price:       req.Price,
			Terms:       req.Terms,
			Conditions:  req.Conditions,
			Images:      images,
			Facilities:  facilities,
		}

		return txRepo.Create(ctx, property)
	})
	if txErr != nil {
		err = fmt.Errorf(errCreate, txErr)
	}
	return
}

func (uc *propertyUsecase) Update(ctx context.Context, id uuid.UUID, req PropertyRequest) (property *entity.Property, err error) {
	defer logger.Track(ctx, "usecase.Update", "id", id)(&err)

	if err = validatePropertyRequest(req); err != nil {
		return
	}

	txErr := uc.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := uc.repo.WithTx(tx)

		p, fErr := txRepo.FindByID(ctx, id)
		if fErr != nil {
			return fmt.Errorf(errUpdate, fErr)
		}

		facilities, fErr := uc.resolveFacilitiesWithRepo(ctx, txRepo, req.FacilityNames)
		if fErr != nil {
			return fErr
		}

		if len(req.ImageURLs) > 0 {
			images := make([]entity.PropertyImage, len(req.ImageURLs))
			for i, url := range req.ImageURLs {
				images[i] = entity.PropertyImage{ID: uuid.New(), PropertyID: p.ID, URL: url}
			}
			p.Images = images
		}

		p.DocumentID = req.DocumentID
		p.Title = req.Title
		p.Description = req.Description
		p.BannerURL = req.BannerURL
		p.Price = req.Price
		p.Terms = req.Terms
		p.Conditions = req.Conditions
		p.Facilities = facilities

		if fErr := txRepo.Update(ctx, p); fErr != nil {
			return fmt.Errorf(errUpdate, fErr)
		}
		property = p
		return nil
	})
	if txErr != nil {
		err = txErr
	}
	return
}

func (uc *propertyUsecase) Delete(ctx context.Context, id uuid.UUID) (err error) {
	defer logger.Track(ctx, "usecase.Delete", "id", id)(&err)

	if err = uc.repo.Delete(ctx, id); err != nil {
		err = fmt.Errorf(errDelete, err)
	}
	return
}

func (uc *propertyUsecase) resolveFacilitiesWithRepo(ctx context.Context, repo repository.PropertyRepository, names []string) ([]entity.Facility, error) {
	facilities := make([]entity.Facility, 0, len(names))
	for _, name := range names {
		f, err := repo.FindOrCreateFacility(ctx, name)
		if err != nil {
			return nil, fmt.Errorf(errResolveFacility, err)
		}
		facilities = append(facilities, *f)
	}
	return facilities, nil
}
