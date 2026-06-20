---
name: go-reference
description: Use when scaffolding a Go project with a Dockerfile, sqlc, Makefile, or GitHub Actions CI. Refers to github.com/lwlee2608/go-reference as the canonical template.
user-invocable: true
argument-hint: [Dockerfile | Makefile | sqlc | CI]
---

# go-reference Template

When the user asks to Dockerize a Go project, add a Makefile, set up sqlc, or add CI, read [`lwlee2608/go-reference`](https://github.com/lwlee2608/go-reference) and refer to it as the canonical template. Do not write these files from scratch or from memory.

## Arguments

The argument (e.g. "Add Dockerfile", "Set up sqlc") scopes the work — copy only the matching files, not the whole template. If no argument is given, ask which parts to scaffold.

## Rules

1. **Clone the live repo, never work from memory.** Source of truth is `https://github.com/lwlee2608/go-reference` (branch `main`). Clone it once with `git clone --depth=1 https://github.com/lwlee2608/go-reference /tmp/go-reference-template`, then `ls` and `Read` files from there. This lets you see the full project layout (`cmd/`, `internal/api`, `internal/db`) — that structure is part of the template too, not just the individual files. The repo evolves; reproducing contents from memory causes drift.

2. **Adapt, don't blind-copy.** The template hardcodes `go-reference` in `Dockerfile`, `Makefile` (`APP := go-reference`), and the `cmd/go-reference/` path — replace with the target project's binary name. Otherwise keep the template's structure and conventions.
