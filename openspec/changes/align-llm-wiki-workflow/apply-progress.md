# Apply Progress: align-llm-wiki-workflow

## Mode

- Implementation mode: Standard
- Delivery strategy: auto-forecast
- Chain strategy: feature-branch-chain (`Rama integradora`)
- Current slice: PR 2 — lint/templates/tests
- Size exception: No

## Completed Tasks

- [x] 1.1 Review the existing uncommitted `README.md` changes and only edit around them to keep the workflow-story slice intact.
- [x] 1.2 Update `context.md` with the smallest safe action that removes stale routing/knowledge claims; delete only if rewrite would add more churn.
- [x] 1.3 Refresh `CLAUDE.md` guidance to match the agent-led workflow and current tool/template paths.
- [x] 2.1 Update `internal/templates/assets/commands/wiki-lint.md` to require contradiction, stale-claim, concept/entity, link, citation, and research-gap checks; keep no-repair wording.
- [x] 2.2 Mirror the same lint contract in `internal/templates/assets/schema.md.template`.
- [x] 2.3 Adjust `internal/generator/generator.go` comments only if stale path guidance still exists.
- [x] 2.4 Extend `internal/generator/generator_test.go` and `internal/tools/tools_test.go` to assert the generated backend outputs preserve the lint contract.

## Files Changed

| File | Action | Notes |
|---|---|---|
| `context.md` | Modified | Rewrote stale routing investigation into short, accurate, explicitly-ephemeral repo context. |
| `CLAUDE.md` | Modified | Corrected agent-led workflow framing, current template paths, and release guidance without inventing guardrails not yet implemented. |
| `internal/generator/generator.go` | Modified | Fixed stale `<slug>-wiki/` comment to match current output path. |
| `internal/templates/assets/commands/wiki-lint.md` | Modified | Added semantic lint categories for contradictions, stale claims, concept/entity consistency, link coverage, citations, and research gaps while keeping the guidance report-only. |
| `internal/templates/assets/schema.md.template` | Modified | Mirrored the semantic lint contract in generated schema instructions and preserved the no-auto-repair boundary. |
| `internal/generator/generator_test.go` | Modified | Added assertions that generated CLAUDE/AGENTS/wiki-lint outputs keep the semantic lint contract. |
| `internal/tools/tools_test.go` | Modified | Added backend installer assertions for Claude, OpenCode, and Pi lint/schema outputs. |
| `openspec/changes/align-llm-wiki-workflow/tasks.md` | Modified | Marked completed tasks and recorded the resolved chain strategy for apply. |

## Verification

- Previous slice ran: `git diff --check`
- This slice ran: `go test ./internal/generator -run Test`
- This slice ran: `go test ./internal/tools -run Test`

## Remaining Tasks

- [ ] 3.1 Create `internal/doctor/doctor.go` with bounded, read-only checks for manifest, core files, tool outputs, wikilinks, and index coverage.
- [ ] 3.2 Add `internal/doctor/doctor_test.go` using `t.TempDir()` to cover pass/fail findings and no-write behavior.
- [ ] 3.3 Wire `internal/cmd/doctor.go` into Cobra and add `internal/cmd/doctor_test.go` for wiki and non-wiki execution paths.
- [ ] 4.1 Add `.github/workflows/ci.yml` to run `go test ./...` on push/PR.
- [ ] 4.2 Update `Makefile` with a `preflight` target and make `release` depend on it before tag/push.
- [ ] 4.3 Verify release/help text still matches the new preflight flow.
- [ ] 5.1 Run focused Go tests for generator, tools, doctor, and cmd packages.
- [ ] 5.2 Run `go test ./...` and confirm no stale docs or generated outputs were unintentionally changed.

## Notes

- `README.md` was reviewed and intentionally left untouched to preserve the user’s in-progress workflow-story edits.
- `doc/inception-review.md` and `.atl/skill-registry.md` were not modified.
- Semantic lint guidance remains agent-led evaluation only; this slice did not add automatic repair, deterministic ingest/query, vector search, or RAG behavior.
