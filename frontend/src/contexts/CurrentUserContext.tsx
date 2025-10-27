/* eslint-disable react-refresh/only-export-components */
import { createContext, useContext, useState, useEffect, useMemo, useCallback } from 'react';
import type { ReactNode } from 'react';
import type { User } from '../types/user';
import { fetchUsers } from '../services/api';

/**
 * localStorage key for persisting the selected user's ID across sessions.
 * Stored as a string representation of the user ID number.
 */
const LOCAL_STORAGE_KEY = 'currentUserId';

/**
 * Context type definition for user selection state management.
 *
 * @property {User[]} users - Complete list of all available users fetched from the API
 * @property {User | null} currentUser - The currently selected user, or null if no user is selected
 * @property {function} setCurrentUserById - Function to change the current user by their ID
 * @property {boolean} isLoading - Loading state indicator for the initial user fetch
 * @property {string | null} error - Error message if user fetch fails, null otherwise
 */
interface CurrentUserContextType {
  users: User[];
  currentUser: User | null;
  setCurrentUserById: (userId: number) => void;
  updateUser: (user: User) => void;
  refreshUser: () => Promise<void>;
  isLoading: boolean;
  error: string | null;
}

const CurrentUserContext = createContext<CurrentUserContextType | undefined>(undefined);

interface CurrentUserProviderProps {
  children: ReactNode;
}

/**
 * Provider component for user selection state management with localStorage persistence.
 *
 * This component fetches all available users from the API on mount and manages the currently
 * selected user. User selection is persisted to localStorage and restored on page refresh.
 *
 * **localStorage Persistence Logic:**
 * 1. On mount, fetches all users from the API
 * 2. Attempts to restore previously selected user from localStorage using the stored ID
 * 3. If stored ID exists and matches a user in the fetched list, that user is selected
 * 4. If no valid stored ID, defaults to the first user in the list
 * 5. Whenever user selection changes via setCurrentUserById, the new ID is saved to localStorage
 *
 * **AbortController Pattern:**
 * Uses AbortController to prevent state updates if the component unmounts during the fetch.
 * This prevents "Can't perform a React state update on an unmounted component" warnings.
 *
 * @param {CurrentUserProviderProps} props - Component props
 * @param {ReactNode} props.children - Child components that will have access to the context
 *
 * @example
 * ```tsx
 * <CurrentUserProvider>
 *   <App />
 * </CurrentUserProvider>
 * ```
 */
export function CurrentUserProvider({ children }: CurrentUserProviderProps) {
  const [users, setUsers] = useState<User[]>([]);
  const [currentUser, setCurrentUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const controller = new AbortController();
    const signal = controller.signal;

    async function loadUsers() {
      try {
        setIsLoading(true);
        setError(null);

        const fetchedUsers = await fetchUsers();

        // Check if component unmounted during fetch
        if (signal.aborted) return;

        setUsers(fetchedUsers);

        // Handle edge case: no users in the system
        if (fetchedUsers.length === 0) {
          setCurrentUser(null);
          setIsLoading(false);
          return;
        }

        // STEP 1: Try to restore previously selected user from localStorage
        // This provides persistence across page refreshes and browser sessions
        const storedUserId = localStorage.getItem(LOCAL_STORAGE_KEY);

        if (storedUserId) {
          const parsedId = parseInt(storedUserId, 10);
          // STEP 2: Validate that the stored ID still exists in the fetched user list
          // This handles cases where a user might have been deleted from the database
          const storedUser = fetchedUsers.find(user => user.id === parsedId);

          if (storedUser) {
            // Success: Restored user from localStorage
            setCurrentUser(storedUser);
            setIsLoading(false);
            return;
          }
          // If stored user not found, fall through to default behavior
        }

        // STEP 3: Default to first user if no valid stored user
        // This happens on first visit or when stored user no longer exists
        const firstUser = fetchedUsers[0];
        setCurrentUser(firstUser);
        // Save the default selection to localStorage for future visits
        localStorage.setItem(LOCAL_STORAGE_KEY, firstUser.id.toString());

        setIsLoading(false);
      } catch (err) {
        if (signal.aborted) return;

        const errorMessage = err instanceof Error ? err.message : 'Failed to load users';
        setError(errorMessage);
        setIsLoading(false);
      }
    }

    loadUsers();

    return () => {
      controller.abort();
    };
  }, []);

  const setCurrentUserById = useCallback((userId: number) => {
    const user = users.find(u => u.id === userId);
    if (user) {
      setCurrentUser(user);
      localStorage.setItem(LOCAL_STORAGE_KEY, userId.toString());
    }
  }, [users]);

  const updateUser = useCallback((user: User) => {
    setCurrentUser(user);
    // Also update the user in the users array
    setUsers(prevUsers => prevUsers.map(u => u.id === user.id ? user : u));
  }, []);

  const refreshUser = useCallback(async () => {
    if (!currentUser) return;

    try {
      const fetchedUsers = await fetchUsers();
      setUsers(fetchedUsers);

      const updated = fetchedUsers.find(u => u.id === currentUser.id);
      if (updated) {
        setCurrentUser(updated);
      }
    } catch {
      // Silently handle refresh errors
    }
  }, [currentUser]);

  const value = useMemo(() => ({
    users,
    currentUser,
    setCurrentUserById,
    updateUser,
    refreshUser,
    isLoading,
    error,
  }), [users, currentUser, setCurrentUserById, updateUser, refreshUser, isLoading, error]);

  return (
    <CurrentUserContext.Provider value={value}>
      {children}
    </CurrentUserContext.Provider>
  );
}

/**
 * Custom hook to access the current user context.
 *
 * Provides access to user selection state including the list of all users,
 * the currently selected user, and the ability to change the selection.
 *
 * **IMPORTANT:** This hook must be used within a component that is wrapped by CurrentUserProvider.
 * Attempting to use it outside the provider will throw an error.
 *
 * @returns {CurrentUserContextType} The current user context value
 * @throws {Error} If used outside of CurrentUserProvider
 *
 * @example
 * ```tsx
 * function MyComponent() {
 *   const { currentUser, setCurrentUserById, isLoading } = useCurrentUser();
 *
 *   if (isLoading) return <div>Loading...</div>;
 *
 *   return (
 *     <div>
 *       <p>Current user: {currentUser?.name}</p>
 *       <button onClick={() => setCurrentUserById(2)}>
 *         Switch to user 2
 *       </button>
 *     </div>
 *   );
 * }
 * ```
 */
export function useCurrentUser(): CurrentUserContextType {
  const context = useContext(CurrentUserContext);

  if (context === undefined) {
    throw new Error('useCurrentUser must be used within a CurrentUserProvider');
  }

  return context;
}
