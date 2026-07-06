# Proposal: Align LLM Wiki Workflow

## Intent

Align `llm-wiki-generator` with the canonical LLM Wiki workflow: curated immutable `raw/` sources, agent-led ingest/query/lint, markdown synthesis, and CLI/TUI support. Reduce story drift, stale guidance, weak prompts, and release risk without making the CLI a knowledge engine.

## Scope

### In Scope
- Preserve/review the uncommitted `README.md` workflow-story slice.
- Clean stale knowledge in `context.md`, `CLAUDE.md`, and generator comments.
- Strengthen generated `/wiki-lint` and schema semantic checks.
- Add bounded read-only `llm-wiki doctor` structural/wiki health inspection.
- Add CI and release preflight guardrails.

### Out of Scope
- RAG engine, vector search, deterministic ingest/query engine, or automatic content repair.
- Obsidian, qmd, Dataview, or format expansion.
- Default-path Ollama expansion or backend behavior redesign.

## Capabilities

### New Capabilities
- `workflow-documentation-alignment`: Docs/guidance explain the agent-led workflow and remove stale claims.
- `generated-wiki-linting`: Generated lint prompts require contradiction, stale-claim, concept/entity, link, citation, and research-gap checks.
- `wiki-doctor`: CLI exposes read-only manifest, core-file, tool-output, wikilink, and index checks.
- `release-guardrails`: CI and release preflight verify tests before release actions.

### Modified Capabilities
- None; no existing OpenSpec specs are present.

## Approach

Use review-sliced alignment: docs/story, stale cleanup, semantic lint prompts plus template tests, bounded doctor, then CI/release guardrails. Keep units within 400 changed lines and protect current `README.md` edits.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `README.md`, `doc/inception-review.md` | Modified | Align story; preserve edits. |
| `context.md`, `CLAUDE.md`, `internal/generator/generator.go` | Modified | Remove stale claims/comments. |
| `internal/templates/assets/commands/wiki-lint.md`, `internal/templates/assets/schema.md.template` | Modified | Add semantic lint expectations. |
| `internal/generator/generator_test.go`, `internal/tools/tools_test.go` | Modified | Assert lint contract. |
| `internal/cmd/`, `internal/manifest/`, `internal/tools/` | New/Modified | Add read-only `doctor`. |
| `.github/workflows/`, `Makefile` | New/Modified | Add guardrails. |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Scope creeps into ingest/repair engine | Med | Enforce read-only doctor and non-goals. |
| Prompt regressions across backends | Med | Add template/generator assertions. |
| Existing README edits are overwritten | Med | Treat current diff as source input. |
| Change exceeds review budget | Med | Split by work unit. |

## Rollback Plan

Revert by work unit: docs cleanup, prompt/template tests, doctor, and CI/release guardrails are independent. `doctor` is read-only, so rollback removes only command surface and tests.

## Dependencies

- Exploration artifacts and `doc/inception-review.md`.
- Uncommitted `README.md` slice.

## Success Criteria

- [ ] Docs explain CLI scaffolding vs agent-led knowledge operations without stale contradictions.
- [ ] Generated `/wiki-lint` requires semantic wiki health checks across supported backends.
- [ ] `llm-wiki doctor` reports bounded read-only structural/index/link health.
- [ ] CI runs `go test ./...`; release flow has a preflight before tag/push.
