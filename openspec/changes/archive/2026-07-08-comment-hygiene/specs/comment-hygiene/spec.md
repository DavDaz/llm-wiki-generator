# Comment Hygiene Specification

## Purpose

Define the expected hygiene boundaries for maintainer guidance, schema render placeholders, and source comments without changing generated wiki command behavior.

## Requirements

### Requirement: Root maintainer guidance stays current

The system MUST keep root maintainer guidance accurate for the current CLI/package layout, including `doctor`, `wiki/sources.json`, and the generated file/backend context that contributors rely on.

#### Scenario: Current command and file guidance is documented

- GIVEN the repository root guidance is updated
- WHEN a maintainer reads `CLAUDE.md` or `AGENTS.md`
- THEN the guidance includes the `doctor` command where command/package context is summarized
- AND the generated wiki file/package context reflects current outputs such as `wiki/sources.json` and backend instruction files

#### Scenario: Generated wiki commands are not broadened by root guidance updates

- GIVEN the repository documents maintainer-only commands
- WHEN the hygiene change is applied
- THEN generated wiki command scope remains limited to `/wiki-ingest`, `/wiki-query`, and `/wiki-lint`

### Requirement: Schema placeholder cleanup is dependency-safe

The system MUST remove schema render placeholders or plumbing only after confirming no live template asset depends on them.

#### Scenario: Unused placeholder is removed safely

- GIVEN a schema render field or replacement has no matching usage in template assets
- WHEN the hygiene change removes that field or replacement
- THEN schema rendering still succeeds with the remaining live placeholders only

#### Scenario: Live placeholder blocks removal

- GIVEN any template asset still references a placeholder such as `{{COMMANDS_DIR}}`
- WHEN cleanup is evaluated
- THEN the placeholder plumbing MUST be preserved until the template dependency is removed or replaced

### Requirement: Low-signal comments are pruned selectively

The system SHOULD remove comments that merely restate adjacent code, and MUST preserve comments that explain invariants, boundaries, shared-file rules, or failure modes.

#### Scenario: Redundant generator narration is removed

- GIVEN a comment in `internal/generator/generator.go` only repeats the next statement
- WHEN comment hygiene is applied
- THEN that redundant narration may be removed without changing behavior

#### Scenario: Guardrail comments remain

- GIVEN a comment documents a non-obvious invariant or safety boundary
- WHEN comment hygiene is applied
- THEN that comment remains available to future maintainers

### Requirement: Generated wiki command semantics remain unchanged

The system MUST NOT change the behavior contract of generated `/wiki-ingest`, `/wiki-query`, or `/wiki-lint` as part of this hygiene change.

#### Scenario: Render and installer cleanup preserves command behavior

- GIVEN schema rendering or tool installer call sites are simplified
- WHEN generated instruction files are rendered after the change
- THEN the operational semantics of `/wiki-ingest`, `/wiki-query`, and `/wiki-lint` stay unchanged

#### Scenario: Doctor remains outside generated command semantics

- GIVEN `doctor` is a root CLI structural check
- WHEN maintainers review generated wiki instructions after the hygiene change
- THEN `doctor` is not introduced as a generated wiki command requirement
