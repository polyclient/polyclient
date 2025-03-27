package healthcheck

import (
	"encoding/json"
	"net/http"
)

func NewHandler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /healthcheck", http.HandlerFunc(check))

	return http.StripPrefix("/healthcheck", mux)
}

func check(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
}
