# Generated Wiki Linting Specification

## Purpose

Require generated wiki lint guidance to check semantic wiki health instead of only surface formatting.

## Requirements

### Requirement: Semantic Wiki Lint Contract

Generated lint prompts and schema guidance MUST require contradiction, stale-claim, concept/entity consistency, wikilink, citation, and research-gap checks across supported backends. They MUST NOT prescribe automatic repair or deterministic ingest/query behavior.

#### Scenario: Generate semantic lint guidance

- GIVEN a wiki is generated for a supported backend
- WHEN lint guidance is installed
- THEN it includes all required semantic check categories in the generated output

#### Scenario: Keep lint bounded to evaluation

- GIVEN semantic issues are detected during lint review
- WHEN the generated guidance describes next steps
- THEN it asks for inspection or follow-up instead of automatic repair behavior
