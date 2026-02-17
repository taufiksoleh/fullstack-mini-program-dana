package postgres

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/taufiksoleh/backend-mp/internal/domain/apperror"
	"github.com/taufiksoleh/backend-mp/internal/domain/entity"
	domainrepo "github.com/taufiksoleh/backend-mp/internal/domain/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	propID   = uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	imgID    = uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
	facID    = uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")
	fixedNow = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	propCols = []string{"id", "document_id", "title", "description", "banner_url", "price", "terms", "conditions", "created_at", "updated_at"}
	imgCols  = []string{"id", "property_id", "url", "created_at", "updated_at"}
	facCols  = []string{"id", "name", "created_at", "updated_at"}
)

func newTestRepo(t *testing.T) (domainrepo.PropertyRepository, sqlmock.Sqlmock, *sql.DB) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("gorm.Open: %v", err)
	}
	return NewPropertyRepository(gormDB), mock, db
}

func propRow() *sqlmock.Rows {
	return sqlmock.NewRows(propCols).
		AddRow(propID, "doc-001", "Test", "Desc", "banner.jpg", 500000.0, "Monthly", "Good", fixedNow, fixedNow)
}

func imgRows() *sqlmock.Rows {
	return sqlmock.NewRows(imgCols).
		AddRow(imgID, propID, "img.jpg", fixedNow, fixedNow)
}

func facRows() *sqlmock.Rows {
	return sqlmock.NewRows(facCols).
		AddRow(facID, "Pool", fixedNow, fixedNow)
}

