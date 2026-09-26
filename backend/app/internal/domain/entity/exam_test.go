package entity

import "testing"

func TestExamValidate_Success(t *testing.T) {
	e := &Exam{Name: "Glicose", Date: "2026-07-25", Value: 90, Unit: "mg/dL"}
	if err := e.Validate(); err != nil {
		t.Fatalf("esperava sucesso, obtive %v", err)
	}
}

func TestExamValidate_MissingName(t *testing.T) {
	e := &Exam{Date: "2026-07-25"}
	if err := e.Validate(); err != ErrExamNameRequired {
		t.Fatalf("esperava ErrExamNameRequired, obtive %v", err)
	}
}

func TestExamValidate_MissingDate(t *testing.T) {
	e := &Exam{Name: "Glicose"}
	if err := e.Validate(); err != ErrExamDateRequired {
		t.Fatalf("esperava ErrExamDateRequired, obtive %v", err)
	}
}
