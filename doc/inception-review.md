# LLM Wiki Generator Inception Review

Date: 2026-07-06  
Scope: Review `llm-wiki-generator` against the canonical LLM Wiki pattern in `evaluador/llm-wiki.md`.

## Executive verdict

`llm-wiki-generator` mostly complies with the LLM Wiki idea. Its best framing is not “a wiki engine” or “a RAG system”; it is a Go CLI/TUI that scaffolds and maintains the operating environment for AI-agent-maintained markdown wikis.

The direction is valid. The main risk is product drift: onboarding and documentation currently give too much attention to tool management, TUI features, and optional agents before making the core mental model unmistakable.

Recommended direction: keep the CLI/TUI, but re-center the product around the canonical workflow:

1. Curated immutable sources in `raw/`.
2. AI-maintained compiled knowledge in `wiki/`.
3. Generated schema/instructions that discipline the agent.
4. Repeated ingest/query/lint workflows that keep `index.md` and `log.md` current.

Do not build a full RAG/search/ingestion engine yet. The highest-value next step is a clearer core workflow plus deterministic health checks.

## Canonical baseline

The LLM Wiki pattern requires these baseline capabilities:

| Area | Baseline requirement |
| --- | --- |
| Raw sources | Curated, immutable source documents. The AI reads them but does not mutate them. |
| Wiki | LLM-generated markdown pages that accumulate synthesis over time. |
| Schema | Agent instructions that define structure, conventions, and workflows. |
| Ingest | New sources update summaries, entity/concept pages, `index.md`, and `log.md`. |
| Query | The agent reads the wiki first, synthesizes answers with citations, and can file valuable answers back into the wiki. |
| Lint | Periodic health checks for contradictions, stale claims, orphans, weak links, missing concepts, and research gaps. |
| Navigation | `index.md` is the content catalog; `log.md` is the chronological trace. |

Optional enhancements include qmd/search, Obsidian plugins, Marp, Dataview, image handling, Ollama, and elaborate CLI/TUI workflows.

## Compliance matrix

| Requirement | Status | Evidence / notes |
| --- | --- | --- |
| Three-layer model | Compliant | The generator creates `raw/`, `wiki/`, `wiki.toml`, and agent instructions. |
| Immutable raw sources | Compliant | Generated schema and ingest instructions protect `raw/` as source material. |
| Compounding markdown wiki | Compliant | Templates instruct agents to create/update summaries, concepts, links, and synthesis. |
| Agent schema/instructions | Compliant | `CLAUDE.md` / `AGENTS.md` are generated from `internal/templates/assets/schema.md.template`. |
| Ingest workflow | Mostly compliant | `/wiki-ingest` exists and updates pages/index/log; re-ingest state via `sources.json` is useful. |
| Query workflow | Compliant | `/wiki-query` reads wiki pages and can save valuable answers back into the wiki. |
| Lint workflow | Partially compliant | `/wiki-lint` exists, but should be strengthened around semantic contradictions, stale claims, and missing concepts. |
| `index.md` + `log.md` | Compliant | Both are generated and treated as navigation/traceability anchors. |
| Right-sized optional tooling | Mixed | Multiple tool backends and TUI are useful, but docs should label them as optional support, not the core idea. |

## What is working well

- The Go binary is a good fit for repeatable scaffolding and migration.
- `wiki.toml` provides a clear manifest/source of truth for installed tool support.
- Tool adapters are isolated behind backend-specific code for Claude Code, OpenCode, and Pi.
- Generated templates preserve the central LLM Wiki workflows: ingest, query, and lint.
- Tests cover important generator, manifest, tools, dashboard, and raw-note behavior.
- `go test ./...` and `make test` passed during review; `make test` uses the race detector.
- The manual raw-note manager supports the real workflow of capturing new raw material without leaving the dashboard.

## Gaps and risks

### 1. Documentation does not lead with the core workflow strongly enough

The README should make it clear that `llm-wiki` scaffolds the workspace and installs agent instructions. The AI agent performs ingest/query/lint; the Go binary does not currently execute those knowledge operations itself.

Impact: new users may expect a deterministic ingestion engine or RAG tool instead of an agent-operated wiki workflow.

### 2. Some context/docs are stale

Observed stale or misleading areas:

- `context.md` describes older no-argument behavior that no longer matches current command routing.
- `CLAUDE.md` references a top-level `assets/` directory, while editable templates now live under `internal/templates/assets/`.
- README manage documentation omits the current raw-note creation flow.
- A generator comment still implies `<slug>-wiki/` while current behavior creates `<slug>/`.

