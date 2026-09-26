package repository

import (
	"context"

	"dose-certa-backend/internal/domain/entity"
)

// ExamRepository é o contrato de persistência dos exames.
type ExamRepository interface {
	Create(ctx context.Context, e *entity.Exam) (*entity.Exam, error)
	FindAll(ctx context.Context) ([]*entity.Exam, error) // ordenado por data desc
	FindByID(ctx context.Context, id string) (*entity.Exam, error)
	Update(ctx context.Context, id string, e *entity.Exam) (*entity.Exam, error)
	Delete(ctx context.Context, id string) error
}
