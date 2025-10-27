package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"modernfi-treasury-app/internal/services"
)

// YieldHandler handles HTTP requests for yield data
type YieldHandler struct {
	treasuryService *services.TreasuryService
}

// NewYieldHandler creates a new YieldHandler with the provided TreasuryService
func NewYieldHandler(treasuryService *services.TreasuryService) *YieldHandler {
	return &YieldHandler{
		treasuryService: treasuryService,
	}
}

// GetYields handles GET requests to fetch the latest treasury yields
func (h *YieldHandler) GetYields(w http.ResponseWriter, r *http.Request) {
	// Fetch latest yields from the treasury service
	yieldData, err := h.treasuryService.GetLatestYields()
	if err != nil {
		// Log the error for debugging
		log.Printf("Error fetching treasury yields: %v", err)

		// Return 500 Internal Server Error with error message
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to fetch treasury data",
		})
		return
	}

	// Set content type and return successful response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(yieldData)
}

// GetHistoricalYields handles GET requests to /api/yields/historical
// Query parameter: period (1W, 1M, 3M, 6M, 1Y, 5Y, 10Y, 30Y) - defaults to 3M
func (h *YieldHandler) GetHistoricalYields(w http.ResponseWriter, r *http.Request) {
	// Parse query parameter
	period := r.URL.Query().Get("period")
	if period == "" {
		period = "3M" // Default to 3 months
	}

	// Validate period
	validPeriods := map[string]bool{
		"1W":  true,
		"1M":  true,
		"3M":  true,
		"6M":  true,
		"1Y":  true,
		"5Y":  true,
		"10Y": true,
		"30Y": true,
	}
	if !validPeriods[period] {
		log.Printf("Invalid period requested: %s", period)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid period. Must be one of: 1W, 1M, 3M, 6M, 1Y, 5Y, 10Y, 30Y",
		})
		return
	}

	// Fetch historical yields
	data, err := h.treasuryService.GetHistoricalYields(period)
	if err != nil {
		log.Printf("Error fetching historical yields: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to fetch historical treasury data",
		})
		return
	}

	// Return successful response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
