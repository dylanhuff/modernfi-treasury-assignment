import type { User } from '../types/user';
import type { Transaction, TransactionRequest, TransactionResponse, BuyRequest } from '../types/transaction';
import type { Holding, SellRequest } from '../types/holding';

// Re-export types for convenience
export type { Holding, SellRequest } from '../types/holding';

/**
 * Represents a single point on the treasury yield curve.
 *
 * @property {string} term - The maturity term (e.g., "1M", "3M", "6M", "1Y", "2Y", "5Y", "10Y", "30Y")
 * @property {number} rate - The yield rate as a percentage (e.g., 4.5 for 4.5%)
 */
export interface YieldPoint {
  term: string;
  rate: number;
}

/**
 * Contains the complete yield curve data for a specific date.
 *
 * @property {string} date - The date this yield data is from (ISO 8601 format)
 * @property {YieldPoint[]} yields - Array of yield points for all treasury terms
 */
export interface YieldData {
  date: string;
  yields: YieldPoint[];
}

/**
 * Represents a single data point in the historical yield time series.
 * This format is optimized for Tremor LineChart with flattened structure.
 *
 * @property {string} date - The date in YYYY-MM-DD format
 * @property {number} 10Y - 10-year Treasury Note yield (percentage)
 * @property {number} 5Y - 5-year Treasury Note yield (percentage)
 * @property {number} 2Y - 2-year Treasury Note yield (percentage)
 */
export interface HistoricalDataPoint {
  date: string;        // YYYY-MM-DD format
  "10Y": number;       // 10-year yield (Treasury Note)
  "5Y": number;        // 5-year yield (Treasury Note)
  "2Y": number;        // 2-year yield (Treasury Note)
}

/**
 * Contains historical yield data for a specific time period.
 * Matches the backend HistoricalYieldData Go struct exactly.
 * The data format is optimized for Tremor LineChart - no transformation needed on frontend.
 *
 * @property {string} period - Time period ("1M", "3M", "6M", or "1Y")
 * @property {string} startDate - Start date of the period (YYYY-MM-DD format)
 * @property {string} endDate - End date of the period (YYYY-MM-DD format)
 * @property {string[]} terms - Array of maturity terms included (e.g., ["10Y", "5Y", "2Y"])
 * @property {HistoricalDataPoint[]} data - Array of historical data points
 */
export interface HistoricalYieldData {
  period: string;      // "1M", "3M", "6M", or "1Y"
  startDate: string;   // YYYY-MM-DD format
  endDate: string;     // YYYY-MM-DD format
  terms: string[];     // ["10Y", "5Y", "2Y"]
  data: HistoricalDataPoint[];
}

// Use relative URLs when in production (served by nginx that proxies to backend)
// Use full URL for local development (Vite dev server on port 5173)
const API_BASE_URL = import.meta.env.DEV
  ? (import.meta.env.VITE_API_URL || 'http://localhost:8080')
  : ''; // Empty string means relative URLs (same origin)

/**
 * Fetches current US Treasury yield curve data from the backend API.
 *
 * The backend caches this data from Treasury.gov for 1 hour to improve performance.
 * Returns yield rates for all standard treasury terms (1M through 30Y).
 *
 * @returns {Promise<YieldData>} Promise resolving to the current yield curve data
 * @throws {Error} If the API request fails or returns a non-OK status
 *
 * @example
 * ```tsx
 * const yieldData = await fetchYields();
 * console.log(yieldData.yields); // [{term: "1M", rate: 4.5}, ...]
 * ```
 */
