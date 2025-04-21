package db

import (
	"context"
	"fmt"
)

// SDK provides high-level access to database operations.
type SDK struct {
	manager *ConnectionManager
}

// NewDatabaseSDK creates a new DatabaseSDK instance.
func NewDatabaseSDK(manager *ConnectionManager) *SDK {
	return &SDK{
		manager: manager,
	}
}

// GetManager returns the connection manager.
func (s *SDK) GetManager() *ConnectionManager {
	return s.manager
}

// Ping checks the liveness of the named connection.
func (s *SDK) Ping(ctx context.Context, connectionName string) error {
	active, err := s.manager.GetActiveConnection(ctx, connectionName)
	if err != nil {
		return fmt.Errorf("failed to get active connection: %w", err)
	}

	return active.handle.Ping(ctx)
}

// InfoOperations provides high-level access to information about the database connection.
type InfoOperations struct {
	sdk *SDK
}

// Info provides access to info-related operations.
func (s *SDK) Info() *InfoOperations {
	return &InfoOperations{sdk: s}
}

// CurrentDatabase returns the name of the current database used by the connection.
func (o *InfoOperations) CurrentDatabase(ctx context.Context, connectionName string) (string, error) {
	active, err := o.sdk.manager.GetActiveConnection(ctx, connectionName)
	if err != nil {
		return "", fmt.Errorf("failed to get active connection: %w", err)
	}

	info, ok := active.handle.(ConnectionInfo)
	if !ok {
		return "unknown", nil
	}

	return info.CurrentDatabase(ctx), nil
}

// ServerVersion returns the version of the database server used by the connection.
func (o *InfoOperations) ServerVersion(ctx context.Context, connectionName string) (string, error) {
	active, err := o.sdk.manager.GetActiveConnection(ctx, connectionName)
	if err != nil {
		return "", fmt.Errorf("failed to get active connection: %w", err)
	}

	info, ok := active.handle.(ConnectionInfo)
	if !ok {
		return "unknown", nil
	}

	return info.ServerVersion(ctx), nil
}

// InspectorOperations provides access to inspector-related operations.
type InspectorOperations struct {
	sdk *SDK
}

// Inspector provides access to inspector-related operations.
func (s *SDK) Inspector() *InspectorOperations {
	return &InspectorOperations{sdk: s}
}

// ListSchemas returns a list of all schemas in the connection.
func (o *InspectorOperations) ListSchemas(ctx context.Context, connectionName string) ([]SchemaSummary, error) {
	active, err := o.sdk.manager.GetActiveConnection(ctx, connectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to get active connection: %w", err)
	}

	inspector, ok := active.handle.(SchemaLister)
	if !ok {
		return nil, fmt.Errorf("driver %s does not support schema listing", active.driverName)
	}

	return inspector.ListSchemas(ctx)
}

// GetSchema returns the schema with the given name.
func (o *InspectorOperations) GetSchema(ctx context.Context, connectionName, schemaName string) (*SchemaDetail, error) {
	active, err := o.sdk.manager.GetActiveConnection(ctx, connectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to get active connection: %w", err)
	}

	inspector, ok := active.handle.(SchemaGetter)
	if !ok {
		return nil, fmt.Errorf("driver %s does not support schema listing", active.driverName)
	}

	return inspector.GetSchema(ctx, schemaName)
}

// ListTables returns a list of all tables in the connection.
func (o *InspectorOperations) ListTables(ctx context.Context, connectionName string, options ...ListTablesOption) ([]TableSummary, error) {
	active, err := o.sdk.manager.GetActiveConnection(ctx, connectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to get active connection: %w", err)
	}

	inspector, ok := active.handle.(TableLister)
	if !ok {
		return nil, fmt.Errorf("driver %s does not support table listing", active.driverName)
	}

	return inspector.ListTables(ctx, options...)
}

// GetTable returns the table with the given name.
func (o *InspectorOperations) GetTable(ctx context.Context, connectionName, tableName string, options ...GetTableOption) (*TableDetail, error) {
	active, err := o.sdk.manager.GetActiveConnection(ctx, connectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to get active connection: %w", err)
	}

	inspector, ok := active.handle.(TableGetter)
	if !ok {
		return nil, fmt.Errorf("driver %s does not support table listing", active.driverName)
	}

	return inspector.GetTable(ctx, tableName, options...)
}

// QueryOperations provides access to query-related operations.
type QueryOperations struct {
	sdk *SDK
}

// Query provides access to query-related operations.
func (s *SDK) Query() *QueryOperations {
	return &QueryOperations{sdk: s}
}

// Execute runs a query against the connection.
// The 'query' parameter is driver-specific (e.g., SQL string, BSON document, Cypher).
// Use driver-specific query builder objects or simple types like string/map.
func (o *QueryOperations) Execute(ctx context.Context, connectionName, query string) (Result, error) {
	active, err := o.sdk.manager.GetActiveConnection(ctx, connectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to get active connection: %w", err)
	}

	executor, ok := active.handle.(QueryExecutor)
	if !ok {
		return nil, fmt.Errorf("driver %s does not support query execution", active.driverName)
	}

	return executor.Execute(ctx, query)
}
