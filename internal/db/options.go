// Order is the order to sort the list of tables.

package db

// ListTablesOptions defines options for the SchemaLister.ListTables method.
type ListTablesOptions struct {
	// Schema is the name of the schema to list tables from.
	Schema string
	// Filter is a filter to apply to the list of tables.
	Filter string
	// Limit is the maximum number of tables to return.
	Limit int
	// Offset is the offset of the first table to return.
	Offset int
}

// ListTablesOption is a functional option for configuring ListTablesOptions.
type ListTablesOption func(*ListTablesOptions)

// WithTablesSchema sets the schema for listing tables.
func WithTablesSchema(schema string) ListTablesOption {
	return func(opts *ListTablesOptions) {
		if schema != "" {
			opts.Schema = schema
		}
	}
}

// WithTablesFilter sets the filter for listing tables.
func WithTablesFilter(filter string) ListTablesOption {
	return func(opts *ListTablesOptions) {
		if filter != "" {
			opts.Filter = filter
		}
	}
}

// WithTablesLimit sets the limit for listing tables.
func WithTablesLimit(limit int) ListTablesOption {
	return func(opts *ListTablesOptions) {
		if limit > 0 {
			opts.Limit = limit
		}
	}
}

// WithTablesOffset sets the offset for listing tables.
func WithTablesOffset(offset int) ListTablesOption {
	return func(opts *ListTablesOptions) {
		if offset > 0 {
			opts.Offset = offset
		}
	}
}

// NewListTablesOptions creates a ListTablesOptions with the given options applied.
func NewListTablesOptions(opts ...ListTablesOption) *ListTablesOptions {
	options := &ListTablesOptions{}

	for _, opt := range opts {
		opt(options)
	}

	return options
}
