/**
 * Custom Joyride styles to match Tremor/Tailwind aesthetic.
 *
 * Design principles:
 * - Follows Tremor design tokens exactly (colors, shadows, border radius, typography)
 * - Primary color: Tremor brand blue (#3B82F6 / blue-500)
 * - Background: Clean white with tremor-dropdown shadow (elevated UI)
 * - Text: Tremor content colors (gray-700 emphasis, gray-900 strong)
 * - Border radius: tremor-default (0.5rem) for cards, tremor-small (0.375rem) for buttons
 * - Typography: tremor-default (0.875rem / 1.25rem line-height)
 * - z-index: High enough to overlay all content (10000)
 * - Responsive: Adapts width for mobile, tablet, and desktop
 * - Spotlight: Enhanced padding for better visual hierarchy
 */
export const TOUR_STYLES = {
  options: {
    // Primary color - Tremor brand blue-500
    primaryColor: '#3B82F6',

    // Background - Clean white (tremor-background-DEFAULT)
    backgroundColor: '#FFFFFF',

    // Text color - Tremor content-emphasis (gray-700)
    textColor: '#374151',

    // Arrow color - Match background for seamless appearance
    arrowColor: '#FFFFFF',

    // Overlay color - Semi-transparent dark for focus
    overlayColor: 'rgba(0, 0, 0, 0.5)',

    // z-index - High enough to overlay all content
    zIndex: 10000,

    // Width - Responsive (mobile-friendly default, max-width handled by CSS)
    width: 420,

    // Beacon size - Standard Joyride size
    beaconSize: 36,
  },

  // Tooltip container styling
  tooltip: {
    // Border radius - tremor-default (matches Cards)
    borderRadius: '0.5rem',

    // Box shadow - tremor-dropdown (elevated UI like modals/dropdowns)
    // This is more prominent than tremor-card for better visibility
    boxShadow: '0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1)',

    // No padding on container (padding applied to content/footer)
    padding: 0,

    // Responsive width (max 90% of viewport on mobile)
    maxWidth: '90vw',

    // Smooth transitions
    transition: 'all 0.3s ease-in-out',
  },

  // Tooltip content area
  tooltipContent: {
    // Padding - comfortable spacing (tremor standard)
    padding: '1.25rem',

    // Typography - tremor-default
    fontSize: '0.875rem',
    lineHeight: '1.25rem',

    // Text color - Tremor content-emphasis
    color: '#374151',

    // Background color - ensure white for contrast
    backgroundColor: '#FFFFFF',
  },

  // Tooltip title (if used)
  tooltipTitle: {
    // Font size - tremor-title
    fontSize: '1.125rem',
    lineHeight: '1.75rem',

    // Font weight - semibold
    fontWeight: '600',

    // Color - tremor-content-strong
    color: '#111827', // gray-900

    // Margin bottom for spacing
    marginBottom: '0.5rem',
  },

  // Tooltip footer (button area)
  tooltipFooter: {
    // Spacing
    marginTop: '0.75rem',
    padding: '0 1.25rem 1.25rem 1.25rem',

    // Display flex for button layout
    display: 'flex',
    gap: '0.75rem',
    justifyContent: 'space-between',
    alignItems: 'center',
  },

  // Button - Next/Primary action
  buttonNext: {
    // Background - Tremor brand blue-500
    backgroundColor: '#3B82F6',

    // Border radius - tremor-small
    borderRadius: '0.375rem',

    // Text color - white
    color: '#FFFFFF',

    // Typography - tremor-default, medium weight
    fontSize: '0.875rem',
    fontWeight: '500',
    lineHeight: '1.25rem',

    // Spacing
    padding: '0.5rem 1rem',

    // Reset default button styles
    outline: 'none',
    border: 'none',

    // Cursor
    cursor: 'pointer',

    // Smooth transitions
    transition: 'background-color 0.2s ease-in-out, transform 0.1s ease-in-out',

    // Hover state (inline style, but documented for reference)
    // :hover { backgroundColor: '#2563EB' } // blue-600
  },

  // Button - Back/Secondary action
  buttonBack: {
    // Text color - Tremor content-DEFAULT (gray-500)
    color: '#6B7280',

    // Typography - tremor-default, medium weight
    fontSize: '0.875rem',
    fontWeight: '500',
    lineHeight: '1.25rem',

    // Spacing
    marginRight: '0.5rem',
    padding: '0.5rem 0.75rem',

    // Reset default button styles
    outline: 'none',
    border: 'none',
    backgroundColor: 'transparent',

    // Cursor
    cursor: 'pointer',

    // Smooth transitions
    transition: 'color 0.2s ease-in-out',

    // Hover state (inline style, but documented for reference)
    // :hover { color: '#374151' } // gray-700
  },

  // Button - Skip
  buttonSkip: {
    // Text color - Tremor content-DEFAULT (gray-500)
    color: '#6B7280',

    // Typography - tremor-default, medium weight
    fontSize: '0.875rem',
    fontWeight: '500',
    lineHeight: '1.25rem',

    // Spacing
    padding: '0.5rem 0.75rem',

    // Reset default button styles
    outline: 'none',
    border: 'none',
    backgroundColor: 'transparent',

    // Cursor
    cursor: 'pointer',

    // Smooth transitions
    transition: 'color 0.2s ease-in-out',

    // Hover state (inline style, but documented for reference)
    // :hover { color: '#374151' } // gray-700
  },

  // Beacon (pulsing circle indicator)
  beacon: {
    outline: 'none',
  },

  beaconInner: {
    // Background - Tremor brand blue-500
    backgroundColor: '#3B82F6',
  },

  beaconOuter: {
    // Background - blue-500 with 20% opacity
    backgroundColor: 'rgba(59, 130, 246, 0.2)',

    // Border - solid blue-500
    border: '2px solid #3B82F6',
  },

  // Spotlight (highlight area around target element)
  spotlight: {
    // Border radius - tremor-default (matches cards)
    borderRadius: '0.5rem',

    // Enhanced padding for better visual hierarchy
    // Adds more breathing room around highlighted elements
    padding: 8, // 8px padding around spotlight area

    // Smooth transitions
    transition: 'all 0.3s ease-in-out',
  },

  // Overlay (dark background behind spotlight)
  overlay: {
    // Mix blend mode for better contrast
    mixBlendMode: 'hard-light' as const,

    // Smooth transitions
    transition: 'opacity 0.3s ease-in-out',
  },
};
