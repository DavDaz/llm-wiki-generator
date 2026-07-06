## Exploration: align-llm-wiki-workflow

### Current State
`llm-wiki-generator` already implements the canonical LLM Wiki foundation: `llm-wiki init` creates `raw/`, `wiki/`, `wiki/index.md`, `wiki/log.md`, `wiki/sources.json`, `wiki.toml`, and tool-specific agent instructions. The knowledge operations remain agent-led through `/wiki-ingest`, `/wiki-query`, and `/wiki-lint`; the Go binary scaffolds, migrates, manages, and can safely add deterministic inspection without becoming an ingestion or RAG engine.

The first documentation slice is in progress in `README.md`: it now leads with the canonical workflow, names feature tiers, and documents `New raw note` as a convenience. The remaining alignment work is mostly stale context cleanup, semantic prompt hardening, a bounded `doctor` command, and release guardrails.

### Affected Areas
- `README.md` — product story has been re-centered, but should be reviewed for final wording and consistency after downstream changes.
- `doc/inception-review.md` — source plan for the change; already captures scope, PR slices, and non-goals.
- `context.md` — stale generated context still describes obsolete no-argument routing and launcher behavior that current code/tests have already fixed.
- `CLAUDE.md` — repository guidance still references a top-level `assets/` directory; editable generated content now lives under `internal/templates/assets/`.
- `internal/generator/generator.go` — comment still says wikis are created at `ParentDir/<slug>-wiki/`, while code creates `ParentDir/<slug>/`.
- `internal/templates/assets/commands/wiki-lint.md` — current lint prompt is useful but mostly structural; it needs explicit semantic checks for contradictions, stale claims, missing concepts/entities, weak links, missing citations, and research gaps.
- `internal/templates/assets/schema.md.template` — schema-level lint rules mirror the weaker structural lint expectations and should be kept consistent with the generated command prompt.
- `internal/generator/generator_test.go` and/or `internal/tools/tools_test.go` — best place to add template contract assertions so generated lint expectations do not regress across Claude/OpenCode/Pi installs.
- `internal/cmd/` — natural home for a bounded `doctor` subcommand that loads the manifest and performs filesystem/index/link checks.
- `internal/manifest/manifest.go` — existing validation can be reused by `doctor` for manifest semantics.
- `internal/tools/*.go` — tool registry and `IsInstalled` checks can support doctor validation of enabled backend files.
- `.github/workflows/release.yml` — release currently runs only on tags and delegates to GoReleaser; there is no observed PR/push CI test workflow.
- `Makefile` — `release` tags and pushes immediately after semver validation; it should gain or call a preflight target before production release actions.

### Approaches
1. **Review-sliced alignment** — Implement the change as small work units: docs/story, stale context cleanup, semantic lint prompt with tests, doctor, then CI/release guardrails.
   - Pros: Keeps each PR under the 400-line review budget, respects the existing implementation plan, and lets docs establish the mental model before adding runtime checks.
   - Cons: Requires multiple PRs or commits and a little coordination across phases.
   - Effort: Medium

2. **Single-pass implementation** — Update docs, prompts, doctor, tests, and CI in one broad change.
   - Pros: Fastest path to visible completion.
   - Cons: Likely exceeds the review budget, mixes product messaging with runtime behavior, and makes rollback/review harder.
   - Effort: Medium-High

3. **Runtime-first implementation** — Add `doctor` and CI before finishing stale docs and prompt semantics.
   - Pros: Quickly adds deterministic safety checks.
   - Cons: Leaves the repo teaching stale concepts while new behavior is added; this does not address the primary product-drift risk first.
   - Effort: Medium

### Recommendation
Use **Review-sliced alignment**. Finish the in-progress README slice, then clean stale repository knowledge before modifying generated prompts or adding `doctor`. The recommended implementation chain is:

1. Product-story docs: finish/review `README.md` canonical workflow and feature-tier messaging.
2. Stale context cleanup: replace or remove `context.md`, fix `CLAUDE.md`, and update the stale `<slug>-wiki/` generator comment.
3. Semantic lint strengthening: update `wiki-lint.md` and `schema.md.template`, then add generator/tool tests asserting the required lint concepts are installed for every backend.
4. Deterministic `doctor`: add a Cobra command that validates manifest, core directories/files, enabled tool outputs, broken wikilinks, and index coverage only. It must not ingest, query, summarize, repair, or mutate knowledge content.
5. CI/release guardrails: add test CI for `go test ./...`; add `release-preflight` or equivalent to check clean tree and tests before tag/push; decide whether lint is required or documented as optional based on `golangci-lint` availability.

### Risks
- `context.md` is misleading enough that future agents may reintroduce already-fixed no-argument routing bugs if it remains.
- Prompt changes without tests can silently regress across generated Claude/OpenCode/Pi outputs.
- `doctor` scope can creep into deterministic ingest/repair behavior; keep it read-only and structural/consistency-focused.
- CI/release guardrails may expose environment assumptions, especially around `golangci-lint` not being installed locally.
- The current README change is already uncommitted, so later SDD apply work should avoid overwriting that direct documentation slice.

### Ready for Proposal
Yes — proceed to `sdd-propose`. The proposal should frame this as alignment and guardrails, not as a new knowledge engine. Explicitly keep RAG, vector search, deterministic ingest/query, automatic repairs, Obsidian/Dataview/qmd, and default-path Ollama expansion out of scope.
