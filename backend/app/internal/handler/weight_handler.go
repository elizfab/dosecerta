package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"dose-certa-backend/internal/domain/entity"
	"dose-certa-backend/pkg/response"
)

// weightUseCase é a dependência do handler (interface → testável com fakes).
type weightUseCase interface {
	Create(ctx context.Context, w *entity.WeightRecord) (*entity.WeightRecord, error)
	List(ctx context.Context) ([]*entity.WeightRecord, error)
	Get(ctx context.Context, id string) (*entity.WeightRecord, error)
	Update(ctx context.Context, id string, w *entity.WeightRecord) (*entity.WeightRecord, error)
	Delete(ctx context.Context, id string) error
}

type WeightHandler struct {
	uc weightUseCase
}

func NewWeightHandler(uc weightUseCase) *WeightHandler {
	return &WeightHandler{uc: uc}
}

func (h *WeightHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/weight-records", h.list)
	mux.HandleFunc("POST /api/v1/weight-records", h.create)
	mux.HandleFunc("GET /api/v1/weight-records/{id}", h.get)
	mux.HandleFunc("PUT /api/v1/weight-records/{id}", h.update)
	mux.HandleFunc("DELETE /api/v1/weight-records/{id}", h.delete)
}

type weightInput struct {
	Date         string  `json:"date"`
	WeightKg     float64 `json:"weightKg"`
	Waist        float64 `json:"waist"`
	Hip          float64 `json:"hip"`
	Glute        float64 `json:"glute"`
	Chest        float64 `json:"chest"`
	ArmRight     float64 `json:"armRight"`
	ArmLeft      float64 `json:"armLeft"`
	ForearmRight float64 `json:"forearmRight"`
	ForearmLeft  float64 `json:"forearmLeft"`
	ThighRight   float64 `json:"thighRight"`
	ThighLeft    float64 `json:"thighLeft"`
	CalfRight    float64 `json:"calfRight"`
	CalfLeft     float64 `json:"calfLeft"`
	Neck         float64 `json:"neck"`
	Shoulders    float64 `json:"shoulders"`
	BodyFatPct   float64 `json:"bodyFatPct"`
}

func (in *weightInput) toEntity() *entity.WeightRecord {
	return &entity.WeightRecord{
		Date: in.Date, WeightKg: in.WeightKg,
		Waist: in.Waist, Hip: in.Hip, Glute: in.Glute, Chest: in.Chest,
		ArmRight: in.ArmRight, ArmLeft: in.ArmLeft,
		ForearmRight: in.ForearmRight, ForearmLeft: in.ForearmLeft,
		ThighRight: in.ThighRight, ThighLeft: in.ThighLeft,
		CalfRight: in.CalfRight, CalfLeft: in.CalfLeft,
		Neck: in.Neck, Shoulders: in.Shoulders, BodyFatPct: in.BodyFatPct,
	}
}

func (h *WeightHandler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.uc.List(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, items)
}

func (h *WeightHandler) create(w http.ResponseWriter, r *http.Request) {
	var in weightInput
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

func (h *WeightHandler) get(w http.ResponseWriter, r *http.Request) {
	item, err := h.uc.Get(r.Context(), r.PathValue("id"))
	if err != nil {
		writeDomainError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, item)
}

func (h *WeightHandler) update(w http.ResponseWriter, r *http.Request) {
	var in weightInput
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

func (h *WeightHandler) delete(w http.ResponseWriter, r *http.Request) {
	if err := h.uc.Delete(r.Context(), r.PathValue("id")); err != nil {
		writeDomainError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"id": r.PathValue("id")})
}
