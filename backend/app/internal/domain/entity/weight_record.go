package entity

import (
	"errors"
	"time"
)


// WeightRecord representa uma medição corporal do paciente numa data.
// Além do peso, guarda circunferências opcionais (cm). É apresentado na UI
// como "Medidas corporais".
type WeightRecord struct {
	ID       string  `bson:"_id,omitempty" json:"id"`
	Date     string  `bson:"date" json:"date"` // "YYYY-MM-DD"
	WeightKg float64 `bson:"weightKg" json:"weightKg"`

	// Circunferências em cm (opcionais).
	Waist        float64 `bson:"waist,omitempty" json:"waist,omitempty"`               // cintura/barriga
	Hip          float64 `bson:"hip,omitempty" json:"hip,omitempty"`                   // quadril
	Glute        float64 `bson:"glute,omitempty" json:"glute,omitempty"`               // glúteo
	Chest        float64 `bson:"chest,omitempty" json:"chest,omitempty"`               // peito
	ArmRight     float64 `bson:"armRight,omitempty" json:"armRight,omitempty"`         // braço direito
	ArmLeft      float64 `bson:"armLeft,omitempty" json:"armLeft,omitempty"`           // braço esquerdo
	ForearmRight float64 `bson:"forearmRight,omitempty" json:"forearmRight,omitempty"` // antebraço direito
	ForearmLeft  float64 `bson:"forearmLeft,omitempty" json:"forearmLeft,omitempty"`   // antebraço esquerdo
	ThighRight   float64 `bson:"thighRight,omitempty" json:"thighRight,omitempty"`     // coxa direita
	ThighLeft    float64 `bson:"thighLeft,omitempty" json:"thighLeft,omitempty"`       // coxa esquerda
	CalfRight    float64 `bson:"calfRight,omitempty" json:"calfRight,omitempty"`       // panturrilha direita
	CalfLeft     float64 `bson:"calfLeft,omitempty" json:"calfLeft,omitempty"`         // panturrilha esquerda
	Neck         float64 `bson:"neck,omitempty" json:"neck,omitempty"`                 // pescoço
	Shoulders    float64 `bson:"shoulders,omitempty" json:"shoulders,omitempty"`       // ombros
	BodyFatPct   float64 `bson:"bodyFatPct,omitempty" json:"bodyFatPct,omitempty"`     // % de gordura

	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}

var (
	ErrWeightDateRequired = errors.New("data é obrigatória")
	ErrWeightDateInvalid  = errors.New("data deve estar no formato YYYY-MM-DD")
	ErrWeightValueInvalid = errors.New("peso deve ser maior que zero")
)

// Validate garante as invariantes mínimas da medida.
func (w *WeightRecord) Validate() error {
	if w.Date == "" {
		return ErrWeightDateRequired
	}
	if _, err := time.Parse(dateLayout, w.Date); err != nil {
		return ErrWeightDateInvalid
	}
	if w.WeightKg <= 0 {
		return ErrWeightValueInvalid
	}
	return nil
}
