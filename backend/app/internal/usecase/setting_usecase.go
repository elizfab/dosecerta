package usecase

import (
	"context"
	"fmt"

	"dose-certa-backend/internal/domain/entity"
	"dose-certa-backend/internal/domain/repository"
)

// SettingUseCase concentra as regras das configurações (key/value).
type SettingUseCase struct {
	repo repository.SettingRepository
}

func NewSettingUseCase(repo repository.SettingRepository) *SettingUseCase {
	return &SettingUseCase{repo: repo}
}

// Get busca uma configuração pela chave.
func (uc *SettingUseCase) Get(ctx context.Context, key string) (*entity.Setting, error) {
	if key == "" {
		return nil, entity.ErrSettingKeyRequired
	}
	return uc.repo.Get(ctx, key)
}

// Set cria/atualiza uma configuração.
func (uc *SettingUseCase) Set(ctx context.Context, key, value string) (*entity.Setting, error) {
	s := &entity.Setting{Key: key, Value: value}
	if err := s.Validate(); err != nil {
		return nil, err
	}
	saved, err := uc.repo.Set(ctx, s)
	if err != nil {
		return nil, fmt.Errorf("falha ao salvar configuração: %w", err)
	}
	return saved, nil
}
