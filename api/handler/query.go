package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/polyclient/polyclient/internal/engine"
)

// QueryHandler is the HTTP handler for the query resource.
type QueryHandler struct {
	engine *engine.Engine
}

// NewQueryHandler creates a new instance of QueryHandler.
func NewQueryHandler(e *engine.Engine) *QueryHandler {
	return &QueryHandler{engine: e}
}

// RegisterRoutes registers the routes for the QueryHandler.
func (h *QueryHandler) RegisterRoutes(r chi.Router) {}
