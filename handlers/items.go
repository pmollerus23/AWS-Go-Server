package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
)

// Item represents an item in our system.
type Item struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// In-memory store for demo purposes
var (
	items     = make(map[int64]Item)
	itemsMux  sync.RWMutex // Protects items map and nextID
	nextID    int64 = 1
)

// HandleItemsGet returns a handler that retrieves all items.
func HandleItemsGet(logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		itemsMux.RLock()
		itemsCount := len(items)
		// Convert map to slice
		itemsList := make([]Item, 0, itemsCount)
		for _, item := range items {
			itemsList = append(itemsList, item)
		}
		itemsMux.RUnlock()

		logger.Info("retrieving all items", "count", itemsCount)

		if err := encode(w, r, http.StatusOK, itemsList); err != nil {
			logger.Error("failed to encode response", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})
}

// CreateItemRequest represents the request to create an item.
type CreateItemRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CreateItemResponse represents the response after creating an item.
type CreateItemResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// HandleItemsCreate returns a handler that creates a new item.
func HandleItemsCreate(logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, problems, err := decodeValid[CreateItemRequest](r)
		if err != nil {
			logger.Error("failed to decode request", "error", err)
			if len(problems) > 0 {
				encode(w, r, http.StatusBadRequest, map[string]interface{}{
					"error":    "validation failed",
					"problems": problems,
				})
				return
			}
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// Create the item (protected by write lock)
		itemsMux.Lock()
		id := nextID
		nextID++
		item := Item{
			ID:          id,
			Name:        req.Name,
			Description: req.Description,
		}
		items[id] = item
		itemsMux.Unlock()

		logger.Info("item created", "id", id, "name", req.Name)

		resp := CreateItemResponse{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
		}

		if err := encode(w, r, http.StatusCreated, resp); err != nil {
			logger.Error("failed to encode response", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})
}

// Valid implements the Validator interface for CreateItemRequest.
func (r CreateItemRequest) Valid(ctx context.Context) map[string]string {
	problems := make(map[string]string)

	if r.Name == "" {
		problems["name"] = "name is required and cannot be empty"
	}
	if len(r.Name) > 100 {
		problems["name"] = "name must be 100 characters or less"
	}
	if len(r.Description) > 500 {
		problems["description"] = "description must be 500 characters or less"
	}

	return problems
}

// Helper functions for encoding/decoding
// These are duplicated here for the handlers package
// In a real app, you might put these in a shared package

func encode[T any](w http.ResponseWriter, r *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func decodeValid[T Validator](r *http.Request) (T, map[string]string, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, nil, err
	}
	if problems := v.Valid(r.Context()); len(problems) > 0 {
		return v, problems, fmt.Errorf("invalid %T: %d problems", v, len(problems))
	}
	return v, nil, nil
}

type Validator interface {
	Valid(ctx context.Context) map[string]string
}
