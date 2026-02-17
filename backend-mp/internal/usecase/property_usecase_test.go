package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/taufiksoleh/backend-mp/internal/domain/apperror"
	"github.com/taufiksoleh/backend-mp/internal/domain/entity"
	"github.com/taufiksoleh/backend-mp/internal/domain/repository"
	mockrepo "github.com/taufiksoleh/backend-mp/internal/domain/repository/mock"
	"github.com/taufiksoleh/backend-mp/internal/testdata/fixture"
	"github.com/taufiksoleh/backend-mp/internal/usecase"
	"go.uber.org/mock/gomock"
)

func newTestUsecase(t *testing.T) (usecase.PropertyUsecase, *mockrepo.MockPropertyRepository) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mockRepo := mockrepo.NewMockPropertyRepository(ctrl)
	// db is nil; Create/Update transaction path requires integration testing
	uc := usecase.NewPropertyUsecase(mockRepo, nil)
	return uc, mockRepo
}

// ---------------------------------------------------------------------------
// GetAll
// ---------------------------------------------------------------------------

func TestGetAll(t *testing.T) {
	filter := repository.PropertyFilter{Keyword: "", Page: 1, PageSize: 10}

	tests := []struct {
		name      string
		filter    repository.PropertyFilter
		mockSetup func(*mockrepo.MockPropertyRepository)
		wantTotal int64
		wantLen   int
		wantErr   bool
	}{
		{
			name:   "returns paginated results",
			filter: filter,
			mockSetup: func(m *mockrepo.MockPropertyRepository) {
				m.EXPECT().
					FindAllWithFilter(gomock.Any(), filter).
					Return(repository.PaginatedProperties{
						Properties: []entity.Property{*fixture.Property()},
						Total:      1,
					}, nil)
			},
			wantTotal: 1,
			wantLen:   1,
		},
		{
			name:   "propagates repo error",
			filter: filter,
			mockSetup: func(m *mockrepo.MockPropertyRepository) {
				m.EXPECT().
					FindAllWithFilter(gomock.Any(), filter).
					Return(repository.PaginatedProperties{}, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			uc, mockRepo := newTestUsecase(t)
			tc.mockSetup(mockRepo)

			result, err := uc.GetAll(context.Background(), tc.filter)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result.Total != tc.wantTotal {
				t.Errorf("total: want %d, got %d", tc.wantTotal, result.Total)
			}
			if len(result.Properties) != tc.wantLen {
				t.Errorf("len(Properties): want %d, got %d", tc.wantLen, len(result.Properties))
			}
		})
	}
}

// ---------------------------------------------------------------------------
// GetByID
// ---------------------------------------------------------------------------

func TestGetByID(t *testing.T) {
	prop := fixture.Property()

	tests := []struct {
		name      string
		id        uuid.UUID
		mockSetup func(*mockrepo.MockPropertyRepository)
		wantErr   error
	}{
		{
			name: "found",
			id:   prop.ID,
			mockSetup: func(m *mockrepo.MockPropertyRepository) {
				m.EXPECT().FindByID(gomock.Any(), prop.ID).Return(prop, nil)
			},
		},
		{
			name: "not found",
			id:   uuid.New(),
			mockSetup: func(m *mockrepo.MockPropertyRepository) {
				m.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(nil, apperror.ErrNotFound)
			},
			wantErr: apperror.ErrNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			uc, mockRepo := newTestUsecase(t)
			tc.mockSetup(mockRepo)

			result, err := uc.GetByID(context.Background(), tc.id)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("want error %v, got %v", tc.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result.ID != tc.id {
				t.Errorf("ID: want %s, got %s", tc.id, result.ID)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// GetByDocumentID
// ---------------------------------------------------------------------------

func TestGetByDocumentID(t *testing.T) {
	prop := fixture.Property()

	tests := []struct {
		name       string
		documentID string
		mockSetup  func(*mockrepo.MockPropertyRepository)
		wantErr    error
	}{
		{
			name:       "found",
			documentID: prop.DocumentID,
			mockSetup: func(m *mockrepo.MockPropertyRepository) {
				m.EXPECT().FindByDocumentID(gomock.Any(), prop.DocumentID).Return(prop, nil)
			},
		},
		{
			name:       "not found",
			documentID: "missing-doc",
			mockSetup: func(m *mockrepo.MockPropertyRepository) {
				m.EXPECT().FindByDocumentID(gomock.Any(), "missing-doc").Return(nil, apperror.ErrNotFound)
			},
			wantErr: apperror.ErrNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			uc, mockRepo := newTestUsecase(t)
			tc.mockSetup(mockRepo)

			result, err := uc.GetByDocumentID(context.Background(), tc.documentID)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("want error %v, got %v", tc.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result.DocumentID != tc.documentID {
				t.Errorf("DocumentID: want %s, got %s", tc.documentID, result.DocumentID)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Create — validation only (transaction path requires integration testing)
// ---------------------------------------------------------------------------

func TestCreate_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     usecase.PropertyRequest
		wantErr error
	}{
		{
			name:    "empty title",
			req:     usecase.PropertyRequest{Title: "", Price: 100000},
			wantErr: apperror.ErrValidation,
		},
		{
			name:    "whitespace title",
			req:     usecase.PropertyRequest{Title: "   ", Price: 100000},
			wantErr: apperror.ErrValidation,
		},
		{
			name:    "zero price",
			req:     usecase.PropertyRequest{Title: "Valid", Price: 0},
			wantErr: apperror.ErrValidation,
		},
		{
			name:    "negative price",
			req:     usecase.PropertyRequest{Title: "Valid", Price: -1},
			wantErr: apperror.ErrValidation,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			uc, _ := newTestUsecase(t)
			_, err := uc.Create(context.Background(), tc.req)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("want error %v, got %v", tc.wantErr, err)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Update — validation only (transaction path requires integration testing)
// ---------------------------------------------------------------------------

func TestUpdate_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     usecase.PropertyRequest
		wantErr error
	}{
		{
			name:    "empty title",
			req:     usecase.PropertyRequest{Title: "", Price: 100000},
			wantErr: apperror.ErrValidation,
		},
		{
			name:    "zero price",
			req:     usecase.PropertyRequest{Title: "Valid", Price: 0},
			wantErr: apperror.ErrValidation,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			uc, _ := newTestUsecase(t)
			_, err := uc.Update(context.Background(), uuid.New(), tc.req)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("want error %v, got %v", tc.wantErr, err)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Delete
// ---------------------------------------------------------------------------

func TestDelete(t *testing.T) {
	id := fixture.PropertyID

	tests := []struct {
		name      string
		mockSetup func(*mockrepo.MockPropertyRepository)
		wantErr   error
	}{
		{
			name: "success",
			mockSetup: func(m *mockrepo.MockPropertyRepository) {
				m.EXPECT().Delete(gomock.Any(), id).Return(nil)
			},
		},
		{
			name: "not found",
			mockSetup: func(m *mockrepo.MockPropertyRepository) {
				m.EXPECT().Delete(gomock.Any(), id).Return(apperror.ErrNotFound)
			},
			wantErr: apperror.ErrNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			uc, mockRepo := newTestUsecase(t)
			tc.mockSetup(mockRepo)

			err := uc.Delete(context.Background(), id)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("want error %v, got %v", tc.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
