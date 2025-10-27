# Mobile Responsiveness Testing Guide
## Tour Component - Treasury Management Dashboard

This guide provides comprehensive testing procedures to ensure the onboarding tour works flawlessly on mobile devices.

---

## Quick Reference: Test Viewports

| Device Category | Width | Example Devices |
|----------------|-------|-----------------|
| Small Phone | 375px | iPhone SE, iPhone 12 Mini |
| Standard Phone | 390px - 428px | iPhone 13/14/15, Pixel 6 |
| Large Phone | 414px - 430px | iPhone 14 Pro Max, Galaxy S23+ |
| Tablet | 768px - 1024px | iPad, iPad Pro, Android tablets |
| Desktop | 1025px+ | Laptops, desktops |

---

## Testing Tools

### Browser DevTools (Recommended for Quick Testing)

**Chrome DevTools:**
1. Open DevTools (F12 or Cmd+Option+I on Mac)
2. Click "Toggle device toolbar" (Cmd+Shift+M on Mac)
3. Select device presets or set custom dimensions
4. Test orientation: Portrait and Landscape

**Preset Devices to Test:**
- iPhone SE (375 x 667)
- iPhone 12 Pro (390 x 844)
- iPhone 14 Pro Max (430 x 932)
- iPad (768 x 1024)
- iPad Pro (1024 x 1366)

**Safari Responsive Design Mode:**
1. Open Web Inspector (Cmd+Option+I)
2. Click "Responsive Design Mode" (Cmd+Shift+R)
3. Choose iOS device presets

### Real Device Testing (Recommended for Final Validation)

**iOS Devices:**
- iPhone SE (3rd gen) - smallest common viewport
- iPhone 13/14/15 - standard size
- iPhone 14 Pro Max - largest phone
- iPad - tablet experience
- Safari browser

**Android Devices:**
- Pixel 6/7 - standard Android
- Galaxy S23 - Samsung experience
- Chrome browser

### Tunnel for Real Device Testing

For testing on actual devices on your local network:

```bash
# Start dev server with network access
npm run dev -- --host

# Or use ngrok for external access
npx ngrok http 5173

# Or use localtunnel
npx localtunnel --port 5173
```

---

## Test Checklist

### 1. Tooltip Layout (All Viewports)

**Small Phone (375px)**
- [ ] Tooltips do NOT overflow screen horizontally
- [ ] Maximum width is 95vw (leaves 5% margin on sides)
- [ ] Content text is readable at 13px (0.8125rem)
- [ ] Line height (1.3rem) provides adequate spacing
- [ ] All text wraps properly without horizontal scroll

**Standard Phone (390px - 640px)**
- [ ] Tooltips do NOT overflow screen
- [ ] Maximum width is 90vw
- [ ] Content text is readable at 13px (0.8125rem)
- [ ] Adequate padding around content (1.25rem)

**Tablet (768px - 1024px)**
- [ ] Tooltips display at 420px max width
- [ ] Tooltips are centered and don't feel cramped
- [ ] Content is comfortable to read

**Desktop (1025px+)**
- [ ] Tooltips display at 420px max width
- [ ] Tooltips position correctly relative to targets
- [ ] Content has proper spacing

### 2. Touch Targets (Mobile & Tablet)

Touch targets should meet iOS Human Interface Guidelines (44x44px minimum).

**All Mobile Devices**
- [ ] "Next" button is at least 44px tall
- [ ] "Back" button is at least 44px tall
- [ ] "Skip" button is at least 44px tall
- [ ] All buttons have adequate padding (0.75rem+)
- [ ] Buttons don't overlap or crowd each other
- [ ] Minimum 8px gap between buttons

**Tap Testing**
- [ ] All buttons respond to single tap
- [ ] No accidental double-taps required
- [ ] Tap areas are generous (not pixel-perfect)
- [ ] Buttons provide visual feedback on tap

### 3. Tooltip Positioning

**Test Each Step on Mobile:**

For each of the 9 tour steps, verify:

| Step | Target | Expected Placement | Mobile Check |
|------|--------|-------------------|--------------|
| 1 | Welcome | Center (modal) | [ ] Centers on screen |
| 2 | User Selector | Bottom | [ ] Below header, visible |
| 3 | Balance | Bottom | [ ] Below element, doesn't overlap |
| 4 | Fund/Withdraw | Bottom | [ ] Visible, doesn't overflow |
| 5 | Yield Curve | Bottom | [ ] Below chart, readable |
| 6 | Buy Form | Left (auto on mobile) | [ ] Adapts placement, visible |
| 7 | Holdings | Top (auto on mobile) | [ ] Adapts placement, visible |
| 8 | Sell Form | Left (auto on mobile) | [ ] Adapts placement, visible |
| 9 | History | Top (auto on mobile) | [ ] Adapts placement, visible |

