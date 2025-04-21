// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package api

import (
	"net/http"

	"github.com/polyclient/polyclient/api/resource"
	"github.com/polyclient/polyclient/gui"
)

// Router is the HTTP router for the API.
type Router struct {
	mux *http.ServeMux
}

// NewRouter creates a new HTTP router for the API.
func NewRouter() *Router {
	mux := http.NewServeMux()

	mux.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

			return
		}

		http.FileServer(http.FS(gui.DistDirFS)).ServeHTTP(w, r)
	}))

	mux.Handle("/api/connections", resource.NewConnectionHandler(nil))

	return &Router{mux: mux}
}

// ServeHTTP serves the HTTP router.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
