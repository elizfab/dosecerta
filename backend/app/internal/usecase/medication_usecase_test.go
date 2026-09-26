package usecase

import (
	"context"
	"errors"
	"testing"

	"dose-certa-backend/internal/domain/entity"
	"dose-certa-backend/internal/domain/repository"
)

// mock em memória que satisfaz repository.MedicationRepository.
type medRepoMock struct {
	createFn   func(ctx context.Context, m *entity.Medication) (*entity.Medication, error)
	findAllFn  func(ctx context.Context) ([]*entity.Medication, error)
	findByIDFn func(ctx context.Context, id string) (*entity.Medication, error)
	updateFn   func(ctx context.Context, id string, m *entity.Medication) (*entity.Medication, error)
	deleteFn   func(ctx context.Context, id string) error
}

func (r *medRepoMock) Create(ctx context.Context, m *entity.Medication) (*entity.Medication, error) {
	return r.createFn(ctx, m)
}
func (r *medRepoMock) FindAll(ctx context.Context) ([]*entity.Medication, error) {
	return r.findAllFn(ctx)
}
func (r *medRepoMock) FindByID(ctx context.Context, id string) (*entity.Medication, error) {
	return r.findByIDFn(ctx, id)
}
func (r *medRepoMock) Update(ctx context.Context, id string, m *entity.Medication) (*entity.Medication, error) {
	return r.updateFn(ctx, id, m)
}
func (r *medRepoMock) Delete(ctx context.Context, id string) error {
	return r.deleteFn(ctx, id)
}

func TestMedicationCreate_Success(t *testing.T) {
	var saved *entity.Medication
	repo := &medRepoMock{
		createFn: func(ctx context.Context, m *entity.Medication) (*entity.Medication, error) {
			m.ID = "abc"
			saved = m
			return m, nil
		},
	}
	uc := NewMedicationUseCase(repo)

	out, err := uc.Create(context.Background(), &entity.Medication{
		Name: "Omeprazol", Dosage: "20 mg", Schedule: "08:00",
	})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if out.ID != "abc" {
		t.Fatalf("esperava id 'abc', obtive %q", out.ID)
	}
	if !saved.Active {
		t.Fatalf("Create deveria marcar Active=true")
	}
	if saved.CreatedAt.IsZero() || saved.UpdatedAt.IsZero() {
		t.Fatalf("Create deveria preencher timestamps")
	}
}

func TestMedicationCreate_ValidationError(t *testing.T) {
	repo := &medRepoMock{
		createFn: func(ctx context.Context, m *entity.Medication) (*entity.Medication, error) {
			t.Fatal("repo.Create não deveria ser chamado quando a validação falha")
			return nil, nil
		},
	}
	uc := NewMedicationUseCase(repo)

	_, err := uc.Create(context.Background(), &entity.Medication{Dosage: "20 mg", Schedule: "08:00"})
	if err != entity.ErrMedicationNameRequired {
		t.Fatalf("esperava ErrMedicationNameRequired, obtive %v", err)
	}
}

func TestMedicationGet_NotFound(t *testing.T) {
	repo := &medRepoMock{
		findByIDFn: func(ctx context.Context, id string) (*entity.Medication, error) {
			return nil, repository.ErrNotFound
		},
	}
	uc := NewMedicationUseCase(repo)

	_, err := uc.Get(context.Background(), "missing")
	if !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("esperava ErrNotFound, obtive %v", err)
	}
}

func TestMedicationList_Success(t *testing.T) {
	repo := &medRepoMock{
		findAllFn: func(ctx context.Context) ([]*entity.Medication, error) {
			return []*entity.Medication{{ID: "1"}, {ID: "2"}}, nil
		},
	}
	uc := NewMedicationUseCase(repo)

	out, err := uc.List(context.Background())
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("esperava 2 itens, obtive %d", len(out))
	}
}
