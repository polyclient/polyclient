package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/polyclient/polyclient/internal/constant"
	"github.com/polyclient/polyclient/internal/db"
	"github.com/polyclient/polyclient/internal/engine"
)

// TableHandler is the HTTP handler for the table resource.
type TableHandler struct {
	engine *engine.Engine
}

// NewTableHandler creates a new instance of TableHandler.
func NewTableHandler(e *engine.Engine) *TableHandler {
	return &TableHandler{engine: e}
}

// RegisterRoutes registers the routes for the TableHandler.
func (h *TableHandler) RegisterRoutes(r chi.Router) {
	r.Get("/", h.handleListTables)
	r.Get("/{name}", h.handleGetTable)
}

func (h *TableHandler) handleListTables(w http.ResponseWriter, r *http.Request) {
	connName := r.Context().Value(constant.ContextKeyConnectionName).(string)

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

	tables, err := h.engine.SDK.Inspector().ListTables(r.Context(), connName,
		db.WithTablesSchema(schema),
		db.WithTablesFilter(filter),
		db.WithTablesLimit(limitInt),
		db.WithTablesOffset(offsetInt),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(tables); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *TableHandler) handleGetTable(w http.ResponseWriter, r *http.Request) {
	connName := r.Context().Value(constant.ContextKeyConnectionName).(string)
	name := chi.URLParam(r, "name")

	table, err := h.engine.SDK.Inspector().GetTable(r.Context(), connName, name)
	if err != nil {
		// TODO: implement error constants
		// if errors.Is(err, db.ErrTableNotFound) {
		// 	http.Error(w, err.Error(), http.StatusNotFound)
		// 	return
		// }
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(table); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
