-- ============================================================================
-- Treasury Management Application - Consolidated Database Schema
-- ============================================================================
-- This schema consolidates all migrations into a single, clean definition.
-- For ModernFi take-home assignment.
-- ============================================================================

-- Drop existing objects if they exist (for clean recreation)
DROP TABLE IF EXISTS holdings CASCADE;
DROP TABLE IF EXISTS transactions CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TYPE IF EXISTS transaction_type CASCADE;

-- ============================================================================
-- ENUMS
-- ============================================================================

-- Transaction types: fund (deposit), withdraw, buy (treasury), sell (treasury)
CREATE TYPE transaction_type AS ENUM ('fund', 'withdraw', 'buy', 'sell');

-- ============================================================================
-- TABLES
-- ============================================================================

-- Users Table
-- Stores user account information and current balance
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    balance NUMERIC(12, 2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT users_balance_non_negative CHECK (balance >= 0)
);

-- Transactions Table
-- Records all financial transactions (deposits, withdrawals, buys, sells)
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    timestamp TIMESTAMP NOT NULL DEFAULT NOW(),
    type transaction_type NOT NULL,
    term VARCHAR(10),  -- Treasury term (e.g., "1M", "6M", "2Y") - nullable for fund/withdraw
    amount DECIMAL(12, 2) NOT NULL,
    yield_at_transaction DECIMAL(5, 2),  -- Yield % at time of buy/sell - nullable for fund/withdraw
    balance_after DECIMAL(12, 2) NOT NULL,
    holding_id INTEGER,  -- References holding for sell transactions - nullable

    -- Constraints
    CONSTRAINT transactions_amount_positive CHECK (amount > 0)
);

-- Holdings Table
-- Tracks active treasury holdings (bills, notes, bonds)
CREATE TABLE holdings (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    term VARCHAR(10) NOT NULL,  -- Treasury term (e.g., "1M", "6M", "2Y", "10Y", "30Y")
    amount DECIMAL(12, 2) NOT NULL,  -- Original purchase amount
    yield_at_purchase DECIMAL(5, 2) NOT NULL,  -- Yield % at time of purchase
    purchase_date TIMESTAMP NOT NULL DEFAULT NOW(),
    remaining_amount DECIMAL(12, 2) NOT NULL,  -- Amount remaining after partial sells
    face_value DECIMAL(12, 2),  -- Maturity value (for T-Bills with discount pricing)
    purchase_price DECIMAL(12, 2),  -- Actual price paid (discounted for T-Bills)
    security_type VARCHAR(10),  -- 'bill' (≤1Y), 'note' (2Y-10Y), 'bond' (30Y)

    -- Constraints
    CONSTRAINT holdings_amount_positive CHECK (amount > 0),
    CONSTRAINT holdings_remaining_non_negative CHECK (remaining_amount >= 0),
    CONSTRAINT holdings_remaining_lte_amount CHECK (remaining_amount <= amount)
);

-- ============================================================================
-- INDEXES
-- ============================================================================

-- Users table indexes
CREATE INDEX idx_users_name ON users(name);

-- Transactions table indexes
CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE INDEX idx_transactions_timestamp ON transactions(timestamp DESC);
CREATE INDEX idx_transactions_type ON transactions(type);

-- Holdings table indexes
CREATE INDEX idx_holdings_user_id ON holdings(user_id);
CREATE INDEX idx_holdings_purchase_date ON holdings(purchase_date DESC);

-- ============================================================================
-- COMMENTS
-- ============================================================================

COMMENT ON TABLE users IS 'User accounts with current balance';
COMMENT ON TABLE transactions IS 'All financial transactions (deposits, withdrawals, treasury trades)';
COMMENT ON TABLE holdings IS 'Active treasury holdings (bills, notes, bonds)';

COMMENT ON COLUMN holdings.security_type IS 'Type of treasury security: bill (≤1Y), note (2Y-10Y), bond (30Y)';
COMMENT ON COLUMN holdings.face_value IS 'Amount received at maturity (par value for T-Bills)';
COMMENT ON COLUMN holdings.purchase_price IS 'Actual discounted price paid (for T-Bills)';
COMMENT ON COLUMN transactions.holding_id IS 'References the holding being sold (for sell transactions)';
