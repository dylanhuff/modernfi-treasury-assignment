package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"modernfi-treasury-app/internal/database"
	"modernfi-treasury-app/internal/utils"
)

type TransactionService struct {
	queries *database.Queries
	pool    *pgxpool.Pool
}

func NewTransactionService(queries *database.Queries, pool *pgxpool.Pool) *TransactionService {
	return &TransactionService{
		queries: queries,
		pool:    pool,
	}
}

// FundAccount adds funds to user account atomically
func (s *TransactionService) FundAccount(ctx context.Context, userID int32, amount pgtype.Numeric) (*database.User, error) {
	// Validate amount > 0
	amountFloat, err := amount.Float64Value()
	if err != nil {
		return nil, fmt.Errorf("invalid amount format: %w", err)
	}
	if !amountFloat.Valid || amountFloat.Float64 <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	var updatedUser *database.User

	// Use database transaction for atomicity
	err = pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		qtx := s.queries.WithTx(tx)

		// Update user balance
		user, err := qtx.UpdateUserBalance(ctx, database.UpdateUserBalanceParams{
			Balance: amount,
			ID:      userID,
		})
		if err != nil {
			return fmt.Errorf("failed to update balance: %w", err)
		}

		// Create transaction record
		_, err = qtx.CreateTransaction(ctx, database.CreateTransactionParams{
			UserID:             userID,
			Type:               database.TransactionTypeFund,
			Term:               pgtype.Text{Valid: false},
			Amount:             amount,
			YieldAtTransaction: pgtype.Numeric{Valid: false},
			BalanceAfter:       user.Balance,
			HoldingID:          pgtype.Int4{Valid: false},
		})
		if err != nil {
			return fmt.Errorf("failed to create transaction record: %w", err)
		}

		updatedUser = &user
		return nil
	})

	return updatedUser, err
}

// WithdrawAccount withdraws funds from user account atomically
func (s *TransactionService) WithdrawAccount(ctx context.Context, userID int32, amount pgtype.Numeric) (*database.User, error) {
	// Validate amount > 0
	amountFloat, err := amount.Float64Value()
	if err != nil {
		return nil, fmt.Errorf("invalid amount format: %w", err)
	}
	if !amountFloat.Valid || amountFloat.Float64 <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	// Get current user to check balance (quick pre-check for better UX)
	user, err := s.queries.GetUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Validate sufficient balance
	balanceFloat, err := user.Balance.Float64Value()
	if err != nil {
		return nil, fmt.Errorf("invalid balance format: %w", err)
	}
	if !balanceFloat.Valid {
		return nil, errors.New("user balance is invalid")
	}
	if balanceFloat.Float64 < amountFloat.Float64 {
		return nil, errors.New("insufficient balance")
	}

	var updatedUser *database.User

	// Use database transaction for atomicity
	err = pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		qtx := s.queries.WithTx(tx)

		// Re-check balance inside transaction to prevent race conditions
		// Use FOR UPDATE to lock the row until transaction completes
		currentUser, err := qtx.GetUserForUpdate(ctx, userID)
		if err != nil {
			return fmt.Errorf("failed to get user in transaction: %w", err)
		}

		currentBalanceFloat, err := currentUser.Balance.Float64Value()
		if err != nil {
			return fmt.Errorf("invalid current balance format: %w", err)
		}
		if !currentBalanceFloat.Valid {
			return errors.New("current user balance is invalid")
		}
		if currentBalanceFloat.Float64 < amountFloat.Float64 {
			return errors.New("insufficient balance")
		}

		// Create negative amount for withdrawal
		negativeAmount := pgtype.Numeric{}
		err = negativeAmount.Scan(fmt.Sprintf("-%.2f", amountFloat.Float64))
		if err != nil {
			return fmt.Errorf("failed to create negative amount: %w", err)
		}

		// Update user balance (negative amount to subtract)
		user, err := qtx.UpdateUserBalance(ctx, database.UpdateUserBalanceParams{
			Balance: negativeAmount,
			ID:      userID,
		})
		if err != nil {
			// Check if error is due to balance constraint violation (SQLSTATE 23514)
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23514" {
				return errors.New("insufficient balance")
			}
			return fmt.Errorf("failed to update balance: %w", err)
		}

		// Create transaction record
		_, err = qtx.CreateTransaction(ctx, database.CreateTransactionParams{
			UserID:             userID,
			Type:               database.TransactionTypeWithdraw,
			Term:               pgtype.Text{Valid: false},
			Amount:             amount,
			YieldAtTransaction: pgtype.Numeric{Valid: false},
			BalanceAfter:       user.Balance,
			HoldingID:          pgtype.Int4{Valid: false},
		})
		if err != nil {
			return fmt.Errorf("failed to create transaction record: %w", err)
		}

		updatedUser = &user
		return nil
	})

	return updatedUser, err
}

