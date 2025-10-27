# CurrentUserContext Usage Guide

## Overview

The `CurrentUserContext` provides a centralized way to manage user selection across your application. It handles user fetching, selection persistence via localStorage, and provides an easy-to-use hook for accessing current user data in any component.

---

## Quick Start

### 1. Using the `useCurrentUser` Hook

The `useCurrentUser` hook is the primary way to access user data in your components.

**Basic Usage:**
```tsx
import { useCurrentUser } from '../contexts/CurrentUserContext';

function MyComponent() {
  const { currentUser, isLoading, error } = useCurrentUser();

  if (isLoading) {
    return <div>Loading user data...</div>;
  }

  if (error) {
    return <div>Error: {error}</div>;
  }

  if (!currentUser) {
    return <div>No user selected</div>;
  }

  return (
    <div>
      <h1>Welcome, {currentUser.name}!</h1>
      <p>Balance: ${currentUser.balance.toLocaleString()}</p>
    </div>
  );
}
```

---

## Common Use Cases

### Accessing Current User in New Components

Any component within the `CurrentUserProvider` tree can access the current user:

```tsx
import { useCurrentUser } from '../contexts/CurrentUserContext';

function TransactionHistory() {
  const { currentUser } = useCurrentUser();

  // Use currentUser.id to fetch user-specific data
  useEffect(() => {
    if (currentUser) {
      fetchTransactions(currentUser.id);
    }
  }, [currentUser]);

  return (
    <div>
      <h2>Transactions for {currentUser?.name}</h2>
      {/* Your transaction list */}
    </div>
  );
}
```

### Changing the Selected User

Use `setCurrentUserById` to programmatically change the current user:

```tsx
import { useCurrentUser } from '../contexts/CurrentUserContext';

function UserSwitcher() {
  const { users, currentUser, setCurrentUserById } = useCurrentUser();

  return (
    <div>
      <h3>Current: {currentUser?.name}</h3>
      <select
        value={currentUser?.id || ''}
        onChange={(e) => setCurrentUserById(Number(e.target.value))}
      >
        {users.map(user => (
          <option key={user.id} value={user.id}>
            {user.name}
          </option>
        ))}
      </select>
    </div>
  );
}
```

### Displaying User Balance

```tsx
import { useCurrentUser } from '../contexts/CurrentUserContext';

function BalanceDisplay() {
  const { currentUser } = useCurrentUser();

  const formatBalance = (balance: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 2,
    }).format(balance);
  };

  return (
    <div className="balance-display">
      <span className="label">Balance:</span>
      <span className="amount">
        {currentUser ? formatBalance(currentUser.balance) : 'N/A'}
      </span>
    </div>
  );
}
```

### User-Specific API Calls

```tsx
import { useCurrentUser } from '../contexts/CurrentUserContext';
import { useState, useEffect } from 'react';

function UserOrders() {
  const { currentUser } = useCurrentUser();
  const [orders, setOrders] = useState([]);

  useEffect(() => {
    if (!currentUser) return;

    // Fetch orders for the current user
    fetch(`/api/v1/users/${currentUser.id}/orders`)
      .then(res => res.json())
      .then(setOrders)
      .catch(console.error);
  }, [currentUser?.id]); // Re-fetch when user changes

  return (
    <div>
      <h2>Orders for {currentUser?.name}</h2>
      {orders.map(order => (
        <div key={order.id}>{/* Order details */}</div>
      ))}
    </div>
  );
}
```

### Conditional Rendering Based on User

```tsx
import { useCurrentUser } from '../contexts/CurrentUserContext';

function AdminPanel() {
  const { currentUser } = useCurrentUser();

  // Example: Check if user has sufficient balance for a feature
  const canAccessPremiumFeatures = currentUser && currentUser.balance > 100000;

  if (!canAccessPremiumFeatures) {
    return <div>Insufficient balance for premium features</div>;
  }

  return <div>{/* Premium features */}</div>;
}
```

---

## Complete Component Example

Here's a full example showing multiple aspects of the context:

```tsx
import { useCurrentUser } from '../contexts/CurrentUserContext';
import { useState } from 'react';

function CompleteExample() {
  const {
    users,
    currentUser,
    setCurrentUserById,
    isLoading,
    error
  } = useCurrentUser();

  const [localState, setLocalState] = useState('');

  // Handle loading state
  if (isLoading) {
    return <div className="loading">Loading users...</div>;
  }

  // Handle error state
  if (error) {
    return <div className="error">Error: {error}</div>;
  }

  // Handle no user selected
  if (!currentUser) {
    return <div>No user selected</div>;
  }

  return (
    <div className="user-dashboard">
      {/* User selector */}
      <header>
        <select
          value={currentUser.id}
          onChange={(e) => setCurrentUserById(Number(e.target.value))}
          className="user-select"
        >
          {users.map(user => (
            <option key={user.id} value={user.id}>
              {user.name}
            </option>
          ))}
        </select>
      </header>

      {/* Display user info */}
      <main>
        <h1>Welcome, {currentUser.name}!</h1>
        <p>Account Balance: ${currentUser.balance.toLocaleString()}</p>
        <p>Account created: {new Date(currentUser.created_at).toLocaleDateString()}</p>

        {/* User-specific action */}
        <button onClick={() => console.log('Action for user:', currentUser.id)}>
          Perform Action
        </button>
      </main>
    </div>
  );
}

export default CompleteExample;
```

