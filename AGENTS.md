# AGENTS - Quo2 Guardrails

## A) Non-Negotiables

1. Keep `package main` + page-oriented files (`0.*`, `page.*`).
2. Keep `chi` + `net/http`; explicit route lines.
3. Keep server-rendered HTML via `klec/lib/htmlHelper`.
4. Never use templates, JSON-first app flow, REST/SPA layering, ORM/DI/DTO architecture.
5. Never run `gofmt`/`go fmt`/`goimports`/IDE auto-formatters, and never do style-only "idiomatic Go" rewrites.
6. Never put `panic(...)` in request/page handlers (`page.*.go`).
7. No detailed error plumbing

## B) Code style
 1. Read ~/quo2/docs/CODING_STYLE.md before writing or modifying go code.

## C) Interaction Rule

1. Keep feedback technical and concrete.
2. Cut the bullshit, get to the point, don't polish it up.

## D) Quality Check

- You must achieve code clarity, simplicity, maintainability, and no needless duplication.
- I will ask another AI to criticize your work for all factors.
- Anticipate all criticisms, do not rush to a solution.
- Minimal code changes.
- After each function you add, (a) check for duplication and (b) assess bloat.

## E) Enforcement

- `./scripts/check-guardrails.sh`
- `./scripts/check-all.sh`
- `./scripts/run-dev.sh -port 7777`

Before any completion claim, run `./scripts/check-all.sh`.
`run-dev.sh` refuses to start if guardrails fail.

After any Go code recompilation (not only SQL and/or CSS), start and keep two persistent dev terminals running:
- `./scripts/run-dev.sh -watch -port 3333 -layout desktop`
- `./scripts/run-dev.sh -watch -port 4444 -layout phone`
