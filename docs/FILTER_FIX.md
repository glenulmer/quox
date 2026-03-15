# FILTER_FIX

## Goal

Remove `LoadFilterLookups` and move all static lookup loading into `0.boots.go`, with lookup data stored on `App.lookup` and normalized toward `IdMap_t` usage.

`CurrentDBDate()` is dynamic and out of scope for static lookup loading.

## Guardrails To Preserve

1. Keep `package main` and page-oriented files.
2. Keep `chi` + `net/http` and server-rendered HTML patterns.
3. Do not add `panic(...)` in `page.*.go` handlers.
4. Keep rewrite contracts untouched (`data-post`, `data-record`, rewrite IDs/selectors).
5. Keep DB static lookups in bootstrap only, not request handlers.

## SQL Execution Rule

If SQL changes are needed, propose SQL snippets in this plan and hand them off for manual user execution and testing. Do not execute SQL migrations directly from this workflow.

## Function Replacement Order (Recommended)

1. **Data type migration first**
   - Edit `FilterLookups_t` in `page.2.filters.types.go`.
   - Convert lookup sets to `IdMap_t[...]` where possible:
     - `hospitalLevels`, `dentalLevels`
     - `priorCoverOptions`, `examOptions`, `specialistOptions`
   - Keep explicit defaults (`priorCoverDefault`, `examDefault`, `specialistDefault`).
   - Remove `...Allowed map[int]bool` fields once consumers can validate via `IdMap_t.byId`.

2. **Bootstrap loaders in `0.boots.go`**
   - Move filter lookup orchestration from `LoadFilterLookups()` into `0.boots.go`.
   - Replace slice-returning query helpers with bootstrap loaders that return `IdMap_t`:
     - `QueryPriorCoverOptions` -> `LoadPriorCoverIdMap` (or equivalent)
     - `QueryReferralOptions` -> `LoadReferralIdMap` (for specialist validation)
     - `QueryLevelChooser` -> `LoadLevelNameIdMap(categ int)` (or equivalent)
   - Keep `CategIDByName` (or renamed equivalent) in bootstrap path.
   - Update `LoadStaticData()` to set all lookup fields once at startup.

3. **Filter-page consumer migration**
   - Update `LoadFiltersPageState`, `DefaultFilterState`, `NormalizeFilterState` in `page.2.filters.go` to use `IdMap_t` reads (`sort` + `byId`).
   - Replace validators:
     - `PickLevel(... []LevelName_t ...)` -> `PickLevel(... IdMap_t[LevelName_t] ...)`
     - `PickCodebook(... map[int]bool ...)` -> `PickCodebook(... IdMap_t[CodebookOption_t] ...)`
   - Update render helpers:
     - `LevelSelect(... []LevelName_t)` -> `LevelSelect(... IdMap_t[LevelName_t])`
     - `CodebookSelect(... []CodebookOption_t)` -> `CodebookSelect(... IdMap_t[CodebookOption_t])`
   - Keep form `name=` keys and rewrite IDs exactly unchanged.

4. **Remove obsolete filter bootstrap wrapper**
   - Delete `LoadFilterLookups` and no-longer-needed helper functions in `0.bootstrap.filters.go`.
   - If constants are still needed (`spPriorCovQuery`, `spReferralsQuery`, `spLevelChooser`, specialist/exam codes), keep them in a bootstrap file used by `0.boots.go`.

5. **Validation hardening**
   - Keep startup validation for empties and required codes:
     - prior-cover options non-empty
     - hospital/dental levels non-empty
     - specialist fixed codes compatible with referral lookup
   - Ensure all failures remain bootstrap-time panics, never request-time handler panics.

6. **Checks**
   - Run `./scripts/check-all.sh` after code changes.

7. **Final docs sync (required final step)**
   - Update `docs/STATIC_LOOKUPS.md` to reflect final symbols, cache fields, and loader function names in `0.boots.go`.
   - Update `docs/QUERY_REGISTRY.md` only if symbols/proc wiring changed.
   - This doc-update step is intentionally last so docs match final code exactly.

## Expected End State

1. `LoadFilterLookups` no longer exists.
2. Static lookup loading happens once at startup in `0.boots.go`.
3. Filter lookups are represented with `IdMap_t` where applicable.
4. Filter page behavior and rewrite protocol remain unchanged.
5. Guardrail scripts pass, including static/query registry checks.
