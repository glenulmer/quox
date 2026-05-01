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
3. Any path starting with an alphabetic is relative to the repository root

## C) Code style
 1. Read docs/CODING_STYLE.md before writing or modifying go code. If you can't find it, stop and ask me.

## D) Quality Check

- Keep scope of changes as minimal as possible
- Do not refactor or build new architecture unless I explicitly tell you to do so
- Code must be human-readable and simple
- If I ask for code analysis, stating generalized assumptions means you are wrong - do analysis and report succinctly with specifics
- Before adding any helper, confirm no equivalent exists in the codebase. Propose it with justification and wait for explicit approval before proceeding.
- Do not add const values unless you know will use them repeatedly or they are longer than 25 characters. (Instead add a comment that explains the purpose of the hardcoded constant.)
