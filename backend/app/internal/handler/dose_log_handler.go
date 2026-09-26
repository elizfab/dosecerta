package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"dose-certa-backend/internal/domain/entity"
	"dose-certa-backend/pkg/response"
)

// doseLogUseCase é a dependência do handler (interface → testável com fakes).
type doseLogUseCase interface {
	MarkTaken(ctx context.Context, medicationID, date string) (*entity.DoseLog, error)
	Unmark(ctx context.Context, medicationID, date string) error
	ListByDate(ctx context.Context, date string) ([]*entity.DoseLog, error)
}

type DoseLogHandler struct {
	uc doseLogUseCase
}

func NewDoseLogHandler(uc doseLogUseCase) *DoseLogHandler {
	return &DoseLogHandler{uc: uc}
}

// RegisterRoutes registra as rotas de doses tomadas.
//
//	GET    /api/v1/doses?date=YYYY-MM-DD          → lista as doses tomadas na data
//	POST   /api/v1/doses  {medicationId, date}    → marca dose como tomada
//	DELETE /api/v1/doses?medicationId=..&date=..  → desmarca
func (h *DoseLogHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/doses", h.list)
	mux.HandleFunc("POST /api/v1/doses", h.mark)
	mux.HandleFunc("DELETE /api/v1/doses", h.unmark)
}

type doseInput struct {
	MedicationID string `json:"medicationId"`
	Date         string `json:"date"`
}

func (h *DoseLogHandler) list(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	items, err := h.uc.ListByDate(r.Context(), date)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, items)
}

func (h *DoseLogHandler) mark(w http.ResponseWriter, r *http.Request) {
	var in doseInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response.Error(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}
	created, err := h.uc.MarkTaken(r.Context(), in.MedicationID, in.Date)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, created)
}

func (h *DoseLogHandler) unmark(w http.ResponseWriter, r *http.Request) {
	medicationID := r.URL.Query().Get("medicationId")
	date := r.URL.Query().Get("date")
	if err := h.uc.Unmark(r.Context(), medicationID, date); err != nil {
		writeDomainError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"medicationId": medicationID, "date": date})
}
