# OrderForm Component Tests

This document outlines comprehensive tests for the OrderForm component. These tests validate form behavior, validation logic, and user interaction flows.

## Test Setup Required

To run these tests, install the following packages:
```bash
npm install --save-dev vitest @testing-library/react @testing-library/user-event @testing-library/jest-dom
```

## Test Suite: OrderForm Component

### 1. Rendering Tests

#### Test: Component renders with all required elements
**Given**: OrderForm component is rendered
**When**: Component mounts
**Then**: Should display:
- "Buy Treasury" title
- Available balance text
- Treasury term dropdown
- Amount input field
- Buy Treasury button

#### Test: Balance displays correctly from context
**Given**: User has balance of $1,000,000
**When**: OrderForm renders
**Then**: Should display "Available Balance: $1,000,000"

#### Test: All 8 treasury terms are available in dropdown
**Given**: Dropdown is opened
**When**: User views options
**Then**: Should show all terms: 1M, 3M, 6M, 1Y, 2Y, 5Y, 10Y, 30Y

---

### 2. Form Validation Tests

#### Test: Form validation - no term selected
**Given**: Amount is entered
**When**: User clicks Buy without selecting term
**Then**: Should display error "Please select a treasury term"

#### Test: Form validation - no amount entered
**Given**: Term is selected
**When**: User clicks Buy without entering amount
**Then**: Should display error "Please enter an amount greater than $0"

#### Test: Form validation - zero amount
**Given**: Term is selected
**When**: User enters amount "0" and clicks Buy
**Then**: Should display error "Please enter an amount greater than $0"

#### Test: Form validation - negative amount
**Given**: Term is selected
**When**: User enters negative amount and clicks Buy
**Then**: Should display error "Please enter an amount greater than $0"

#### Test: Form validation - insufficient balance
**Given**: User has balance $100,000
**And**: Term "6M" is selected
**When**: User enters amount "150000" and clicks Buy
**Then**: Should display error "Insufficient balance. You have $100,000 available"

#### Test: Buy button disabled when form invalid
**Given**: Form has invalid state (no term or invalid amount)
**When**: User views the button
**Then**: Buy button should be disabled

---

### 3. Form Interaction Tests

#### Test: Amount input formats numbers with commas
**Given**: Amount input is focused
**When**: User types "100000"
**Then**: Input should display "100,000"

#### Test: Amount input removes non-numeric characters
**Given**: Amount input is focused
**When**: User types "100,000.50"
**Then**: Input value should be 100000 (integer only)

#### Test: Balance preview updates as amount changes
**Given**: User has balance $500,000
**And**: User enters amount "100000"
**When**: Amount is entered
**Then**: Should display "Balance after purchase: $400,000"

#### Test: Balance preview hidden when amount is zero
**Given**: Amount input is empty or zero
**When**: Component renders
**Then**: Balance preview should not be visible

---

### 4. Successful Buy Flow Tests

#### Test: Successful buy updates balance
**Given**: Valid term and amount are entered
**When**: User clicks Buy and API succeeds
**Then**:
- User balance should update in context
- Success message should display
- Form should clear (term and amount reset)

#### Test: Success message auto-hides after 3 seconds
**Given**: Buy operation succeeds
**When**: Success message displays
**Then**: Message should disappear after 3 seconds

#### Test: Transactions refresh after successful buy
**Given**: Buy operation succeeds
**When**: API call completes
**Then**: refreshTransactions() should be called

#### Test: Form clears after successful buy
**Given**: Term "6M" and amount "100000" are entered
**When**: Buy succeeds
**Then**:
- Term should be empty
- Amount should be empty
- Display value should be empty

---

### 5. Error Handling Tests

#### Test: API error displays error message
**Given**: Valid form data
**When**: buyTreasury API call fails
**Then**: Error message should display with API error text

#### Test: Network error displays generic message
**Given**: Valid form data
**When**: Network request fails (no response)
**Then**: Should display "Failed to purchase treasury"

#### Test: Error message clears on next submission
**Given**: Previous buy attempt failed
**When**: User clicks Buy again
**Then**: Previous error message should clear before new attempt

---

### 6. Loading State Tests

#### Test: Button shows loading state during API call
**Given**: Valid form data
**When**: User clicks Buy (during API call)
**Then**:
- Button should show "Processing..." text
- Button should display loading spinner
- Button should be disabled

#### Test: Button disabled during loading
**Given**: Buy operation in progress
**When**: User tries to click button
**Then**: Button should be disabled (prevents double submission)

---

### 7. Context Integration Tests

