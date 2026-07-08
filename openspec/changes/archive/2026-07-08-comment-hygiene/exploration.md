## Exploration: comment-hygiene

### Current State
The prior read-only review findings still reproduce in the current tree:

- `llm-wiki doctor` is implemented and registered through Cobra (`internal/cmd/doctor.go:13-20`) and command registration is covered by `internal/cmd/doctor_test.go:15-19`. README also documents it in the headless command list and behavior notes (`README.md:97-107`).
- Root maintainer guidance is stale: `CLAUDE.md` lists `internal/cmd/` commands without `doctor` (`CLAUDE.md:57`) and omits `llm-wiki doctor` from available commands (`CLAUDE.md:74-86`). `AGENTS.md` has the same architecture-summary drift by describing generation output without `wiki/sources.json` or the doctor package/command context (`AGENTS.md:16-27`).
- `SchemaData.CommandsDir` is still present and passed from tool renderers (`internal/templates/render.go:39-40`, `internal/tools/claude.go:59-78`), and `RenderSchema` still replaces `{{COMMANDS_DIR}}` (`internal/templates/render.go:54-67`). Current `schema.md.template` only references `{{COMMANDS_TREE}}` (`internal/templates/assets/schema.md.template:22-34`); grep found no `{{COMMANDS_DIR}}` usage in template assets.
- `internal/generator/generator.go` has several section comments that mostly repeat the following statement (`Build manifest`, `Create directory structure`, `Write manifest`, `Write .gitignore`, `Install enabled tools` at `internal/generator/generator.go:49`, `69`, `76`, `93`, `98`). The comment at `81` is more useful because it groups three generated files.
- Good explanatory comments should remain: Bubble Tea/huh pointer invariant (`internal/tui/wizard/wizard.go:21-22`, `37`), shared `AGENTS.md` uninstall rules (`internal/tools/opencode.go:48`, `internal/tools/pi.go:48`), read-only doctor boundary (`internal/doctor/doctor.go:1`, `52`, `73`), malformed page skip behavior and atomic writes (`internal/pages/pages.go:31-32`, `96`), and embedded FS panic justification (`internal/templates/embed.go:16-17`).

### Affected Areas
- `CLAUDE.md` — update root maintainer command architecture and available command list so `doctor` is not hidden from Claude-oriented contributors.
- `AGENTS.md` — update root agent guidance to keep OpenCode/Pi-oriented context aligned with the same current command architecture.
- `internal/templates/render.go` — remove or correct stale `CommandsDir` schema plumbing if no template placeholder needs it.
- `internal/tools/claude.go` — simplify `renderSchema` and `SchemaData` construction after removing unused `commandsDir` input.
- `internal/tools/opencode.go` — update `renderSchema` call site if the `commandsDir` parameter is removed.
- `internal/tools/pi.go` — update `renderSchema` call site if the `commandsDir` parameter is removed.
- `internal/generator/generator.go` — prune low-signal comments while preserving exported API comments and helpful grouping comments.
- `internal/generator/generator_test.go` — behavior spec for generated structure and rendered instructions; run after template/render plumbing changes.
- `internal/tools/tools_test.go` — behavior spec for tool install output and rendered instruction content; run after render signature changes.
- `internal/cmd/doctor_test.go` — confirms doctor command registration; no change expected, but useful if docs/architecture claims are questioned.

### Approaches
1. **Targeted hygiene fix** — Align stale maintainer guidance, remove unused `CommandsDir` render plumbing, and prune only obvious generator comments.
   - Pros: Fixes all reproduced findings with a small review footprint; preserves high-value comments; keeps code and docs consistent.
   - Cons: Requires touching both docs and render call sites; tests must catch accidental generated output drift.
   - Effort: Low

2. **Docs-and-comments only** — Update `CLAUDE.md` / `AGENTS.md` and prune generator comments, but leave `CommandsDir` in place.
   - Pros: Smallest diff and lowest behavioral risk.
   - Cons: Leaves confirmed stale code-comment/plumbing noise behind; future readers will still see a placeholder path that is not used.
   - Effort: Low

3. **Broader generated command documentation** — Add doctor references into generated wiki instructions/templates as well as root guidance.
   - Pros: Makes doctor more visible to generated wiki users.
   - Cons: Expands scope beyond comment hygiene; `doctor` is a CLI structural check, while generated `/wiki-*` operations are agent-led knowledge workflows.
   - Effort: Medium

### Recommendation
Use **Targeted hygiene fix**. Scope should include only root maintainer guidance, dead schema render plumbing, and low-signal comments. Keep generated wiki operational instructions focused on `/wiki-ingest`, `/wiki-query`, and `/wiki-lint`; do not add `doctor` to generated command trees unless a later product decision expands generated docs.

Recommended later verification:

- `go test ./internal/generator -run Test`
- `go test ./internal/tools -run Test`
- `go test ./internal/cmd -run Test`
- `go test ./...`

### Risks
- Removing `CommandsDir` is safe based on current grep results, but any future template reintroducing `{{COMMANDS_DIR}}` would need a test or compile-visible data path.
- Editing root guidance can drift from repo language conventions if rewritten too broadly; keep changes surgical.
- Over-pruning comments could erase domain-specific guardrails. Preserve comments that explain non-obvious invariants, boundaries, or failure modes.

### Non-Goals
- Do not implement product behavior changes.
- Do not alter generated `/wiki-ingest`, `/wiki-query`, or `/wiki-lint` semantics.
- Do not add RAG/vector/deterministic ingest/query/auto-repair behavior.
- Do not rewrite Spanish project docs beyond the stale command/guidance lines needed for this hygiene change.

### Ready for Proposal
Yes — propose a narrow cleanup change that updates root guidance, removes unused `CommandsDir` render plumbing, and prunes redundant comments while preserving the explicitly valuable explanatory comments.
