package repository

import (
	"context"

	"dose-certa-backend/internal/domain/entity"
)

// MedicationRepository é o contrato de persistência para medicamentos.
// A implementação concreta (MongoDB) fica em internal/repository/mongodb.
type MedicationRepository interface {
	Create(ctx context.Context, m *entity.Medication) (*entity.Medication, error)
	FindAll(ctx context.Context) ([]*entity.Medication, error)
	FindByID(ctx context.Context, id string) (*entity.Medication, error)
	Update(ctx context.Context, id string, m *entity.Medication) (*entity.Medication, error)
	Delete(ctx context.Context, id string) error
}
