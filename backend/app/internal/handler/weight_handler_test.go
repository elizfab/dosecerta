package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"dose-certa-backend/internal/domain/entity"
)

type weightUCFake struct {
	createFn func(ctx context.Context, w *entity.WeightRecord) (*entity.WeightRecord, error)
	listFn   func(ctx context.Context) ([]*entity.WeightRecord, error)
	getFn    func(ctx context.Context, id string) (*entity.WeightRecord, error)
	updateFn func(ctx context.Context, id string, w *entity.WeightRecord) (*entity.WeightRecord, error)
	deleteFn func(ctx context.Context, id string) error
}

func (f *weightUCFake) Create(ctx context.Context, w *entity.WeightRecord) (*entity.WeightRecord, error) {
	return f.createFn(ctx, w)
}
func (f *weightUCFake) List(ctx context.Context) ([]*entity.WeightRecord, error) {
	return f.listFn(ctx)
}
func (f *weightUCFake) Get(ctx context.Context, id string) (*entity.WeightRecord, error) {
	return f.getFn(ctx, id)
}
func (f *weightUCFake) Update(ctx context.Context, id string, w *entity.WeightRecord) (*entity.WeightRecord, error) {
	return f.updateFn(ctx, id, w)
}
func (f *weightUCFake) Delete(ctx context.Context, id string) error { return f.deleteFn(ctx, id) }

func TestWeightHandler_Create201(t *testing.T) {
	uc := &weightUCFake{
		createFn: func(ctx context.Context, w *entity.WeightRecord) (*entity.WeightRecord, error) {
			w.ID = "1"
			return w, nil
		},
	}
	mux := http.NewServeMux()
	NewWeightHandler(uc).RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/weight-records", strings.NewReader(`{"date":"2026-07-25","weightKg":70.5}`))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("esperava 201, obtive %d (%s)", rec.Code, rec.Body.String())
	}
}

func TestWeightHandler_Create400OnInvalid(t *testing.T) {
	uc := &weightUCFake{
		createFn: func(ctx context.Context, w *entity.WeightRecord) (*entity.WeightRecord, error) {
			return nil, entity.ErrWeightValueInvalid
		},
	}
	mux := http.NewServeMux()
	NewWeightHandler(uc).RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/weight-records", strings.NewReader(`{"date":"2026-07-25","weightKg":0}`))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("esperava 400, obtive %d", rec.Code)
	}
}
