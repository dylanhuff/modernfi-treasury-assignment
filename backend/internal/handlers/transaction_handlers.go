package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"modernfi-treasury-app/internal/database"
	"modernfi-treasury-app/internal/services"
	"modernfi-treasury-app/internal/utils"
)

// TransactionHandlers handles HTTP requests for fund and withdraw operations.
// It uses the TransactionService for atomic database operations.
type TransactionHandlers struct {
	txService       *services.TransactionService
	queries         *database.Queries
	treasuryService *services.TreasuryService
}

// NewTransactionHandlers creates and returns a new TransactionHandlers instance.
func NewTransactionHandlers(
	txService *services.TransactionService,
	queries *database.Queries,
	treasuryService *services.TreasuryService,
) *TransactionHandlers {
	return &TransactionHandlers{
		txService:       txService,
		queries:         queries,
		treasuryService: treasuryService,
	}
}

// TransactionRequest represents the incoming JSON request for fund/withdraw operations
type TransactionRequest struct {
	UserID int32   `json:"user_id"`
	Amount float64 `json:"amount"`
}

// BuyRequest represents the incoming JSON request for buy operations
type BuyRequest struct {
	UserID    int32   `json:"user_id"`
	Term      string  `json:"term"`
	FaceValue float64 `json:"face_value"`
}

// SellRequest represents the incoming JSON request for sell operations
type SellRequest struct {
	UserID    int32   `json:"user_id"`
	HoldingID int32   `json:"holding_id"`
	Amount    float64 `json:"amount"`
}