**Placement Adaptation:**
- [ ] Tooltips automatically reposition if target is near edge
- [ ] Arrow points to correct target element
- [ ] Tooltips never go off-screen
- [ ] Spotlight highlights target element correctly

### 4. Scrolling Behavior

**Auto-Scroll on Step Change:**
- [ ] App scrolls smoothly to target element (300ms duration)
- [ ] 20px offset prevents header overlap
- [ ] Target element is fully visible after scroll
- [ ] Tooltip appears after scroll completes
- [ ] No jarring jumps or instant scrolling

**Manual Scroll During Tour:**
- [ ] User can scroll manually while tour is active
- [ ] Tooltips reposition when page scrolls
- [ ] Overlay remains consistent during scroll
- [ ] Spotlight follows target element

### 5. Text Readability

**Font Sizes:**

| Viewport | Content Text | Button Text | Progress Text |
|----------|-------------|-------------|---------------|
| ≤375px | 13px (0.8125rem) | 14px (0.875rem) | 11px (0.6875rem) |
| 376-640px | 13px (0.8125rem) | 14px (0.875rem) | 12px (0.75rem) |
| 641px+ | 14px (0.875rem) | 14px (0.875rem) | 12px (0.75rem) |

**Readability Checks:**
- [ ] All text is sharp and anti-aliased
- [ ] Text has adequate contrast (WCAG AA minimum 4.5:1)
- [ ] Line length doesn't exceed 60-70 characters
- [ ] No text truncation or ellipsis
- [ ] Multi-line text has comfortable line height

### 6. Interaction & Accessibility

