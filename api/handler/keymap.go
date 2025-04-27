package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/polyclient/polyclient/internal/engine"
)

// KeymapHandler is the HTTP handler for the keymap resource.
type KeymapHandler struct {
	engine *engine.Engine
}

// NewKeymapHandler creates a new instance of KeymapHandler.
func NewKeymapHandler(e *engine.Engine) *KeymapHandler {
	return &KeymapHandler{engine: e}
}

// RegisterRoutes registers the routes for the KeymapHandler.
func (h *KeymapHandler) RegisterRoutes(r chi.Router) {
	r.Get("/", h.handleGetKeymap)
}

func (h *KeymapHandler) handleGetKeymap(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