// TransactionResponse represents the JSON response for fund/withdraw operations
type TransactionResponse struct {
	Success bool        `json:"success"`
	User    interface{} `json:"user,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// FundHandler handles POST /api/v1/fund requests.
// Expects JSON body with user_id and amount fields.
// Returns updated user object on success, or error message on failure.
func (h *TransactionHandlers) FundHandler(w http.ResponseWriter, r *http.Request) {
	var req TransactionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding fund request: %v", err)
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Convert float64 to pgtype.Numeric using string representation
	amount := pgtype.Numeric{}
	if err := amount.Scan(fmt.Sprintf("%.2f", req.Amount)); err != nil {
		log.Printf("Error converting amount to numeric: %v", err)
		respondWithError(w, http.StatusBadRequest, "invalid amount format")
		return
	}

	user, err := h.txService.FundAccount(r.Context(), req.UserID, amount)
	if err != nil {
		log.Printf("Error funding account for user %d: %v", req.UserID, err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, TransactionResponse{
		Success: true,
		User:    user,
	})
}

// WithdrawHandler handles POST /api/v1/withdraw requests.
// Expects JSON body with user_id and amount fields.
// Validates sufficient balance before withdrawal.
// Returns updated user object on success, or error message on failure.
func (h *TransactionHandlers) WithdrawHandler(w http.ResponseWriter, r *http.Request) {
	var req TransactionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding withdraw request: %v", err)
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Convert float64 to pgtype.Numeric using string representation
	amount := pgtype.Numeric{}
	if err := amount.Scan(fmt.Sprintf("%.2f", req.Amount)); err != nil {
		log.Printf("Error converting amount to numeric: %v", err)
		respondWithError(w, http.StatusBadRequest, "invalid amount format")
		return
	}

	user, err := h.txService.WithdrawAccount(r.Context(), req.UserID, amount)
	if err != nil {
		log.Printf("Error withdrawing from account for user %d: %v", req.UserID, err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, TransactionResponse{
		Success: true,
		User:    user,
	})
}

// GetUserTransactions handles GET /api/v1/users/{userId}/transactions requests.
// Returns all transactions for the specified user, ordered by timestamp DESC.
// Supports fund, withdraw, buy, and sell transaction types.
// Used by frontend TransactionHistory component to display transaction table.
// Returns HTTP 400 if user ID is invalid, HTTP 500 for database errors.
func (h *TransactionHandlers) GetUserTransactions(w http.ResponseWriter, r *http.Request) {
	// Parse user ID from URL parameter
	userIDStr := chi.URLParam(r, "userId")
	userID, err := strconv.ParseInt(userIDStr, 10, 32)
	if err != nil {
		log.Printf("Invalid user ID: %s", userIDStr)
		respondWithError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	// Fetch transactions using existing sqlc query
	transactions, err := h.queries.GetTransactionsByUser(r.Context(), int32(userID))
	if err != nil {
		log.Printf("Error fetching transactions for user %d: %v", userID, err)
		respondWithError(w, http.StatusInternalServerError, "failed to fetch transactions")
		return
	}

	// Return transactions (empty array if no transactions)
	respondWithJSON(w, http.StatusOK, transactions)
}

// respondWithJSON is a helper function to send JSON responses with proper headers and status code
func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// respondWithError is a helper function to send error responses in a consistent format
func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	respondWithJSON(w, statusCode, TransactionResponse{
		Success: false,
		Error:   message,
	})
}

// BuyHandler handles POST /api/v1/buy requests.
// Expects JSON body with user_id, term, and face_value fields.
// Fetches current yield data, validates the term, calculates purchase price, and executes the buy operation atomically.
// Returns updated user object with purchase details on success, or error message on failure.
func (h *TransactionHandlers) BuyHandler(w http.ResponseWriter, r *http.Request) {
	var req BuyRequest

	// Decode JSON request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding buy request: %v", err)
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	log.Printf("Buy request received: user_id=%d, term=%s, face_value=%.2f", req.UserID, req.Term, req.FaceValue)

	// Validate term is in allowed list
	validTerms := map[string]bool{
		"1M":  true,
		"3M":  true,
		"6M":  true,
		"1Y":  true,
		"2Y":  true,
		"5Y":  true,
		"10Y": true,
		"30Y": true,
	}

	if !validTerms[req.Term] {
		log.Printf("Invalid term provided: %s", req.Term)
		respondWithError(w, http.StatusBadRequest, "invalid term: must be one of 1M, 3M, 6M, 1Y, 2Y, 5Y, 10Y, 30Y")
		return
	}

	// Fetch current yield data from treasury service
	yieldData, err := h.treasuryService.GetLatestYields()
	if err != nil {
		log.Printf("Error fetching yield data: %v", err)
		respondWithError(w, http.StatusInternalServerError, "failed to fetch current yield data")
		return
	}

	// Extract yield rate for selected term
	var yieldRate float64
	found := false
	for _, yieldPoint := range yieldData.Yields {
		if yieldPoint.Term == req.Term {
			yieldRate = yieldPoint.Rate
			found = true
			break
		}
	}

	if !found {
		log.Printf("Yield not found for term: %s", req.Term)
		respondWithError(w, http.StatusInternalServerError, "yield data not available for selected term")
		return
	}

	log.Printf("Current yield for %s: %.2f%%", req.Term, yieldRate)

	// Calculate purchase price using T-Bill discount pricing
	purchasePrice, err := utils.CalculateBillPrice(req.FaceValue, yieldRate, req.Term)
	if err != nil {
		// If term is not a valid T-Bill term, fall back to par pricing
		purchasePrice = req.FaceValue
		log.Printf("Using par pricing for term %s: purchase_price=%.2f", req.Term, purchasePrice)
	} else {
		discount := req.FaceValue - purchasePrice
		log.Printf("T-Bill discount pricing: face_value=%.2f, purchase_price=%.2f, discount=%.2f", req.FaceValue, purchasePrice, discount)
	}

	// Convert face value to pgtype.Numeric
	faceValueNumeric := pgtype.Numeric{}
	if err := faceValueNumeric.Scan(fmt.Sprintf("%.2f", req.FaceValue)); err != nil {
		log.Printf("Error converting face value to numeric: %v", err)
		respondWithError(w, http.StatusBadRequest, "invalid face value format")
		return
	}

	// Convert yield to pgtype.Numeric
	currentYield := pgtype.Numeric{}
	if err := currentYield.Scan(fmt.Sprintf("%.2f", yieldRate)); err != nil {
		log.Printf("Error converting yield to numeric: %v", err)
		respondWithError(w, http.StatusInternalServerError, "invalid yield format")
		return
	}

	// Call txService.BuyTreasury() with face value (service will calculate purchase price again)
	user, err := h.txService.BuyTreasury(r.Context(), req.UserID, req.Term, faceValueNumeric, currentYield)
	if err != nil {
		log.Printf("Error executing buy order for user %d: %v", req.UserID, err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("Buy order successful: user_id=%d, term=%s, face_value=%.2f, purchase_price=%.2f, yield=%.2f%%",
		req.UserID, req.Term, req.FaceValue, purchasePrice, yieldRate)

	// Return success response with updated user and purchase details
	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success":        true,
		"user":           user,
		"face_value":     req.FaceValue,
		"purchase_price": purchasePrice,
		"discount":       req.FaceValue - purchasePrice,
	})
}

// SellHandler handles POST /api/v1/sell requests.
// Expects JSON body with user_id, holding_id, and amount fields.
// Validates holding ownership, calculates yield, and processes the sell atomically.
// Returns updated user object on success, or error message on failure.
func (h *TransactionHandlers) SellHandler(w http.ResponseWriter, r *http.Request) {
	var req SellRequest

	// Decode JSON request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding sell request: %v", err)
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	log.Printf("Sell request received: user_id=%d, holding_id=%d, amount=%.2f", req.UserID, req.HoldingID, req.Amount)

	// Convert amount to pgtype.Numeric
	amount := pgtype.Numeric{}
	if err := amount.Scan(fmt.Sprintf("%.2f", req.Amount)); err != nil {
		log.Printf("Error converting amount to numeric: %v", err)
		respondWithError(w, http.StatusBadRequest, "invalid amount format")
		return
	}

	// Call txService.SellTreasury()
	user, err := h.txService.SellTreasury(r.Context(), req.UserID, req.HoldingID, amount)
	if err != nil {
		log.Printf("Error executing sell order for user %d: %v", req.UserID, err)

		// Map specific errors to appropriate HTTP status codes
		errMsg := err.Error()

		// Not found errors (404)
		if errMsg == "holding not found: no rows in result set" {
			respondWithError(w, http.StatusNotFound, "holding not found")
			return
		}

		// Forbidden errors (403) - holding doesn't belong to user
		if errMsg == "unauthorized: holding does not belong to user" {
			respondWithError(w, http.StatusForbidden, "unauthorized: holding does not belong to user")
			return
		}

		// All other errors (400) - bad request
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("Sell order successful: user_id=%d, holding_id=%d, amount=%.2f", req.UserID, req.HoldingID, req.Amount)

	// Return success response with updated user
	respondWithJSON(w, http.StatusOK, TransactionResponse{
		Success: true,
		User:    user,
	})
}
