package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/polyclient/polyclient/internal/engine"
)

// SettingsHandler is the HTTP handler for the settings resource.
type SettingsHandler struct {
	engine *engine.Engine
}

// NewSettingsHandler creates a new instance of SettingsHandler.
func NewSettingsHandler(e *engine.Engine) *SettingsHandler {
	return &SettingsHandler{engine: e}
}

// RegisterRoutes registers the routes for the SettingsHandler.
func (h *SettingsHandler) RegisterRoutes(r chi.Router) {
	r.Get("/settings", h.handleGetSettings)
}

func (h *SettingsHandler) handleGetSettings(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(h.engine.Settings); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
