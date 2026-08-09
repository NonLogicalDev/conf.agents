# Comments And Documentation

## Put The Why Near The Code

Use comments generously when a reader would otherwise have to read a lot of code to understand what is happening. Good comments compress a complicated section into a clear idea, explain why the code works a certain way, and summarize what a block does when that helps the reader follow it.

Treat comments as part of the implementation. Mechanical readability does not
explain history, permissions, risk, the order of phases, side effects, or
conditions for safe changes.

A future engineer should not need hidden task context or a search across the
repository to understand why code is shaped the way it is. Keep shared behavior
and API context in
module or API documentation and local rationale beside the decision it explains.

Internal task plans, task numbers, conversation history, coordinator
instructions, and agent reasoning are not code rationale. Do not mention them
in code comments or documentation. State the lasting behavior, compatibility
rule, protocol, risk, or condition for a safe change directly. Keep
temporary orchestration and private reasoning in the internal plan instead.
When work the user has permitted finds this leak in a public artifact, treat it as
a real review defect. Do not preserve it because the branch is stacked, the
deadline is close, or fixing wording changes a commit SHA.

## Lead With Purpose And Observable Outcome

For a module, API, workflow phase, or substantial code block, start with the
purpose before the mechanics. Explain what capability, user outcome, or
operator outcome the code enables and how a reader can recognize that it is
working. Then explain the important implementation detail, constraint, or
tradeoff.

A local comment can be shorter: name the risk or outcome that makes the nearby
decision necessary. Do not open with file layout, control-flow narration, or a
list of implementation steps when the reader still does not know why the code
exists.

## Separate Reusable Documentation From Local Reasons

Documentation for an item should make sense when the reader has not seen any of
its callers. State the item's primary purpose, its role in the architecture,
why the item needs to exist, and why its responsibility or shape is useful.
For a type, explain why its fields, variants, ownership, or lifetime travel
together. For a function, explain its responsibility and the behavior it preserves,
not the sequence one particular caller happens to use.

Write enough context for a reviewer with little knowledge of the
repository or task conversation to understand each changed item on its own.
Name the relevant domain, lifecycle or workflow phase, input and permission
source, output, mutation or persistence, important relationship, failure or
recovery behavior, and conditions for a safe change when those facts affect
review. Match the explanation to the item. Use several sentences or paragraphs
when one sentence cannot explain it. Do not impose a limit of one line, and do
not rely on a distant module comment, commit message, or
prior conversation to supply context required to review a function safely.

Keep rationale for a particular invocation beside that invocation. Ordering,
fallback choice, retry placement, phase transitions, and the risk prevented at
one call site normally belong in a local comment, not in reusable item docs.
An item contract may still document a true universal precondition or ordering
rule; distinguish that from advice that is correct for only one caller.

Bad:

```rust
/// Carries the fields used while rebuilding an archive segment.
struct SegmentRebuildPlan { /* ... */ }

/// Run this after loading the manifest and before publishing the replacement.
async fn verify_replacement(/* ... */) -> Result<()> { /* ... */ }
```

Good:

```rust
/// Records the exact manifest generation used to rebuild one archive
/// segment safely.
///
/// The source generation prevents applying a plan to newer state, the expected
/// digest identifies the content to preserve, and replacement
/// objects travel with that expectation so planning and publication cannot
/// accidentally refer to different candidates.
struct SegmentRebuildPlan { /* ... */ }

/// Verifies that the replacement objects in a rebuild plan preserve the
/// segment's expected digest before the plan becomes eligible for publication.
///
/// Reads the candidate objects without publishing them and returns an error when their
/// digest differs, allowing callers to abandon or rebuild a stale candidate.
async fn verify_replacement(/* ... */) -> Result<()> { /* ... */ }

// Verify immediately before publication so this rebuild cannot publish an
// invalid replacement after planning has succeeded.
verify_replacement(&plan, store).await?;
store.publish(plan).await
```

The good item docs explain architectural purpose and shape in isolation. The
local comment explains why this caller places verification at that point.

## What To Explain

Use names, documentation, and comments to explain:

- the responsibility owned here and responsibilities deliberately excluded;
- who may change, delete, access, or save data, and the limits on
  persistence, or protocol decisions;
- rules to preserve, state changes, ordering, and side effects;
- why an important guard or branch exists, what information is absent, and
  when the behavior may change safely;
- why operations happen in a particular order; and
- the purpose of each phase in a larger workflow.

Documentation for the operator should name prerequisites, supported invocation
modes, required credentials or tools, and which resources may change. Module
documentation can state what a script or component reads, validates, and
writes.

## Long Orchestration

Prefer specifically named helpers for coherent phases. When a function remains long because phases share state or need to be read together, use comments to introduce each phase, summarize the work it performs, and explain why each transition happens.

Length alone does not require comments. A short, declarative, or already
decomposed function may be clear without them. A long coordinator is not clear
merely because every statement and type is mechanically readable.

## Avoid Syntax Narration

Do not restate assignments, loops, conditions, assertions, obvious literals, or
raw comparisons. Explain why the choice matters, what state may change, what
risk the guard prevents, or what evidence makes the behavior safe.

Give comments enough detail to help the reader and use exact function, type, field, and concept names. Remove stale comments rather than preserving a confident explanation that no longer matches the code.

## Preserve Honest Auditability

Optimize for auditability and behavior preservation, not minimum line count. A
branch, guard, or comment that looks repetitive can be the clearest honest
representation when paths have different types, permissions, lifetimes, side
effects, or failure behavior.

Do not collapse real distinctions merely because expressions look similar. If
deliberate repetition protects a meaningful distinction, keep the explicit shape
and explain the reason to keep it explicit.

Remove a comment only after checking that it is genuinely duplicate, stale, or
misleading. Adjacent comments may look repetitive while documenting different
contracts. Avoid unrelated comment cleanup when the task does not require it.

## Keep Rationale Local

A commit message or review description cannot be the only place that explains
permissions, persistence, protocols, security, state changes, fixtures, or
negative proof. Summarize a decision that affects the whole system in review
prose when useful, but
keep the implementation rationale beside the code that depends on it.

Before handoff, a junior engineer should be able to answer:

- What inputs does this code require?
- What may change during execution?
- What data can it modify?
- What validation happens before modification?
- Why does each important guard exist?
- Why does each unusual type or function exist, and why is it shaped this
  way instead of placing that responsibility elsewhere?
- Which rationale is a reusable item contract, and which belongs beside a
  particular call site?
- What are the major phases, and why are they ordered this way?
- What gives and limits each permission?
