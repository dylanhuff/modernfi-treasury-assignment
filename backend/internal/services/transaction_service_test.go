package services

import (
	"context"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"modernfi-treasury-app/internal/database"
)

// TestBuyTreasury_Success tests successful treasury purchase
func TestBuyTreasury_Success(t *testing.T) {
	ctx := context.Background()

	// Connect to test database
	// Note: This requires a running PostgreSQL instance for integration testing
	connString := "postgres://postgres:postgres@localhost:5432/treasury_db?sslmode=disable"
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		t.Skipf("Skipping integration test: database not available: %v", err)
		return
	}
	defer pool.Close()

	queries := database.New(pool)
	service := NewTransactionService(queries, pool)

	// Create test user with sufficient balance
	testUser, err := queries.CreateUser(ctx, database.CreateUserParams{
		Name:    "Test User - Buy Success",
		Balance: mustNumeric("500000.00"),
	})
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	defer cleanupUser(t, ctx, queries, testUser.ID)

	// Execute buy order
	amount := mustNumeric("100000.00")
	currentYield := mustNumeric("4.50")
	updatedUser, err := service.BuyTreasury(ctx, testUser.ID, "6M", amount, currentYield)

	// Verify success
	if err != nil {
		t.Fatalf("BuyTreasury failed: %v", err)
	}
	if updatedUser == nil {
		t.Fatal("Expected updated user, got nil")
	}

	// Verify balance decreased by purchase price (discount pricing for 6M T-Bill)
	// Purchase price = $100,000 × (1 - (4.50 / 100 × 180) / 360) = $97,750
	expectedBalance := 402250.00 // $500,000 - $97,750
	actualBalance := mustFloat64(updatedUser.Balance)
	if actualBalance != expectedBalance {
		t.Errorf("Expected balance %f, got %f", expectedBalance, actualBalance)
	}

	// Verify holding was created
	holdings, err := queries.GetHoldingsByUser(ctx, testUser.ID)
	if err != nil {
		t.Fatalf("Failed to get holdings: %v", err)
	}
	if len(holdings) != 1 {
		t.Errorf("Expected 1 holding, got %d", len(holdings))
	}
	if len(holdings) > 0 {
		holding := holdings[0]
		if holding.Term != "6M" {
			t.Errorf("Expected term '6M', got '%s'", holding.Term)
		}
		if mustFloat64(holding.Amount) != 100000.00 {
			t.Errorf("Expected holding amount 100000.00, got %f", mustFloat64(holding.Amount))
		}
		if mustFloat64(holding.YieldAtPurchase) != 4.50 {
			t.Errorf("Expected yield 4.50, got %f", mustFloat64(holding.YieldAtPurchase))
		}
		if mustFloat64(holding.RemainingAmount) != 100000.00 {
			t.Errorf("Expected remaining amount 100000.00, got %f", mustFloat64(holding.RemainingAmount))
		}
	}

	// Verify transaction was created
	transactions, err := queries.GetTransactionsByUser(ctx, testUser.ID)
	if err != nil {
		t.Fatalf("Failed to get transactions: %v", err)
	}
	if len(transactions) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(transactions))
	}
	if len(transactions) > 0 {
		tx := transactions[0]
		if tx.Type != database.TransactionTypeBuy {
			t.Errorf("Expected transaction type 'buy', got '%s'", tx.Type)
		}
		if !tx.Term.Valid || tx.Term.String != "6M" {
			t.Errorf("Expected term '6M', got '%s'", tx.Term.String)
		}
		// Transaction amount should be the purchase price (discounted) for T-Bills
		expectedTxAmount := 97750.00
		if mustFloat64(tx.Amount) != expectedTxAmount {
			t.Errorf("Expected transaction amount %f, got %f", expectedTxAmount, mustFloat64(tx.Amount))
		}
	}
}

// TestBuyTreasury_InsufficientBalance tests buy with insufficient balance
func TestBuyTreasury_InsufficientBalance(t *testing.T) {
	ctx := context.Background()

	connString := "postgres://postgres:postgres@localhost:5432/treasury_db?sslmode=disable"
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		t.Skipf("Skipping integration test: database not available: %v", err)
		return
	}
	defer pool.Close()

	queries := database.New(pool)
	service := NewTransactionService(queries, pool)

	// Create test user with low balance
	testUser, err := queries.CreateUser(ctx, database.CreateUserParams{
		Name:    "Test User - Insufficient Balance",
		Balance: mustNumeric("50000.00"),
	})
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	defer cleanupUser(t, ctx, queries, testUser.ID)

	// Attempt to buy more than available balance
	amount := mustNumeric("100000.00")
	currentYield := mustNumeric("4.50")
	_, err = service.BuyTreasury(ctx, testUser.ID, "6M", amount, currentYield)

	// Verify error returned
	if err == nil {
		t.Fatal("Expected insufficient balance error, got nil")
	}
	if !strings.Contains(err.Error(), "insufficient balance") {
		t.Errorf("Expected error containing 'insufficient balance', got: %v", err)
	}

	// Verify no holding was created
	holdings, err := queries.GetHoldingsByUser(ctx, testUser.ID)
	if err != nil {
		t.Fatalf("Failed to get holdings: %v", err)
	}
	if len(holdings) != 0 {
		t.Errorf("Expected 0 holdings, got %d", len(holdings))
	}

	// Verify no transaction was created
	transactions, err := queries.GetTransactionsByUser(ctx, testUser.ID)
	if err != nil {
		t.Fatalf("Failed to get transactions: %v", err)
	}
	if len(transactions) != 0 {
		t.Errorf("Expected 0 transactions, got %d", len(transactions))
	}

	// Verify balance unchanged
	user, err := queries.GetUser(ctx, testUser.ID)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}
	if mustFloat64(user.Balance) != 50000.00 {
		t.Errorf("Expected balance unchanged at 50000.00, got %f", mustFloat64(user.Balance))
	}
}

