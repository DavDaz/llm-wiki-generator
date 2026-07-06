# Workflow Documentation Alignment Specification

## Purpose

Keep product-story documentation aligned with the agent-led LLM Wiki workflow and remove stale repository guidance.

## Requirements

### Requirement: Workflow Story Alignment

The system documentation MUST describe `llm-wiki` as scaffolding and migration support for wiki structure, while agent-led knowledge ingest, query, and lint remain external workflow responsibilities. Repository guidance MUST remove stale or contradictory claims and MUST preserve in-progress README workflow edits during alignment updates.

#### Scenario: Align workflow docs

- GIVEN repository docs and guidance files describe the product story
- WHEN the workflow alignment change is applied
- THEN the docs state the agent-led workflow without claiming CLI knowledge-engine behavior

#### Scenario: Clean stale claims safely

- GIVEN stale guidance exists beside an uncommitted README workflow slice
- WHEN the docs are updated
- THEN stale claims are removed without overwriting the in-progress README edits
