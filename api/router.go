// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	pMiddleware "github.com/polyclient/polyclient/api/middleware"
	"github.com/polyclient/polyclient/api/resources/connection"
	"github.com/polyclient/polyclient/api/resources/table"
	"github.com/polyclient/polyclient/gui"
	"github.com/polyclient/polyclient/internal/application"
)

// NewRouter creates a new router for the API.
func NewRouter(app *application.Application) (http.Handler, error) {
	r := chi.NewRouter()

	// Middleware setup
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(func(next http.Handler) http.Handler {
		return http.MaxBytesHandler(next, 1<<20) // 1MB body limit
	})

	// Remove trailing slash
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" && r.URL.Path[len(r.URL.Path)-1] == '/' {
				http.Redirect(w, r, r.URL.Path[:len(r.URL.Path)-1], http.StatusMovedPermanently)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	registerStaticRoutes(r)
	registerAPIRoutes(r, app)

	return r, nil
}

// registerStaticRoutes registers the static routes for the API.
func registerStaticRoutes(r *chi.Mux) {
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		fs := http.FileServer(http.FS(gui.DistDirFS))
		fs.ServeHTTP(w, r)
	})
}

// registerAPIRoutes registers the API routes.
func registerAPIRoutes(router *chi.Mux, app *application.Application) {
	router.Route("/api", func(api chi.Router) {
		// Connection routes
		api.Route("/connections", func(r chi.Router) {
			connectionHandler := connection.NewHandler(app)
			connectionHandler.RegisterRoutes(r)
		})

		// Table routes with middleware
		api.Route("/tables", func(r chi.Router) {
			r.Use(pMiddleware.ConnectionName)
			tableHandler := table.NewHandler(app)
			tableHandler.RegisterRoutes(r)
		})
	})
}
