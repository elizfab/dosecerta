package entity

import (
	"errors"
	"time"
)

// DoseLog registra que uma dose de um medicamento foi tomada numa data.
// Um registro por (medicationId, date) — a marcação é idempotente.
type DoseLog struct {
	ID           string    `bson:"_id,omitempty" json:"id"`
	MedicationID string    `bson:"medicationId" json:"medicationId"`
	Date         string    `bson:"date" json:"date"` // "YYYY-MM-DD"
	CreatedAt    time.Time `bson:"createdAt" json:"createdAt"`
}

var (
	ErrDoseMedicationRequired = errors.New("medicationId é obrigatório")
	ErrDoseDateRequired       = errors.New("data é obrigatória")
)

// Validate garante as invariantes mínimas do registro de dose.
func (d *DoseLog) Validate() error {
	if d.MedicationID == "" {
		return ErrDoseMedicationRequired
	}
	if d.Date == "" {
		return ErrDoseDateRequired
	}
	return nil
}
