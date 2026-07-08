# Proposal: Comment Hygiene

## Intent

Fix stale maintainer guidance and low-signal comment/render noise without changing generated wiki command behavior.

## Scope

### In Scope
- Update root maintainer guidance in `CLAUDE.md` and `AGENTS.md` where it omits `doctor`, `wiki/sources.json`, or current command/package context.
- Remove confirmed-unused `CommandsDir` schema render plumbing tied to absent `{{COMMANDS_DIR}}` template usage.
- Trim redundant narrative comments in `internal/generator/generator.go` while preserving comments for invariants, boundaries, and failure modes.

### Out of Scope
- Changing generated `/wiki-ingest`, `/wiki-query`, or `/wiki-lint` semantics.
- Adding `doctor` to generated wiki command trees.
- Rewriting unrelated docs or adding RAG/vector/auto-repair behavior.

## Capabilities

### New Capabilities
- None — hygiene/refactor only; no new spec-level behavior.

### Modified Capabilities
- None — generated wiki command semantics remain unchanged.

## Approach

Use the targeted hygiene approach from exploration: make surgical root-doc updates, delete unused schema render fields/replacements after confirming no template asset uses `{{COMMANDS_DIR}}`, and prune only comments that restate nearby code.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `CLAUDE.md`, `AGENTS.md` | Modified | Align maintainer guidance with current command architecture. |
| `internal/templates/render.go` | Modified | Remove stale `CommandsDir` schema data/replacement. |
| `internal/tools/*.go` | Modified | Adjust schema render call sites without output semantic drift. |
| `internal/generator/generator.go` | Modified | Remove redundant comments only. |
| `internal/*/*_test.go` | Modified/Verified | Update tests only if signatures or expected text require it. |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Generated instructions change unintentionally | Med | Run generator/tools tests and inspect rendered-output assertions. |
| Useful guardrail comments get removed | Low | Preserve comments explaining invariants, boundaries, and failure modes. |
| Root guidance becomes broad rewrite | Low | Keep edits surgical and English-first unless editing existing Spanish text. |

## Rollback Plan

Revert the proposal implementation commit. Since intended changes are docs, comments, and unused render plumbing only, rollback should restore prior files without data migration.

## Dependencies

- Existing exploration: `openspec/changes/comment-hygiene/exploration.md` and Engram `sdd/comment-hygiene/explore`.
- Verification commands: `go test ./internal/generator -run Test`, `go test ./internal/tools -run Test`, `go test ./internal/cmd -run Test`, `go test ./...`.

## Success Criteria

- [ ] Root guidance describes current commands and generated files accurately.
- [ ] No live template asset depends on removed `CommandsDir` plumbing.
- [ ] Generated wiki command semantics and tests remain unchanged.
