package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/polyclient/polyclient/internal/validator"
)

// ConnectionProfile defines a saved database connection profile with metadata and config.
type ConnectionProfile struct {
	Driver               string           `json:"driver" validate:"required"`                  // e.g., "postgres", "mongodb", "mysql", etc.
	Name                 string           `json:"name" validate:"required,max=30,namePattern"` // user-defined name for the connection
	ColorTag             string           `json:"colorTag" validate:"omitempty,iscolor"`       // optional UI color tag to distinguish profiles
	SaveCreds            bool             `json:"saveCreds"`                                   // whether to persist credentials (encrypted)
	Pinned               bool             `json:"pinned"`                                      // whether the profile is pinned in UI
	ConfirmBeforeConnect bool             `json:"confirmBeforeConnect"`                        // whether to confirm before connecting (prevents accidental connects to prod)
	CreatedAt            time.Time        `json:"createdAt"`                                   // profile creation time
	LastUsedAt           time.Time        `json:"lastUsedAt"`                                  // last time profile was used
	Config               ConnectionConfig `json:"config" validate:"required,dive"`             // actual driver-specific config
}

// ConnectionStore defines the interface for persisting connection configurations.
// Implementations MUST handle encryption/decryption of sensitive config values (e.g., passwords in DSNs).
type ConnectionStore interface {
	// SaveProfile stores or updates the configuration for a given name.
	SaveProfile(ctx context.Context, profile *ConnectionProfile) error

	// ListProfiles returns a list of all saved connection profiles.
	ListProfiles(ctx context.Context) ([]*ConnectionProfile, error)

	// ListRecentProfiles returns a list of recently used connection profiles.
	// The threshold specifies the maximum age of profiles to return.
	ListRecentProfiles(ctx context.Context, threshold time.Duration) ([]*ConnectionProfile, error)

	// GetProfile returns the configuration for a given name.
	GetProfile(ctx context.Context, name string) (*ConnectionProfile, error)

	// DeleteProfile removes the configuration for a given name.
	DeleteProfile(ctx context.Context, name string) error
}

// FileConnectionStore is a file-based connection store implementation.
type FileConnectionStore struct {
	path string
}

// NewFileConnectionStore creates a new instance of FileConnectionStore with the specified file path.
func NewFileConnectionStore(path string) *FileConnectionStore {
	return &FileConnectionStore{path: path}
}

// SaveProfile implements ConnectionStore.SaveProfile.
func (s *FileConnectionStore) SaveProfile(ctx context.Context, profile *ConnectionProfile) error {
	v, err := validator.NewCustomValidator()
	if err != nil {
		return fmt.Errorf("failed to create validator: %w", err)
	}

	if err := v.Validate(profile); err != nil {
		return fmt.Errorf("invalid profile: %w", err)
	}

	profilePath := filepath.Join(s.path, profile.Name+".json")

	absPath, err := filepath.Abs(profilePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	cleanPath := filepath.Clean(absPath)

	if err := os.MkdirAll(filepath.Dir(cleanPath), 0o0700); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	profileBytes, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal profile: %w", err)
	}

	// Write to temp file first
	tempFile := cleanPath + ".tmp"
	if err = os.WriteFile(tempFile, profileBytes, 0o0600); err != nil {
		return fmt.Errorf("failed to write profile: %w", err)
	}

	// Atomically rename to final destination
	if err = os.Rename(tempFile, cleanPath); err != nil {
		_ = os.Remove(tempFile)
		return fmt.Errorf("failed to rename profile: %w", err)
	}

	return nil
}

// ListProfiles implements ConnectionStore.List.
func (s *FileConnectionStore) ListProfiles(ctx context.Context) ([]*ConnectionProfile, error) {
	absPath, err := filepath.Abs(s.path)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	cleanPath := filepath.Clean(absPath)

	entries, err := os.ReadDir(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []*ConnectionProfile{}, nil
		}

		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	profiles := make([]*ConnectionProfile, 0, len(entries))

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		profilePath := filepath.Join(cleanPath, entry.Name())

		profile, err := s.readProfileFile(profilePath)
		if err != nil {
			log.Printf("Failed to read profile %s: %s", entry.Name(), err)
			continue
		}

		profiles = append(profiles, profile)
	}

	return profiles, nil
}

// ListRecentProfiles implements ConnectionStore.ListRecentProfiles.
func (s *FileConnectionStore) ListRecentProfiles(ctx context.Context, threshold time.Duration) ([]*ConnectionProfile, error) {
	profiles, err := s.ListProfiles(ctx)
	if err != nil {
		return nil, err
	}

	cutoff := time.Now().Add(-threshold)

	var recent []*ConnectionProfile

	for _, profile := range profiles {
		if profile.LastUsedAt.After(cutoff) {
			recent = append(recent, profile)
		}
	}

	return recent, nil
}

// GetProfile implements ConnectionStore.GetProfile.
func (s *FileConnectionStore) GetProfile(ctx context.Context, name string) (*ConnectionProfile, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("name cannot be empty")
	}

	profilePath := filepath.Join(s.path, name+".json")

	absPath, err := filepath.Abs(profilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	profile, err := s.readProfileFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read profile: %w", err)
	}

	return profile, nil
}

// DeleteProfile implements ConnectionStore.DeleteProfile.
func (s *FileConnectionStore) DeleteProfile(ctx context.Context, name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("name cannot be empty")
	}

	profilePath := filepath.Join(s.path, name+".json")

	absPath, err := filepath.Abs(profilePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	if err := os.Remove(filepath.Clean(absPath)); err != nil {
		if os.IsNotExist(err) {
			return nil // No error if file doesn't exist
		}

		return fmt.Errorf("failed to delete profile: %w", err)
	}

	return nil
}

// readProfileFile reads and unmarshals a profile from a file.
func (s *FileConnectionStore) readProfileFile(path string) (*ConnectionProfile, error) {
	profileBytes, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("failed to read profile: %w", err)
	}

	var profile ConnectionProfile
	if err := json.Unmarshal(profileBytes, &profile); err != nil {
		return nil, fmt.Errorf("failed to unmarshal profile: %w", err)
	}

	return &profile, nil
}