// mockPreloads mocks the 3 preload queries GORM generates:
// 1. property_facilities join table
// 2. facilities
// 3. property_images
func mockPreloads(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "property_facilities" WHERE`)).
		WillReturnRows(sqlmock.NewRows([]string{"property_id", "facility_id"}).AddRow(propID, facID))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "facilities" WHERE`)).
		WillReturnRows(facRows())
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "property_images" WHERE`)).
		WillReturnRows(imgRows())
}

func mockEmptyPreloads(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "property_facilities" WHERE`)).
		WillReturnRows(sqlmock.NewRows([]string{"property_id", "facility_id"}))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "property_images" WHERE`)).
		WillReturnRows(sqlmock.NewRows(imgCols))
}

// ---------------------------------------------------------------------------
// FindByID
// ---------------------------------------------------------------------------

func TestFindByID(t *testing.T) {
	dbErr := errors.New("connection refused")

	tests := []struct {
		name      string
		id        uuid.UUID
		mockSetup func(sqlmock.Sqlmock)
		wantErr   error
		wantTitle string
	}{
		{
			name: "found",
			id:   propID,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "properties" WHERE id = $1`)).
					WillReturnRows(propRow())
				mockPreloads(mock)
			},
			wantTitle: "Test",
		},
		{
			name: "not found",
			id:   uuid.New(),
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "properties" WHERE id = $1`)).
					WillReturnRows(sqlmock.NewRows(propCols))
			},
			wantErr: apperror.ErrNotFound,
		},
		{
			name: "db error",
			id:   propID,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "properties" WHERE id = $1`)).
					WillReturnError(dbErr)
			},
			wantErr: dbErr,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo, mock, db := newTestRepo(t)
			defer db.Close()
			tc.mockSetup(mock)

			p, err := repo.FindByID(context.Background(), tc.id)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("want error %v, got %v", tc.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if p.Title != tc.wantTitle {
				t.Errorf("Title: want %q, got %q", tc.wantTitle, p.Title)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// FindByDocumentID
// ---------------------------------------------------------------------------

func TestFindByDocumentID(t *testing.T) {
	tests := []struct {
		name       string
		documentID string
		mockSetup  func(sqlmock.Sqlmock)
		wantErr    error
	}{
		{
			name:       "found",
			documentID: "doc-001",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "properties" WHERE document_id = $1`)).
					WillReturnRows(propRow())
				mockPreloads(mock)
			},
		},
		{
			name:       "not found",
			documentID: "missing",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "properties" WHERE document_id = $1`)).
					WillReturnRows(sqlmock.NewRows(propCols))
			},
			wantErr: apperror.ErrNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo, mock, db := newTestRepo(t)
			defer db.Close()
			tc.mockSetup(mock)

			_, err := repo.FindByDocumentID(context.Background(), tc.documentID)
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

// ---------------------------------------------------------------------------
// FindAllWithFilter
// ---------------------------------------------------------------------------

func TestFindAllWithFilter(t *testing.T) {
	tests := []struct {
		name      string
		filter    domainrepo.PropertyFilter
		mockSetup func(sqlmock.Sqlmock)
		wantTotal int64
		wantLen   int
		wantErr   bool
	}{
		{
			name:   "returns results",
			filter: domainrepo.PropertyFilter{Page: 1, PageSize: 10},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "properties"`)).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "properties" LIMIT`)).
					WillReturnRows(propRow())
				mockPreloads(mock)
			},
			wantTotal: 1,
			wantLen:   1,
		},
		{
			name:   "empty result",
			filter: domainrepo.PropertyFilter{Page: 1, PageSize: 10},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "properties"`)).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "properties" LIMIT`)).
					WillReturnRows(sqlmock.NewRows(propCols))
				mockEmptyPreloads(mock)
			},
			wantTotal: 0,
			wantLen:   0,
		},
		{
			name:   "count error",
			filter: domainrepo.PropertyFilter{Page: 1, PageSize: 10},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "properties"`)).
					WillReturnError(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo, mock, db := newTestRepo(t)
			defer db.Close()
			tc.mockSetup(mock)

			result, err := repo.FindAllWithFilter(context.Background(), tc.filter)
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
				t.Errorf("Total: want %d, got %d", tc.wantTotal, result.Total)
			}
			if len(result.Properties) != tc.wantLen {
				t.Errorf("len(Properties): want %d, got %d", tc.wantLen, len(result.Properties))
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Delete
// ---------------------------------------------------------------------------

func TestDelete(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "success",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "properties" WHERE id = $1`)).
					WithArgs(propID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
		},
		{
			name: "db error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "properties" WHERE id = $1`)).
					WithArgs(propID).
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo, mock, db := newTestRepo(t)
			defer db.Close()
			tc.mockSetup(mock)

			err := repo.Delete(context.Background(), propID)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// FindOrCreateFacility
// ---------------------------------------------------------------------------

func TestFindOrCreateFacility(t *testing.T) {
	tests := []struct {
		name      string
		facName   string
		mockSetup func(sqlmock.Sqlmock)
		wantName  string
		wantErr   bool
	}{
		{
			name:    "found existing",
			facName: "Pool",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "facilities" WHERE`)).
					WillReturnRows(facRows())
			},
			wantName: "Pool",
		},
		{
			name:    "db error",
			facName: "Pool",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "facilities" WHERE`)).
					WillReturnError(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo, mock, db := newTestRepo(t)
			defer db.Close()
			tc.mockSetup(mock)

			f, err := repo.FindOrCreateFacility(context.Background(), tc.facName)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if f.Name != tc.wantName {
				t.Errorf("Name: want %q, got %q", tc.wantName, f.Name)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------

func TestCreate(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "success",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "properties"`)).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(propID))
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "property_images"`)).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(imgID))
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "facilities"`)).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(facID))
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "property_facilities"`)).
					WillReturnRows(sqlmock.NewRows([]string{"property_id", "facility_id"}).AddRow(propID, facID))
				mock.ExpectCommit()
			},
		},
		{
			name: "db error on insert",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "properties"`)).
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo, mock, db := newTestRepo(t)
			defer db.Close()
			tc.mockSetup(mock)

			prop := &entity.Property{
				ID:         propID,
				DocumentID: "doc-001",
				Title:      "Test",
				Price:      500000,
				Images:     []entity.PropertyImage{{ID: imgID, PropertyID: propID, URL: "img.jpg"}},
				Facilities: []entity.Facility{{ID: facID, Name: "Pool"}},
				CreatedAt:  fixedNow,
				UpdatedAt:  fixedNow,
			}

			err := repo.Create(context.Background(), prop)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
