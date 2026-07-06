# Release Guardrails Specification

## Purpose

Prevent releases from skipping essential verification.

## Requirements

### Requirement: Test-Based Preflight Guardrails

CI and release preflight flows MUST verify `go test ./...` succeeds before tag or push-based release actions continue. Guardrails SHOULD fail fast with a clear preflight result when verification has not passed.

#### Scenario: Allow verified release flow

- GIVEN the full Go test suite succeeds
- WHEN CI or release preflight runs
- THEN the workflow allows subsequent release actions to continue

#### Scenario: Block unverified release flow

- GIVEN tests fail or were not run successfully
- WHEN CI or release preflight runs
- THEN release actions stop before tag or push completion