Impact: future agents and maintainers can make wrong changes because the repo teaches outdated behavior.

### 3. Lint is too structural compared with the canonical baseline

The canonical `llm-wiki.md` expects wiki health checks to find contradictions, stale claims superseded by newer sources, missing concepts, weak cross-references, and research gaps.

The current lint prompt is useful, but it should more explicitly check semantic health, not only structure.

### 4. Release tooling lacks guardrails

Current state is workable for a small tool, but release risk is higher than necessary:

- `make lint` is declared but failed in this environment because `golangci-lint` is not installed.
- There is a release workflow, but no observed CI workflow that runs tests/lint before release.
- `make release` tags and pushes; it should be treated as a production action, not a dry run.

### 5. Optional features are too prominent

TUI status management, Ollama, multiple backend migration, raw-note forms, and markdown guide rendering are all useful. None of them is the LLM Wiki core.

Impact: the product can feel like a tool manager instead of a disciplined wiki generator.

## What to keep

- Go CLI/TUI as the installer, migrator, and maintenance cockpit.
- `wiki.toml` manifest.
- Generated schema plus agent command templates.
- `raw/`, `wiki/`, `index.md`, `log.md`, and `sources.json`.
- Multi-agent support if the project genuinely targets Claude Code, OpenCode, and Pi.
- Raw-note creation, but document it as a convenience workflow.

## What to demote or simplify

| Feature | Recommendation |
| --- | --- |
| Ollama modelfile | Move out of the primary happy path; keep as advanced optional setup. |
| Page status management | Keep as optional review/publishing workflow, not core LLM Wiki behavior. |
| Multi-backend tooling | Keep if actively used; otherwise consider a simpler generic `AGENTS.md` backend plus optional adapters. |
| Glamour guide viewer | Acceptable, but optional; consider plain markdown fallback if dependency weight matters. |
| `context.md` | Remove, regenerate, or clearly mark as ephemeral if it can become stale. |

## Recommended right-sized plan

This plan should be implemented as a sequence of small, reviewable changes. The first objective is not to add more machinery; it is to make the Karpathy-style LLM Wiki workflow obvious, reliable, and hard to misuse.

### Phase 1 — Re-center the product story

- Rewrite the README quick path around the canonical flow:
  1. `llm-wiki init`
  2. Add files to `raw/`
  3. Ask the selected AI agent to run `/wiki-ingest`
  4. Query with `/wiki-query`
  5. Health-check with `/wiki-lint`
- Add a clear sentence: “The CLI scaffolds and manages the wiki; the AI agent maintains the knowledge.”
- Add a feature-tier table: core, convenience, advanced/optional.

### Phase 2 — Fix stale repo knowledge

- Update or remove stale `context.md`.
- Fix `CLAUDE.md` references from `assets/` to `internal/templates/assets/`.
- Update README manage section to include raw-note creation.
- Fix comments/docs that still imply `<slug>-wiki/` if current behavior is `<slug>/`.

### Phase 3 — Strengthen wiki health behavior

- Update `wiki-lint.md` to explicitly check:
  - contradictions between pages,
  - stale claims superseded by newer sources,
  - orphan pages and weakly linked pages,
  - missing concept/entity pages,
  - missing source citations,
  - research gaps worth filling.
- Add golden/contract tests that prevent command-template drift.

### Phase 4 — Add deterministic safety checks

Add a small `llm-wiki doctor` command before considering heavier automation.

Minimum useful checks:

- `wiki.toml` exists and validates.
- `raw/` and `wiki/` exist.
- `wiki/index.md`, `wiki/log.md`, and `wiki/sources.json` exist.
- Generated tool files match enabled tools in `wiki.toml`.
- Basic broken-link scan for wiki links.
- Basic index coverage check against actual wiki pages.

This gives users confidence without turning the project into a full ingestion engine.

### Phase 5 — Add CI/release gates

- Add CI for `go test ./...`.
- Decide whether `golangci-lint` is required; if yes, install it in CI and document setup.
- Add a release preflight target that checks clean tree + tests + optional lint before tagging.

## Implementation plan

### Guiding principle

Keep the Go binary focused on scaffolding, migration, inspection, and safety checks. Keep knowledge synthesis in the AI-agent workflow. This distinction is the product boundary and should guide every implementation decision.

### Delivery sequence

