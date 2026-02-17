package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/taufiksoleh/backend-mp/internal/delivery/http/handler"
	"github.com/taufiksoleh/backend-mp/internal/domain/apperror"
	"github.com/taufiksoleh/backend-mp/internal/domain/entity"
	"github.com/taufiksoleh/backend-mp/internal/domain/repository"
	"github.com/taufiksoleh/backend-mp/internal/testdata/fixture"
	mockusecase "github.com/taufiksoleh/backend-mp/internal/usecase/mock"
	"go.uber.org/mock/gomock"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupRouter(h *handler.PropertyHandler) *gin.Engine {
	r := gin.New()
	v1 := r.Group("/api/v1")
	props := v1.Group("/properties")
	props.GET("", h.GetAll)
	props.GET("/doc/:documentId", h.GetByDocumentID)
	props.GET("/:id", h.GetByID)
	props.POST("", h.Create)
	props.PUT("/:id", h.Update)
	props.DELETE("/:id", h.Delete)
	return r
}

func newMockHandler(t *testing.T) (*handler.PropertyHandler, *mockusecase.MockPropertyUsecase) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mock := mockusecase.NewMockPropertyUsecase(ctrl)
	return handler.NewPropertyHandler(mock), mock
}

// ---------------------------------------------------------------------------
// GET /api/v1/properties
// ---------------------------------------------------------------------------

func TestGetAll(t *testing.T) {
	prop := fixture.Property()

	tests := []struct {
		name       string
		url        string
		mockSetup  func(*mockusecase.MockPropertyUsecase)
		wantCode   int
		assertBody func(*testing.T, []byte)
	}{
		{
			name: "200 default params",
			url:  "/api/v1/properties",
			mockSetup: func(m *mockusecase.MockPropertyUsecase) {
				m.EXPECT().
					GetAll(gomock.Any(), repository.PropertyFilter{Keyword: "", Page: 1, PageSize: 10}).
					Return(repository.PaginatedProperties{Properties: []entity.Property{*prop}, Total: 1}, nil)
			},
			wantCode: http.StatusOK,
			assertBody: func(t *testing.T, b []byte) {
				var body map[string]any
				if err := json.Unmarshal(b, &body); err != nil {
					t.Fatalf("invalid JSON: %v", err)
				}
				listings := body["data"].([]any)
				if len(listings) != 1 {
					t.Errorf("expected 1 listing, got %d", len(listings))
				}
			},
		},
		{
			name: "200 keyword and pageSize",
			url:  "/api/v1/properties?keyword=villa&pageSize=5",
			mockSetup: func(m *mockusecase.MockPropertyUsecase) {
				m.EXPECT().
					GetAll(gomock.Any(), repository.PropertyFilter{Keyword: "villa", Page: 1, PageSize: 5}).
					Return(repository.PaginatedProperties{Properties: []entity.Property{}, Total: 0}, nil)
			},
			wantCode: http.StatusOK,
		},
		{
			name: "200 pageSize capped at 100",
			url:  "/api/v1/properties?pageSize=999",
			mockSetup: func(m *mockusecase.MockPropertyUsecase) {
				m.EXPECT().
					GetAll(gomock.Any(), repository.PropertyFilter{Keyword: "", Page: 1, PageSize: 100}).
					Return(repository.PaginatedProperties{}, nil)
			},
			wantCode: http.StatusOK,
		},
		{
			name: "500 on repo error",
			url:  "/api/v1/properties",
			mockSetup: func(m *mockusecase.MockPropertyUsecase) {
				m.EXPECT().
					GetAll(gomock.Any(), gomock.Any()).
					Return(repository.PaginatedProperties{}, fmt.Errorf("db failure"))
			},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h, mock := newMockHandler(t)
			tc.mockSetup(mock)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, tc.url, nil)
			setupRouter(h).ServeHTTP(w, req)

			if w.Code != tc.wantCode {
				t.Fatalf("want %d, got %d — body: %s", tc.wantCode, w.Code, w.Body.String())
			}
			if tc.assertBody != nil {
				tc.assertBody(t, w.Body.Bytes())
			}
		})
	}
}

// ---------------------------------------------------------------------------
// GET /api/v1/properties/:id
// ---------------------------------------------------------------------------

