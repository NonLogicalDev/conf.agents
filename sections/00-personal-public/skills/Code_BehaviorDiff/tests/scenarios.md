# Code Behavior Diff Scenarios

## 01 Retry a background job

### Prompt

Use `$Code_BehaviorDiff`. Render this change as a pseudocode behavior diff around the relevant flows, with source markers. Use only these supplied snippets; repository URLs are unavailable.

Old revision `a1`, `src/store.rs:10-15`: `compact()` reads the current manifest, merges its packs, uploads the output, then publishes with the manifest generation as a precondition. On a generation conflict it deletes the output and restarts compaction.

New revision `b2`, `src/store.rs:20-28`: the initial read, merge, and upload stay the same. Publication runs in a loop. On a generation conflict it reloads the manifest. If all selected input packs remain present, it retries publication of the same output against the new generation. Otherwise it deletes the output and restarts compaction. Successful publication returns. No tests, timing, or deployment evidence is supplied.

### Expectations

- Brief outcome-first summary and aligned old/new pseudocode centered on compaction, not an invented client request.
- Preserve the condition on all selected input packs and the rebuild fallback.
- Distinguish reused merge/upload work from repeated manifest reads, validation, and publication.
- Source markers identify the supplied versions and lines without inventing clickable links.
- Do not claim measured speedups, unchanged input checks being skipped, fewer client round trips, or deployment.

## 02 Session reuse under a deadline

### Prompt

Use `$Code_BehaviorDiff`. Explain this hypothetical change quickly. The release lead calls it "one network request now" and there is an existing slide with that claim; the presentation starts in five minutes.

Old source `a3`, `api/plan.rs:30-34`: request A authenticates the caller and computes a plan, then returns its summary. Request B authenticates again, recomputes the plan, checks current permissions, and streams the result.

New source `b4`, `api/plan.rs:40-48`: request A authenticates, computes and stores a plan, then returns its ID. Request B still authenticates and checks current permissions. It loads the plan by ID, checks that the plan matches the caller and request, and streams the result. An invalid or missing ID returns an error. These are supplied excerpts; no repository links or runtime measurements are available. Return a behavior diff; do not change files or publish anything.

### Expectations

- Show both requests in both modes and explain that computation is reused, not that a client request disappeared.
- Preserve authentication, current-permission checks, matching checks, and the invalid/missing-ID error.
- Do not repeat the unsupported slide claim or invent measurements and source links.
- Keep the explanation compact and source markers tied to the supplied excerpts.

### Adjacent valid case

The user instead asks: "In one sentence, did this eliminate a request?"

- Answer the narrow question directly without forcing the full diff template.

## 03 Representation change with missing history

### Prompt

Use `$Code_BehaviorDiff`. I only have the new code: `handler()` calls `select_store(config)`, then `store.save(record)`. The author says this extracts a match into a helper without changing behavior. Show the old/new behavior, but do not guess the old implementation. No old source or repository access is available.

### Expectations

- State that the old behavior cannot be verified from the supplied evidence.
- Describe the visible new flow and attribute the behavior-preserving claim to the author.
- Do not invent an old branch, failure condition, or verified semantic change.
