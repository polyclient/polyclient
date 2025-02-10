package utils

import (
	"fmt"
	"time"
)

// Retry repeatedly attempts to execute a specified function up to a given number of times.
// It pauses for a specified duration between each attempt. The function stops retrying
// if it succeeds (returns nil) or if the maximum number of attempts is reached.
// Returns an error if the operation fails after the specified attempts.
func Retry(fn func() error, maxAttempts int, delayDuration time.Duration) error {
	if maxAttempts <= 0 {
		return fmt.Errorf("invalid number of attempts: %d. Must be greater than 0", maxAttempts)
	}

	var err error

	for i := 0; i < maxAttempts; i++ {
		err = fn()

		if err == nil {
			return nil
		}

		if i < maxAttempts-1 {
			time.Sleep(delayDuration)
		}
	}

	return fmt.Errorf("operation failed after %d attempts: %w", maxAttempts, err)
}

// RetryWithBackoff is a utility function that retries a given function with an exponential backoff.
// It starts with an initial delay and multiplies it by a backoff factor after each failed attempt.
// The function stops retrying if it succeeds or if the maximum number of attempts is reached.
// Returns an error if the operation fails after the specified attempts.
func RetryWithBackoff(fn func() error, maxAttempts int, initialDelay time.Duration, backoffFactor float64) error {
	return Retry(fn, maxAttempts, initialDelay*time.Duration(backoffFactor))
}
