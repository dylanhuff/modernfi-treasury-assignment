/**
 * Treasury security types based on maturity term.
 *
 * - Bills: Short-term (≤1 year) zero-coupon securities purchased at discount
 * - Notes: Medium-term (2-10 years) securities that earn interest
 * - Bonds: Long-term (30 years) securities that earn interest
 */
export const SecurityType = {
  Bill: 'bill',
  Note: 'note',
  Bond: 'bond',
} as const;

export type SecurityType = typeof SecurityType[keyof typeof SecurityType];

/**
 * Valid treasury terms supported by the application
 */
export type TreasuryTerm = '1M' | '3M' | '6M' | '1Y' | '2Y' | '5Y' | '10Y' | '30Y';

/**
 * Map treasury terms to their duration in days (matches backend TermDurationDays)
 */
export const TERM_DAYS: Record<TreasuryTerm, number> = {
  '1M': 30,
  '3M': 90,
  '6M': 180,
  '1Y': 365,
  '2Y': 730,
  '5Y': 1825,
  '10Y': 3650,
  '30Y': 10950,
};

/**
 * Determine security type based on treasury term (mirrors backend GetSecurityType)
 *
 * @param term - Treasury term code
 * @returns Security type or null if invalid term
 */
export function getSecurityType(term: string): SecurityType | null {
  switch (term) {
    case '1M':
    case '3M':
    case '6M':
    case '1Y':
      return SecurityType.Bill;
    case '2Y':
    case '5Y':
    case '10Y':
      return SecurityType.Note;
    case '30Y':
      return SecurityType.Bond;
    default:
      return null;
  }
}

/**
 * Get human-readable name for security type
 */
export function getSecurityTypeName(type: SecurityType): string {
  switch (type) {
    case SecurityType.Bill:
      return 'Treasury Bill';
    case SecurityType.Note:
      return 'Treasury Note';
    case SecurityType.Bond:
      return 'Treasury Bond';
  }
}

/**
 * Get badge color for security type display
 */
export function getSecurityTypeBadgeColor(type: SecurityType): 'blue' | 'purple' | 'indigo' {
  switch (type) {
    case SecurityType.Bill:
      return 'blue';
    case SecurityType.Note:
      return 'purple';
    case SecurityType.Bond:
      return 'indigo';
  }
}

/**
 * Get short label for security type
 */
export function getSecurityTypeLabel(type: SecurityType): string {
  switch (type) {
    case SecurityType.Bill:
      return 'Bill';
    case SecurityType.Note:
      return 'Note';
    case SecurityType.Bond:
      return 'Bond';
  }
}

/**
 * Calculate purchase price for Treasury securities (mirrors backend pricing logic)
 *
 * Treasury Bills: purchased at discount using 360-day money market convention
 * Treasury Notes/Bonds: purchased at par (face value)
 *
 * @param faceValue - Amount at maturity
 * @param yieldRate - Annual yield rate as percentage (0-100)
 * @param term - Treasury term code
 * @param securityType - Type of security
 * @returns Purchase price or null if invalid inputs
 */
export function calculatePurchasePrice(
  faceValue: number,
  yieldRate: number,
  term: TreasuryTerm,
  securityType: SecurityType | null
): number | null {
  // Validate inputs
  if (faceValue <= 0 || yieldRate < 0 || yieldRate > 100) {
    return null;
  }

  const days = TERM_DAYS[term];
  if (!days || !securityType) {
    return null;
  }

  // For Treasury Bills: calculate discounted purchase price
  if (securityType === SecurityType.Bill) {
    // Formula: price = faceValue × (1 - (yieldRate / 100 × days) / 360)
    // 360-day money market convention
    const discountFactor = (yieldRate / 100.0 * days) / 360.0;
    const price = faceValue * (1.0 - discountFactor);
    return Math.round(price * 100) / 100;
  }

  // For Treasury Notes/Bonds: purchase at par (face value)
  if (securityType === SecurityType.Note || securityType === SecurityType.Bond) {
    return Math.round(faceValue * 100) / 100;
  }

  return null;
}

/**
 * Calculate maturity value for Treasury Notes and Bonds
 *
 * Uses simplified interest model: principal + (principal × yield × years)
 *
 * @param principal - Face value amount
 * @param yieldRate - Annual yield rate as percentage (0-100)
 * @param daysHeld - Number of days held
 * @returns Total value at maturity
 */
export function calculateMaturityValue(
  principal: number,
  yieldRate: number,
  daysHeld: number
): number {
  // Simple interest formula: principal × (yieldRate / 100) × (daysHeld / 365)
  const simpleInterest = principal * (yieldRate / 100) * (daysHeld / 365);
  const maturityValue = principal + simpleInterest;
  return Math.round(maturityValue * 100) / 100;
}

/**
 * Get tooltip text explaining the security type
 */
export function getSecurityTypeTooltip(type: SecurityType): string {
  switch (type) {
    case SecurityType.Bill:
      return 'Treasury Bills are zero-coupon securities purchased at a discount. You pay less than face value and receive the full amount at maturity.';
    case SecurityType.Note:
    case SecurityType.Bond:
      return 'Treasury Notes and Bonds are purchased at face value and earn interest over time. This simplified model calculates total return at maturity.';
  }
}
