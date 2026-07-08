# Design: Comment Hygiene

## Technical Approach

Apply a surgical hygiene refactor: update root maintainer guidance, remove confirmed-unused schema render plumbing for `{{COMMANDS_DIR}}`, and trim only comments that restate adjacent code. Generated `/wiki-ingest`, `/wiki-query`, and `/wiki-lint` command files remain byte-for-byte semantically unchanged.

## Architecture Decisions

| Decision | Choice | Alternatives considered | Rationale |
|---|---|---|---|
| Keep change as docs/refactor only | Do not edit `internal/templates/assets/commands/*` or command behavior sections in `schema.md.template` | Broader generated instruction cleanup | The proposal explicitly excludes generated command semantic changes. |
| Remove `CommandsDir` only | Delete `SchemaData.CommandsDir`, its replacer entry, and the `renderSchema` parameter | Keep field as future-proof placeholder | Grep confirms no `{{COMMANDS_DIR}}` usage in template assets, so keeping it is misleading. |
| Preserve invariant comments | Remove narrative comments in `generator.go`; keep package/exported docs and comments explaining shared-file/legacy/migration boundaries | Delete all comments aggressively | Useful comments protect behavior; hygiene should reduce noise, not erase intent. |

## Data Flow

No runtime data-flow change.

```text
manifest.Manifest ──→ tools.renderSchema ──→ templates.RenderSchema
        │                       │                    │
        └──────── tool tree + instructions file ─────┘
```

After the refactor, `commandsDir` no longer flows into `SchemaData`; `commandsTree` still renders backend-specific layout.

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `CLAUDE.md` | Modify | Add `doctor`, `wiki/sources.json`, and current command/package context to maintainer guidance with minimal Spanish edits. |
| `AGENTS.md` | Modify | Align command list and architecture notes with `doctor`, `wiki/sources.json`, and current test guidance. |
| `internal/templates/render.go` | Modify | Remove `SchemaData.CommandsDir` and the `{{COMMANDS_DIR}}` replacement. Keep `CommandsTree` and `InstructionsFile`. |
| `internal/tools/claude.go` | Modify | Change `renderSchema(m, commandsDir, commandsTree, instructionsFile)` to `renderSchema(m, commandsTree, instructionsFile)`; remove `CommandsDir` assignment. Keep `claudeSkillsDir` for install path and legacy cleanup. |
| `internal/tools/opencode.go` | Modify | Update `renderSchema` call signature only; keep `.opencode/commands` constants and install path. |
| `internal/tools/pi.go` | Modify | Update `renderSchema` call signature only; keep `.pi/prompts` constants and install path. |
| `internal/generator/generator.go` | Modify | Remove redundant inline comments such as "Build manifest", "Create directory structure", "Write manifest", and similar code-restating markers. Preserve exported docs and the re-render comment in `Migrate`. |

## Interfaces / Contracts

Only an internal helper signature changes:

```go
func renderSchema(m *manifest.Manifest, commandsTree, instructionsFile string) (string, error)
```

`templates.SchemaData` contract keeps required render inputs: wiki identity, language/date, entities, page types, conventions, `CommandsTree`, and `InstructionsFile`.

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit | Template render no longer requires unused `CommandsDir` | Run existing generator/tools tests; add/update assertions only if compile or expected text changes require it. |
| Integration | Tool install output paths and lint contract remain stable | `go test ./internal/generator -run Test` and `go test ./internal/tools -run Test`. |
| Command wiring | `doctor` registration and behavior stay documented and passing | `go test ./internal/cmd -run Test`. |
| Full suite | No regression across CLI/TUI/generator packages | `go test ./...`. |

## Risk Controls

- Do not edit `internal/templates/assets/commands/wiki-*.md`.
- Do not edit the INGEST/QUERY/LINT behavior text in `schema.md.template` unless tests force a marker-only update.
- Before removing comments, classify them: keep comments for invariants, shared-file deletion rules, legacy migrations, exported API docs, and failure modes.
- Inspect any test snapshot/assertion changes to ensure they reflect plumbing removal or root guidance only.
- Keep the review budget under 400 changed lines by avoiding broad rewrites.

## Migration / Rollout

No data migration required. Existing generated wikis are unaffected until users run `llm-wiki migrate`; even then, command semantics should remain unchanged.

## Rollback Plan

Revert the implementation commit. Because the change is limited to root guidance, comments, and unused internal render plumbing, rollback restores prior files without data migration or generated-wiki repair.

## Open Questions

None.
