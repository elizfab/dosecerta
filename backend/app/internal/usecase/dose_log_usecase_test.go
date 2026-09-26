package usecase

import (
	"context"
	"errors"
	"testing"

	"dose-certa-backend/internal/domain/entity"
)

// mock em memória que satisfaz repository.DoseLogRepository.
type doseRepoMock struct {
	markFn       func(ctx context.Context, medID, date string) (*entity.DoseLog, error)
	unmarkFn     func(ctx context.Context, medID, date string) error
	findByDateFn func(ctx context.Context, date string) ([]*entity.DoseLog, error)
}

func (r *doseRepoMock) MarkTaken(ctx context.Context, medID, date string) (*entity.DoseLog, error) {
	return r.markFn(ctx, medID, date)
}
func (r *doseRepoMock) Unmark(ctx context.Context, medID, date string) error {
	return r.unmarkFn(ctx, medID, date)
}
func (r *doseRepoMock) FindByDate(ctx context.Context, date string) ([]*entity.DoseLog, error) {
	return r.findByDateFn(ctx, date)
}

func TestDoseMarkTaken_Success(t *testing.T) {
	repo := &doseRepoMock{
		markFn: func(ctx context.Context, medID, date string) (*entity.DoseLog, error) {
			return &entity.DoseLog{ID: "1", MedicationID: medID, Date: date}, nil
		},
	}
	uc := NewDoseLogUseCase(repo)

	out, err := uc.MarkTaken(context.Background(), "med1", "2026-07-28")
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if out.MedicationID != "med1" || out.Date != "2026-07-28" {
		t.Fatalf("registro inesperado: %+v", out)
	}
}

func TestDoseMarkTaken_ValidationError(t *testing.T) {
	repo := &doseRepoMock{
		markFn: func(ctx context.Context, medID, date string) (*entity.DoseLog, error) {
			t.Fatal("repo não deveria ser chamado com validação falha")
			return nil, nil
		},
	}
	uc := NewDoseLogUseCase(repo)

	if _, err := uc.MarkTaken(context.Background(), "", "2026-07-28"); !errors.Is(err, entity.ErrDoseMedicationRequired) {
		t.Fatalf("esperava ErrDoseMedicationRequired, obtive %v", err)
	}
	if _, err := uc.MarkTaken(context.Background(), "med1", ""); !errors.Is(err, entity.ErrDoseDateRequired) {
		t.Fatalf("esperava ErrDoseDateRequired, obtive %v", err)
	}
}

func TestDoseListByDate_RequiresDate(t *testing.T) {
	uc := NewDoseLogUseCase(&doseRepoMock{})
	if _, err := uc.ListByDate(context.Background(), ""); !errors.Is(err, entity.ErrDoseDateRequired) {
		t.Fatalf("esperava ErrDoseDateRequired, obtive %v", err)
	}
}

func TestDoseUnmark_Success(t *testing.T) {
	called := false
	repo := &doseRepoMock{
		unmarkFn: func(ctx context.Context, medID, date string) error {
			called = true
			return nil
		},
	}
	uc := NewDoseLogUseCase(repo)
	if err := uc.Unmark(context.Background(), "med1", "2026-07-28"); err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if !called {
		t.Fatal("repo.Unmark deveria ter sido chamado")
	}
}
