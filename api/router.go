package api

import (
	"net/http"

	"github.com/polyclient/polyclient/api/features/healthcheck"
	"github.com/polyclient/polyclient/gui"
)

type Router struct {
	mux *http.ServeMux
}

func NewRouter() *Router {
	mux := http.NewServeMux()

	mux.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		http.FileServer(http.FS(gui.DistDirFS)).ServeHTTP(w, r)
	}))

	// API v1 routes
	v1Router := http.NewServeMux()
	v1Router.Handle("/healthcheck/", healthcheck.NewHandler())
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", v1Router))

	return &Router{mux: mux}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