| Step | Change | Primary files | Acceptance criteria |
| --- | --- | --- | --- |
| 1 | Rewrite the README happy path around the canonical LLM Wiki workflow. | `README.md` | README explains that the CLI scaffolds/manages the wiki while the AI agent performs ingest/query/lint; quick start follows init → raw sources → ingest → query → lint. |
| 2 | Add feature tiers to reduce product drift. | `README.md` | Features are grouped as core, convenience, and advanced/optional; TUI, Ollama, backend migration, and page status are not presented as the core value proposition. |
| 3 | Fix stale repository knowledge. | `context.md`, `CLAUDE.md`, source comments, README manage docs | No docs claim obsolete no-arg behavior, obsolete `assets/` paths, or old `<slug>-wiki/` output. Raw-note creation is documented as a convenience flow. |
| 4 | Strengthen the generated wiki lint prompt. | `internal/templates/assets/commands/wiki-lint.md` or equivalent template path | `/wiki-lint` explicitly checks contradictions, stale claims, weak links, missing concepts/entities, missing citations, and research gaps. |
| 5 | Protect prompt-template behavior with tests. | `internal/generator/*_test.go`, `internal/tools/*_test.go`, or focused template tests | Tests fail if generated command templates lose the semantic lint expectations or required workflow anchors. |
| 6 | Add deterministic `doctor` checks. | `internal/cmd/`, `internal/manifest/`, `internal/generator/` as needed | `llm-wiki doctor` validates manifest, required directories/files, enabled tool outputs, basic wiki links, and index coverage without performing ingest/query itself. |
| 7 | Add CI and release preflight guardrails. | `.github/workflows/`, `Makefile`, release docs | CI runs `go test ./...`; release preflight checks clean tree and tests before any tag/push action. Lint is either installed/documented or clearly optional. |

### Suggested PR slicing

| PR | Scope | Why this boundary works |
| --- | --- | --- |
| PR 1 | Product-story docs only: README happy path, feature tiers, manage/raw-note docs. | Low risk, immediately aligns the project with the canonical model. |
| PR 2 | Stale context cleanup: `context.md`, `CLAUDE.md`, comments, obsolete path references. | Keeps maintainer/agent guidance correct before deeper code changes. |
| PR 3 | Semantic lint prompt plus template/generator tests. | Strengthens the core wiki health loop without adding a deterministic ingestion engine. |
| PR 4 | `llm-wiki doctor` command. | Adds deterministic trust checks as a focused CLI capability. |
| PR 5 | CI and release preflight. | Hardens delivery after product behavior and docs are clearer. |

### Non-goals for this implementation round

- Do not build a full deterministic ingest/query engine.
- Do not add vector search, qmd, Dataview, Obsidian plugins, or RAG infrastructure yet.
- Do not remove the TUI or multi-backend support unless later usage data shows they are dead weight.
- Do not make optional agent/tool integrations part of the default happy path.

### Verification plan

- For README and context-only changes, review rendered markdown and check terminology against the canonical workflow.
- For template changes, run `go test ./internal/generator -run Test` and any focused tool/template tests.
- For command wiring or `doctor`, run `go test ./...`.
- Before release-related changes land, verify `make test` locally when available because it runs the race detector.

### Recommended first implementation target

Start with PR 1 and PR 2 before adding `doctor`. The project currently needs sharper product framing more than new runtime behavior. Once the docs and agent-facing context teach the right mental model, semantic lint and deterministic checks will compound in the right direction.

## Alternatives considered

| Alternative | Fit | Tradeoff |
| --- | --- | --- |
| Keep current CLI/TUI and add `doctor` | Best | Preserves current value while adding trust and clarity. |
| Strip to shell/templates only | Too small now | Easier to maintain, but loses migration, manifest, and TUI value already built. |
| Build full deterministic ingest/query engine | Too heavy | Would fight the core idea: the LLM maintains the wiki, not a rigid parser/RAG pipeline. |
| Add vector search/qmd now | Premature | Useful later at scale, but optional in the canonical model. |

## Decision

Proceed with `llm-wiki-generator` as the right foundation, but tighten it:

1. Make the core LLM Wiki workflow obvious.
2. Demote optional tools from the primary path.
3. Fix stale docs/context.
4. Strengthen lint semantics.
5. Add deterministic `doctor` and CI/release guardrails.

This keeps the project aligned with the LLM Wiki pattern without overbuilding a separate knowledge engine.
