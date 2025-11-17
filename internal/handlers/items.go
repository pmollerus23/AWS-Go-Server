package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
)

// Item represents an item in our system.
type Item struct {
	ID          int64  `json:"id" example:"1"`
	Name        string `json:"name" example:"Sample Item"`
	Description string `json:"description" example:"This is a sample item description"`
}

// In-memory store for demo purposes
var (
	items    = make(map[int64]Item)
	itemsMux sync.RWMutex // Protects items map and nextID
	nextID   int64        = 1
)

// HandleItemsGet returns a handler that retrieves all items.
//
//	@Summary		List all items
//	@Description	Get a list of all items in the system
//	@Tags			items
//	@Produce		json
//	@Success		200	{array}		Item
//	@Failure		401	{string}	string	"Unauthorized"
//	@Failure		500	{string}	string	"Internal Server Error"
//	@Security		BearerAuth
//	@Router			/api/v1/items [get]
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
	Name        string `json:"name" example:"New Item" minLength:"1" maxLength:"100"`
	Description string `json:"description" example:"Item description" maxLength:"500"`
}

// CreateItemResponse represents the response after creating an item.
type CreateItemResponse struct {
	ID          int64  `json:"id" example:"1"`
	Name        string `json:"name" example:"New Item"`
	Description string `json:"description" example:"Item description"`
}

// ValidationError represents validation error response
type ValidationError struct {
	Error    string            `json:"error" example:"validation failed"`
	Problems map[string]string `json:"problems"`
}

// HandleItemsCreate returns a handler that creates a new item.
//
//	@Summary		Create a new item
//	@Description	Create a new item with name and description
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Param			item	body		CreateItemRequest	true	"Item to create"
//	@Success		201		{object}	CreateItemResponse
//	@Failure		400		{object}	ValidationError	"Validation error"
//	@Failure		401		{string}	string			"Unauthorized"
//	@Failure		500		{string}	string			"Internal Server Error"
//	@Security		BearerAuth
//	@Router			/api/v1/items [post]
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
