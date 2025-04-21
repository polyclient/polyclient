package resource

import (
	"net/http"

	"github.com/polyclient/polyclient/internal/db"
)

// ConnectionNameHeader is the HTTP header for the connection name.
const ConnectionNameHeader = "X-Connection-Name"

// ConnectionCreateRequest is the expected JSON body for the connection endpoint.
type ConnectionCreateRequest struct {
	Driver string              `json:"driver"`
	Config db.ConnectionConfig `json:"config"`
}

// ConnectionCreateResponse represents a response to create a new database connection.
type ConnectionCreateResponse struct {
	SessionID string `json:"sessionId"`
}

// ConnectionHandler is the HTTP handler for the connection endpoint.
type ConnectionHandler struct {
	sdk *db.SDK
	mux *http.ServeMux
}

// NewConnectionHandler creates a new HTTP handler for the connection endpoint.
func NewConnectionHandler(sdk *db.SDK) *ConnectionHandler {
	h := &ConnectionHandler{
		sdk: sdk,
		mux: http.NewServeMux(),
	}

	h.registerRoutes(h.mux)

	return h
}

// ServeHTTP serves the HTTP handler for the connection endpoint.
func (h *ConnectionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

// HandleListConnections handles the list connections endpoint.
func (h *ConnectionHandler) HandleListConnections(w http.ResponseWriter, r *http.Request) {
}

// HandleListRecentConnections handles the list recent connections endpoint.
func (h *ConnectionHandler) HandleListRecentConnections(w http.ResponseWriter, r *http.Request) {
}

// HandleGetConnection handles the get connection endpoint.
func (h *ConnectionHandler) HandleGetConnection(w http.ResponseWriter, r *http.Request) {
}

// HandleCreateConnection handles the create connection endpoint.
func (h *ConnectionHandler) HandleCreateConnection(w http.ResponseWriter, r *http.Request) {
}

// HandleDeleteConnection handles the delete connection endpoint.
func (h *ConnectionHandler) HandleDeleteConnection(w http.ResponseWriter, r *http.Request) {
}

func (h *ConnectionHandler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", h.HandleListConnections)
	mux.HandleFunc("GET /recent", h.HandleListRecentConnections)
	mux.HandleFunc("GET /{name}", h.HandleGetConnection)
	mux.HandleFunc("POST /", h.HandleCreateConnection)
	mux.HandleFunc("DELETE /{name}", h.HandleDeleteConnection)
}
