# Tasks: Comment Hygiene

## Review Workload Forecast

| Field | Value |
|-------|-------|
| Estimated changed lines | 120-220 |
| 400-line budget risk | Low |
| Chained PRs recommended | No |
| Suggested split | Single PR |
| Delivery strategy | auto-forecast |
| Chain strategy | pending |

Decision needed before apply: No
Chained PRs recommended: No
Chain strategy: pending
400-line budget risk: Low

### Suggested Work Units

| Unit | Goal | Likely PR | Notes |
|------|------|-----------|-------|
| 1 | Update maintainer guidance and remove dead render plumbing | PR 1 | Covers root docs + schema helper cleanup together; verify with focused tests. |
| 2 | Prune redundant generator comments and verify regressions | PR 1 | Keep only `internal/generator/generator.go` comment cleanup; no behavior changes. |

## Phase 1: Root Guidance and Render Plumbing

- [x] 1.1 Update `CLAUDE.md` to mention `doctor`, `wiki/sources.json`, and current backend/package context without broadening generated wiki command scope.
- [x] 1.2 Update `AGENTS.md` to align command guidance, architecture notes, and test guidance with `doctor` and `wiki/sources.json`.
- [x] 1.3 Remove `SchemaData.CommandsDir` and the `{{COMMANDS_DIR}}` replacer from `internal/templates/render.go`.
- [x] 1.4 Update `internal/tools/claude.go` to drop the `commandsDir` argument from `renderSchema(...)` and remove the unused field assignment.
- [x] 1.5 Update `internal/tools/opencode.go` and `internal/tools/pi.go` to call the new `renderSchema(...)` signature.

## Phase 2: Comment Hygiene in Generator

- [x] 2.1 Prune redundant inline comments in `internal/generator/generator.go` that only restate adjacent code.
- [x] 2.2 Keep comments that explain invariants, shared-file rules, migration behavior, or other non-obvious boundaries in `internal/generator/generator.go`.

## Phase 3: Verification

- [x] 3.1 Run `go test ./internal/generator -run Test` to confirm generator behavior stays stable after comment and render cleanup.
- [x] 3.2 Run `go test ./internal/tools -run Test` to verify tool render wiring and install-path behavior still compile and pass.
- [x] 3.3 Run `go test ./internal/cmd -run Test` to confirm command registration/docs remain consistent.
- [x] 3.4 Run `go test ./...` for full regression coverage.

## Phase 4: Cleanup

- [x] 4.1 Review the diff to confirm `CommandsDir` plumbing is fully removed and no generated command semantics changed.
- [x] 4.2 Ensure the final wording in `CLAUDE.md` and `AGENTS.md` stays concise and review-friendly.