#### Test: Uses currentUser from context
**Given**: CurrentUserContext provides user with id=8
**When**: Buy is submitted
**Then**: Should call buyTreasury(8, term, amount)

#### Test: Calls updateUser after successful buy
**Given**: Buy succeeds with updated user data
**When**: API returns new user object
**Then**: updateUser() should be called with new user data

#### Test: No user in context disables form
**Given**: currentUser is null
**When**: Component renders
**Then**: Buy button should be disabled

---

### 8. Edge Cases

#### Test: Very large amount formatting
**Given**: User enters "99999999"
**When**: Input updates
**Then**: Should display "99,999,999"

#### Test: Empty amount after typing and deleting
**Given**: User types "1000" then deletes all digits
**When**: Input becomes empty
**Then**:
- amount should be undefined
- displayValue should be empty
- Balance preview should hide

#### Test: Rapid term changes don't break state
**Given**: User rapidly selects different terms
**When**: Dropdown value changes multiple times
**Then**: Component state should update correctly

---

## Validation Rules Summary

The OrderForm validates the following:

1. **Term validation**: Must select one of 8 valid terms
2. **Amount validation**: Must be > 0
3. **Balance validation**: Amount must be <= available balance
4. **User validation**: currentUser must exist

## Test Coverage Goals

- ✅ All validation rules tested
- ✅ Success flow tested
- ✅ Error flow tested
- ✅ Loading states tested
- ✅ Context integration tested
- ✅ Form clearing tested
- ✅ Edge cases tested

## Example Test Implementation

Below is example code showing how these tests would be implemented with Vitest and React Testing Library:

\`\`\`typescript
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { OrderForm } from './OrderForm';
import { CurrentUserProvider } from '../contexts/CurrentUserContext';
import { TransactionRefreshProvider } from '../contexts/TransactionRefreshContext';
import * as api from '../services/api';

// Mock the API
vi.mock('../services/api', () => ({
  buyTreasury: vi.fn(),
}));

describe('OrderForm', () => {
  const mockUser = {
    id: 1,
    name: 'Test User',
    balance: 1000000,
    created_at: '2024-01-01',
  };

  const mockUpdateUser = vi.fn();
  const mockRefreshTransactions = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should display insufficient balance error when amount exceeds balance', async () => {
    const user = userEvent.setup();

    render(
      <CurrentUserProvider value={{ currentUser: mockUser, updateUser: mockUpdateUser }}>
        <TransactionRefreshProvider value={{ refreshTransactions: mockRefreshTransactions }}>
          <OrderForm />
        </TransactionRefreshProvider>
      </CurrentUserProvider>
    );

    // Select term
    const termSelect = screen.getByLabelText('Treasury Term');
    await user.selectOptions(termSelect, '6M');

    // Enter amount greater than balance
    const amountInput = screen.getByPlaceholderText('Enter amount');
    await user.type(amountInput, '2000000');

    // Click buy
    const buyButton = screen.getByRole('button', { name: /buy treasury/i });
    await user.click(buyButton);

    // Verify error message
    expect(screen.getByText(/insufficient balance/i)).toBeInTheDocument();
  });

  it('should call API and update user on successful buy', async () => {
    const user = userEvent.setup();
    const updatedUser = { ...mockUser, balance: 900000 };
    vi.mocked(api.buyTreasury).mockResolvedValue(updatedUser);

    render(
      <CurrentUserProvider value={{ currentUser: mockUser, updateUser: mockUpdateUser }}>
        <TransactionRefreshProvider value={{ refreshTransactions: mockRefreshTransactions }}>
          <OrderForm />
        </TransactionRefreshProvider>
      </CurrentUserProvider>
    );

    // Select term and amount
    await user.selectOptions(screen.getByLabelText('Treasury Term'), '6M');
    await user.type(screen.getByPlaceholderText('Enter amount'), '100000');

    // Submit form
    await user.click(screen.getByRole('button', { name: /buy treasury/i }));

    // Verify API called
    await waitFor(() => {
      expect(api.buyTreasury).toHaveBeenCalledWith(1, '6M', 100000);
    });

    // Verify user updated
    expect(mockUpdateUser).toHaveBeenCalledWith(updatedUser);

    // Verify transactions refreshed
    expect(mockRefreshTransactions).toHaveBeenCalled();

    // Verify success message
    expect(screen.getByText(/treasury purchased successfully/i)).toBeInTheDocument();
  });
});
\`\`\`

## Running Tests

Once testing dependencies are installed, run tests with:
```bash
npm test                  # Run all tests
npm test OrderForm       # Run OrderForm tests only
npm test -- --coverage   # Run with coverage report
```
