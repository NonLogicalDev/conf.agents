---
name: Code_BehaviorDiff
description: Explain changes as old-versus-new behavior in diff-style pseudocode, organized around affected request flows, background jobs, or state transitions. Use for behavior diffs, pseudocode diffs, side-by-side flow comparisons, and PR behavior breakdowns.
---

# Code Behavior Diff

Show what happens differently, why it matters, and where the source proves it. Compare behavior rather than translating changed lines into prose.

## Establish the comparison

- Identify the two versions or modes being compared: before/after revisions, a setting off/on, or current/proposed behavior. Label them clearly.
- Read the relevant source on both sides and enough caller context to trace the affected flow. For a PR, use its intended base and reviewed head; do not accidentally attribute inherited stack changes to this diff.
- Distinguish code behavior from observed execution, test results, and deployment. If one side is unavailable, state the gap and limit the comparison to what is supported. Label illustrative or proposed behavior.
- Keep the work explanatory unless the user separately asks for changes, experiments, or publication.

## Shape the explanation

1. Start with short TL;DR bullets about the behavioral outcome.
2. Explain the motivation briefly when it is not already clear. Distinguish a documented reason from your inference.
3. Show the relevant old/new flows in fenced `diff` pseudocode. Number separate flows when that helps the reader follow them.
4. Finish with material consequences, exceptions, and evidence limits. For a PR breakdown, add relevant Code_Pedantics callouts on correctness, scope, organization, or comments after the explanation. Keep optional improvements separate from defects; a narrow flow question does not require a full review.

Scale this structure to the request. A narrow follow-up may need only a small diff and one explanation. If the user asks for literal side-by-side columns or another format, preserve that choice. Use plain language; simplify the tone further when asked without losing technical meaning.

## Anchor the diff around the flow

- Choose the flow the change affects: a request, negotiation, push, fetch, compaction, refresh, retry, or publication. Do not force a request-flow story onto a background job or a type-only change.
- Follow the operation from its relevant entry point through the changed decisions to the result. For several affected flows, use separate compact comparisons instead of one giant block or a file-by-file inventory.
- Align corresponding old and new steps. Use `-` for old behavior, `+` for new behavior, and unmarked context for unchanged steps. Keep enough shared context to explain ordering.
- Name the actor at each important boundary: client, server, worker, peer, or storage. Distinguish local calls from network requests. Count actual requests or round trips only when the source establishes them; numbered pseudocode steps are not request counts.
- Preserve the branches and ordering that change the outcome: authorization, cache hits, missing state, retries, errors, and publication when relevant. Do not hide a consequential condition inside `...`.
- Say what work is reused, what is recomputed, and what is checked again. Reusing a session or cached result does not by itself eliminate authorization, freshness checks, or client requests.
- Keep identifiers that help map the explanation to source. Use descriptive pseudocode for incidental syntax, wrappers, and plumbing.
- If the change only moves code or changes representation, say which observable behavior stays the same. Do not invent a behavioral difference to fill the format.

## Attach source markers

Place short markers such as `[O1]` and `[N1]` beside the relevant old and new steps. Immediately below the block, map them to clickable GitHub or Sourcegraph file-and-line links, preferably pinned to the corresponding revisions. Annotate each link with what it establishes.

Use sources you actually inspected. Keep old and new line numbers tied to their own versions. One marker may cover several adjacent steps when the source supports all of them. For supplied snippets without a repository URL, cite their supplied file/revision/line labels and state the link limitation; do not invent a permalink.

## Example shape

Illustrative pseudocode; replace the labels and markers with the inspected behavior and sources:

```diff
 worker.compact():
     output = merge(selected_inputs)
     upload(output)
     on publication_conflict:
-        discard(output)
-        restart_merge()                         # [O1]
+        current = reload_manifest()
+        if selected_inputs_still_present(current):
+            retry_publication(output, current)  # [N1]
+        else:
+            discard(output)
+            restart_merge()                     # [N2]
```

Explain the distinction: the new path reuses merged and uploaded bytes, but reloads metadata and checks whether reuse is still valid. The fallback still rebuilds the output.

## Tests

When changing this skill, read [tests/README.md](tests/README.md). Run the relevant scenarios with fresh subagents that have empty context windows.
