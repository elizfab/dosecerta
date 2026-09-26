package entity

import (
	"errors"
	"time"
)

// Exam representa o resultado de um exame laboratorial numa data.
type Exam struct {
	ID           string    `bson:"_id,omitempty" json:"id"`
	Name         string    `bson:"name" json:"name"`
	Value        float64   `bson:"value,omitempty" json:"value,omitempty"`
	Unit         string    `bson:"unit,omitempty" json:"unit,omitempty"`
	ReferenceMin float64   `bson:"referenceMin,omitempty" json:"referenceMin,omitempty"`
	ReferenceMax float64   `bson:"referenceMax,omitempty" json:"referenceMax,omitempty"`
	Date         string    `bson:"date" json:"date"` // "YYYY-MM-DD"
	Notes        string    `bson:"notes,omitempty" json:"notes,omitempty"`
	CreatedAt    time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time `bson:"updatedAt" json:"updatedAt"`
}

var (
	ErrExamNameRequired = errors.New("nome do exame é obrigatório")
	ErrExamDateRequired = errors.New("data do exame é obrigatória")
	ErrExamDateInvalid  = errors.New("data do exame deve estar no formato YYYY-MM-DD")
)

func (e *Exam) Validate() error {
	if e.Name == "" {
		return ErrExamNameRequired
	}
	if e.Date == "" {
		return ErrExamDateRequired
	}
	if _, err := time.Parse(dateLayout, e.Date); err != nil {
		return ErrExamDateInvalid
	}
	return nil
}
