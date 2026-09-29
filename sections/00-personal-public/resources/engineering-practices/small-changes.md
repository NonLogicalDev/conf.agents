# Small, reviewable changes

Read this reference when preparing a code change, deciding whether to split a pull request, or reducing reviewer effort.

Source: [Google Engineering Practices: Small CLs](https://google.github.io/eng-practices/review/developer/small-cls.html).

A change list, or CL, is a proposed code change. These principles also apply to GitHub pull requests and other work submitted for review.

## Prefer one coherent change

- Give each change one understandable purpose.
- Include the related tests, the context the reviewer needs, and enough implementation for the result to make sense on its own.
- Keep it complete, understandable, and safe to merge. Judge its size by the effort needed to review it, not an arbitrary limit on lines or files.

## Why this helps

Focused changes are usually easier to:

- Read and review without setting aside a large block of time.
- Understand and test thoroughly.
- Merge, debug, and roll back.
- Correct before too much work depends on the wrong approach.
- Explain without asking reviewers to reconstruct hidden context.

## Split work when it helps

Consider these ways to divide work:

- Separate unrelated refactors from a bug fix or feature.
- Separate changes with different owners or reviewers.
- Introduce shared groundwork before the feature that uses it.
- Split independent features into complete, reviewable steps.
- Stack dependent changes when the repository supports that workflow.

Each step must build and include its related tests. Keep enough of a behavior together that reviewers do not need several incomplete pull requests to understand it.

## Apply judgment

Some large changes are still simple to review:

- A generated mechanical change.
- A straightforward file deletion.
- A migration whose parts must change together.

When a broad change is needed:

- Explain what it covers.
- Point out the parts that matter to the reviewer.
- Provide useful verification.