export async function fetchYields(): Promise<YieldData> {
  try {
    const response = await fetch(`${API_BASE_URL}/api/yields`, {
      method: 'GET',
      headers: {
        'Accept': 'application/json',
      },
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data: YieldData = await response.json();
    return data;
  } catch (error) {
    if (error instanceof Error) {
      throw new Error(`Failed to fetch yields: ${error.message}`);
    }
    throw new Error('Failed to fetch yields: Unknown error');
  }
}

/**
 * Fetches historical Treasury yield data for a given time period.
 *
 * Returns time-series data for Treasury Note yields (10Y, 5Y, 2Y) over the specified period.
 * The backend caches historical data for 24 hours as it doesn't change retroactively.
 * Data is returned in a format optimized for Tremor LineChart - no frontend transformation needed.
 *
 * @param {string} period - Time frame: "1M", "3M", "6M", or "1Y" (defaults to "3M" if not provided)
 * @returns {Promise<HistoricalYieldData>} Promise resolving to historical yield data
 * @throws {Error} If fetch fails or period is invalid
 *
 * @example
 * ```tsx
 * const historicalData = await fetchHistoricalYields("3M");
 * console.log(historicalData.data); // [{date: "2025-07-26", "10Y": 4.25, "5Y": 4.10, "2Y": 4.05}, ...]
 * ```
 */
export async function fetchHistoricalYields(
  period: string
): Promise<HistoricalYieldData> {
  try {
    const response = await fetch(
      `${API_BASE_URL}/api/yields/historical?period=${period}`,
      {
        method: 'GET',
        headers: {
          'Accept': 'application/json',
        },
      }
    );

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data: HistoricalYieldData = await response.json();
    return data;
  } catch (error) {
    if (error instanceof Error) {
      throw new Error(`Failed to fetch historical yields: ${error.message}`);
    }
    throw new Error('Failed to fetch historical yields: Unknown error');
  }
}

/**
 * Fetches all users from the backend API.
 *
 * Returns a list of all users in the system with their current balances.
 * Used by CurrentUserProvider to populate the user selection dropdown and
 * restore the previously selected user from localStorage.
 *
 * @returns {Promise<User[]>} Promise resolving to array of all users (empty array if no users exist)
 * @throws {Error} If the API request fails or returns a non-OK status
 *
 * @example
 * ```tsx
 * const users = await fetchUsers();
 * const firstUser = users[0];
 * console.log(`${firstUser.name} has $${firstUser.balance}`);
 * ```
 */
export async function fetchUsers(): Promise<User[]> {
  try {
    const response = await fetch(`${API_BASE_URL}/api/v1/users`, {
      headers: {
        'Accept': 'application/json',
      },
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data: User[] = await response.json();
    return data;
  } catch (error) {
    if (error instanceof Error) {
      throw new Error(`Failed to fetch users: ${error.message}`);
    }
    throw new Error('Failed to fetch users: Unknown error');
  }
}

/**
 * Adds funds to a user's account.
 *
 * Makes a POST request to the fund endpoint to add the specified amount
 * to the user's balance. The operation is atomic - both the user balance
 * and transaction record are updated together in a database transaction.
 *
 * @param {number} userId - The ID of the user to fund
 * @param {number} amount - The amount to add to the account (must be > 0)
 * @returns {Promise<User>} Promise resolving to the updated user with new balance
 * @throws {Error} If the API request fails, amount is invalid, or server returns error
 *
 * @example
 * ```tsx
 * const updatedUser = await fundAccount(1, 50000);
 * console.log(`New balance: $${updatedUser.balance}`);
 * ```
 */
export async function fundAccount(userId: number, amount: number): Promise<User> {
  const response = await fetch(`${API_BASE_URL}/api/v1/fund`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      user_id: userId,
      amount: amount,
    } as TransactionRequest),
  });

  const data: TransactionResponse = await response.json();

  if (!response.ok || !data.success) {
    throw new Error(data.error || 'Fund operation failed');
  }

  if (!data.user) {
    throw new Error('No user data returned');
  }

  return {
    id: data.user.id,
    name: data.user.name,
    balance: parseFloat(data.user.balance),
    created_at: data.user.created_at,
  };
}

/**
 * Withdraws funds from a user's account.
 *
 * Makes a POST request to the withdraw endpoint to subtract the specified amount
 * from the user's balance. The operation validates sufficient balance and is atomic -
 * both the user balance and transaction record are updated together in a database transaction.
 *
 * @param {number} userId - The ID of the user to withdraw from
 * @param {number} amount - The amount to withdraw from the account (must be > 0 and <= balance)
 * @returns {Promise<User>} Promise resolving to the updated user with new balance
 * @throws {Error} If the API request fails, amount is invalid, insufficient balance, or server returns error
 *
 * @example
 * ```tsx
 * try {
 *   const updatedUser = await withdrawAccount(1, 25000);
 *   console.log(`New balance: $${updatedUser.balance}`);
 * } catch (error) {
 *   console.error('Withdrawal failed:', error.message); // e.g., "insufficient balance"
 * }
 * ```
 */
export async function withdrawAccount(userId: number, amount: number): Promise<User> {
  const response = await fetch(`${API_BASE_URL}/api/v1/withdraw`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      user_id: userId,
      amount: amount,
    } as TransactionRequest),
  });

  const data: TransactionResponse = await response.json();

  if (!response.ok || !data.success) {
    throw new Error(data.error || 'Withdraw operation failed');
  }

  if (!data.user) {
    throw new Error('No user data returned');
  }

  return {
    id: data.user.id,
    name: data.user.name,
    balance: parseFloat(data.user.balance),
    created_at: data.user.created_at,
  };
}

/**
 * Buys a treasury security for a user.
 *
 * Makes a POST request to purchase a treasury with the specified term and face value.
 * For T-Bills, the user pays a discounted price below face value and receives the full
 * face value at maturity. The operation validates sufficient balance for the purchase price
 * and is atomic - the holding record, transaction record, and user balance are all updated
 * together in a database transaction.
 *
 * @param {number} userId - The ID of the user making the purchase
 * @param {string} term - The treasury term (1M, 3M, 6M, 1Y, 2Y, 5Y, 10Y, 30Y)
 * @param {number} faceValue - The face value amount (amount at maturity, not purchase cost)
 * @returns {Promise<User>} Promise resolving to the updated user with new balance
 * @throws {Error} If the API request fails, face value is invalid, insufficient balance for purchase cost, or server returns error
 *
 * @example
 * ```tsx
 * try {
 *   // User pays ~$97,750 for a $100,000 face value 6M T-Bill
 *   const updatedUser = await buyTreasury(1, "6M", 100000);
 *   console.log(`New balance: $${updatedUser.balance}`);
 * } catch (error) {
 *   console.error('Purchase failed:', error.message);
 * }
 * ```
 */
