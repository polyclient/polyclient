package validator

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
)

// CustomValidator is a wrapper for go-playground/validator that includes custom validators.
// Useful to reuse the same validators across the different parts of the application (e.g. internals, API, CLI).
type CustomValidator struct {
	validate *validator.Validate
}

// NewCustomValidator returns a new CustomValidator instance.
func NewCustomValidator() (*CustomValidator, error) {
	v := validator.New(validator.WithRequiredStructEnabled())

	custom := &CustomValidator{
		validate: v,
	}

	if err := custom.includeCustomValidators(); err != nil {
		return nil, fmt.Errorf("failed to include custom validators: %w", err)
	}

	return custom, nil
}

// Validate validates a struct using the included custom validators.
func (v *CustomValidator) Validate(data any) error {
	return v.validate.Struct(data)
}

func (v *CustomValidator) includeCustomValidators() error {
	customValidators := []struct {
		tag string
		fn  validator.Func
	}{
		{
			tag: "idPattern",
			fn: func(fl validator.FieldLevel) bool {
				return regexp.MustCompile(`^[a-z0-9_-]+$`).MatchString(fl.Field().String())
			},
		},
		{
			tag: "namePattern",
			fn: func(fl validator.FieldLevel) bool {
				return regexp.MustCompile(`^[a-zA-Z0-9._-]+$`).MatchString(fl.Field().String())
			},
		},
		{
			tag: "driverNamePattern",
			fn: func(fl validator.FieldLevel) bool {
				supportedDrivers := []string{"postgres", "mysql", "mongodb", "sqlite"}
				driver := strings.ToLower(fl.Field().String())

				return slices.Contains(supportedDrivers, driver)
			},
		},
		{
			tag: "entrypointPattern",
			fn: func(fl validator.FieldLevel) bool {
				return regexp.MustCompile(`^[a-zA-Z0-9._/-]+\.wasm$`).MatchString(fl.Field().String())
			},
		},
	}

	for _, cv := range customValidators {
		if err := v.validate.RegisterValidation(cv.tag, cv.fn); err != nil {
			return fmt.Errorf("failed to register custom validator: %w", err)
		}
	}

	return nil
}
