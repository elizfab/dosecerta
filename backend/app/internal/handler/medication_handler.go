package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"dose-certa-backend/internal/domain/entity"
	"dose-certa-backend/internal/domain/repository"
	"dose-certa-backend/pkg/response"
)

// medicationUseCase é a dependência do handler (interface → testável com fakes).
type medicationUseCase interface {
	Create(ctx context.Context, m *entity.Medication) (*entity.Medication, error)
	List(ctx context.Context) ([]*entity.Medication, error)
	Get(ctx context.Context, id string) (*entity.Medication, error)
	Update(ctx context.Context, id string, m *entity.Medication) (*entity.Medication, error)
	Delete(ctx context.Context, id string) error
}

type MedicationHandler struct {
	uc medicationUseCase
}

func NewMedicationHandler(uc medicationUseCase) *MedicationHandler {
	return &MedicationHandler{uc: uc}
}

// RegisterRoutes registra as rotas REST do recurso (net/http nativo, Go 1.22).
func (h *MedicationHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/medications", h.list)
	mux.HandleFunc("POST /api/v1/medications", h.create)
	mux.HandleFunc("GET /api/v1/medications/{id}", h.get)
	mux.HandleFunc("PUT /api/v1/medications/{id}", h.update)
	mux.HandleFunc("DELETE /api/v1/medications/{id}", h.delete)
}

type medicationInput struct {
	Name      string `json:"name"`
	Dosage    string `json:"dosage"`
	Schedule  string `json:"schedule"`
	Notes     string `json:"notes"`
	Period    string `json:"period"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Stock     int    `json:"stock"`
	Active    *bool  `json:"active"` // opcional; ausente = ativo (não pausa ao editar)
}

// toEntity converte o input em entidade. activeDefault define o active quando
// o cliente não envia o campo (ex.: criação → sempre ativo).
func (in medicationInput) toEntity(activeDefault bool) *entity.Medication {
	active := activeDefault
	if in.Active != nil {
		active = *in.Active
	}
	return &entity.Medication{
		Name:      in.Name,
		Dosage:    in.Dosage,
		Schedule:  in.Schedule,
		Notes:     in.Notes,
		Period:    in.Period,
		StartDate: in.StartDate,
		EndDate:   in.EndDate,
		Stock:     in.Stock,
		Active:    active,
	}
}

func (h *MedicationHandler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.uc.List(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, items)
}

func (h *MedicationHandler) create(w http.ResponseWriter, r *http.Request) {
	var in medicationInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response.Error(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}
	created, err := h.uc.Create(r.Context(), in.toEntity(true))
	if err != nil {
		writeDomainError(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, created)
}

func (h *MedicationHandler) get(w http.ResponseWriter, r *http.Request) {
	item, err := h.uc.Get(r.Context(), r.PathValue("id"))
	if err != nil {
		writeDomainError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, item)
}

func (h *MedicationHandler) update(w http.ResponseWriter, r *http.Request) {
	var in medicationInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response.Error(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}
	updated, err := h.uc.Update(r.Context(), r.PathValue("id"), in.toEntity(true))
	if err != nil {
		writeDomainError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, updated)
}

func (h *MedicationHandler) delete(w http.ResponseWriter, r *http.Request) {
	if err := h.uc.Delete(r.Context(), r.PathValue("id")); err != nil {
		writeDomainError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"id": r.PathValue("id")})
}

// writeDomainError mapeia erros de domínio para os códigos HTTP corretos.
func writeDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		response.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, repository.ErrInvalidID):
		response.Error(w, http.StatusBadRequest, err.Error())
	case isValidationError(err):
		response.Error(w, http.StatusBadRequest, err.Error())
	default:
		response.Error(w, http.StatusInternalServerError, err.Error())
	}
}

func isValidationError(err error) bool {
	return errors.Is(err, entity.ErrMedicationNameRequired) ||
		errors.Is(err, entity.ErrMedicationDosageRequired) ||
		errors.Is(err, entity.ErrMedicationScheduleRequired) ||
		errors.Is(err, entity.ErrMedicationScheduleInvalid) ||
		errors.Is(err, entity.ErrWeightDateRequired) ||
		errors.Is(err, entity.ErrWeightValueInvalid) ||
		errors.Is(err, entity.ErrWeightDateInvalid) ||
		errors.Is(err, entity.ErrExamNameRequired) ||
		errors.Is(err, entity.ErrExamDateRequired) ||
		errors.Is(err, entity.ErrExamDateInvalid) ||
		errors.Is(err, entity.ErrDoseMedicationRequired) ||
		errors.Is(err, entity.ErrDoseDateRequired) ||
		errors.Is(err, entity.ErrSettingKeyRequired)
}
