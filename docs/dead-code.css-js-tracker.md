# Dead Code Tracker (CSS + JS)

Date: 2026-04-04
Scope: `/home/glen/quox` `static/css` + `static/js` with usage traced through Go render code and JS selectors/class operations.

## Completed In This Pass

### Removed files not referenced/loaded

- `static/css/page.1.quote.css`
- `static/js/date.parts.js`

### Unused selectors now commented out in active CSS files

- `static/css/page.1.quote.desktop.css`
  - `.qbuy-head-title`
  - `.quote-desk-results-title`
  - `.quote-desk-plan-toolbar`
  - `.quote-plan-sort-label`
- `static/css/page.1.quote.phone.css`
  - `.qbuy-head-title`
  - `.quote-plan-sort-label`
  - `.quote-plan-sums`
- `static/css/page.2.editq.desktop.css`
  - `.editq-cust-preview`
  - `.editq-cust-empty`
  - `.editq-main`
  - `.editq-condition-list`
  - `.editq-prime-plan-total`
  - `.editq-prime-plan-charges`
  - `.editq-prime-plan-sum`
  - `.editq-prime-plan-op`
  - `.editq-field`
  - `.editq-label`
  - `.editq-dependent-head`
  - `.editq-dependent-title`
  - `.editq-dependent-head .editq-del-btn` (rule dead because `.editq-dependent-head` is unused)
- `static/css/page.2.editq.phone.css`
  - `.editq-cust-preview`
  - `.editq-cust-empty`
  - `.editq-condition-list`
  - `.editq-prime-plan-total`
  - `.editq-prime-plan-charges`
  - `.editq-prime-plan-sum`
  - `.editq-prime-plan-op`
  - `.editq-field`
  - `.editq-label`
  - `.editq-dependent-head`
  - `.editq-dependent-head-text`
  - `.editq-dependent-title`
  - `.editq-dependent-head .editq-del-btn` (rule dead because `.editq-dependent-head` is unused)

## JS Dead Code Status

- No dead internal code found in loaded JS files:
  - `static/js/page.1.quote.js`
  - `static/js/page.1.quote.buy.js`
  - `static/js/page.2.editq.js`
- Dead JS found was only the whole unreferenced file `static/js/date.parts.js` (already removed).

## Keep: Not Dead

- `quote-span-1` .. `quote-span-12`
  - Generated dynamically by `QuoteSpanClass(span)` in Go (`"quote-span-" + span`), so literal class search can miss them.
  - These are active and should not be deleted.

## Verification commands used

- Find unreferenced CSS classes:
  - `classes=$(rg -o --no-filename '\\.[_a-zA-Z][-_a-zA-Z0-9]*' static/css/*.css | sed 's/^\\.//' | sort -u); for c in $classes; do if ! rg -q --glob '*.go' --glob 'static/js/*.js' "(^|[^A-Za-z0-9_-])${c}([^A-Za-z0-9_-]|$)"; then echo "$c"; fi; done | sort`
- Check asset-path references:
  - `rg -n "page\\.1\\.quote\\.css|date\\.parts\\.js" -S`
