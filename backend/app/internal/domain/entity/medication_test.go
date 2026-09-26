package entity

import "testing"

func TestMedicationValidate_Success(t *testing.T) {
	m := &Medication{Name: "Omeprazol", Dosage: "20 mg", Schedule: "08:00"}
	if err := m.Validate(); err != nil {
		t.Fatalf("esperava sucesso, obtive erro: %v", err)
	}
}

func TestMedicationValidate_MissingName(t *testing.T) {
	m := &Medication{Dosage: "20 mg", Schedule: "08:00"}
	if err := m.Validate(); err != ErrMedicationNameRequired {
		t.Fatalf("esperava ErrMedicationNameRequired, obtive: %v", err)
	}
}

func TestMedicationValidate_MissingDosage(t *testing.T) {
	m := &Medication{Name: "Omeprazol", Schedule: "08:00"}
	if err := m.Validate(); err != ErrMedicationDosageRequired {
		t.Fatalf("esperava ErrMedicationDosageRequired, obtive: %v", err)
	}
}

func TestMedicationValidate_MissingSchedule(t *testing.T) {
	m := &Medication{Name: "Omeprazol", Dosage: "20 mg"}
	if err := m.Validate(); err != ErrMedicationScheduleRequired {
		t.Fatalf("esperava ErrMedicationScheduleRequired, obtive: %v", err)
	}
}
