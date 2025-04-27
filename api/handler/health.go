package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/polyclient/polyclient/internal/engine"
)

// HealthHandler is the HTTP handler for the health resource.
type HealthHandler struct {
	engine *engine.Engine
}

// NewHealthHandler creates a new instance of HealthHandler.
func NewHealthHandler(e *engine.Engine) *HealthHandler {
	return &HealthHandler{engine: e}
}

// RegisterRoutes registers the routes for the HealthHandler.
func (h *HealthHandler) RegisterRoutes(r chi.Router) {
	r.Get("/", h.handleGetHealth)
}

func (h *HealthHandler) handleGetHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
