/**
 * Interface for a single tour step
 * Matches the structure expected by React Joyride
 */
interface TourStep {
  target: string;
  content: string;
  placement?: 'top' | 'bottom' | 'left' | 'right' | 'center';
  disableBeacon?: boolean;
}

/**
 * Tour step definitions for the Treasury Management Dashboard walkthrough.
 * Comprehensive tour covering all features of the application.
 */
export const TOUR_STEPS: TourStep[] = [
  {
    // Step 1: Welcome - Sets context and expectations
    target: 'body',
    content: 'Welcome to the Treasury Management Dashboard! This guided tour will walk you through all features including user management, treasury investments, portfolio tracking, and transaction history.',
    placement: 'center',
    disableBeacon: true, // Show immediately without beacon
  },
  {
    // Step 2: User Selector - Switch between different users
    target: '[data-tour-id="user-selector"]',
    content: 'Switch between different user accounts to view their individual portfolios and balances. Each user has their own treasury holdings and transaction history.',
    placement: 'bottom',
    disableBeacon: true,
  },
  {
    // Step 3: Balance Display - Shows current account balance
    target: '[data-tour-id="balance-display"]',
    content: 'Your current account balance is displayed here. This balance updates in real-time as you buy or sell treasuries, deposit funds, or withdraw money.',
    placement: 'bottom',
    disableBeacon: true,
  },
  {
    // Step 4: Fund & Withdraw - Account management
    target: '[data-tour-id="fund-withdraw-buttons"]',
    content: 'Add funds to your account or withdraw money at any time. These actions are tracked in your transaction history for complete transparency.',
    placement: 'bottom',
    disableBeacon: true,
  },
  {
    // Step 5: Yield Curve - Demonstrates live API data fetching + visualization
    target: '[data-tour-id="yield-curve-chart"]',
    content: 'Live treasury yield data is fetched from the U.S. Treasury API and visualized here. The curve updates daily and shows current rates across all maturity terms (1 month to 30 years).',
    placement: 'bottom',
    disableBeacon: true,
  },
  {
    // Step 6: Buy Form - Purchase treasury securities
    target: '[data-tour-id="order-form"]',
    content: 'Purchase T-Bills, T-Notes, or T-Bonds by selecting a term and entering the face value amount. Purchase amounts are validated against your available balance to ensure you have sufficient funds.',
    placement: 'left',
    disableBeacon: true,
  },
  {
    // Step 7: Current Holdings - View your portfolio
    target: '[data-tour-id="current-holdings"]',
    content: 'View all your active treasury holdings with details like purchase price, face value, yield rates, and maturity dates. Holdings are updated in real-time as you buy or sell.',
    placement: 'top',
    disableBeacon: true,
  },
  {
    // Step 8: Sell Form - Liquidate holdings
    target: '[data-tour-id="sell-form"]',
    content: 'Sell your treasury holdings before maturity. The system calculates your proceeds including principal and prorated yield based on how long you\'ve held the security.',
    placement: 'left',
    disableBeacon: true,
  },
  {
    // Step 9: Transaction History - Complete audit trail
    target: '[data-tour-id="transaction-history"]',
    content: 'All transactions are permanently recorded here, including buys, sells, deposits, and withdrawals. This provides a complete audit trail of all account activity.',
    placement: 'top',
    disableBeacon: true,
  },
];
