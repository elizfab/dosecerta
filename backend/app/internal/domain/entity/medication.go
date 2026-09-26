package entity

import (
	"errors"
	"time"
)

const dateLayout = "2006-01-02"
const timeLayout = "15:04"

// Medication representa um medicamento cadastrado pelo paciente.
type Medication struct {
	ID       string `bson:"_id,omitempty" json:"id"`
	Name     string `bson:"name" json:"name"`
	Dosage   string `bson:"dosage" json:"dosage"`     // ex: "20 mg"
	Schedule string `bson:"schedule" json:"schedule"` // "HH:mm"
	Notes    string `bson:"notes,omitempty" json:"notes,omitempty"`
	Active   bool   `bson:"active" json:"active"`

	// Campos opcionais adicionais.
	Period    string `bson:"period,omitempty" json:"period,omitempty"`       // manhã/tarde/noite/madrugada (derivável do horário)
	StartDate string `bson:"startDate,omitempty" json:"startDate,omitempty"` // "YYYY-MM-DD"
	EndDate   string `bson:"endDate,omitempty" json:"endDate,omitempty"`     // "YYYY-MM-DD"
	Stock     int    `bson:"stock,omitempty" json:"stock,omitempty"`         // quantidade em estoque

	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

// Erros de validação do domínio.
var (
	ErrMedicationNameRequired      = errors.New("nome do medicamento é obrigatório")
	ErrMedicationDosageRequired    = errors.New("dose/posologia é obrigatória")
	ErrMedicationScheduleRequired  = errors.New("horário é obrigatório")
	ErrMedicationScheduleInvalid   = errors.New("horário deve estar no formato HH:mm")
)

// Validate garante as invariantes mínimas do medicamento.
func (m *Medication) Validate() error {
	if m.Name == "" {
		return ErrMedicationNameRequired
	}
	if m.Dosage == "" {
		return ErrMedicationDosageRequired
	}
	if m.Schedule == "" {
		return ErrMedicationScheduleRequired
	}
	if _, err := time.Parse(timeLayout, m.Schedule); err != nil {
		return ErrMedicationScheduleInvalid
	}
	return nil
}
