package repository

import (
	"context"

	"dose-certa-backend/internal/domain/entity"
)

// WeightRecordRepository é o contrato de persistência dos registros de peso.
type WeightRecordRepository interface {
	Create(ctx context.Context, w *entity.WeightRecord) (*entity.WeightRecord, error)
	FindAll(ctx context.Context) ([]*entity.WeightRecord, error) // ordenado por data desc
	FindByID(ctx context.Context, id string) (*entity.WeightRecord, error)
	Update(ctx context.Context, id string, w *entity.WeightRecord) (*entity.WeightRecord, error)
	Delete(ctx context.Context, id string) error
}
