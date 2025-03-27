package middleware

import (
	"log"
	"net/http"
	"time"
)

// statusWriter wraps the http.ResponseWriter to capture the status code.
type statusWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code of the response when set.
func (w *statusWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// Logger is a middleware that logs the request details.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(sw, r)
		log.Println(sw.statusCode, r.Method, r.URL.Path, time.Since(start))
	})
}
