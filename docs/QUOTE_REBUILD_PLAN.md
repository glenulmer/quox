# Quote Rebuild Plan

This file is for implementation tracking only.

## Current Read

- Shared quote state/default parsing exists, but the current form renderer still hardcodes fields inline.
- Shared plan/filter/pricing logic exists and is the strongest part of the rebuild so far.
- The current page is still one responsive page, not a true phone layout and desktop layout split.
- The old live-update path is only half-carried forward: server rewrite support exists, but the page currently has no JS client wired back in.

## Goal

Build the quote page around three explicit layers:

- Shared quote/control definitions.
- Shared plan/filter/pricing logic.
- Separate render targets for phone and desktop.

The phone and desktop pages must share the same state and plan results, while being free to present them differently.

## Layer 1: Shared Quote Definitions

Create one shared quote-definition layer that is the source of truth for:

- field name
- field label
- control kind
- choice source / chooser proc
- default value rule
- bool/date/number handling
- layout hints for phone and desktop groups

This layer replaces repeated field knowledge spread across:

- defaults
- form apply logic
- state summary rendering
- layout-specific control rendering

`QuoteFieldDefs()` should evolve from a label list into a real control-definition table.

## Layer 2: Shared Quote / Plan Logic

Keep all quote parsing, filter decisions, and pricing decisions out of layout files.

The shared logic layer must own:

- quote defaults
- quote state merge/apply rules
- value parsing helpers
- age calculation
- plan filter decisions
- addon/default choice decisions
- total/base/surcharge calculation
- visible vs filtered plan lists

`QuotePlans(...)` should remain the source of truth for plan results. Layout files should consume data from it, not recompute business rules while rendering.

## Layer 3: Separate Render Targets

Introduce explicit render files:

- `z.render.phone.go`
- `z.render.desktop.go`

These files should render the same quote state and the same `QuotePlans(...)` result in different ways.

Phone target:

- compact stacked quote panel
- card-first visible plans
- expandable plan detail inside each card
- filtered plans folded below the main plan list

Desktop target:

- denser quote workbench
- real desktop plan table, not phone cards stretched wider
- strong column alignment for total, plan row, and addon/category cells
- filtered plans folded below the main table

Do not rely on “render both and hide one with CSS” as the main architecture. CSS can still help with breakpoints inside each renderer, but phone and desktop should each have their own deliberate markup path.

## Control Rendering Plan

Render controls from the shared quote-definition layer instead of hardcoding each field in one big form function.

That means:

- one shared helper to render an individual control from its definition
- one phone control layout function
- one desktop control layout function
- one shared apply path that accepts the same field names regardless of layout

The control-definition layer must be good enough that adding or changing a field happens once, not in several files.

## Interaction Plan

Restore live quote updates as part of the rebuild, not as a later polish task.

Required behavior:

- quote field changes post to `/quote-info-change`
- server rewrites the quote-dependent sections
- plan and filtered-plan sections update without full-page reload
- desktop and phone both use the same POST contract

If plan-level addon selectors are interactive, they must also reuse the same state/update contract instead of creating a separate business-logic path.

## Implementation Order

1. Rewrite the rebuild plan and lock the target architecture.
2. Expand quote field definitions into real control definitions.
3. Make defaults/apply/state-summary consume that shared definition layer.
4. Keep `QuotePlans(...)` as the shared plan engine and trim any layout leakage from it.
5. Extract shared render fragments that are neutral across layouts.
6. Build `z.render.phone.go`.
7. Build `z.render.desktop.go`.
8. Reconnect live JS updates to the rebuilt page.
9. Verify that both layouts show the same quote state and the same plan/filter results for the same input.

## Done When

- Quote fields are defined once in one shared definition layer.
- Defaults, form apply logic, and rendering all consume that one definition layer.
- Plan/filter/pricing rules live outside the layout files.
- `z.render.phone.go` exists and is the primary phone renderer.
- `z.render.desktop.go` exists and is the primary desktop renderer.
- Desktop is not just phone markup widened by CSS.
- Visible-plan totals show plan plus addons.
- Base plan row is shown separately so totals are auditable.
- Vision defaults on in the UI.
- Filtered plans remain folded below the main results in both layouts.
- Live update behavior works again.

## Guardrails

- Keep routes explicit.
- Keep `package main`.
- Keep server-rendered HTML.
- Keep rendering separate from quote/plan/filter logic.
- Keep phone and desktop renderers explicit rather than hiding architecture inside CSS.
