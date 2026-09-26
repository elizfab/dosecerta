package usecase

import (
	"context"
	"testing"

	"dose-certa-backend/internal/domain/entity"
)

type weightRepoMock struct {
	createFn   func(ctx context.Context, w *entity.WeightRecord) (*entity.WeightRecord, error)
	findAllFn  func(ctx context.Context) ([]*entity.WeightRecord, error)
	findByIDFn func(ctx context.Context, id string) (*entity.WeightRecord, error)
	updateFn   func(ctx context.Context, id string, w *entity.WeightRecord) (*entity.WeightRecord, error)
	deleteFn   func(ctx context.Context, id string) error
}

func (r *weightRepoMock) Create(ctx context.Context, w *entity.WeightRecord) (*entity.WeightRecord, error) {
	return r.createFn(ctx, w)
}
func (r *weightRepoMock) FindAll(ctx context.Context) ([]*entity.WeightRecord, error) {
	return r.findAllFn(ctx)
}
func (r *weightRepoMock) FindByID(ctx context.Context, id string) (*entity.WeightRecord, error) {
	return r.findByIDFn(ctx, id)
}
func (r *weightRepoMock) Update(ctx context.Context, id string, w *entity.WeightRecord) (*entity.WeightRecord, error) {
	return r.updateFn(ctx, id, w)
}
func (r *weightRepoMock) Delete(ctx context.Context, id string) error {
	return r.deleteFn(ctx, id)
}

func TestWeightCreate_Success(t *testing.T) {
	var saved *entity.WeightRecord
	repo := &weightRepoMock{
		createFn: func(ctx context.Context, w *entity.WeightRecord) (*entity.WeightRecord, error) {
			w.ID = "1"
			saved = w
			return w, nil
		},
	}
	uc := NewWeightUseCase(repo)

	out, err := uc.Create(context.Background(), &entity.WeightRecord{Date: "2026-07-25", WeightKg: 70.5})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if out.ID != "1" {
		t.Fatalf("esperava id '1', obtive %q", out.ID)
	}
	if saved.CreatedAt.IsZero() {
		t.Fatalf("Create deveria preencher CreatedAt")
	}
}

func TestWeightCreate_InvalidValue(t *testing.T) {
	repo := &weightRepoMock{
		createFn: func(ctx context.Context, w *entity.WeightRecord) (*entity.WeightRecord, error) {
			t.Fatal("repo.Create não deveria ser chamado")
			return nil, nil
		},
	}
	uc := NewWeightUseCase(repo)

	_, err := uc.Create(context.Background(), &entity.WeightRecord{Date: "2026-07-25", WeightKg: -1})
	if err != entity.ErrWeightValueInvalid {
		t.Fatalf("esperava ErrWeightValueInvalid, obtive %v", err)
	}
}
