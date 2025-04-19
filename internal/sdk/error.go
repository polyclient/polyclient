package sdk

import "errors"

var (
	// ErrDriverNotFound is returned when a driver is not found in the registry.
	ErrDriverNotFound = errors.New("driver not found in the registry")

	// ErrDriverAlreadyRegistered is returned when a driver is already registered in the registry.
	ErrDriverAlreadyRegistered = errors.New("driver already registered in the registry")

	// ErrDriverDoesNotSupportDatabaseListing is returned when a driver does not support database listing.
	ErrDriverDoesNotSupportDatabaseListing = errors.New("driver does not support database listing")

	// ErrDriverDoesNotSupportSchemaListing is returned when a driver does not support schema listing.
	ErrDriverDoesNotSupportSchemaListing = errors.New("driver does not support schema listing")

	// ErrDriverDoesNotSupportTableListing is returned when a driver does not support table listing.
	ErrDriverDoesNotSupportTableListing = errors.New("driver does not support table listing")

	// ErrDriverDoesNotSupportViewListing is returned when a driver does not support view listing.
	ErrDriverDoesNotSupportViewListing = errors.New("driver does not support view listing")

	// ErrDriverDoesNotSupportMaterializedViewListing is returned when a driver does not support materialized view listing.
	ErrDriverDoesNotSupportMaterializedViewListing = errors.New("driver does not support materialized view listing")

	// ErrDriverDoesNotSupportColumnListing is returned when a driver does not support column listing.
	ErrDriverDoesNotSupportColumnListing = errors.New("driver does not support column listing")

	// ErrDriverDoesNotSupportIndexListing is returned when a driver does not support index listing.
	ErrDriverDoesNotSupportIndexListing = errors.New("driver does not support index listing")

	// ErrDriverDoesNotSupportQueryExecution is returned when a driver does not support query execution.
	ErrDriverDoesNotSupportQueryExecution = errors.New("driver does not support query execution")
)
