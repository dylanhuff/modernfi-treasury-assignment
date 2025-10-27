export type TransactionType = 'fund' | 'withdraw' | 'buy' | 'sell';

/**
 * Represents a complete transaction record from the database.
 * Includes all fields needed for fund, withdraw, buy, and sell transactions.
 */
export interface Transaction {
  id: number;
  user_id: number;
  timestamp: string; // ISO 8601 format from backend
  type: TransactionType;
  term: string | null; // Only populated for buy/sell
  amount: string; // Decimal as string to preserve precision
  yield_at_transaction: string | null; // Only populated for buy/sell
  balance_after: string; // Decimal as string
  holding_id: number | null; // Only populated for sell
}

export interface TransactionRequest {
  user_id: number;
  amount: number;
}

/**
 * Request payload for buying a treasury security.
 * Sent to POST /api/v1/buy endpoint.
 */
export interface BuyRequest {
  user_id: number;
  term: string; // Treasury term: 1M, 3M, 6M, 1Y, 2Y, 5Y, 10Y, 30Y
  amount?: number; // Deprecated: Use face_value instead (kept for backward compatibility)
  face_value: number; // Amount at maturity (for T-Bills, this is the face value)
}

export interface TransactionResponse {
  success: boolean;
  user?: {
    id: number;
    name: string;
    balance: string;
    created_at: string;
  };
  error?: string;
}
