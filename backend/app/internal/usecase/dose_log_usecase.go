package usecase

import (
	"context"
	"fmt"

	"dose-certa-backend/internal/domain/entity"
	"dose-certa-backend/internal/domain/repository"
)

// DoseLogUseCase concentra as regras de negócio das doses tomadas.
type DoseLogUseCase struct {
	repo repository.DoseLogRepository
}

func NewDoseLogUseCase(repo repository.DoseLogRepository) *DoseLogUseCase {
	return &DoseLogUseCase{repo: repo}
}

// MarkTaken registra que a dose foi tomada (idempotente).
func (uc *DoseLogUseCase) MarkTaken(ctx context.Context, medicationID, date string) (*entity.DoseLog, error) {
	d := &entity.DoseLog{MedicationID: medicationID, Date: date}
	if err := d.Validate(); err != nil {
		return nil, err
	}
	created, err := uc.repo.MarkTaken(ctx, medicationID, date)
	if err != nil {
		return nil, fmt.Errorf("falha ao marcar dose: %w", err)
	}
	return created, nil
}

// Unmark desfaz a marcação (idempotente).
func (uc *DoseLogUseCase) Unmark(ctx context.Context, medicationID, date string) error {
	d := &entity.DoseLog{MedicationID: medicationID, Date: date}
	if err := d.Validate(); err != nil {
		return err
	}
	if err := uc.repo.Unmark(ctx, medicationID, date); err != nil {
		return fmt.Errorf("falha ao desmarcar dose: %w", err)
	}
	return nil
}

// ListByDate retorna as doses tomadas numa data.
func (uc *DoseLogUseCase) ListByDate(ctx context.Context, date string) ([]*entity.DoseLog, error) {
	if date == "" {
		return nil, entity.ErrDoseDateRequired
	}
	items, err := uc.repo.FindByDate(ctx, date)
	if err != nil {
		return nil, fmt.Errorf("falha ao listar doses: %w", err)
	}
	return items, nil
}
