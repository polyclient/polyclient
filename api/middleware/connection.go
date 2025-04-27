package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/polyclient/polyclient/internal/constant"
)

// ConnectionName is a Chi middleware that extracts the connection name from the request header
// and makes it available in the request context under the key ContextKeyConnectionName.
// If the header is missing, it returns a 400 error to the client.
func ConnectionName(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		connName := strings.TrimSpace(r.Header.Get(constant.HTTPHeaderConnectionName))
		if connName == "" {
			http.Error(w, fmt.Sprintf(
				"missing connection name in '%s' header", constant.HTTPHeaderConnectionName,
			), http.StatusBadRequest)

			return
		}

		ctx := context.WithValue(r.Context(), constant.ContextKeyConnectionName, connName)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
