package repository

import (
	"context"

	"dose-certa-backend/internal/domain/entity"
)

// DoseLogRepository é o contrato de persistência para doses tomadas.
type DoseLogRepository interface {
	// MarkTaken cria (ou mantém) o registro de dose tomada — idempotente por (medicationId, date).
	MarkTaken(ctx context.Context, medicationID, date string) (*entity.DoseLog, error)
	// Unmark remove o registro. Retorna nil mesmo se não existir (idempotente).
	Unmark(ctx context.Context, medicationID, date string) error
	// FindByDate lista as doses tomadas numa data.
	FindByDate(ctx context.Context, date string) ([]*entity.DoseLog, error)
}
