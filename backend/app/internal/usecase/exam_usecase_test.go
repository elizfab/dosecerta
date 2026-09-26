package usecase

import (
	"context"
	"testing"

	"dose-certa-backend/internal/domain/entity"
)

type examRepoMock struct {
	createFn   func(ctx context.Context, e *entity.Exam) (*entity.Exam, error)
	findAllFn  func(ctx context.Context) ([]*entity.Exam, error)
	findByIDFn func(ctx context.Context, id string) (*entity.Exam, error)
	updateFn   func(ctx context.Context, id string, e *entity.Exam) (*entity.Exam, error)
	deleteFn   func(ctx context.Context, id string) error
}

func (r *examRepoMock) Create(ctx context.Context, e *entity.Exam) (*entity.Exam, error) {
	return r.createFn(ctx, e)
}
func (r *examRepoMock) FindAll(ctx context.Context) ([]*entity.Exam, error) { return r.findAllFn(ctx) }
func (r *examRepoMock) FindByID(ctx context.Context, id string) (*entity.Exam, error) {
	return r.findByIDFn(ctx, id)
}
func (r *examRepoMock) Update(ctx context.Context, id string, e *entity.Exam) (*entity.Exam, error) {
	return r.updateFn(ctx, id, e)
}
func (r *examRepoMock) Delete(ctx context.Context, id string) error { return r.deleteFn(ctx, id) }

func TestExamCreate_Success(t *testing.T) {
	var saved *entity.Exam
	repo := &examRepoMock{
		createFn: func(ctx context.Context, e *entity.Exam) (*entity.Exam, error) {
			e.ID = "1"
			saved = e
			return e, nil
		},
	}
	uc := NewExamUseCase(repo)

	out, err := uc.Create(context.Background(), &entity.Exam{Name: "Glicose", Date: "2026-07-25", Value: 90})
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

func TestExamCreate_MissingName(t *testing.T) {
	repo := &examRepoMock{
		createFn: func(ctx context.Context, e *entity.Exam) (*entity.Exam, error) {
			t.Fatal("repo.Create não deveria ser chamado")
			return nil, nil
		},
	}
	uc := NewExamUseCase(repo)

	_, err := uc.Create(context.Background(), &entity.Exam{Date: "2026-07-25"})
	if err != entity.ErrExamNameRequired {
		t.Fatalf("esperava ErrExamNameRequired, obtive %v", err)
	}
}
