# Design: Align LLM Wiki Workflow

## Technical Approach

Implement the alignment as reviewable slices that keep `llm-wiki` focused on scaffolding, migration, inspection, and release safety. Knowledge ingest/query/lint remains agent-led through generated instructions. The design explicitly excludes RAG, vector search, deterministic ingest/query engines, and automatic repair.

## Architecture Decisions

| Decision | Choice | Alternatives considered | Rationale |
|---|---|---|---|
| Preserve README slice | Treat current `README.md` diff as input and edit around it | Rewrite README from scratch | `README.md` already contains the new workflow story; preserving it avoids losing uncommitted work. |
| Semantic lint location | Update `internal/templates/assets/commands/wiki-lint.md` and `internal/templates/assets/schema.md.template` | Add runtime Go lint engine | Specs require generated guidance, not deterministic knowledge analysis. |
| Backend coverage | Assert generated Claude skill, OpenCode command, and Pi prompt outputs contain the lint contract | Test only template source | Tool installers copy/render templates differently; backend output tests catch drift. |
| Doctor command | Add read-only Cobra command backed by a small internal diagnostic package | Fold into `status` or add repair mode | `status` summarizes manifest/tool state; `doctor` reports bounded health without mutation. |
| Release guardrails | Add CI workflow for `go test ./...` and Makefile preflight before tagging | Rely on GoReleaser hooks | GoReleaser starts after tag push; release must fail before tag/push. |

## Data Flow

```text
Docs slice: current README diff ──→ targeted edits ──→ review-safe markdown

Templates: assets/wiki-lint + schema ──→ tool installers ──→ generated backend files

Doctor: cwd ──→ wiki.toml ──→ bounded file scans ──→ stdout findings + exit code

Release: make release ──→ preflight(test + clean tree) ──→ tag ──→ push ──→ GoReleaser
```

## File Changes

| File | Action | Description |
|---|---|---|
| `README.md` | Modify | Preserve current workflow-story edits; refine stale wording only. |
| `context.md` | Modify/Delete | Remove stale root-routing and obsolete guidance. |
| `CLAUDE.md` | Modify | Fix template path, command list, and release preflight guidance. |
| `internal/generator/generator.go` | Modify | Fix stale `<slug>-wiki/` comment if still present. |
| `internal/templates/assets/commands/wiki-lint.md` | Modify | Add contradiction, stale-claim, entity/concept, wikilink, citation, and research-gap checks; keep no-repair rule. |
| `internal/templates/assets/schema.md.template` | Modify | Mirror semantic lint contract in generated schema. |
| `internal/generator/generator_test.go` | Modify | Assert generated wiki output includes semantic lint expectations. |
| `internal/tools/tools_test.go` | Modify | Assert Claude/OpenCode/Pi installed outputs preserve lint contract. |
| `internal/doctor/doctor.go` | Create | Pure read-only checks and result model. |
| `internal/doctor/doctor_test.go` | Create | Table-driven filesystem checks using `t.TempDir()`. |
| `internal/cmd/doctor.go` | Create | Cobra command wiring and output formatting. |
| `internal/cmd/doctor_test.go` | Create | Command behavior and non-wiki error coverage. |
| `.github/workflows/ci.yml` | Create | Run `go test ./...` on PRs/pushes. |
| `Makefile` | Modify | Add `preflight` and make `release` depend on it before tag/push. |

## Interfaces / Contracts

```go
type Finding struct { Check, Severity, Path, Message string }
type Report struct { WikiRoot string; Findings []Finding }
func Run(root string) (Report, error) // read-only; no writes, no repairs
```

Doctor checks are bounded to: manifest load/validate, `raw/`, `wiki/`, `wiki/index.md`, `wiki/log.md`, `wiki/sources.json`, enabled tool outputs, `[[wikilink]]` targets in `wiki/*.md`, and basic index coverage for wiki pages.

## Testing Strategy

| Layer | What to Test | Approach |
|---|---|---|
| Unit | Doctor checks and exit-worthy findings | Table-driven tests with `t.TempDir()`. |
| Integration | Cobra `doctor` command output | Stub temp wiki, run command, assert stdout and errors. |
| Template | Semantic lint contract across generated backends | Existing generator/tools tests assert installed files. |
| CI | Guardrail execution | `ci.yml` runs `go test ./...`; `make release` runs preflight first. |

## Migration / Rollout

No data migration required. Deliver as chained/reviewable work units: docs preservation, stale cleanup, semantic lint + tests, doctor command, CI/release guardrails. Each unit should remain under the 400-line review budget or become its own PR slice.

## Open Questions

- [ ] Should `context.md` be deleted or rewritten as explicitly ephemeral generated context?
