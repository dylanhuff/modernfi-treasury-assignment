package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"modernfi-treasury-app/internal/database"
)

// UserHandler handles HTTP requests related to user operations.
// It uses sqlc-generated queries for type-safe database access.
type UserHandler struct {
	queries *database.Queries
}

// NewUserHandler creates and returns a new UserHandler instance.
// The queries parameter should be initialized with a database connection pool.
func NewUserHandler(queries *database.Queries) *UserHandler {
	return &UserHandler{queries: queries}
}

// GetAllUsers handles GET /api/v1/users requests.
// Returns a JSON array of all users in the system, ordered by name.
// Returns an empty array ([]) if no users exist, never null.
// Returns HTTP 500 with error message if database query fails.
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.queries.ListUsers(r.Context())
	if err != nil {
		log.Printf("Error fetching users: %v", err)
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	// sqlc with emit_empty_slices ensures users is [] not nil
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Printf("Error encoding users: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
