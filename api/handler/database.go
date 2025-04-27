package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/polyclient/polyclient/internal/engine"
)

// DatabaseHandler is the HTTP handler for the table resource.
type DatabaseHandler struct {
	engine *engine.Engine
}

// NewDatabaseHandler creates a new instance of DatabaseHandler.
func NewDatabaseHandler(e *engine.Engine) *DatabaseHandler {
	return &DatabaseHandler{engine: e}
}

// RegisterRoutes registers the routes for the DatabaseHandler.
func (h *DatabaseHandler) RegisterRoutes(r chi.Router) {
	r.Get("/", h.handleListDatabases)
	r.Get("/{name}", h.handleGetDatabase)
}

func (h *DatabaseHandler) handleListDatabases(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *DatabaseHandler) handleGetDatabase(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
