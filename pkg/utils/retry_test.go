package utils_test

import (
	"testing"
	"time"

	"github.com/polyclient/polyclient/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRetry(t *testing.T) {
	tests := []struct {
		name          string
		fn            func() error
		maxAttempts   int
		delayDuration time.Duration
		expectedErr   error
		errorMsg      string
	}{
		{
			name: "success on first attempt",
			fn: func() error {
				return nil
			},
			maxAttempts:   5,
			delayDuration: time.Millisecond,
			expectedErr:   nil,
			errorMsg:      "",
		},
		{
			name: "failure after max attempts",
			fn: func() error {
				return assert.AnError
			},
			maxAttempts:   3,
			delayDuration: time.Millisecond,
			expectedErr:   assert.AnError,
			errorMsg:      "operation failed after 3 attempts: " + assert.AnError.Error(),
		},
		{
			name: "zero attempts should fail immediately",
			fn: func() error {
				return nil
			},
			maxAttempts:   0,
			delayDuration: time.Millisecond,
			expectedErr:   assert.AnError,
			errorMsg:      "invalid number of attempts: 0. Must be greater than 0",
		},
		{
			name: "negative attempts should fail immediately",
			fn: func() error {
				return nil
			},
			maxAttempts:   -1,
			delayDuration: time.Millisecond,
			expectedErr:   assert.AnError,
			errorMsg:      "invalid number of attempts: -1. Must be greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := utils.Retry(tt.fn, tt.maxAttempts, tt.delayDuration)

			if err != tt.expectedErr {
				require.Error(t, err)
				assert.EqualError(t, err, tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRetryWithBackoff(t *testing.T) {
	tests := []struct {
		name          string
		fn            func() error
		maxAttempts   int
		initialDelay  time.Duration
		backoffFactor float64
		expectedErr   error
		errorMsg      string
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := utils.RetryWithBackoff(tt.fn, tt.maxAttempts, tt.initialDelay, tt.backoffFactor)

			if err != tt.expectedErr {
				require.Error(t, err)
				assert.EqualError(t, err, tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
