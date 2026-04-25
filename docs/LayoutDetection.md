# Layout Detection Plan

## Goal

Move layout choice from global `App.layout` to per-session device mode (`mobile` or `desktop`) with server bootstrap detection and client confirmation.

## Constraints

- Keep page-oriented handlers and existing HTML helper flow.
- No template engine.
- Cookie (`device`) takes priority over UA.
- Redirect script must be first element emitted inside `<head>`.

## Implementation Steps

1. Session model updates
- Extend `SessionVars_t` with:
  - `device string`
  - `deviceConfirmed bool`
- Default new sessions to `device=desktop`, `deviceConfirmed=false`.
- Update state/session conversion so device fields are preserved when `SetState` is called.

2. Session store/device helpers
- Add mode constants (`mobile`, `desktop`) and validation helper.
- Add store getters/setters for device + confirmed flag by token.
- Expose request helpers:
  - `SessionDeviceMode(*http.Request) string`
  - `SetSessionDeviceMode(*http.Request, string)`

3. Middleware detection order
- In `SessionMiddleware`:
  - Ensure session token.
  - Check `device` cookie first; if valid, write to session and mark confirmed.
  - If not confirmed, run UA detection (`github.com/mssola/useragent`) for newly created/uninitialized sessions.

4. Runtime layout selection
- Add helper to map session device -> render layout (`phone`/`desktop`).
- Replace `App.layout` checks in quote/editq GET+POST render paths with request-scoped layout helper.

5. Head script positioning
- Extend `lib/htmlHelper/Head_t` with a method for a pre-head inline script rendered immediately after `<head>` and before meta tags.

6. Client confirmation script
- Build inline script with server mode embedded as literal string.
- Script logic:
  - detect client mode (`innerWidth < 768 || pointer:coarse`)
  - compare with server mode
  - on mismatch set `device` cookie and reload via `location.replace(location.href)`
  - on match do nothing

7. Inject script on full-page quote/editq renders
- Add pre-head script injection in:
  - `Page1Quote`
  - `Page2EditQ`

8. Dependency and verification
- Add `github.com/mssola/useragent` to module deps.
- Run `go test ./...` and confirm no regressions.

## Risks

- Existing `SetState` currently rewrites all session vars; without preservation, device metadata would be lost.
- If cookies are blocked, mismatch reload cannot persist preference; script should avoid forced infinite retry.
- Width-only classification is insufficient for modern iPhone landscape; coarse-pointer fallback is required.

## Expected Result

- First request gets server best-guess layout.
- First rendered page quickly confirms/corrects layout client-side.
- Subsequent requests stay stable via cookie + session-confirmed mode.