// BuyTreasury purchases a treasury security for a user atomically
// For T-Bills (1M, 3M, 6M, 1Y): faceValue is the amount at maturity, purchasePrice is calculated using discount pricing
// For Notes/Bonds (2Y, 5Y, 10Y, 30Y): uses par pricing (purchase price = face value)
func (s *TransactionService) BuyTreasury(
	ctx context.Context,
	userID int32,
	term string,
	faceValue pgtype.Numeric,
	currentYield pgtype.Numeric,
) (*database.User, error) {
	// Determine security type (bill, note, or bond)
	securityType, err := utils.GetSecurityType(term)
	if err != nil {
		return nil, fmt.Errorf("invalid term: %w", err)
	}

	// Validate face value > 0
	faceValueFloat, err := faceValue.Float64Value()
	if err != nil {
		return nil, fmt.Errorf("invalid face value format: %w", err)
	}
	if !faceValueFloat.Valid || faceValueFloat.Float64 <= 0 {
		return nil, errors.New("face value must be greater than zero")
	}

	// Extract yield rate for pricing calculation
	yieldRateFloat, err := currentYield.Float64Value()
	if err != nil {
		return nil, fmt.Errorf("invalid yield rate format: %w", err)
	}
	if !yieldRateFloat.Valid {
		return nil, errors.New("yield rate is required")
	}
	// Edge case validation: yield rate must be non-negative
	if yieldRateFloat.Float64 < 0 {
		return nil, errors.New("yield rate must be greater than or equal to zero")
	}

	// Calculate purchase price based on security type
	var purchasePriceFloat float64

	if securityType == utils.SecurityTypeBill {
		// Treasury Bills: Use discount pricing
		// price = faceValue × (1 - (yield × days) / 360)
		purchasePriceFloat, err = utils.CalculateBillPrice(faceValueFloat.Float64, yieldRateFloat.Float64, term)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate bill price: %w", err)
		}
	} else {
		// Treasury Notes/Bonds: Use par pricing
		purchasePriceFloat, err = utils.CalculateNoteBondPrice(faceValueFloat.Float64, yieldRateFloat.Float64, term)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate note/bond price: %w", err)
		}
	}

	// Convert purchase price to pgtype.Numeric
	purchasePrice := pgtype.Numeric{}
	err = purchasePrice.Scan(fmt.Sprintf("%.2f", purchasePriceFloat))
	if err != nil {
		return nil, fmt.Errorf("failed to create purchase price: %w", err)
	}

	// Get current user to check balance
	user, err := s.queries.GetUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Validate sufficient balance for purchase price (NOT face value!)
	balanceFloat, err := user.Balance.Float64Value()
	if err != nil {
		return nil, fmt.Errorf("invalid balance format: %w", err)
	}
	if !balanceFloat.Valid {
		return nil, errors.New("user balance is invalid")
	}
	if balanceFloat.Float64 < purchasePriceFloat {
		// Create friendly security type name for error message
		securityTypeName := "Treasury Bill"
		if securityType == utils.SecurityTypeNote {
			securityTypeName = "Treasury Note"
		} else if securityType == utils.SecurityTypeBond {
			securityTypeName = "Treasury Bond"
		}
		return nil, fmt.Errorf("insufficient balance: need %.2f for %s (face value: %.2f)",
			purchasePriceFloat, securityTypeName, faceValueFloat.Float64)
	}

	var updatedUser *database.User

	// Use database transaction for atomicity
	err = pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		qtx := s.queries.WithTx(tx)

		// Re-check balance inside transaction to prevent race conditions
		// Use FOR UPDATE to lock the row until transaction completes
		currentUser, err := qtx.GetUserForUpdate(ctx, userID)
		if err != nil {
			return fmt.Errorf("failed to get user in transaction: %w", err)
		}

		currentBalanceFloat, err := currentUser.Balance.Float64Value()
		if err != nil {
			return fmt.Errorf("invalid current balance format: %w", err)
		}
		if !currentBalanceFloat.Valid {
			return errors.New("current user balance is invalid")
		}
		// Check against purchase price (NOT face value!)
		if currentBalanceFloat.Float64 < purchasePriceFloat {
			return errors.New("insufficient balance")
		}

		// Create holding record with security type, face_value, and purchase_price
		// amount column is set to face_value for backward compatibility
		holding, err := qtx.CreateHolding(ctx, database.CreateHoldingParams{
			UserID:          userID,
			Term:            term,
			Amount:          faceValue, // Set to face value for backward compatibility
			YieldAtPurchase: currentYield,
			PurchaseDate:    pgtype.Timestamp{Time: time.Now(), Valid: true},
			RemainingAmount: faceValue,                                      // Initially, remaining amount equals face value
			FaceValue:       faceValue,                                      // Amount at maturity
			PurchasePrice:   purchasePrice,                                  // Actual discounted price paid (or par for notes/bonds)
			SecurityType:    pgtype.Text{String: securityType, Valid: true}, // bill, note, or bond
		})
		if err != nil {
			return fmt.Errorf("failed to create holding: %w", err)
		}

		// Create negative purchase price for withdrawal (subtract from balance)
		// Deduct purchase price, NOT face value!
		negativePurchasePrice := pgtype.Numeric{}
		err = negativePurchasePrice.Scan(fmt.Sprintf("-%.2f", purchasePriceFloat))
		if err != nil {
			return fmt.Errorf("failed to create negative purchase price: %w", err)
		}

		// Update user balance (deduct purchase price)
		user, err := qtx.UpdateUserBalance(ctx, database.UpdateUserBalanceParams{
			Balance: negativePurchasePrice,
			ID:      userID,
		})
		if err != nil {
			// Check if error is due to balance constraint violation (SQLSTATE 23514)
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23514" {
				return errors.New("insufficient balance")
			}
			return fmt.Errorf("failed to update balance: %w", err)
		}

		// Create transaction record (amount stores purchase price for buy transactions)
		_, err = qtx.CreateTransaction(ctx, database.CreateTransactionParams{
			UserID:             userID,
			Type:               database.TransactionTypeBuy,
			Term:               pgtype.Text{String: term, Valid: true},
			Amount:             purchasePrice, // Record the actual amount deducted (purchase price)
			YieldAtTransaction: currentYield,
			BalanceAfter:       user.Balance,
			HoldingID:          pgtype.Int4{Int32: holding.ID, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("failed to create transaction record: %w", err)
		}

		updatedUser = &user
		return nil
	})

	return updatedUser, err
}