// TestBuyTreasury_InvalidAmount tests buy with invalid amounts
func TestBuyTreasury_InvalidAmount(t *testing.T) {
	ctx := context.Background()

	connString := "postgres://postgres:postgres@localhost:5432/treasury_db?sslmode=disable"
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		t.Skipf("Skipping integration test: database not available: %v", err)
		return
	}
	defer pool.Close()

	queries := database.New(pool)
	service := NewTransactionService(queries, pool)

	// Create test user
	testUser, err := queries.CreateUser(ctx, database.CreateUserParams{
		Name:    "Test User - Invalid Amount",
		Balance: mustNumeric("500000.00"),
	})
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	defer cleanupUser(t, ctx, queries, testUser.ID)

	testCases := []struct {
		name   string
		amount string
	}{
		{"Zero amount", "0.00"},
		{"Negative amount", "-1000.00"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			amount := mustNumeric(tc.amount)
			currentYield := mustNumeric("4.50")
			_, err := service.BuyTreasury(ctx, testUser.ID, "6M", amount, currentYield)

			// Verify error returned
			if err == nil {
				t.Fatalf("Expected error for %s, got nil", tc.name)
			}
			if !strings.Contains(err.Error(), "face value must be greater than zero") {
				t.Errorf("Expected error containing 'face value must be greater than zero', got: %v", err)
			}
		})
	}
}

// TestBuyTreasury_AtomicTransaction tests that failed operations don't leave partial state
func TestBuyTreasury_AtomicTransaction(t *testing.T) {
	ctx := context.Background()

	connString := "postgres://postgres:postgres@localhost:5432/treasury_db?sslmode=disable"
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		t.Skipf("Skipping integration test: database not available: %v", err)
		return
	}
	defer pool.Close()

	queries := database.New(pool)
	service := NewTransactionService(queries, pool)

	// Create test user with low balance
	testUser, err := queries.CreateUser(ctx, database.CreateUserParams{
		Name:    "Test User - Atomic Test",
		Balance: mustNumeric("100000.00"),
	})
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	defer cleanupUser(t, ctx, queries, testUser.ID)

	// Attempt buy that should fail due to balance constraint
	// For 6M T-Bill at 4.50% yield, face value of $102,500 costs ~$100,194 (exceeds $100,000 balance)
	amount := mustNumeric("102500.00")
	currentYield := mustNumeric("4.50")
	_, err = service.BuyTreasury(ctx, testUser.ID, "6M", amount, currentYield)

	// Should fail due to insufficient balance
	if err == nil {
		t.Fatal("Expected insufficient balance error, got nil")
	}

	// Verify database is in consistent state (no partial records)
	holdings, _ := queries.GetHoldingsByUser(ctx, testUser.ID)
	if len(holdings) != 0 {
		t.Errorf("Expected 0 holdings after failed transaction, got %d", len(holdings))
	}

	transactions, _ := queries.GetTransactionsByUser(ctx, testUser.ID)
	if len(transactions) != 0 {
		t.Errorf("Expected 0 transactions after failed transaction, got %d", len(transactions))
	}

	// Verify balance unchanged
	user, _ := queries.GetUser(ctx, testUser.ID)
	if mustFloat64(user.Balance) != 100000.00 {
		t.Errorf("Expected balance unchanged at 100000.00, got %f", mustFloat64(user.Balance))
	}
}

// Helper functions

func mustNumeric(s string) pgtype.Numeric {
	var n pgtype.Numeric
	if err := n.Scan(s); err != nil {
		panic(err)
	}
	return n
}

func mustFloat64(n pgtype.Numeric) float64 {
	v, err := n.Float64Value()
	if err != nil {
		panic(err)
	}
	if !v.Valid {
		panic("invalid numeric value")
	}
	return v.Float64
}

func cleanupUser(t *testing.T, ctx context.Context, queries *database.Queries, userID int32) {
	// Clean up test user (cascade will clean up holdings and transactions)
	if err := queries.DeleteUser(ctx, userID); err != nil {
		t.Logf("Warning: failed to cleanup test user %d: %v", userID, err)
	}
}
