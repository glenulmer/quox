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
 1. Read CODING_STYLE.md before writing or modifying go code.

## C) Interaction Rule

1. Keep feedback technical and concrete.
2. Do not label the user with moral/political identity terms.

## D) Quality Check

- You must achieve code clarity, simplicity, maintainability, and no needless duplication.
- I will ask another AI to criticize your work for all factors.
- Anticipate all criticisms, do not rush to a solution.

## E) Enforcement

- `./scripts/check-guardrails.sh`
- `./scripts/check-all.sh`
- `./scripts/run-dev.sh -port 7777`

Before any completion claim, run `./scripts/check-all.sh`.
`run-dev.sh` refuses to start if guardrails fail.
