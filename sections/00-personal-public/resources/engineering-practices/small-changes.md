# Small, reviewable changes

Read this reference when preparing a code change, deciding whether to split a pull request, or reducing reviewer effort.

Source: [Google Engineering Practices: Small CLs](https://google.github.io/eng-practices/review/developer/small-cls.html).

A change list, or CL, is a proposed code change. The same principles apply to a GitHub pull request or another reviewable unit of work.

## Prefer one coherent change

Keep a change focused on one understandable purpose. Include its related tests, the context the reviewer needs, and enough implementation for the result to make sense on its own.

This is a judgment about review effort, not a strict limit on lines or files. A change should remain complete, understandable, and safe to merge.

## Why this helps

Focused changes are usually easier to:

- Read and review without setting aside a large block of time.
- Reason about and test thoroughly.
- Merge, debug, and roll back.
- Correct before too much work depends on the wrong approach.
- Explain without asking reviewers to reconstruct hidden context.

## Split work when it genuinely helps

Useful boundaries include:

- Separate unrelated refactors from a bug fix or feature.
- Separate changes with different owners or reviewers.
- Introduce shared groundwork before the feature that uses it.
- Split independent features into complete, reviewable steps.
- Stack dependent changes when the repository supports that workflow.

Keep each step buildable and preserve the related tests. Avoid splitting a change so aggressively that reviewers need several incomplete pull requests to understand a single behavior.

## Apply judgment

A generated mechanical change, a straightforward file deletion, or a necessarily coupled migration can be large while remaining simple to review. When a broad change is genuinely needed, explain its boundaries, highlight the meaningful parts, and provide useful verification.

The principle is to respect the reviewer's time, not to optimize for an arbitrary line count.
