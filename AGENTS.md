# AGENTS - Guardrails

## A) Non-Negotiables

1. Keep `package main` + page-oriented files (`0.*`, `page.*`).
2. Keep `chi` + `net/http`; explicit route lines.
3. Keep server-rendered HTML via `klec/lib/htmlHelper`.
4. Never use templates, JSON-first app flow, REST/SPA layering, ORM/DI/DTO architecture.
5. Never run `gofmt`/`go fmt`/`goimports`/IDE auto-formatters, and never do style-only "idiomatic Go" rewrites.
6. No detailed error plumbing

## B) Interaction Rule

1. Keep feedback technical and concrete.
2. State conclusions and recent actions first. Details afterward.
3. If I say "shutdown everything", that means kill all persistent terminals & resources that can be released
4. Any path starting with an alphabetic is relative to the repository root

## C) Code style
 1. Read docs/CODING_STYLE.md before writing or modifying go code.

## D) Quality Check

- Keep scope of changes as minimal as possible
- Code must be human-readable and simple
- If I ask for code analysis, stating generalized assumptions means you are wrong and unhelpful. Do analysis and report succinctly and with specifics.

## E) Enforcement

- `./scripts/check-guardrails.sh`
- `./scripts/check-all.sh`
- `./scripts/run-dev.sh -port 7777`

Before any completion claim, run `./scripts/check-all.sh`.
`run-dev.sh` refuses to start if guardrails fail.

After any Go code recompilation (not only SQL and/or CSS), start and keep two persistent & privileged dev terminals running:
- `./scripts/run-dev.sh -watch -port 3333 -layout desktop`
- `./scripts/run-dev.sh -watch -port 4444 -layout phone`
