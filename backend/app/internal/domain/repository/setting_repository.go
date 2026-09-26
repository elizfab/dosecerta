package repository

import (
	"context"

	"dose-certa-backend/internal/domain/entity"
)

// SettingRepository é o contrato de persistência para configurações (key/value).
type SettingRepository interface {
	// Get retorna a configuração pela chave. Retorna ErrNotFound se não existir.
	Get(ctx context.Context, key string) (*entity.Setting, error)
	// Set cria ou atualiza (upsert) a configuração.
	Set(ctx context.Context, s *entity.Setting) (*entity.Setting, error)
}