**Touch Interactions:**
- [ ] Buttons respond immediately to tap
- [ ] No hover-only interactions (mobile doesn't have hover)
- [ ] Active states provide visual feedback
- [ ] No accidental clicks on overlay

**Keyboard/Assistive Tech (Tablet):**
- [ ] Tour can be controlled via keyboard (Tab, Enter, Esc)
- [ ] Focus indicators are visible (2px blue outline)
- [ ] Screen reader announces tour content
- [ ] Skip button is easily reachable

**Gesture Support:**
- [ ] Swipe gestures don't interfere with tour
- [ ] Pinch-to-zoom disabled during tour (prevents UI breaks)
- [ ] Tour overlay prevents background touches

### 7. Performance

**Loading & Animation:**
- [ ] Tour starts within 500ms of trigger
- [ ] Target validation completes quickly (< 5 seconds)
- [ ] Transitions are smooth (60fps)
- [ ] No jank during overlay fade-in
- [ ] Tooltip animations are smooth (250ms)
- [ ] Scrolling animations don't stutter

**Memory & Battery:**
- [ ] Tour doesn't cause memory leaks
- [ ] Closing tour frees resources
- [ ] No background polling or unnecessary re-renders

### 8. Orientation Changes

**Rotate Device During Tour:**
- [ ] Tour adapts to new orientation
- [ ] Tooltips reposition correctly
- [ ] No layout breaks
- [ ] Content remains readable
- [ ] Buttons remain accessible

**Portrait vs Landscape:**
- [ ] Both orientations work equally well
- [ ] Tooltips don't feel cramped in landscape
- [ ] Portrait doesn't cause excessive scrolling

### 9. Edge Cases

**Small Content Targets:**
- [ ] Tour highlights small elements correctly
- [ ] Spotlight doesn't obscure nearby content
- [ ] Tooltip doesn't overlap target

**Long Content:**
- [ ] Long tour text wraps properly
- [ ] Multi-paragraph tooltips are readable
- [ ] Footer buttons remain accessible

**Interruptions:**
- [ ] Tour handles screen lock gracefully
- [ ] Tour resumes after app backgrounding
- [ ] Browser tab switching doesn't break tour
- [ ] Notifications don't interfere

**Network Conditions:**
- [ ] Tour works offline (no external dependencies)
- [ ] Slow network doesn't delay tour start
- [ ] Tour CSS loads before tour activates

---

## Test Procedures

### Procedure 1: Fresh User Experience (Mobile)

**Setup:**
1. Clear localStorage: `localStorage.clear()`
2. Reload page on mobile device/viewport
3. Wait for app to load fully

**Steps:**
1. [ ] Tour auto-starts within 1 second
2. [ ] Welcome step appears centered
3. [ ] Tap "Next" - transitions to User Selector
4. [ ] Tap "Next" through all 9 steps
5. [ ] Each step highlights correct element
6. [ ] Progress shows "Step X of 9"
7. [ ] Final step shows "Finish" button
8. [ ] Tap "Finish" - tour closes
9. [ ] Refresh page - tour does NOT restart

**Expected Behavior:**
- Smooth transitions between steps
- All targets visible and highlighted
- No layout shifts or breaks
- Tour completion persists

### Procedure 2: Skip Flow (Mobile)

**Setup:**
1. Clear localStorage
2. Reload page on mobile device

**Steps:**
1. [ ] Tour auto-starts
2. [ ] Tap "Skip" button immediately
3. [ ] Tour closes without showing remaining steps
4. [ ] Refresh page
5. [ ] Tour does NOT restart

**Expected Behavior:**
- Skip button is easily tappable (44px touch target)
- Tour closes immediately
- Skip is treated as completion (doesn't restart)

### Procedure 3: Manual Restart (Mobile)

**Setup:**
1. Complete tour (or skip it)
2. Locate help/tour button in header

**Steps:**
1. [ ] Tap help/tour button
2. [ ] Tour restarts from Step 1
3. [ ] Navigate through tour
4. [ ] Close or complete tour
5. [ ] Tour can be restarted again

**Expected Behavior:**
- Help button is visible and accessible
- Tour always restarts from beginning
- Unlimited manual restarts allowed

### Procedure 4: Landscape Orientation (Mobile)

**Setup:**
1. Clear localStorage
2. Rotate device to landscape before loading

**Steps:**
1. [ ] Load page in landscape mode
2. [ ] Tour starts successfully
3. [ ] All tooltips fit on screen
4. [ ] Navigate through all steps
5. [ ] Rotate to portrait mid-tour
6. [ ] Tour adapts to new orientation
7. [ ] Complete tour in portrait

**Expected Behavior:**
- Tour works in both orientations
- Smooth transition between orientations
- No layout breaks or overflow

### Procedure 5: Tablet Experience

**Setup:**
1. Use tablet device or 768px+ viewport
2. Clear localStorage

**Steps:**
1. [ ] Tour auto-starts
2. [ ] Tooltips are sized appropriately (420px max)
3. [ ] Touch targets are comfortable (40-44px)
4. [ ] All 9 steps work correctly
5. [ ] Keyboard navigation works (if available)
6. [ ] Landscape and portrait both work

**Expected Behavior:**
- Desktop-like tooltips with mobile touch targets
- Comfortable spacing and sizing
- Works with both touch and keyboard

---

## Common Issues & Solutions

### Issue: Tooltips Overflow Screen

**Symptoms:**
- Tooltip extends beyond viewport width
- Horizontal scrolling required
- Content cut off

**Solution:**
- Verify `max-width: 90vw` or `95vw` is applied
- Check `.__floater__body` CSS is not overridden
- Ensure content doesn't have fixed width elements

**File:** `frontend/src/tour/tour.css` (lines 15-60)

### Issue: Buttons Too Small on Mobile

**Symptoms:**
- Difficult to tap buttons
- Accidental taps on wrong buttons
- User frustration

**Solution:**
- Verify `min-height: 44px` for mobile buttons
- Check padding is adequate (0.75rem+)
- Ensure gap between buttons (0.5rem+)

**File:** `frontend/src/tour/tour.css` (lines 27-31, 56-59)

### Issue: Text Too Small to Read

**Symptoms:**
- Users squinting or zooming
- Complaints about readability
- Accessibility issues

**Solution:**
- Verify mobile font-size is 13px+ (0.8125rem)
- Check line-height is adequate (1.3rem+)
- Ensure contrast ratio meets WCAG AA (4.5:1)

**File:** `frontend/src/tour/tour.css` (lines 21-24)

### Issue: Tour Doesn't Start on Mobile

**Symptoms:**
- Tour never appears
- Console errors about missing targets
- Validation fails

**Solution:**
- Check all `data-tour-id` attributes exist
- Verify components render before tour starts
- Review validation logic in Tour.tsx

**Files:**
- `frontend/src/components/Tour.tsx` (lines 70-110)
- Component files with `data-tour-id` attributes

### Issue: Poor Touch Response

**Symptoms:**
- Buttons don't respond to first tap
- Need multiple taps
- Delay in interaction

**Solution:**
- Remove any hover-only styles
- Verify `:active` states are defined
- Check for event handler issues
- Ensure no competing touch handlers

**File:** `frontend/src/tour/tour.css` (lines 86-104)

---

## Validation Criteria

### ✅ Pass Criteria

Tour is considered mobile-ready when:

1. **Layout:** All tooltips fit on screen at all tested viewports (375px - 1024px+)
2. **Touch Targets:** All buttons meet 44x44px minimum (iOS HIG)
3. **Readability:** Text is comfortable to read (13px+ on mobile)
4. **Performance:** Smooth 60fps animations, no jank
5. **Accessibility:** Keyboard navigation works, focus visible
6. **Orientation:** Works in both portrait and landscape
7. **Functionality:** All 9 steps work correctly on mobile
8. **Persistence:** Tour completion is saved and respected
9. **Manual Control:** Help button allows unlimited restarts
10. **Edge Cases:** No crashes, errors, or layout breaks

### ❌ Fail Criteria

Tour requires fixes if:

- Any tooltip overflows viewport
- Buttons smaller than 40px tall
- Text smaller than 12px (0.75rem)
- Animations stutter or drop frames
- Tour doesn't adapt to orientation change
- Any step fails to display correctly
- Console errors during tour
- Touch interactions don't respond
- Accessibility violations (WCAG AA)

---

## Browser Testing Matrix

| Browser | iOS | Android | Desktop |
|---------|-----|---------|---------|
| Safari | ✅ Required | N/A | ⚪ Optional |
| Chrome | ⚪ Optional | ✅ Required | ✅ Required |
| Firefox | N/A | ⚪ Optional | ⚪ Optional |
| Edge | N/A | N/A | ⚪ Optional |

**Legend:**
- ✅ Required for release
- ⚪ Optional but recommended

---

## Automated Testing (Future Enhancement)

While manual testing is required for MVP, consider adding automated tests:

**Playwright/Cypress Mobile Tests:**
```javascript
// Example: Test tour on mobile viewport
test('tour displays correctly on mobile', async () => {
  await page.setViewportSize({ width: 375, height: 667 });
  await page.goto('/');

  // Verify tour starts
  await expect(page.locator('.react-joyride__tooltip')).toBeVisible();

  // Check tooltip width
  const tooltip = await page.locator('.__floater__body');
  const box = await tooltip.boundingBox();
  expect(box.width).toBeLessThanOrEqual(375 * 0.95); // 95vw

  // Verify buttons are large enough
  const nextButton = await page.locator('button[aria-label*="Next"]');
  const buttonBox = await nextButton.boundingBox();
  expect(buttonBox.height).toBeGreaterThanOrEqual(44);
});
```

**Visual Regression Testing:**
- Percy.io for screenshot comparisons
- Chromatic for Storybook visual testing
- BackstopJS for DIY visual regression

---

## Success Metrics

Track these metrics to validate mobile experience:

1. **Completion Rate:** % of users who complete tour on mobile
2. **Skip Rate:** % of users who skip tour on mobile
3. **Time to Complete:** Average time to finish all 9 steps
4. **Error Rate:** Console errors during tour
5. **Device Coverage:** % of mobile devices tested
6. **User Feedback:** Qualitative feedback from real users

**Target Metrics (MVP):**
- Completion Rate: >60%
- Skip Rate: <30%
- Error Rate: 0%
- Device Coverage: 100% (all common devices)

---

## Testing Schedule

**Pre-Release (MVP):**
- [ ] Day 1: Browser DevTools testing (all viewports)
- [ ] Day 2: Real device testing (iPhone, Android)
- [ ] Day 3: Edge case testing and bug fixes

**Post-Release (Monitoring):**
- Week 1: Monitor analytics for completion rates
- Week 2: Collect user feedback
- Month 1: Analyze metrics and plan improvements

---

## Additional Resources

**React Joyride Docs:**
- https://docs.react-joyride.com/
- Styling: https://docs.react-joyride.com/styling
- Props: https://docs.react-joyride.com/props

**Mobile Testing Tools:**
- BrowserStack: https://www.browserstack.com/
- Sauce Labs: https://saucelabs.com/
- Chrome DevTools Device Mode: https://developer.chrome.com/docs/devtools/device-mode/

**Accessibility Guidelines:**
- iOS Human Interface Guidelines: https://developer.apple.com/design/human-interface-guidelines/
- Material Design Touch Targets: https://m3.material.io/foundations/accessible-design/accessibility-basics
- WCAG 2.1 (AA): https://www.w3.org/WAI/WCAG21/quickref/

---

## Conclusion

This guide ensures comprehensive mobile testing coverage for the onboarding tour. Follow all procedures and check all items before considering the tour mobile-ready.

**Remember:** Real device testing is irreplaceable. Always test on actual devices before release.

**Next Steps After Testing:**
1. Document any issues found
2. Fix critical issues immediately
3. Prioritize nice-to-have improvements
4. Update this guide with lessons learned
5. Mark section 5.3 as complete in implementation plan
