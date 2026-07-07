# Code Context

> This file is ephemeral working context, not a source of truth. If it conflicts with code, tests, `README.md`, `AGENTS.md`, or active OpenSpec artifacts, trust those sources and refresh or discard this file.

## Current Product Story

- `llm-wiki` scaffolds and migrates wiki workspaces for Claude Code, OpenCode, and Pi.
- The CLI/TUI manages structure, tool installation, and convenience flows.
- Knowledge operations remain agent-led through generated `/wiki-ingest`, `/wiki-query`, and `/wiki-lint` instructions.
- The project is intentionally not a RAG engine or deterministic knowledge processor.

## High-Signal Files

1. `README.md` — user-facing workflow story and daily usage.
2. `AGENTS.md` — maintainer guidance, repo commands, and architecture notes.
3. `internal/cmd/root.go` — no-arg routing and command entry behavior.
4. `internal/generator/generator.go` — wiki creation and migration flow.
5. `internal/templates/assets/` — editable source for generated instructions and schema content.

## Safe Reading Order

1. Read `README.md` for the product contract.
2. Read `AGENTS.md` for maintainer-level constraints.
3. Read `Makefile` and `.goreleaser.yaml` for release/build behavior.
4. Then read only the package you are changing.

## Staleness Rule

- Do not preserve investigative notes here as long-lived repository knowledge.
- Prefer OpenSpec artifacts for change-specific reasoning.
- Rewrite or delete this file when it stops matching the current codebase.
