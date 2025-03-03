// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package plugin

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"

	"github.com/go-playground/validator/v10"
)

var supportedManifestFiles = [...]string{"manifest.json"}

type Manifest struct {
	Id          string   `json:"id"          validate:"required,idPattern"`
	Name        string   `json:"name"        validate:"required,gte=1,lte=30"`
	Scope       string   `json:"scope"       validate:"required,oneof=sql nosql ai theme"`
	Version     string   `json:"version"     validate:"required,semver"`
	Description string   `json:"Description" validate:"lte=300"`
	Authors     []Author `json:"authors"     validate:"min=1,dive"`
	License     string   `json:"license"`
	Repository  string   `json:"repository"  validate:"uri"`
	Homepage    string   `json:"homepage"    validate:"uri"`
	Keywords    []string `json:"keywords"    validate:"max=5,unique"`
	Entrypoint  string   `json:"entrypoint"  validate:"required,entrypointPattern"`
}

type Author struct {
	Name     string   `json:"name"     validate:"required"`
	Email    string   `json:"email"    validate:"omitempty,email"`
	Websites []string `json:"websites" validate:"omitempty,min=1,max=5,unique,dive,uri"`
}

// ValidateManifest verifies that a give manifest configuration
// complies with all the defined validation rules.
func ValidateManifest(m *Manifest) error {
	v := validator.New(validator.WithRequiredStructEnabled())

	// Register custom validation for id pattern
	if err := v.RegisterValidation("idPattern", func(fl validator.FieldLevel) bool {
		return regexp.MustCompile(`^[a-z0-9_-]+$`).MatchString(fl.Field().String())
	}); err != nil {
		return fmt.Errorf("invalid plugin manifest: %w", err)
	}

	// Register custom validation for entrypoint pattern
	if err := v.RegisterValidation("entrypointPattern", func(fl validator.FieldLevel) bool {
		return regexp.MustCompile(`^[a-zA-Z0-9._/-]+\.wasm$`).MatchString(fl.Field().String())
	}); err != nil {
		return fmt.Errorf("invalid plugin manifest: %w", err)
	}

	if err := v.Struct(m); err != nil {
		return fmt.Errorf("invalid plugin manifest: %w", err)
	}

	return nil
}

// LoadManifest reads and validates a plugin manifest from the
// specified file path.
func LoadManifest(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load manifest: %w", err)
	}

	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to load manifest: %w", err)
	}

	if err := ValidateManifest(&manifest); err != nil {
		return nil, fmt.Errorf("failed to load manifest: %w", err)
	}

	return &manifest, nil
}

// FindManifestPath searches for the first plugin manifest file
// under a lookupPath and returns its path. If lookupPath itself
// is a plugin manifest, it returns it immediately.
func FindManifestPath(lookupPath string) (string, error) {
	info, err := os.Stat(lookupPath)
	if err != nil {
		return "", fmt.Errorf("failed to find plugin's manifest: %w", err)
	}

	if isPluginManifest(info) {
		return lookupPath, nil
	}

	var foundPath string

	err = filepath.Walk(lookupPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if isPluginManifest(info) {
			foundPath = path
			return filepath.SkipAll
		}

		return nil
	})
	if err != nil || foundPath == "" {
		return "", fmt.Errorf("failed to find plugin manifest in %s", lookupPath)
	}

	return foundPath, nil
}

// FindManifestPaths searches for all the plugin manifest files
// under a lookupPath and returns its paths. If lookupPath itself
// is a plugin manifest, it returns it immediately in a slice.
func FindManifestPaths(lookupPath string) ([]string, error) {
	info, err := os.Stat(lookupPath)
	if err != nil {
		return nil, fmt.Errorf("failed to find plugins manifests: %w", err)
	}

	if isPluginManifest(info) {
		return []string{lookupPath}, nil
	}

	foundPaths := []string{}

	err = filepath.Walk(lookupPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if isPluginManifest(info) {
			foundPaths = append(foundPaths, path)
		}

		return nil
	})
	if err != nil || len(foundPaths) == 0 {
		return []string{}, fmt.Errorf("failed to find plugins manifests in %s", err)
	}

	return foundPaths, nil
}

func isPluginManifest(info os.FileInfo) bool {
	return !os.ModeAppend.IsDir() && slices.Contains(supportedManifestFiles[:], filepath.Base(info.Name()))
}
