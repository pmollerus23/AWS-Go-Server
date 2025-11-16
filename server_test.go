package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

// waitForReady polls the health endpoint until it returns 200 or times out.
func waitForReady(
	ctx context.Context,
	timeout time.Duration,
	endpoint string,
) error {
	client := http.Client{}
	startTime := time.Now()
	for {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}
		resp, err := client.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return nil
		}
		if resp != nil {
			resp.Body.Close()
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if time.Since(startTime) >= timeout {
				return fmt.Errorf("timeout reached while waiting for endpoint")
			}
			time.Sleep(250 * time.Millisecond)
		}
	}
}

func TestServer(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)

	// Start the server in a goroutine
	go func() {
		if err := run(ctx, io.Discard, []string{"testapp"}); err != nil {
			t.Logf("server error: %v", err)
		}
	}()

	// Wait for server to be ready
	if err := waitForReady(ctx, 5*time.Second, "http://localhost:8080/healthz"); err != nil {
		t.Fatalf("server never became ready: %v", err)
	}

	t.Run("health check", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/healthz")
		if err != nil {
			t.Fatalf("failed to call healthz: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("failed to read response body: %v", err)
		}

		if string(body) != "OK" {
			t.Errorf("expected body 'OK', got %q", string(body))
		}
	})

	t.Run("create and get items", func(t *testing.T) {
		// Create an item
		createReq := map[string]string{
			"name":        "Test Item",
			"description": "This is a test item",
		}
		body, err := json.Marshal(createReq)
		if err != nil {
			t.Fatalf("failed to marshal request: %v", err)
		}

		resp, err := http.Post(
			"http://localhost:8080/api/v1/items",
			"application/json",
			bytes.NewReader(body),
		)
		if err != nil {
			t.Fatalf("failed to create item: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("expected status 201, got %d", resp.StatusCode)
		}

		var createResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if createResp["name"] != "Test Item" {
			t.Errorf("expected name 'Test Item', got %v", createResp["name"])
		}

		// Get all items
		resp, err = http.Get("http://localhost:8080/api/v1/items")
		if err != nil {
			t.Fatalf("failed to get items: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		var items []map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if len(items) < 1 {
			t.Error("expected at least 1 item")
		}
	})

	t.Run("validation error", func(t *testing.T) {
		// Try to create an item with empty name
		createReq := map[string]string{
			"name":        "",
			"description": "This should fail",
		}
		body, err := json.Marshal(createReq)
		if err != nil {
			t.Fatalf("failed to marshal request: %v", err)
		}

		resp, err := http.Post(
			"http://localhost:8080/api/v1/items",
			"application/json",
			bytes.NewReader(body),
		)
		if err != nil {
			t.Fatalf("failed to create item: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", resp.StatusCode)
		}

		var errResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if errResp["error"] != "validation failed" {
			t.Errorf("expected validation error, got %v", errResp)
		}
	})

	t.Run("not found", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/nonexistent")
		if err != nil {
			t.Fatalf("failed to call endpoint: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", resp.StatusCode)
		}
	})
}
