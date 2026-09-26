package usecase

import (
	"context"
	"fmt"
	"time"

	"dose-certa-backend/internal/domain/entity"
	"dose-certa-backend/internal/domain/repository"
)

// MedicationUseCase concentra as regras de negócio dos medicamentos.
// Depende da interface do repositório (não da implementação).
type MedicationUseCase struct {
	repo repository.MedicationRepository
}

func NewMedicationUseCase(repo repository.MedicationRepository) *MedicationUseCase {
	return &MedicationUseCase{repo: repo}
}

// Create valida e persiste um novo medicamento.
func (uc *MedicationUseCase) Create(ctx context.Context, in *entity.Medication) (*entity.Medication, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	in.Active = true
	in.CreatedAt = now
	in.UpdatedAt = now

	created, err := uc.repo.Create(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar medicamento: %w", err)
	}
	return created, nil
}

// List retorna todos os medicamentos.
func (uc *MedicationUseCase) List(ctx context.Context) ([]*entity.Medication, error) {
	items, err := uc.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("falha ao listar medicamentos: %w", err)
	}
	return items, nil
}

// Get busca um medicamento por id.
func (uc *MedicationUseCase) Get(ctx context.Context, id string) (*entity.Medication, error) {
	item, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// Update valida e atualiza um medicamento existente.
func (uc *MedicationUseCase) Update(ctx context.Context, id string, in *entity.Medication) (*entity.Medication, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}
	in.UpdatedAt = time.Now().UTC()

	updated, err := uc.repo.Update(ctx, id, in)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

// Delete remove um medicamento por id.
func (uc *MedicationUseCase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}
