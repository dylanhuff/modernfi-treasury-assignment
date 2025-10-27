/**
 * LocalStorage utility functions for tour completion tracking
 *
 * These functions handle persisting the tour completion state to prevent
 * repeated tours while allowing manual restarts via the help button.
 */

const TOUR_STORAGE_KEY = 'treasuryApp_tourCompleted';

/**
 * Retrieves the tour completion status from localStorage
 *
 * @returns {boolean} true if tour has been completed, false otherwise
 *
 * Error Handling:
 * - Returns false if localStorage is unavailable (private browsing, etc.)
 * - Returns false if the stored value is invalid
 * - Returns false on any unexpected errors
 */
export function getTourCompleted(): boolean {
  try {
    // Check if localStorage is available
    if (typeof window === 'undefined' || !window.localStorage) {
      console.warn('localStorage is not available');
      return false;
    }

    const value = localStorage.getItem(TOUR_STORAGE_KEY);

    // If no value stored, tour hasn't been completed
    if (value === null) {
      return false;
    }

    // Parse the stored value as boolean
    return value === 'true';
  } catch (error) {
    // Handle localStorage access errors (e.g., SecurityError in private browsing)
    console.error('Error reading tour completion status from localStorage:', error);
    return false;
  }
}

/**
 * Saves the tour completion status to localStorage
 *
 * @param {boolean} completed - Whether the tour has been completed
 *
 * Error Handling:
 * - Logs warning if localStorage is unavailable
 * - Catches and logs any storage errors (quota exceeded, etc.)
 */
export function setTourCompleted(completed: boolean): void {
  try {
    // Check if localStorage is available
    if (typeof window === 'undefined' || !window.localStorage) {
      console.warn('localStorage is not available, cannot save tour completion status');
      return;
    }

    // Store the boolean value as a string
    localStorage.setItem(TOUR_STORAGE_KEY, String(completed));
  } catch (error) {
    // Handle localStorage write errors (e.g., QuotaExceededError)
    console.error('Error saving tour completion status to localStorage:', error);
  }
}

