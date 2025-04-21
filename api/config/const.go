package config

const (
	// HTTPHeaderConnectionName is the HTTP header for the connection name.
	HTTPHeaderConnectionName = "X-PolyClient-Connection-Name"
)

// ContextKey is the type for context keys.
type ContextKey string

const (
	// ContextKeyConnectionName is the context key for the connection name.
	ContextKeyConnectionName ContextKey = "connectionName"
)
