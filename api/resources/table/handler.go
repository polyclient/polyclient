package table

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/polyclient/polyclient/api/config"
	"github.com/polyclient/polyclient/internal/application"
	"github.com/polyclient/polyclient/internal/db"
)

// Handler is the HTTP handler for the table endpoint.
type Handler struct {
	app *application.Application
}

// NewHandler creates a new HTTP handler for the table endpoint.
func NewHandler(app *application.Application) *Handler {
	return &Handler{app: app}
}

// RegisterRoutes registers the routes for the table endpoint.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/", h.handleListTables)
	r.Get("/{name}", h.handleGetTable)
}

func (h *Handler) handleListTables(w http.ResponseWriter, r *http.Request) {
	connName := r.Context().Value(config.ContextKeyConnectionName).(string)

	schema := r.URL.Query().Get("schema")
	filter := r.URL.Query().Get("filter")
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	limitInt, err := strconv.Atoi(limit)
	if err != nil && limit != "" {
		http.Error(w, "invalid limit parameter", http.StatusBadRequest)
		return
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil && offset != "" {
		http.Error(w, "invalid offset parameter", http.StatusBadRequest)
		return
	}

	tables, err := h.app.SDK.Inspector().ListTables(r.Context(), connName,
		db.WithTablesSchema(schema),
		db.WithTablesFilter(filter),
		db.WithTablesLimit(limitInt),
		db.WithTablesOffset(offsetInt),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(tables); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) handleGetTable(w http.ResponseWriter, r *http.Request) {
	connName := r.Context().Value(config.ContextKeyConnectionName).(string)
	name := chi.URLParam(r, "name")

	table, err := h.app.SDK.Inspector().GetTable(r.Context(), connName, name)
	if err != nil {
		// TODO: implement error constants
		// if errors.Is(err, db.ErrTableNotFound) {
		// 	http.Error(w, err.Error(), http.StatusNotFound)
		// 	return
		// }
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(table); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
