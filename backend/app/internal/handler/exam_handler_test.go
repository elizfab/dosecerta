package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"dose-certa-backend/internal/domain/entity"
	"dose-certa-backend/internal/domain/repository"
)

type examUCFake struct {
	createFn func(ctx context.Context, e *entity.Exam) (*entity.Exam, error)
	listFn   func(ctx context.Context) ([]*entity.Exam, error)
	getFn    func(ctx context.Context, id string) (*entity.Exam, error)
	updateFn func(ctx context.Context, id string, e *entity.Exam) (*entity.Exam, error)
	deleteFn func(ctx context.Context, id string) error
}

func (f *examUCFake) Create(ctx context.Context, e *entity.Exam) (*entity.Exam, error) {
	return f.createFn(ctx, e)
}
func (f *examUCFake) List(ctx context.Context) ([]*entity.Exam, error) { return f.listFn(ctx) }
func (f *examUCFake) Get(ctx context.Context, id string) (*entity.Exam, error) {
	return f.getFn(ctx, id)
}
func (f *examUCFake) Update(ctx context.Context, id string, e *entity.Exam) (*entity.Exam, error) {
	return f.updateFn(ctx, id, e)
}
func (f *examUCFake) Delete(ctx context.Context, id string) error { return f.deleteFn(ctx, id) }

func examRouter(uc examUseCase) *http.ServeMux {
	mux := http.NewServeMux()
	NewExamHandler(uc).RegisterRoutes(mux)
	return mux
}

func TestExamHandler_Create201(t *testing.T) {
	uc := &examUCFake{
		createFn: func(ctx context.Context, e *entity.Exam) (*entity.Exam, error) {
			e.ID = "1"
			return e, nil
		},
	}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/exams",
		strings.NewReader(`{"name":"Glicose","date":"2026-07-25","value":90,"unit":"mg/dL"}`))
	rec := httptest.NewRecorder()
	examRouter(uc).ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("esperava 201, obtive %d (%s)", rec.Code, rec.Body.String())
	}
}

func TestExamHandler_Create400OnMissingName(t *testing.T) {
	uc := &examUCFake{
		createFn: func(ctx context.Context, e *entity.Exam) (*entity.Exam, error) {
			return nil, entity.ErrExamNameRequired
		},
	}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/exams",
		strings.NewReader(`{"date":"2026-07-25"}`))
	rec := httptest.NewRecorder()
	examRouter(uc).ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("esperava 400, obtive %d", rec.Code)
	}
}

func TestExamHandler_Get404(t *testing.T) {
	uc := &examUCFake{
		getFn: func(ctx context.Context, id string) (*entity.Exam, error) {
			return nil, repository.ErrNotFound
		},
	}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/exams/xyz", nil)
	rec := httptest.NewRecorder()
	examRouter(uc).ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("esperava 404, obtive %d", rec.Code)
	}
}
