// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package polling

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

// UntilReady checks if the server at targetURL is available by sending a GET request
// and retrying until success or timeout.
func UntilReady(ctx context.Context, targetURL string) error {
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}

	retryInterval := 100 * time.Millisecond
	retries := 0

	for {
		retries++

		slog.Info("Polling for server readiness", "url", targetURL, "retries", retries)

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			req, err := http.NewRequestWithContext(ctx, http.MethodHead, targetURL, http.NoBody)
			if err != nil {
				return err
			}

			resp, err := httpClient.Do(req)
			if err != nil {
				time.Sleep(retryInterval)
				continue
			}

			if resp.StatusCode == http.StatusOK {
				return nil
			}

			time.Sleep(retryInterval)

			_ = resp.Body.Close()
		}
	}
}
