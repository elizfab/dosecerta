package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"dose-certa-backend/internal/domain/entity"
	"dose-certa-backend/pkg/response"
)

// settingUseCase é a dependência do handler (interface → testável com fakes).
type settingUseCase interface {
	Get(ctx context.Context, key string) (*entity.Setting, error)
	Set(ctx context.Context, key, value string) (*entity.Setting, error)
}

type SettingHandler struct {
	uc settingUseCase
}

func NewSettingHandler(uc settingUseCase) *SettingHandler {
	return &SettingHandler{uc: uc}
}

// RegisterRoutes registra as rotas de configuração (key/value).
//
//	GET /api/v1/settings/{key}          → lê a configuração
//	PUT /api/v1/settings/{key} {value}  → cria/atualiza a configuração
func (h *SettingHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/settings/{key}", h.get)
	mux.HandleFunc("PUT /api/v1/settings/{key}", h.set)
}

type settingInput struct {
	Value string `json:"value"`
}

func (h *SettingHandler) get(w http.ResponseWriter, r *http.Request) {
	item, err := h.uc.Get(r.Context(), r.PathValue("key"))
	if err != nil {
		writeDomainError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, item)
}

func (h *SettingHandler) set(w http.ResponseWriter, r *http.Request) {
	var in settingInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response.Error(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}
	saved, err := h.uc.Set(r.Context(), r.PathValue("key"), in.Value)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, saved)
}
