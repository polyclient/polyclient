package pluginsdk

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/go-playground/validator/v10"
)

type Manifest struct {
	Name        string   `json:"name" validate:"required,namePattern"`
	Version     string   `json:"version" validate:"required,semver"`
	DisplayName string   `json:"displayName"`
	Description string   `json:"description"`
	Author      string   `json:"author"`
	License     string   `json:"license"`
	Repository  string   `json:"repository" validate:"url"`
	Categories  []string `json:"categories" validate:"dive,oneof=Database Data Data Data AI Theme Cloud Service Storage Security Development Testing Logger Other"`
	EntryPoint  string   `json:"entryPoint" validate:"required"`
}

// FindManifest searches for the first manifest.json file found under rootPath.
// If rootPath itself is a manifest.json file, it is returned immediately.
func FindManifest(rootPath string) (string, error) {
	info, err := os.Stat(rootPath)
	if err != nil {
		return "", fmt.Errorf("failed to lookup plugin manifest: %w", err)
	}

	// If the given path is a file and it's manifest.json, return it.
	if !info.IsDir() && filepath.Base(rootPath) == "manifest.json" {
		return rootPath, nil
	}

	var foundManifest string
	err = filepath.Walk(rootPath, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// If a file named manifest.json is found, capture it and stop the walk.
		if !info.IsDir() && filepath.Base(p) == "manifest.json" {
			foundManifest = p
			return filepath.SkipAll
		}

		return nil
	})
	if err != nil {
		return "", err
	}
	if foundManifest == "" {
		return "", fmt.Errorf("manifest.json not found in %s", rootPath)
	}

	return foundManifest, nil
}

// FindAllManifests searches for all manifest.json files under rootPath.
// If rootPath itself is a manifest.json file, it returns a slice containing just that file.
func FindAllManifests(rootPath string) ([]string, error) {
	info, err := os.Stat(rootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup plugin manifests: %w", err)
	}

	manifests := []string{}

	// If the given path is a file and it's manifest.json, return it.
	if !info.IsDir() && filepath.Base(rootPath) == "manifest.json" {
		return []string{rootPath}, nil
	}

	err = filepath.Walk(rootPath, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Base(p) == "manifest.json" {
			manifests = append(manifests, p)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(manifests) == 0 {
		return nil, fmt.Errorf("no manifest.json files found in %s", rootPath)
	}

	return manifests, nil
}

// ValidateManifest validates the plugin manifest.
func ValidateManifest(m *Manifest) error {
	v := validator.New(validator.WithRequiredStructEnabled())

	if err := v.RegisterValidation("namePattern", func(fl validator.FieldLevel) bool {
		return regexp.MustCompile(`^[a-z0-9_-]+$`).MatchString(fl.Field().String())
	}); err != nil {
		return fmt.Errorf("invalid plugin manifest: %w", err)
	}

	if err := v.Struct(m); err != nil {
		return fmt.Errorf("invalid plugin manifest: %w", err)
	}

	return nil
}

// LoadManifest loads the plugin manifest from the given path.
// If the manifest is invalid, an error is returned.
func LoadManifest(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}

	if err := ValidateManifest(&manifest); err != nil {
		return nil, err
	}

	return &manifest, nil
}
