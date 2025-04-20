package sdk

import (
	"context"
	"fmt"

	"github.com/polyclient/polyclient/internal/db"
)

// DatabaseSDK provides high-level access to database operations.
// built on top of the driver registry and connection interfaces.
type DatabaseSDK struct {
	driverRegistry *db.Registry[db.Driver]
}

// NewDatabaseSDK creates a new DatabaseSDK instance.
func NewDatabaseSDK(driverRegistry *db.Registry[db.Driver]) *DatabaseSDK {
	return &DatabaseSDK{
		driverRegistry: driverRegistry,
	}
}

// Session represents a database session with an open connection.
type Session struct {
	conn       db.Connection
	driverName string
}

// SchemaOperations provides access to schema-related operations.
type SchemaOperations struct {
	session *Session
}

// QueryOperations provides access to query-related operations.
type QueryOperations struct {
	session *Session
}

// OpenConnection opens a new database connection.
func (s *DatabaseSDK) OpenConnection(ctx context.Context, driverName string, config db.Config) (*Session, error) {
	driver, ok := s.driverRegistry.Get(driverName)
	if !ok {
		return nil, ErrDriverNotFound
	}

	conn, err := driver.Open(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}

	return &Session{
		conn:       conn,
		driverName: driverName,
	}, nil
}

// Info returns information about the current database connection.
func (s *Session) Info() db.ConnectionInfo {
	return s.conn.Info()
}

// Ping verifies that the connection is alive.
func (s *Session) Ping(ctx context.Context) error {
	return s.conn.Ping(ctx)
}

// CloseConnection closes the database connection.
func (s *Session) CloseConnection() error {
	return s.conn.Close()
}

// Schema provides access to schema-related operations.
func (s *Session) Schema() *SchemaOperations {
	return &SchemaOperations{session: s}
}

// ListSchemas returns a list of all schemas in the connection.
func (o *SchemaOperations) ListSchemas(ctx context.Context) ([]db.SchemaSummary, error) {
	inspector, ok := o.session.conn.(db.SchemaLister)
	if !ok {
		return nil, ErrDriverDoesNotSupportSchemaListing
	}

	return inspector.ListSchemas(ctx)
}

// GetSchema returns the schema with the given name.
func (o *SchemaOperations) GetSchema(ctx context.Context, name string) (db.SchemaDetail, error) {
	inspector, ok := o.session.conn.(db.SchemaGetter)
	if !ok {
		return db.SchemaDetail{}, ErrDriverDoesNotSupportSchemaListing
	}

	return inspector.GetSchema(ctx, name)
}

// ListTables returns a list of all tables in the connection.
func (o *SchemaOperations) ListTables(ctx context.Context, options ...db.ListTablesOption) ([]db.TableSummary, error) {
	inspector, ok := o.session.conn.(db.TableLister)
	if !ok {
		return nil, ErrDriverDoesNotSupportTableListing
	}

	return inspector.ListTables(ctx, options...)
}

// GetTable returns the table with the given name.
func (o *SchemaOperations) GetTable(ctx context.Context, name string, options ...db.GetTableOption) (db.TableDetail, error) {
	inspector, ok := o.session.conn.(db.TableGetter)
	if !ok {
		return db.TableDetail{}, ErrDriverDoesNotSupportTableListing
	}

	return inspector.GetTable(ctx, name, options...)
}

// ListViews returns a list of all views in the connection.
func (o *SchemaOperations) ListViews(ctx context.Context) ([]db.ViewSummary, error) {
	inspector, ok := o.session.conn.(db.ViewLister)
	if !ok {
		return nil, ErrDriverDoesNotSupportViewListing
	}

	return inspector.ListViews(ctx)
}

// GetView returns the view with the given name.
func (o *SchemaOperations) GetView(ctx context.Context, name string) (db.ViewDetail, error) {
	inspector, ok := o.session.conn.(db.ViewGetter)
	if !ok {
		return db.ViewDetail{}, ErrDriverDoesNotSupportViewListing
	}

	return inspector.GetView(ctx, name)
}

// ListMaterializedViews returns a list of all materialized views in the connection.
func (o *SchemaOperations) ListMaterializedViews(ctx context.Context) ([]db.MaterializedViewSummary, error) {
	inspector, ok := o.session.conn.(db.MaterializedViewLister)
	if !ok {
		return nil, ErrDriverDoesNotSupportMaterializedViewListing
	}

	return inspector.ListMaterializedViews(ctx)
}

// GetMaterializedView returns the materialized view with the given name.
func (o *SchemaOperations) GetMaterializedView(ctx context.Context, name string) (db.MaterializedViewDetail, error) {
	inspector, ok := o.session.conn.(db.MaterializedViewGetter)
	if !ok {
		return db.MaterializedViewDetail{}, ErrDriverDoesNotSupportMaterializedViewListing
	}

	return inspector.GetMaterializedView(ctx, name)
}

// Query provides access to query-related operations.
func (s *Session) Query() *QueryOperations {
	return &QueryOperations{session: s}
}

// Execute runs a query against the connection.
// The 'query' parameter is driver-specific (e.g., SQL string, BSON document, Cypher).
// Use driver-specific query builder objects or simple types like string/map.
func (o *QueryOperations) Execute(ctx context.Context, query string) (db.Result, error) {
	executor, ok := o.session.conn.(db.QueryExecutor)
	if !ok {
		return nil, ErrDriverDoesNotSupportQueryExecution
	}

	return executor.Execute(ctx, query)
}
