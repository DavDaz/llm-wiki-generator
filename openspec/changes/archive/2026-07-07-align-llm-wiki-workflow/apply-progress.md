# Apply Progress: align-llm-wiki-workflow

## Mode

- Implementation mode: Standard
- Delivery strategy: auto-forecast
- Chain strategy: feature-branch-chain (`Rama integradora`)
- Current slice: final verification
- Size exception: No

## Completed Tasks

- [x] 1.1 Review the existing uncommitted `README.md` changes and only edit around them to keep the workflow-story slice intact.
- [x] 1.2 Update `context.md` with the smallest safe action that removes stale routing/knowledge claims; delete only if rewrite would add more churn.
- [x] 1.3 Refresh `CLAUDE.md` guidance to match the agent-led workflow and current tool/template paths.
- [x] 2.1 Update `internal/templates/assets/commands/wiki-lint.md` to require contradiction, stale-claim, concept/entity, link, citation, and research-gap checks; keep no-repair wording.
- [x] 2.2 Mirror the same lint contract in `internal/templates/assets/schema.md.template`.
- [x] 2.3 Adjust `internal/generator/generator.go` comments only if stale path guidance still exists.
- [x] 2.4 Extend `internal/generator/generator_test.go` and `internal/tools/tools_test.go` to assert the generated backend outputs preserve the lint contract.
- [x] 3.1 Create `internal/doctor/doctor.go` with bounded, read-only checks for manifest, core files, tool outputs, wikilinks, and index coverage.
- [x] 3.2 Add `internal/doctor/doctor_test.go` using `t.TempDir()` to cover pass/fail findings and no-write behavior.
- [x] 3.3 Wire `internal/cmd/doctor.go` into Cobra and add `internal/cmd/doctor_test.go` for wiki and non-wiki execution paths.
- [x] 4.1 Add `.github/workflows/ci.yml` to run `go test ./...` on push/PR.
- [x] 4.2 Update `Makefile` with a `preflight` target and make `release` depend on it before tag/push.
- [x] 4.3 Verify release/help text still matches the new preflight flow.
- [x] 5.1 Run focused Go tests for generator, tools, doctor, and cmd packages.
- [x] 5.2 Run `go test ./...` and confirm no stale docs or generated outputs were unintentionally changed.

## Files Changed

| File | Action | Notes |
|---|---|---|
| `README.md` | Modified | Re-centered the canonical LLM Wiki workflow and clarified the CLI-versus-agent boundary. |
| `context.md` | Modified | Replaced stale routing/debug notes with a short ephemeral context file and clear staleness guard. |
| `CLAUDE.md` | Modified | Corrected workflow framing, current template asset paths, and release-status guidance relevant to this slice. |
| `internal/generator/generator.go` | Modified | Fixed stale `<slug>-wiki/` path wording in the `Init` comment. |
| `internal/templates/assets/commands/wiki-lint.md` | Modified | Added mandatory semantic lint categories for contradictions, stale claims, concept/entity consistency, links, citations, and research gaps while keeping lint as report-only guidance. |
| `internal/templates/assets/schema.md.template` | Modified | Mirrored the semantic lint contract and explicitly kept lint agent-led and non-repairing in generated instructions. |
| `internal/generator/generator_test.go` | Modified | Added assertions that generated instructions and Claude/OpenCode/Pi lint outputs preserve the semantic lint contract. |
| `internal/tools/tools_test.go` | Modified | Added backend installer assertions that schema and lint command outputs keep the semantic lint contract. |
| `internal/doctor/doctor.go` | Created | Added the bounded read-only doctor report with manifest, core file, tool output, wikilink, and index coverage checks. |
| `internal/doctor/doctor_test.go` | Created | Added temp-dir coverage for healthy/failing reports and verified the doctor run stays read-only. |
| `internal/cmd/doctor.go` | Created | Added the read-only Cobra wiring that loads the current wiki, delegates to `internal/doctor`, prints concise findings, and exits non-zero when structural errors are found. |
| `internal/cmd/doctor_test.go` | Created | Added command registration plus wiki/non-wiki execution-path coverage with temp-dir fixtures. |
| `.github/workflows/ci.yml` | Created | Added a push/pull-request CI workflow that runs `go test ./...` on GitHub Actions. |
| `Makefile` | Modified | Added `preflight` and `validate-release-version`, then made `release` depend on clean-tree and full-test verification before tag/push. |
| `CLAUDE.md` | Modified | Updated the release guidance in Spanish to document `make preflight` and the new CI workflow. |
| `AGENTS.md` | Modified | Updated the English repo guidance so release instructions mention `make preflight` and `ci.yml`. |
| `openspec/changes/archive/2026-07-07-align-llm-wiki-workflow/tasks.md` | Archived copy | Records the completed Phase 4 CI/release guardrail tasks and Phase 5 verification tasks in the archived change bundle. |
| `openspec/changes/archive/2026-07-07-align-llm-wiki-workflow/apply-progress.md` | Archived copy | Preserves the merged PR 1/2/3A/3B/4 progress and final verification state after archive. |
| `openspec/changes/archive/2026-07-07-align-llm-wiki-workflow/verify-report.md` | Archived copy | Preserves the final SDD verification evidence, compliance matrix, and PASS verdict after archive. |

## Verification

- Ran: `go clean -testcache` ✅
- Ran: `go test ./internal/generator -run Test` ✅ (`0.600s`)
- Ran: `go test ./internal/tools -run Test` ✅ (`0.424s`)
- Ran: `go test ./internal/doctor -run Test` ✅ (`0.767s`)
- Ran: `go test ./internal/cmd` ✅ (`0.883s`)
- Ran: `make -n release VERSION=v0.3.0` ✅
- Ran: `go test ./...` ✅
- Ran: `git diff --check` ✅
- Not run: `make release` (explicitly out of scope and would tag/push).
- Not run: `make preflight` end-to-end because its clean-tree guard is expected to fail while this branch is still intentionally uncommitted; the dry run confirmed the target order and `go test ./...` passed independently.

## Remaining Tasks

- None.

## Notes

- The revert removed the in-progress README/OpenSpec slice, so this re-apply restored the intended docs changes directly from Engram/source verification.
- To stay within the review budget, the filesystem OpenSpec rehydration for this slice was kept minimal to `tasks.md` and `apply-progress.md`; proposal/spec/design remain available in Engram.
- PR 3A stayed inside the doctor-package boundary: read-only checks plus package tests only; Cobra wiring, CI, release guardrails, and broader docs cleanup remain untouched.
- PR 3B stayed inside the command-wiring boundary: Cobra registration, concise read-only output, and clear non-wiki failure handling only; `internal/doctor` behavior and release guardrails remain untouched.
- PR 4 stayed inside the release-safety boundary: CI test workflow, clean-tree + test preflight, and release/help text alignment only; no tagging/pushing or broader docs cleanup was performed.
- Final verification confirmed tasks 1.1-5.2 complete, focused tests and full suite passing, no whitespace issues, and no RAG/vector/deterministic ingest/query/automatic repair scope introduced by the change.
- `.atl/skill-registry.md` was intentionally not modified.
- The lint guidance remains evaluation-only: no auto-repair, deterministic ingest/query, vector search, or RAG behavior was introduced.
