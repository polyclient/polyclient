package api

import (
	"net/http"

	"github.com/polyclient/polyclient/api/features/health"
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

	mux.Handle("/api/health", health.NewHandler())

	return &Router{mux: mux}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