func TestGetByID(t *testing.T) {
	prop := fixture.Property()

	tests := []struct {
		name      string
		urlID     string
		mockSetup func(*mockusecase.MockPropertyUsecase)
		wantCode  int
	}{
		{
			name:  "200 found",
			urlID: prop.ID.String(),
			mockSetup: func(m *mockusecase.MockPropertyUsecase) {
				m.EXPECT().GetByID(gomock.Any(), prop.ID).Return(prop, nil)
			},
			wantCode: http.StatusOK,
		},
		{
			name:  "404 not found",
			urlID: uuid.New().String(),
			mockSetup: func(m *mockusecase.MockPropertyUsecase) {
				m.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(nil, apperror.ErrNotFound)
			},
			wantCode: http.StatusNotFound,
		},
		{
			name:      "400 invalid UUID",
			urlID:     "not-a-uuid",
			mockSetup: func(m *mockusecase.MockPropertyUsecase) {},
			wantCode:  http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h, mock := newMockHandler(t)
			tc.mockSetup(mock)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/properties/"+tc.urlID, nil)
			setupRouter(h).ServeHTTP(w, req)

			if w.Code != tc.wantCode {
				t.Fatalf("want %d, got %d — body: %s", tc.wantCode, w.Code, w.Body.String())
			}
		})
	}
}

// ---------------------------------------------------------------------------
// POST /api/v1/properties
// ---------------------------------------------------------------------------

func TestCreate(t *testing.T) {
	tests := []struct {
		name      string
		body      string
		mockSetup func(*mockusecase.MockPropertyUsecase)
		wantCode  int
	}{
		{
			name: "201 created",
			body: `{"title":"Test","price":100000}`,
			mockSetup: func(m *mockusecase.MockPropertyUsecase) {
				m.EXPECT().Create(gomock.Any(), gomock.Any()).Return(fixture.Property(), nil)
			},
			wantCode: http.StatusCreated,
		},
		{
			name:      "400 missing title",
			body:      `{"price":100000}`,
			mockSetup: func(m *mockusecase.MockPropertyUsecase) {},
			wantCode:  http.StatusBadRequest,
		},
		{
			name:      "400 missing price",
			body:      `{"title":"Test"}`,
			mockSetup: func(m *mockusecase.MockPropertyUsecase) {},
			wantCode:  http.StatusBadRequest,
		},
		{
			name:      "400 invalid JSON",
			body:      `{invalid}`,
			mockSetup: func(m *mockusecase.MockPropertyUsecase) {},
			wantCode:  http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h, mock := newMockHandler(t)
			tc.mockSetup(mock)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/properties", bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")
			setupRouter(h).ServeHTTP(w, req)

			if w.Code != tc.wantCode {
				t.Fatalf("want %d, got %d — body: %s", tc.wantCode, w.Code, w.Body.String())
			}
		})
	}
}

// ---------------------------------------------------------------------------
// DELETE /api/v1/properties/:id
// ---------------------------------------------------------------------------

func TestDelete(t *testing.T) {
	tests := []struct {
		name      string
		urlID     string
		mockSetup func(*mockusecase.MockPropertyUsecase)
		wantCode  int
	}{
		{
			name:  "204 deleted",
			urlID: fixture.PropertyID.String(),
			mockSetup: func(m *mockusecase.MockPropertyUsecase) {
				m.EXPECT().Delete(gomock.Any(), fixture.PropertyID).Return(nil)
			},
			wantCode: http.StatusNoContent,
		},
		{
			name:  "404 not found",
			urlID: uuid.New().String(),
			mockSetup: func(m *mockusecase.MockPropertyUsecase) {
				m.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(apperror.ErrNotFound)
			},
			wantCode: http.StatusNotFound,
		},
		{
			name:      "400 invalid UUID",
			urlID:     "bad-id",
			mockSetup: func(m *mockusecase.MockPropertyUsecase) {},
			wantCode:  http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h, mock := newMockHandler(t)
			tc.mockSetup(mock)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodDelete, "/api/v1/properties/"+tc.urlID, nil)
			setupRouter(h).ServeHTTP(w, req)

			if w.Code != tc.wantCode {
				t.Fatalf("want %d, got %d — body: %s", tc.wantCode, w.Code, w.Body.String())
			}
		})
	}
}
