package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"dose-certa-backend/internal/domain/entity"
	"dose-certa-backend/internal/domain/repository"
)

// fake do use case de medicamentos.
type medUCFake struct {
	createFn func(ctx context.Context, m *entity.Medication) (*entity.Medication, error)
	listFn   func(ctx context.Context) ([]*entity.Medication, error)
	getFn    func(ctx context.Context, id string) (*entity.Medication, error)
	updateFn func(ctx context.Context, id string, m *entity.Medication) (*entity.Medication, error)
	deleteFn func(ctx context.Context, id string) error
}

func (f *medUCFake) Create(ctx context.Context, m *entity.Medication) (*entity.Medication, error) {
	return f.createFn(ctx, m)
}
func (f *medUCFake) List(ctx context.Context) ([]*entity.Medication, error) { return f.listFn(ctx) }
func (f *medUCFake) Get(ctx context.Context, id string) (*entity.Medication, error) {
	return f.getFn(ctx, id)
}
func (f *medUCFake) Update(ctx context.Context, id string, m *entity.Medication) (*entity.Medication, error) {
	return f.updateFn(ctx, id, m)
}
func (f *medUCFake) Delete(ctx context.Context, id string) error { return f.deleteFn(ctx, id) }

func newRouter(h *MedicationHandler) *http.ServeMux {
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	return mux
}

func TestMedicationHandler_Create201(t *testing.T) {
	uc := &medUCFake{
		createFn: func(ctx context.Context, m *entity.Medication) (*entity.Medication, error) {
			m.ID = "1"
			return m, nil
		},
	}
	mux := newRouter(NewMedicationHandler(uc))

	body := `{"name":"Omeprazol","dosage":"20 mg","schedule":"08:00"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/medications", strings.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("esperava 201, obtive %d", rec.Code)
	}
	var resp struct {
		Success bool `json:"success"`
		Data    struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("json inválido: %v", err)
	}
	if !resp.Success || resp.Data.ID != "1" {
		t.Fatalf("resposta inesperada: %s", rec.Body.String())
	}
}

func TestMedicationHandler_Get404(t *testing.T) {
	uc := &medUCFake{
		getFn: func(ctx context.Context, id string) (*entity.Medication, error) {
			return nil, repository.ErrNotFound
		},
	}
	mux := newRouter(NewMedicationHandler(uc))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/medications/xyz", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("esperava 404, obtive %d", rec.Code)
	}
}

func TestMedicationHandler_Create400OnValidation(t *testing.T) {
	uc := &medUCFake{
		createFn: func(ctx context.Context, m *entity.Medication) (*entity.Medication, error) {
			return nil, entity.ErrMedicationNameRequired
		},
	}
	mux := newRouter(NewMedicationHandler(uc))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/medications", strings.NewReader(`{"dosage":"20 mg","schedule":"08:00"}`))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("esperava 400, obtive %d", rec.Code)
	}
}