// SellTreasury sells a treasury holding (full or partial) and returns proceeds to balance
func (s *TransactionService) SellTreasury(
	ctx context.Context,
	userID int32,
	holdingID int32,
	amount pgtype.Numeric,
) (*database.User, error) {
	// Validate amount > 0
	amountFloat, err := amount.Float64Value()
	if err != nil {
		return nil, fmt.Errorf("invalid amount format: %w", err)
	}
	if !amountFloat.Valid || amountFloat.Float64 <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	// Fetch holding to verify it exists and belongs to user
	holding, err := s.queries.GetHoldingByID(ctx, holdingID)
	if err != nil {
		return nil, fmt.Errorf("holding not found: %w", err)
	}

	// Verify holding belongs to user (security check)
	if holding.UserID != userID {
		return nil, errors.New("unauthorized: holding does not belong to user")
	}

	// Validate amount <= remaining_amount
	remainingFloat, err := holding.RemainingAmount.Float64Value()
	if err != nil {
		return nil, fmt.Errorf("invalid remaining amount format: %w", err)
	}
	if !remainingFloat.Valid {
		return nil, errors.New("holding remaining amount is invalid")
	}
	if amountFloat.Float64 > remainingFloat.Float64 {
		return nil, fmt.Errorf("insufficient remaining amount: requested %.2f, available %.2f",
			amountFloat.Float64, remainingFloat.Float64)
	}

	// Calculate proceeds based on security type
	var totalProceeds float64

	// Determine security type from holding (with legacy fallback)
	var securityType string
	if holding.SecurityType.Valid && holding.SecurityType.String != "" {
		// Use stored security type for new holdings
		securityType = holding.SecurityType.String
	} else {
		// For legacy holdings without security_type, infer from term
		inferredType, err := utils.GetSecurityType(holding.Term)
		if err != nil {
			// Fail-fast: Do not allow selling holdings with invalid/unknown security types
			// This ensures data integrity and prevents silent errors
			return nil, fmt.Errorf("cannot determine security type for holding %d (term: %s): %w", holdingID, holding.Term, err)
		}
		securityType = inferredType
	}

	if securityType == utils.SecurityTypeBill {
		// Treasury Bills: Return face value
		// The yield was already earned as the discount (face_value - purchase_price)
		totalProceeds = amountFloat.Float64
	} else {
		// Treasury Notes/Bonds: Calculate maturity value with simple interest
		// maturityValue = principal + (principal × yieldRate × daysHeld / 365)

		// Calculate days held from purchase date to now
		purchaseTime := holding.PurchaseDate.Time
		currentTime := time.Now()
		daysHeld := int(currentTime.Sub(purchaseTime).Hours() / 24)

		// Edge case validation: ensure days held is non-negative (protects against clock issues)
		if daysHeld < 0 {
			return nil, errors.New("invalid holding: purchase date is in the future")
		}

		// Get yield rate from holding
		yieldRateFloat, err := holding.YieldAtPurchase.Float64Value()
		if err != nil || !yieldRateFloat.Valid {
			return nil, fmt.Errorf("invalid yield rate for note/bond holding: %w", err)
		}
		// Edge case validation: yield rate must be non-negative
		if yieldRateFloat.Float64 < 0 {
			return nil, errors.New("invalid holding: yield rate must be greater than or equal to zero")
		}

		// Calculate maturity value using the helper function
		// principal = amount being sold, yieldRate = yield at purchase, daysHeld = time held
		maturityValue, err := utils.CalculateNoteBondMaturityValue(
			amountFloat.Float64,
			yieldRateFloat.Float64,
			daysHeld,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate note/bond maturity value: %w", err)
		}

		totalProceeds = maturityValue
		log.Printf("Selling %s holding %d: principal=%.2f, yield=%.2f%%, days_held=%d, maturity_value=%.2f",
			securityType, holdingID, amountFloat.Float64, yieldRateFloat.Float64, daysHeld, maturityValue)
	}

	var updatedUser *database.User

	// Use database transaction for atomicity
	err = pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		qtx := s.queries.WithTx(tx)

		// Update holding remaining_amount (subtract sold amount)
		newRemainingAmount := remainingFloat.Float64 - amountFloat.Float64
		newRemaining := pgtype.Numeric{}
		err = newRemaining.Scan(fmt.Sprintf("%.2f", newRemainingAmount))
		if err != nil {
			return fmt.Errorf("failed to create new remaining amount: %w", err)
		}

		_, err = qtx.UpdateHoldingRemainingAmount(ctx, database.UpdateHoldingRemainingAmountParams{
			ID:              holdingID,
			RemainingAmount: newRemaining,
		})
		if err != nil {
			return fmt.Errorf("failed to update holding remaining amount: %w", err)
		}

		// Create proceeds amount
		proceedsAmount := pgtype.Numeric{}
		err = proceedsAmount.Scan(fmt.Sprintf("%.2f", totalProceeds))
		if err != nil {
			return fmt.Errorf("failed to create proceeds amount: %w", err)
		}

		// Add proceeds to user balance
		user, err := qtx.UpdateUserBalance(ctx, database.UpdateUserBalanceParams{
			Balance: proceedsAmount,
			ID:      userID,
		})
		if err != nil {
			return fmt.Errorf("failed to update balance: %w", err)
		}

		// Create transaction record (store principal amount for consistency)
		_, err = qtx.CreateTransaction(ctx, database.CreateTransactionParams{
			UserID:             userID,
			Type:               database.TransactionTypeSell,
			Term:               pgtype.Text{String: holding.Term, Valid: true},
			Amount:             amount, // Principal amount (consistent with buy/fund/withdraw)
			YieldAtTransaction: holding.YieldAtPurchase,
			BalanceAfter:       user.Balance,
			HoldingID:          pgtype.Int4{Int32: holdingID, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("failed to create transaction record: %w", err)
		}

		updatedUser = &user
		return nil
	})

	return updatedUser, err
}
