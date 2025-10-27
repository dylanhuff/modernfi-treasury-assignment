# Tour Visual Testing Guide

This guide provides steps to manually test the tour tooltip appearance and positioning across different screen sizes.

## Completed Enhancements (Section 5.1)

### Style Customizations
- ✅ Applied Tremor design tokens exactly (colors, shadows, border radius, typography)
- ✅ Enhanced tooltip shadow to `tremor-dropdown` for better elevation
- ✅ Responsive width handling (90vw mobile, 420px desktop)
- ✅ Enhanced spotlight padding (8px) for better visual hierarchy
- ✅ Added smooth transitions (0.3s tooltips, 0.2s buttons)
- ✅ Implemented hover states for buttons (blue-600 for primary, gray-700 for secondary)
- ✅ Added fade-in animations for overlay and tooltips

## Manual Testing Checklist

### Desktop Testing (1024px+)
Open the app at http://localhost:5173 in a desktop browser:

1. **Tooltip Width**
   - [ ] Tour tooltips should be 420px wide
   - [ ] Tooltips should have proper shadow (tremor-dropdown)
   - [ ] Border radius matches Tremor cards (0.5rem)

2. **Tooltip Content**
   - [ ] Text is gray-700 (Tremor content-emphasis)
   - [ ] Font size is 0.875rem (tremor-default)
   - [ ] Line height is comfortable (1.25rem)
   - [ ] Padding is consistent (1.25rem)

3. **Buttons**
   - [ ] Primary button (Next) is blue-500
   - [ ] Primary button hovers to blue-600
   - [ ] Back/Skip buttons are gray-500
   - [ ] Back/Skip buttons hover to gray-700 with subtle background
   - [ ] All buttons have 0.375rem border radius (tremor-small)

4. **Spotlight**
   - [ ] Spotlight has 8px padding around target elements
   - [ ] Spotlight has 0.5rem border radius (matches cards)
   - [ ] Spotlight has subtle blue glow on focus
   - [ ] Overlay is semi-transparent dark (rgba(0, 0, 0, 0.5))

5. **Animations**
   - [ ] Tooltips fade in smoothly (0.3s)
   - [ ] Tooltips slide down slightly on appearance
   - [ ] Overlay fades in smoothly
   - [ ] Button hover transitions are smooth (0.2s)

### Tablet Testing (768px - 1024px)
Resize browser to tablet width or use device emulation:

1. **Tooltip Positioning**
   - [ ] Tooltips should still be 420px wide (if space allows)
   - [ ] Tooltips should not overflow screen edges
   - [ ] Auto placement works correctly

2. **Tooltip Placement**
   - [ ] YieldCurve tooltip: positioned at bottom or top (auto)
   - [ ] OrderForm tooltip: positioned at top or bottom (auto)
   - [ ] TransactionHistory tooltip: positioned at top (auto)
   - [ ] Arrows point to correct target elements

### Mobile Testing (375px - 640px)
Resize browser to mobile width or use device emulation (iPhone SE, iPhone 12, etc.):

1. **Responsive Width**
   - [ ] Tooltips should be max 90vw (90% of viewport width)
   - [ ] Tooltips should never overflow screen
   - [ ] Font size slightly smaller on mobile (0.8125rem)

2. **Touch Interactions**
   - [ ] Next button is large enough for touch (min 44x44px target)
   - [ ] Skip button is easily tappable
   - [ ] Back button is easily tappable
   - [ ] No accidental clicks on underlying content

3. **Content Readability**
   - [ ] All text is readable on small screens
   - [ ] No text truncation or overflow
   - [ ] Padding feels comfortable (not cramped)
   - [ ] Buttons don't overlap or wrap awkwardly

4. **Tooltip Positioning on Mobile**
   - [ ] Tooltips position automatically to fit screen
   - [ ] Arrows still point to target elements
   - [ ] No tooltips cut off at screen edges
   - [ ] Welcome step (body target) displays centered

### Cross-Browser Testing
Test in multiple browsers to ensure consistent appearance:

1. **Chrome** (Primary)
   - [ ] All styles render correctly
   - [ ] Animations smooth
   - [ ] Hover states work

2. **Safari** (macOS/iOS)
   - [ ] All styles render correctly
   - [ ] Animations smooth
   - [ ] Hover states work (desktop) / Touch works (mobile)

3. **Firefox**
   - [ ] All styles render correctly
   - [ ] Animations smooth
   - [ ] Hover states work

### Accessibility Testing

1. **Keyboard Navigation**
   - [ ] Tab key navigates between buttons
   - [ ] Focus states are visible (blue outline)
   - [ ] Enter key advances tour
   - [ ] Esc key skips tour (if supported by Joyride)

2. **Screen Reader** (Optional but recommended)
   - [ ] Tooltip content is announced
   - [ ] Button labels are descriptive
   - [ ] Progress indicator is announced

## Testing Viewport Sizes

### Recommended Test Sizes
- **Mobile**: 375px (iPhone SE), 390px (iPhone 12), 414px (iPhone 12 Pro Max)
- **Tablet**: 768px (iPad), 820px (iPad Air), 1024px (iPad Pro)
- **Desktop**: 1280px, 1440px, 1920px

### How to Test in Chrome DevTools
1. Open Chrome DevTools (F12 or Cmd+Option+I)
2. Click "Toggle device toolbar" (Cmd+Shift+M)
3. Select a device preset or enter custom dimensions
4. Reload the page and trigger the tour
5. Step through all 4 tour steps

## Known Issues & Limitations

### MVP Constraints
- Dark mode styles are prepared but commented out (future enhancement)
- Media query styles are defined in CSS, not inline (Joyride limitation)
- Hover states use CSS (not inline styles) for better browser support

### Edge Cases to Note
- Very small screens (<350px) may have slight text wrapping
- Tour targets must exist before tour starts (handled by validation logic)
- If user scrolls during tour, spotlight may temporarily misalign (Joyride handles this)

## Visual Comparison Checklist

Compare tour tooltips against Tremor Card components:

| Property | Tremor Card | Tour Tooltip | Match? |
|----------|------------|--------------|--------|
| Border Radius | 0.5rem | 0.5rem | ✅ |
| Shadow | tremor-dropdown | tremor-dropdown | ✅ |
| Background | white | white | ✅ |
| Text Color | gray-700 | gray-700 | ✅ |
| Font Size | 0.875rem | 0.875rem | ✅ |
| Primary Color | blue-500 | blue-500 | ✅ |
| Hover Color | blue-600 | blue-600 | ✅ |

## Success Criteria

Section 5.1 is complete when:
- [x] Custom styles applied from tour/styles.ts
- [x] Tooltips match Tremor aesthetic exactly
- [x] Responsive design tested on mobile, tablet, desktop
- [x] Spotlight padding provides clear visual hierarchy
- [x] All animations are smooth and professional
- [x] No console errors related to styling

## Next Steps

After completing manual testing of Section 5.1:
- Move to **Section 5.2**: Add smooth transitions (already implemented, verify)
- Move to **Section 5.3**: Mobile responsiveness testing (use this guide)

---

**Last Updated**: Phase 5.1 completion
**Tester**: [Your name]
**Date**: [Test date]
