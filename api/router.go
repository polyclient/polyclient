// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package api

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/polyclient/polyclient/api/handler"
	pMiddleware "github.com/polyclient/polyclient/api/middleware"
	"github.com/polyclient/polyclient/gui"
	"github.com/polyclient/polyclient/internal/engine"
)

// Router is the main API router for the API.
type Router struct {
	*chi.Mux
	engine *engine.Engine
}

// NewRouter creates a new router for the API.
func NewRouter(ctx context.Context, e *engine.Engine) *Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	if e.Settings.API.CORS.Enabled {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins: e.Settings.API.CORS.AllowedOrigins,
			AllowedMethods: e.Settings.API.CORS.AllowedMethods,
			AllowedHeaders: e.Settings.API.CORS.AllowedHeaders,
			MaxAge:         e.Settings.API.CORS.MaxAge,
		}))
	}

	if e.Settings.API.RateLimit.Enabled {
		r.Use(httprate.Limit(
			e.Settings.API.RateLimit.RequestsPerMinute,
			e.Settings.API.RateLimit.WindowLength,
		))
	}

	r.Use(middleware.CleanPath)
	r.Use(middleware.StripSlashes)

	if e.Settings.API.Compression.Enabled {
		r.Use(middleware.Compress(e.Settings.API.Compression.Level))
	}

	r.Use(middleware.AllowContentType("application/json", "application/x-www-form-urlencoded"))

	return &Router{r, e}
}

// RegisterGUIRoutes configures the router for the GUI.
func (r *Router) RegisterGUIRoutes() {
	staticFileServer := http.FileServer(http.FS(gui.DistDirFS))

	r.Mux.Get(r.engine.Settings.GUI.Path, func(w http.ResponseWriter, r *http.Request) {
		// TODO: Need to add caching headers here
		staticFileServer.ServeHTTP(w, r)
	})
}

// RegisterAPIRoutesV1 configures the routing for the v1 API endpoints (/api/v1).
func (r *Router) RegisterAPIV1Routes() {
	r.Mux.Route("/api/v1", func(apiRouter chi.Router) {
		apiRouter.Use(middleware.SetHeader("Content-Type", "application/json"))

		apiRouter.Route("/health", func(apiRouter chi.Router) {
			h := handler.NewHealthHandler(r.engine)
			h.RegisterRoutes(apiRouter)
		})

		apiRouter.Route("/settings", func(apiRouter chi.Router) {
			h := handler.NewSettingsHandler(r.engine)
			h.RegisterRoutes(apiRouter)
		})

		apiRouter.Route("/keymap", func(apiRouter chi.Router) {
			h := handler.NewKeymapHandler(r.engine)
			h.RegisterRoutes(apiRouter)
		})

		apiRouter.Route("/connections", func(apiRouter chi.Router) {
			h := handler.NewConnectionHandler(r.engine)
			h.RegisterRoutes(apiRouter)
		})

		apiRouter.Route("/databases", func(apiRouter chi.Router) {
			apiRouter.Use(pMiddleware.ConnectionName)

			h := handler.NewDatabaseHandler(r.engine)
			h.RegisterRoutes(apiRouter)
		})

		apiRouter.Route("/tables", func(apiRouter chi.Router) {
			apiRouter.Use(pMiddleware.ConnectionName)

			h := handler.NewTableHandler(r.engine)
			h.RegisterRoutes(apiRouter)
		})

		apiRouter.Route("/queries", func(apiRouter chi.Router) {
			// TODO: analyze the best way to implement a custom middleware for query timeouts

			apiRouter.Use(pMiddleware.ConnectionName)

			h := handler.NewQueryHandler(r.engine)
			h.RegisterRoutes(apiRouter)
		})
	})
}