export async function buyTreasury(
  userId: number,
  term: string,
  faceValue: number
): Promise<User> {
  const response = await fetch(`${API_BASE_URL}/api/v1/buy`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      user_id: userId,
      term: term,
      face_value: faceValue,
    } as BuyRequest),
  });

  const data: TransactionResponse = await response.json();

  if (!response.ok || !data.success) {
    // Enhance error message to clarify it's about purchase cost, not face value
    const errorMsg = data.error || 'Buy operation failed';
    throw new Error(errorMsg);
  }

  if (!data.user) {
    throw new Error('No user data returned');
  }

  return {
    id: data.user.id,
    name: data.user.name,
    balance: parseFloat(data.user.balance),
    created_at: data.user.created_at,
  };
}

/**
 * Fetches all transactions for a specific user.
 *
 * Returns transactions ordered by timestamp descending (most recent first).
 * Includes fund, withdraw, and future buy/sell transactions.
 *
 * @param {number} userId - The ID of the user whose transactions to fetch
 * @returns {Promise<Transaction[]>} Promise resolving to array of transactions (empty if no transactions)
 * @throws {Error} If the API request fails or user doesn't exist
 *
 * @example
 * ```tsx
 * const transactions = await fetchUserTransactions(1);
 * console.log(`User has ${transactions.length} transactions`);
 * ```
 */
export async function fetchUserTransactions(userId: number): Promise<Transaction[]> {
  try {
    const response = await fetch(`${API_BASE_URL}/api/v1/users/${userId}/transactions`, {
      method: 'GET',
      headers: {
        'Accept': 'application/json',
      },
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data: Transaction[] = await response.json();
    return data;
  } catch (error) {
    if (error instanceof Error) {
      throw new Error(`Failed to fetch transactions: ${error.message}`);
    }
    throw new Error('Failed to fetch transactions: Unknown error');
  }
}

/**
 * Fetches all active holdings for a specific user.
 *
 * Returns holdings with remaining_amount > 0, ordered by purchase date.
 * Used by SellForm to populate the holdings dropdown for selling.
 *
 * @param {number} userId - The ID of the user whose holdings to fetch
 * @returns {Promise<Holding[]>} Promise resolving to array of holdings (empty if no holdings)
 * @throws {Error} If the API request fails or user doesn't exist
 *
 * @example
 * ```tsx
 * const holdings = await fetchUserHoldings(1);
 * console.log(`User has ${holdings.length} active holdings`);
 * ```
 */
export async function fetchUserHoldings(userId: number): Promise<Holding[]> {
  try {
    const response = await fetch(`${API_BASE_URL}/api/v1/users/${userId}/holdings`, {
      method: 'GET',
      headers: {
        'Accept': 'application/json',
      },
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data: Holding[] = await response.json();
    return data;
  } catch (error) {
    if (error instanceof Error) {
      throw new Error(`Failed to fetch holdings: ${error.message}`);
    }
    throw new Error('Failed to fetch holdings: Unknown error');
  }
}

/**
 * Sells a treasury security (full or partial).
 *
 * Makes a POST request to sell the specified amount from a holding.
 * The operation validates ownership and sufficient remaining amount.
 * The sale is atomic - the holding remaining amount, transaction record,
 * and user balance are all updated together in a database transaction.
 * User receives principal + accrued yield based on time held.
 *
 * @param {number} userId - The ID of the user making the sale
 * @param {number} holdingId - The ID of the holding to sell from
 * @param {number} amount - The amount to sell (must be > 0 and <= remaining_amount)
 * @returns {Promise<User>} Promise resolving to the updated user with new balance
 * @throws {Error} If the API request fails, invalid ownership, insufficient remaining amount, or server returns error
 *
 * @example
 * ```tsx
 * try {
 *   const updatedUser = await sellTreasury(1, 5, 50000);
 *   console.log(`New balance: $${updatedUser.balance}`);
 * } catch (error) {
 *   console.error('Sale failed:', error.message);
 * }
 * ```
 */
export async function sellTreasury(
  userId: number,
  holdingId: number,
  amount: number
): Promise<User> {
  const response = await fetch(`${API_BASE_URL}/api/v1/sell`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      user_id: userId,
      holding_id: holdingId,
      amount: amount,
    } as SellRequest),
  });

  const data: TransactionResponse = await response.json();

  if (!response.ok || !data.success) {
    throw new Error(data.error || 'Sell operation failed');
  }

  if (!data.user) {
    throw new Error('No user data returned');
  }

  return {
    id: data.user.id,
    name: data.user.name,
    balance: parseFloat(data.user.balance),
    created_at: data.user.created_at,
  };
}
