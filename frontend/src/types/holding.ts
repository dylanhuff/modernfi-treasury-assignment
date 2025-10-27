import type { SecurityType } from './treasury';

/**
 * Represents a treasury holding from the database.
 * Holdings track purchased treasuries and their remaining amounts.
 */
export interface Holding {
  id: number;
  user_id: number;
  term: string;
  amount: string; // Original purchase amount (decimal as string) - legacy field
  yield_at_purchase: string; // Yield rate at time of purchase
  purchase_date: string; // ISO 8601 format
  remaining_amount: string; // Amount not yet sold (decimal as string)
  // T-Bill discount pricing fields (added in Phase 1)
  face_value?: string; // Maturity amount - what user receives at maturity (null for legacy holdings)
  purchase_price?: string; // Actual cost - what user paid (null for legacy holdings)
  // Security type field (added in Phase 4.5 - Treasury Notes/Bonds implementation)
  security_type?: SecurityType | null; // SecurityType enum value (null for legacy holdings)
}

/**
 * Request payload for selling a treasury holding.
 * Sent to POST /api/v1/sell endpoint.
 */
export interface SellRequest {
  user_id: number;
  holding_id: number;
  amount: number;
}
