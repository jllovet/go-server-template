package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jllovet/go-server-template/config"
	"github.com/jllovet/go-server-template/internal/server"
	"github.com/jllovet/go-server-template/internal/todo"
	"github.com/jllovet/go-server-template/internal/todo/memory"
)

func TestIntegration_TodoWorkflow(t *testing.T) {
	// 1. Setup dependencies (similar to main.go)
	// We use the real in-memory repository and service.
	cfg := &config.Config{
		Host: "localhost",
		Port: "8080", // Not actually used by httptest, but required by struct
	}
	// Discard logs during tests to keep output clean
	logger := log.New(io.Discard, "", 0)

	repo := memory.New()
	service := todo.NewService(repo)
	srv := server.NewServer(service, cfg, logger)

	// 2. Create a test server
	// httptest.NewServer starts a real HTTP server on a random port
	ts := httptest.NewServer(srv)
	defer ts.Close()

	client := ts.Client()
	baseURL := ts.URL

	// Helper to make requests
	request := func(method, path string, body interface{}) (*http.Response, error) {
		var bodyReader io.Reader
		if body != nil {
			b, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			bodyReader = bytes.NewReader(b)
		}
		req, err := http.NewRequest(method, baseURL+path, bodyReader)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		return client.Do(req)
	}

	var createdID string

	// 3. Run the workflow

	t.Run("1. Create Todo", func(t *testing.T) {
		resp, err := request("POST", "/todos", map[string]string{"title": "Integration Test"})
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("expected 201 Created, got %d", resp.StatusCode)
		}

		var created todo.Todo
		if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if created.ID == "" {
			t.Error("expected non-empty ID")
		}
		if created.Title != "Integration Test" {
			t.Errorf("expected title 'Integration Test', got %q", created.Title)
		}
		createdID = created.ID
	})

	t.Run("2. List Todos", func(t *testing.T) {
		resp, err := request("GET", "/todos", nil)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		var todos []todo.Todo
		if err := json.NewDecoder(resp.Body).Decode(&todos); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if len(todos) != 1 {
			t.Errorf("expected 1 todo, got %d", len(todos))
		}
	})

	t.Run("3. Update Title", func(t *testing.T) {
		resp, err := request("PATCH", "/todos/"+createdID, map[string]string{"title": "Updated Title"})
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200 OK, got %d", resp.StatusCode)
		}

		var updated todo.Todo
		if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if updated.Title != "Updated Title" {
			t.Errorf("expected title 'Updated Title', got %q", updated.Title)
		}
	})

	t.Run("4. Mark Complete", func(t *testing.T) {
		resp, err := request("POST", "/todos/"+createdID+"/complete", nil)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200 OK, got %d", resp.StatusCode)
		}
	})

	t.Run("5. Mark Incomplete", func(t *testing.T) {
		resp, err := request("POST", "/todos/"+createdID+"/incomplete", nil)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200 OK, got %d", resp.StatusCode)
		}
	})

	t.Run("6. Delete Todo", func(t *testing.T) {
		resp, err := request("DELETE", "/todos/"+createdID, nil)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			t.Errorf("expected 204 No Content, got %d", resp.StatusCode)
		}
	})
}
