package testutil

import (
	"testing"

	"github.com/polyclient/polyclient/internal/version"
)

// MockVersionProd mocks the version.Version() function to return "v1.2.3".
func MockVersionProd(t *testing.T) {
	t.Helper()

	origVersion := version.Version()
	version.SetVersion("v1.2.3")

	t.Cleanup(func() {
		version.SetVersion(origVersion)
	})
}

// MockVersionDev mocks the version.Version() function to return "dev".
func MockVersionDev(t *testing.T) {
	t.Helper()

	origVersion := version.Version()
	version.SetVersion("dev")

	t.Cleanup(func() {
		version.SetVersion(origVersion)
	})
}
