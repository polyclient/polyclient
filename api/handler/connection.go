// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/polyclient/polyclient/internal/db"
	"github.com/polyclient/polyclient/internal/engine"
)

// ConnectionHandler is the HTTP handler for the connection resource.
type ConnectionHandler struct {
	engine *engine.Engine
}

// NewConnectionHandler creates a new instance of ConnectionHandler.
func NewConnectionHandler(e *engine.Engine) *ConnectionHandler {
	return &ConnectionHandler{engine: e}
}

// RegisterRoutes registers the routes for the ConnectionHandler.
func (h *ConnectionHandler) RegisterRoutes(r chi.Router) {
	r.Get("/", h.handleListConnections)
	r.Get("/recent", h.handleListRecentConnections)
	r.Get("/{name}", h.handleGetConnection)
	r.Post("/", h.handleCreateConnection)
	r.Delete("/{name}", h.handleDeleteConnection)
}

func (h *ConnectionHandler) handleListConnections(w http.ResponseWriter, r *http.Request) {
	profiles, err := h.engine.SDK.GetManager().GetStore().ListProfiles(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(profiles); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *ConnectionHandler) handleListRecentConnections(w http.ResponseWriter, r *http.Request) {
	const defaultThreshold = "24h"

	threshold := r.URL.Query().Get("threshold")
	if threshold == "" {
		threshold = defaultThreshold
	}

	duration, err := time.ParseDuration(threshold)
	if err != nil {
		http.Error(w, "Invalid threshold duration", http.StatusBadRequest)
		return
	}

	profiles, err := h.engine.SDK.GetManager().GetStore().ListRecentProfiles(r.Context(), duration)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(profiles); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *ConnectionHandler) handleGetConnection(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	profile, err := h.engine.SDK.GetManager().GetStore().GetProfile(r.Context(), name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(profile); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *ConnectionHandler) handleCreateConnection(w http.ResponseWriter, r *http.Request) {
	profile := new(db.ConnectionProfile)
	if err := json.NewDecoder(r.Body).Decode(profile); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.engine.Validator.Validate(profile); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.engine.SDK.GetManager().GetStore().SaveProfile(r.Context(), profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(profile); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *ConnectionHandler) handleDeleteConnection(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	err := h.engine.SDK.GetManager().GetStore().DeleteProfile(r.Context(), name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if _, err := w.Write([]byte(name)); err != nil {
		http.Error(w, "failed to write response", http.StatusInternalServerError)
	}
}