---

## Context API Reference

### Available Properties

When you call `useCurrentUser()`, you get access to:

| Property | Type | Description |
|----------|------|-------------|
| `users` | `User[]` | Array of all available users |
| `currentUser` | `User \| null` | Currently selected user (null if none selected) |
| `setCurrentUserById` | `(userId: number) => void` | Function to change the current user |
| `isLoading` | `boolean` | True while fetching users from API |
| `error` | `string \| null` | Error message if fetch fails, null otherwise |

### User Type Definition

```typescript
interface User {
  id: number;
  name: string;
  balance: number;
  created_at: string;
}
```

---

## Important Notes

### 1. Always Check for null

The `currentUser` can be `null` if no user is loaded yet or if the user list is empty. Always use optional chaining or null checks:

```tsx
// Good
currentUser?.name

// Also good
if (currentUser) {
  // Use currentUser.name
}
```

### 2. Hook Must Be Used Inside Provider

The `useCurrentUser` hook will throw an error if used outside of `CurrentUserProvider`. Make sure your component tree is wrapped:

```tsx
// In App.tsx
<CurrentUserProvider>
  <YourComponents />
</CurrentUserProvider>
```

### 3. localStorage Persistence

User selection is automatically saved to localStorage with key `'currentUserId'`. This persists across page refreshes and browser sessions.

### 4. Re-renders When User Changes

Components using `useCurrentUser` will re-render when:
- The current user changes (via `setCurrentUserById`)
- Users are loaded from the API
- Loading or error state changes

---

## Troubleshooting

### "useCurrentUser must be used within a CurrentUserProvider"

**Problem:** You're using the hook outside the provider.

**Solution:** Make sure your component is wrapped in `<CurrentUserProvider>`:

```tsx
// App.tsx
import { CurrentUserProvider } from './contexts/CurrentUserContext';

function App() {
  return (
    <CurrentUserProvider>
      <YourComponent /> {/* Can now use useCurrentUser */}
    </CurrentUserProvider>
  );
}
```

### User Selection Not Persisting

**Problem:** Selected user doesn't persist across page refreshes.

**Solution:** Check browser localStorage is enabled. The context automatically handles persistence via localStorage key `'currentUserId'`.

### currentUser is null

**Problem:** `currentUser` is null even after loading.

**Possible causes:**
1. No users in the database
2. API call failed (check `error` property)
3. Still loading (check `isLoading` property)

**Solution:**
```tsx
const { currentUser, isLoading, error } = useCurrentUser();

if (isLoading) return <div>Loading...</div>;
if (error) return <div>Error: {error}</div>;
if (!currentUser) return <div>No users available</div>;
```

---

## Best Practices

1. **Always handle loading and error states** in your components
2. **Use optional chaining** when accessing currentUser properties
3. **Don't destructure currentUser properties directly** - check for null first
4. **Depend on currentUser.id in useEffect** when fetching user-specific data
5. **Keep the context focused on user selection** - don't add unrelated state

---

## Integration with Future Features

The CurrentUserContext is designed to support upcoming features:

### Fund/Withdraw
```tsx
const { currentUser } = useCurrentUser();

const handleFund = async (amount: number) => {
  await fetch(`/api/v1/users/${currentUser?.id}/fund`, {
    method: 'POST',
    body: JSON.stringify({ amount }),
  });
  // Refresh user data to show updated balance
};
```

### Buy/Sell Orders
```tsx
const { currentUser } = useCurrentUser();

const handleBuyOrder = async (term: string, amount: number) => {
  await fetch(`/api/v1/users/${currentUser?.id}/orders`, {
    method: 'POST',
    body: JSON.stringify({ type: 'buy', term, amount }),
  });
};
```

### Transaction History
```tsx
const { currentUser } = useCurrentUser();

useEffect(() => {
  if (currentUser) {
    fetch(`/api/v1/users/${currentUser.id}/transactions`)
      .then(res => res.json())
      .then(setTransactions);
  }
}, [currentUser]); // Automatically re-fetches when user switches
```

---

## Summary

The `useCurrentUser` hook provides everything you need for user-aware components:

```tsx
const {
  users,           // All users
  currentUser,     // Selected user
  setCurrentUserById, // Change selection
  isLoading,       // Loading state
  error            // Error message
} = useCurrentUser();
```

Use it in any component to access the current user's data, display their information, or build user-specific features.
