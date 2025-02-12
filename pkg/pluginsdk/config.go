package pluginsdk

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"

	"github.com/go-playground/validator/v10"
)

var (
	supportedConfigFiles = [...]string{"polyclient.json"}
)

type Config struct {
	Manifest *Manifest `json:"manifest" validate:"required"`
	Runtime  *Runtime  `json:"runtime"  validate:"required"`
}

type Manifest struct {
	Name        string   `json:"name"        validate:"required,namePattern"`
	Version     string   `json:"version"     validate:"required,semver"`
	DisplayName string   `json:"displayName"`
	Description string   `json:"description"`
	Author      string   `json:"author"`
	License     string   `json:"license"`
	Repository  string   `json:"repository"  validate:"url"`
	Categories  []string `json:"categories"  validate:"dive,oneof=Database 'Data Management' 'Data Analysis' AI Cloud Security Development Testing 'User Interface' Performance Monitoring Networking Integration Utilities Other"`
	EntryPoint  string   `json:"entryPoint"  validate:"required"`
}

type Runtime struct {
	OS []string `json:"os" validate:"required,dive,oneof=linux darwin windows"`
}

// FindConfigPath searches for the first plugin config file under rootPath
// and returns its path. If rootPath itself is a plugin config file, the
// path is returned immediately.
func FindConfigPath(rootPath string) (string, error) {
	info, err := os.Stat(rootPath)
	if err != nil {
		return "", fmt.Errorf("failed to find plugin config: %w", err)
	}

	if CheckConfigFile(info) {
		return rootPath, nil
	}

	var foundConfig string

	err = filepath.Walk(rootPath, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if CheckConfigFile(info) {
			foundConfig = p
			return filepath.SkipAll
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	if foundConfig == "" {
		return "", fmt.Errorf("failed to find plugin config in %s", rootPath)
	}

	return foundConfig, nil
}

// FindConfigPaths searches for all plugin config files under rootPath and
// returns their paths. If rootPath itself is a plugin config file, it is
// returned immediately in a slice.
func FindConfigPaths(rootPath string) ([]string, error) {
	info, err := os.Stat(rootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to find plugin config: %w", err)
	}

	foundConfigs := []string{}

	if CheckConfigFile(info) {
		return []string{rootPath}, nil
	}

	err = filepath.Walk(rootPath, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if CheckConfigFile(info) {
			foundConfigs = append(foundConfigs, p)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	if len(foundConfigs) == 0 {
		return nil, fmt.Errorf("failed to find plugin configs in %s", rootPath)
	}

	return foundConfigs, nil
}

// CheckConfigFile returns true if the given file info belongs to a plugin config file.
func CheckConfigFile(info os.FileInfo) bool {
	return !info.IsDir() && slices.Contains(supportedConfigFiles[:], filepath.Base(info.Name()))
}

// ValidateConfig checks if the given config is a valid plugin config.
func ValidateConfig(m *Config) error {
	v := validator.New(validator.WithRequiredStructEnabled())

	if err := v.RegisterValidation("namePattern", func(fl validator.FieldLevel) bool {
		return regexp.MustCompile(`^[a-z0-9_-]+$`).MatchString(fl.Field().String())
	}); err != nil {
		return fmt.Errorf("invalid plugin config: %w", err)
	}

	if err := v.Struct(m); err != nil {
		return fmt.Errorf("invalid plugin config: %w", err)
	}

	return nil
}

// LoadConfig loads the plugin config from the given path.
// If the config is invalid, an error is returned.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if err := ValidateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
