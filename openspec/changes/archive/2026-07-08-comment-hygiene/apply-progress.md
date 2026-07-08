# Apply Progress: Comment Hygiene

## Implementation Progress

**Change**: comment-hygiene  
**Mode**: Standard

### Completed Tasks
- [x] 1.1 Update `CLAUDE.md` to mention `doctor`, `wiki/sources.json`, and current backend/package context without broadening generated wiki command scope.
- [x] 1.2 Update `AGENTS.md` to align command guidance, architecture notes, and test guidance with `doctor` and `wiki/sources.json`.
- [x] 1.3 Remove `SchemaData.CommandsDir` and the `{{COMMANDS_DIR}}` replacer from `internal/templates/render.go`.
- [x] 1.4 Update `internal/tools/claude.go` to drop the `commandsDir` argument from `renderSchema(...)` and remove the unused field assignment.
- [x] 1.5 Update `internal/tools/opencode.go` and `internal/tools/pi.go` to call the new `renderSchema(...)` signature.
- [x] 2.1 Prune redundant inline comments in `internal/generator/generator.go` that only restate adjacent code.
- [x] 2.2 Keep comments that explain invariants, shared-file rules, migration behavior, or other non-obvious boundaries in `internal/generator/generator.go`.
- [x] 3.1 Run `go test ./internal/generator -run Test`.
- [x] 3.2 Run `go test ./internal/tools -run Test`.
- [x] 3.3 Run `go test ./internal/cmd -run Test`.
- [x] 3.4 Run `go test ./...`.
- [x] 4.1 Review the diff to confirm `CommandsDir` plumbing is fully removed and no generated command semantics changed.
- [x] 4.2 Ensure the final wording in `CLAUDE.md` and `AGENTS.md` stays concise and review-friendly.

### Files Changed
| File | Action | What Was Done |
|------|--------|---------------|
| `CLAUDE.md` | Modified | Added root CLI/package context for `doctor`, `wiki/sources.json`, and backend output layout while keeping generated wiki command scope unchanged. |
| `AGENTS.md` | Modified | Aligned maintainer command guidance, architecture notes, and verification guidance with `doctor`, `wiki/sources.json`, and focused package tests. |
| `internal/templates/render.go` | Modified | Removed `SchemaData.CommandsDir` and the `{{COMMANDS_DIR}}` replacer. |
| `internal/tools/claude.go` | Modified | Simplified `renderSchema` signature and removed the unused `commandsDir` data flow. |
| `internal/tools/opencode.go` | Modified | Updated the `renderSchema` call to the new signature. |
| `internal/tools/pi.go` | Modified | Updated the `renderSchema` call to the new signature. |
| `internal/generator/generator.go` | Modified | Removed redundant inline narration comments and preserved migration/invariant comments. |
| `openspec/changes/comment-hygiene/tasks.md` | Modified | Marked the assigned implementation, verification, and cleanup tasks complete. |

### Deviations from Design
None — implementation matches design.

### Issues Found
None.

### Remaining Tasks
- [ ] None.

### Workload / PR Boundary
- Mode: single PR
- Current work unit: comment-hygiene
- Boundary: root maintainer guidance, dead render plumbing removal, generator comment cleanup, and verification only
- Estimated review budget impact: 35 changed lines in tracked source/docs, well below the 400-line budget

### Verification
- [x] `go test ./internal/generator -run Test`
- [x] `go test ./internal/tools -run Test`
- [x] `go test ./internal/cmd -run Test`
- [x] `go test ./...`

### Status
13/13 tasks complete. Ready for verify.
