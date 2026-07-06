# Wiki Doctor Specification

## Purpose

Provide a bounded read-only command that reports structural wiki health.

## Requirements

### Requirement: Read-Only Structural Checks

`llm-wiki doctor` MUST run bounded read-only checks for manifest validity, required core files, expected tool outputs, wikilink health, and index coverage. The command MUST report findings without changing repository files or wiki content.

#### Scenario: Report structural health

- GIVEN a wiki repository with inspectable files
- WHEN `llm-wiki doctor` runs
- THEN it reports the status of manifest, core files, tool outputs, wikilinks, and index coverage

#### Scenario: Reject unbounded or mutating behavior

- GIVEN a doctor run encounters missing or inconsistent wiki structure
- WHEN findings are emitted
- THEN the command reports only bounded read-only diagnostics and makes no repairs
