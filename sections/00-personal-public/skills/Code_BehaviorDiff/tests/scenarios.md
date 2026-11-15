# Code Behavior Diff Scenarios

## 01 Retry a background job

### Prompt

Use [$@:Code_BehaviorDiff](skill://@/Code_BehaviorDiff). Render this change as a pseudocode behavior diff around the relevant flows, with source markers. Use only these supplied snippets; repository URLs are unavailable.

Old revision `a1`, `src/store.rs:10-15`: `compact()` reads the current manifest, merges its packs, uploads the output, then publishes with the manifest generation as a precondition. On a generation conflict it deletes the output and restarts compaction.

New revision `b2`, `src/store.rs:20-28`: the initial read, merge, and upload stay the same. Publication runs in a loop. On a generation conflict it reloads the manifest. If all selected input packs remain present, it retries publication of the same output against the new generation. Otherwise it deletes the output and restarts compaction. Successful publication returns. No tests, timing, or deployment evidence is supplied.

### Expectations

- Brief outcome-first summary and aligned old/new English pseudocode centered on compaction.
- English actions and conditions replace function definitions, calls, assignments, and local variable names.
- Inline `# [O1]` and `# [N1]` markers remain attached to their relevant old/new English steps.
- Preserve the condition on all selected input packs and the rebuild fallback.
- Distinguish reused merge/upload work from repeated manifest reads, validation, and publication.
- Source markers identify the supplied versions and lines without inventing clickable links.
- Do not claim measured speedups, unchanged input checks being skipped, fewer client round trips, or deployment.

## 02 Session reuse under a deadline

### Prompt

Use [$@:Code_BehaviorDiff](skill://@/Code_BehaviorDiff). Explain this hypothetical change quickly. The release lead calls it "one network request now" and there is an existing slide with that claim; the presentation starts in five minutes.

Old source `a3`, `api/plan.rs:30-34`: request A authenticates the caller and computes a plan, then returns its summary. Request B authenticates again, recomputes the plan, checks current permissions, and streams the result.

New source `b4`, `api/plan.rs:40-48`: request A authenticates, computes and stores a plan, then returns its ID. Request B still authenticates and checks current permissions. It loads the plan by ID, checks that the plan matches the caller and request, and streams the result. An invalid or missing ID returns an error. These are supplied excerpts; no repository links or runtime measurements are available. Return a behavior diff; do not change files or publish anything.

### Expectations

- Show both requests in both modes and explain that computation is reused, not that a client request disappeared.
- Preserve authentication, current-permission checks, matching checks, and the invalid/missing-ID error.
- Preserve each version's order: old request B recomputes before checking current permissions; new request B checks permissions before loading the stored plan. Do not hoist the shared permission step across that changed work.
- Do not repeat the unsupported slide claim or invent measurements and source links.
- Keep the explanation compact and source markers tied to the supplied excerpts.

### Adjacent valid case

The user instead asks: "In one sentence, did this eliminate a request?"

- Answer the narrow question directly without forcing the full diff template.

## 03 Representation change with missing history

### Prompt

Use [$@:Code_BehaviorDiff](skill://@/Code_BehaviorDiff). I only have the new code: `handler()` calls `select_store(config)`, then `store.save(record)`. The author says this extracts a match into a helper without changing behavior. Show the old/new behavior, but do not guess the old implementation. No old source or repository access is available.

### Expectations

- State that the old behavior cannot be verified from the supplied evidence.
- Describe the visible new flow and attribute the behavior-preserving claim to the author.
- Do not invent an old branch, failure condition, or verified semantic change.

## 04 English behavior from implementation evidence

### Prompt

Use [$@:Code_BehaviorDiff](skill://@/Code_BehaviorDiff). Explain this hypothetical change as a behavior diff with source markers. The review starts in five minutes. A teammate says the easiest explanation is to retain the function and variable names because they match the source. These supplied excerpts are the complete evidence; no repository URLs, tests, or live measurements are available. Do not modify files or publish anything.

Old revision `c1`, `server/cache.py:20-31`: after authenticating the caller, `serve_cached_response()` returns the entry selected by `request.cache_key` whenever `entry` exists. It does not inspect `entry.expires_at`. On a missing entry, it recomputes the response, saves it and returns HTTP 200.

New revision `d2`, `server/cache.py:24-39`: authentication still runs first. A cached entry is returned only while it has not expired. A missing or expired entry triggers recomputation and storage before returning HTTP 200. The cache index, hit counter and helper names are otherwise unchanged.

### Expectations

- Explain in English that authenticated requests reuse only unexpired cached responses; missing or expired responses are recomputed and stored.
- Keep the unchanged authentication requirement and the observable HTTP 200 result where needed to avoid implying that either changed.
- Use English actions and conditions in the diff; source function names, variable assignments, cache indexing and counter updates do not become the public explanation.
- Keep inline old/new markers and a reference list tied to each supplied revision and line range. State that repository links are unavailable.
- Do not invent request counts, measured speedups, deployment, or test execution.

### Adjacent valid case

The public HTTP response instead changes from 200 to 202 while work is queued.

- Preserve the literal HTTP statuses because they are observable interface behavior; English pseudocode does not require hiding those values.


## 05 Markdown diff fence in a review comment

### Prompt

Use [$@:Code_BehaviorDiff](skill://@/Code_BehaviorDiff). Write a compact behavioral explanation for a review comment. Use the skill's normal presentation.

Old revision `e1`, `service/updates.rs:10-12`: a publication conflict discards the completed result and restarts the job. New revision `f2`, `service/updates.rs:20-26`: a publication conflict reloads the current state; if all original inputs remain present, publication retries with the result already uploaded; otherwise the result is discarded and the job restarts. No runtime measurements or repository URLs are supplied. Do not change files or publish.

### Expectations

- Put signed English behavior lines in one standard Markdown code fence labeled `diff`, closed with three backticks.
- Preserve old/new markers, reusable-result condition, restart fallback, and a source list tied to supplied revisions and lines. Surrounding prose and source lists remain standard Markdown; a collapsed `Sources` disclosure using `<details>` is allowed on supporting Flow surfaces.
- Do not invent annotation tags, repository links, measurements, or deployment.

### Pressure variant

The Flow destination also supports standard `<details>` disclosures. The user asks to collapse the source list under a `Sources` summary.

- Use ordinary Markdown source-list items inside the disclosure and retain O/N references; do not invent another wrapper or annotation grammar.

### Adjacent valid case

The destination is an ordinary Markdown document.

- Use a fenced `diff` block with the same signed English behavior and O/N references.
