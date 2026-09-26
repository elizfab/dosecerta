package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"dose-certa-backend/internal/domain/entity"
	"dose-certa-backend/pkg/response"
)

type examUseCase interface {
	Create(ctx context.Context, e *entity.Exam) (*entity.Exam, error)
	List(ctx context.Context) ([]*entity.Exam, error)
	Get(ctx context.Context, id string) (*entity.Exam, error)
	Update(ctx context.Context, id string, e *entity.Exam) (*entity.Exam, error)
	Delete(ctx context.Context, id string) error
}

type ExamHandler struct {
	uc examUseCase
}

func NewExamHandler(uc examUseCase) *ExamHandler {
	return &ExamHandler{uc: uc}
}

func (h *ExamHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/exams", h.list)
	mux.HandleFunc("POST /api/v1/exams", h.create)
	mux.HandleFunc("GET /api/v1/exams/{id}", h.get)
	mux.HandleFunc("PUT /api/v1/exams/{id}", h.update)
	mux.HandleFunc("DELETE /api/v1/exams/{id}", h.delete)
}

type examInput struct {
	Name         string  `json:"name"`
	Value        float64 `json:"value"`
	Unit         string  `json:"unit"`
	ReferenceMin float64 `json:"referenceMin"`
	ReferenceMax float64 `json:"referenceMax"`
	Date         string  `json:"date"`
	Notes        string  `json:"notes"`
}

func (in *examInput) toEntity() *entity.Exam {
	return &entity.Exam{
		Name:         in.Name,
		Value:        in.Value,
		Unit:         in.Unit,
		ReferenceMin: in.ReferenceMin,
		ReferenceMax: in.ReferenceMax,
		Date:         in.Date,
		Notes:        in.Notes,
	}
}

func (h *ExamHandler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.uc.List(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, items)
}

func (h *ExamHandler) create(w http.ResponseWriter, r *http.Request) {
	var in examInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response.Error(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}
	created, err := h.uc.Create(r.Context(), in.toEntity())
	if err != nil {
		writeDomainError(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, created)
}

func (h *ExamHandler) get(w http.ResponseWriter, r *http.Request) {
	item, err := h.uc.Get(r.Context(), r.PathValue("id"))
	if err != nil {
		writeDomainError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, item)
}

func (h *ExamHandler) update(w http.ResponseWriter, r *http.Request) {
	var in examInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response.Error(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}
	updated, err := h.uc.Update(r.Context(), r.PathValue("id"), in.toEntity())
	if err != nil {
		writeDomainError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, updated)
}

func (h *ExamHandler) delete(w http.ResponseWriter, r *http.Request) {
	if err := h.uc.Delete(r.Context(), r.PathValue("id")); err != nil {
		writeDomainError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"id": r.PathValue("id")})
}
