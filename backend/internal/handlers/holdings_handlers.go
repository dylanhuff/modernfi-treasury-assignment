package handlers

import (
	"encoding/json"
	"log"
	"math/big"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"modernfi-treasury-app/internal/database"
)

// HoldingsHandlers handles HTTP requests for holdings operations.
type HoldingsHandlers struct {
	queries *database.Queries
}

// NewHoldingsHandlers creates and returns a new HoldingsHandlers instance.
func NewHoldingsHandlers(queries *database.Queries) *HoldingsHandlers {
	return &HoldingsHandlers{
		queries: queries,
	}
}

// GetUserHoldings handles GET /api/v1/users/{id}/holdings requests.
// Returns all holdings for the specified user where remaining_amount > 0.
// Holdings are ordered by purchase_date DESC (most recent first).
func (h *HoldingsHandlers) GetUserHoldings(w http.ResponseWriter, r *http.Request) {
	// Parse user ID from URL parameter
	userIDStr := chi.URLParam(r, "id")
	userID, err := strconv.ParseInt(userIDStr, 10, 32)
	if err != nil {
		log.Printf("Invalid user ID: %s", userIDStr)
		respondWithError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	// Fetch all holdings for user using existing sqlc query
	holdings, err := h.queries.GetHoldingsByUser(r.Context(), int32(userID))
	if err != nil {
		log.Printf("Error fetching holdings for user %d: %v", userID, err)
		respondWithError(w, http.StatusInternalServerError, "failed to fetch holdings")
		return
	}

	// Filter holdings to only include those with remaining_amount > 0
	// Also handle legacy data by providing fallback values
	activeHoldings := []database.Holding{}
	zero := big.NewInt(0)
	for _, holding := range holdings {
		// Check if remaining_amount is valid and > 0
		if holding.RemainingAmount.Valid && holding.RemainingAmount.Int.Cmp(zero) > 0 {
			activeHoldings = append(activeHoldings, holding)
		}
	}

	// Return active holdings (empty array if no holdings with remaining_amount > 0)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(activeHoldings); err != nil {
		log.Printf("Error encoding holdings response: %v", err)
	}
}
