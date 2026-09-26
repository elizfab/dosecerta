package usecase

import (
	"context"
	"fmt"
	"time"

	"dose-certa-backend/internal/domain/entity"
	"dose-certa-backend/internal/domain/repository"
)

// WeightUseCase concentra as regras de negócio dos registros de peso.
type WeightUseCase struct {
	repo repository.WeightRecordRepository
}

func NewWeightUseCase(repo repository.WeightRecordRepository) *WeightUseCase {
	return &WeightUseCase{repo: repo}
}

func (uc *WeightUseCase) Create(ctx context.Context, in *entity.WeightRecord) (*entity.WeightRecord, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}
	in.CreatedAt = time.Now().UTC()

	created, err := uc.repo.Create(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar registro de peso: %w", err)
	}
	return created, nil
}

func (uc *WeightUseCase) List(ctx context.Context) ([]*entity.WeightRecord, error) {
	items, err := uc.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("falha ao listar registros de peso: %w", err)
	}
	return items, nil
}

func (uc *WeightUseCase) Get(ctx context.Context, id string) (*entity.WeightRecord, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *WeightUseCase) Update(ctx context.Context, id string, in *entity.WeightRecord) (*entity.WeightRecord, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}
	return uc.repo.Update(ctx, id, in)
}

func (uc *WeightUseCase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}
