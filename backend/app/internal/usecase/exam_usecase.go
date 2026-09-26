package usecase

import (
	"context"
	"fmt"
	"time"

	"dose-certa-backend/internal/domain/entity"
	"dose-certa-backend/internal/domain/repository"
)

// ExamUseCase concentra as regras de negócio dos exames.
type ExamUseCase struct {
	repo repository.ExamRepository
}

func NewExamUseCase(repo repository.ExamRepository) *ExamUseCase {
	return &ExamUseCase{repo: repo}
}

func (uc *ExamUseCase) Create(ctx context.Context, in *entity.Exam) (*entity.Exam, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	in.CreatedAt = now
	in.UpdatedAt = now

	created, err := uc.repo.Create(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar exame: %w", err)
	}
	return created, nil
}

func (uc *ExamUseCase) List(ctx context.Context) ([]*entity.Exam, error) {
	items, err := uc.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("falha ao listar exames: %w", err)
	}
	return items, nil
}

func (uc *ExamUseCase) Get(ctx context.Context, id string) (*entity.Exam, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *ExamUseCase) Update(ctx context.Context, id string, in *entity.Exam) (*entity.Exam, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}
	in.UpdatedAt = time.Now().UTC()
	return uc.repo.Update(ctx, id, in)
}

func (uc *ExamUseCase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}
