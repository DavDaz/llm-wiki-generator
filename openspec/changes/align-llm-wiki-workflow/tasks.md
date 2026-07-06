# Tasks: Align LLM Wiki Workflow

## Review Workload Forecast

| Field | Value |
|-------|-------|
| Estimated changed lines | 450-650 |
| 400-line budget risk | High |
| Chained PRs recommended | Yes |
| Suggested split | PR 1 docs/stale cleanup → PR 2 lint/templates/tests → PR 3 doctor/CI-release |
| Delivery strategy | auto-forecast |
| Chain strategy | feature-branch-chain (`Rama integradora`) |

Decision needed before apply: No
Chained PRs recommended: Yes
Chain strategy: feature-branch-chain (`Rama integradora`)
400-line budget risk: High

### Suggested Work Units

| Unit | Goal | Likely PR | Notes |
|------|------|-----------|-------|
| 1 | Preserve workflow-story edits and remove stale claims | PR 1 | Base main; keep README diff intact; resolve `context.md` with smallest safe delete-or-rewrite. |
| 2 | Tighten generated lint/schema guidance and backend assertions | PR 2 | Base PR 1; update templates plus generator/tools tests together. |
| 3 | Add bounded read-only doctor command | PR 3 | Base PR 2; include package, Cobra wiring, and command tests. |
| 4 | Add CI and release preflight guardrails | PR 4 | Base PR 3; keep release safety changes isolated. |

## Phase 1: Documentation Alignment

- [x] 1.1 Review the existing uncommitted `README.md` changes and only edit around them to keep the workflow-story slice intact.
- [x] 1.2 Update `context.md` with the smallest safe action that removes stale routing/knowledge claims; delete only if rewrite would add more churn.
- [x] 1.3 Refresh `CLAUDE.md` guidance to match the agent-led workflow and current tool/template paths.

## Phase 2: Semantic Lint Contract

- [x] 2.1 Update `internal/templates/assets/commands/wiki-lint.md` to require contradiction, stale-claim, concept/entity, link, citation, and research-gap checks; keep no-repair wording.
- [x] 2.2 Mirror the same lint contract in `internal/templates/assets/schema.md.template`.
- [x] 2.3 Adjust `internal/generator/generator.go` comments only if stale path guidance still exists.
- [x] 2.4 Extend `internal/generator/generator_test.go` and `internal/tools/tools_test.go` to assert the generated backend outputs preserve the lint contract.

## Phase 3: Read-Only Doctor Command

- [ ] 3.1 Create `internal/doctor/doctor.go` with bounded, read-only checks for manifest, core files, tool outputs, wikilinks, and index coverage.
- [ ] 3.2 Add `internal/doctor/doctor_test.go` using `t.TempDir()` to cover pass/fail findings and no-write behavior.
- [ ] 3.3 Wire `internal/cmd/doctor.go` into Cobra and add `internal/cmd/doctor_test.go` for wiki and non-wiki execution paths.

## Phase 4: Release Safety

- [ ] 4.1 Add `.github/workflows/ci.yml` to run `go test ./...` on push/PR.
- [ ] 4.2 Update `Makefile` with a `preflight` target and make `release` depend on it before tag/push.
- [ ] 4.3 Verify release/help text still matches the new preflight flow.

## Phase 5: Verification / Cleanup

- [ ] 5.1 Run focused Go tests for generator, tools, doctor, and cmd packages.
- [ ] 5.2 Run `go test ./...` and confirm no stale docs or generated outputs were unintentionally changed.
