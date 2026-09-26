package entity

import "testing"

func TestWeightValidate_Success(t *testing.T) {
	w := &WeightRecord{Date: "2026-07-25", WeightKg: 70.5}
	if err := w.Validate(); err != nil {
		t.Fatalf("esperava sucesso, obtive %v", err)
	}
}

func TestWeightValidate_MissingDate(t *testing.T) {
	w := &WeightRecord{WeightKg: 70.5}
	if err := w.Validate(); err != ErrWeightDateRequired {
		t.Fatalf("esperava ErrWeightDateRequired, obtive %v", err)
	}
}

func TestWeightValidate_InvalidValue(t *testing.T) {
	w := &WeightRecord{Date: "2026-07-25", WeightKg: 0}
	if err := w.Validate(); err != ErrWeightValueInvalid {
		t.Fatalf("esperava ErrWeightValueInvalid, obtive %v", err)
	}
}
