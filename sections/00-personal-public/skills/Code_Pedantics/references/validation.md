# Validation And Real Behavior

Run focused behavior tests, formatting, and the repository's proportional validation. When the completion condition depends on real behavior and a safe representative run exists, execute the real entrypoint in an isolated, reversible environment:

Preserve build caches by default. Do not change global cache keys, action environments, output roots, or cache contents merely to force fresh validation: that can invalidate an entire transitive dependency graph and turn a focused proof into unrelated rebuilds. When fresh execution is actually required, disable caching for the narrowest relevant test and state why it is needed. Prefer focused uncached tests plus ordinary cached proportional validation; discard a broader build cache only when the user or repository contract explicitly requires it and the scope and cost are known. For Bazel, select the exact test targets and prefer `--nocache_test_results` or `--cache_test_results=no`; do not change a global action environment or cache version just to rerun tests.

1. Record the clean starting condition and data to leave unchanged.
2. Perturb only the managed input into a valid stale or incorrect state.
3. Invoke the normal entrypoint with its normal defaults.
4. Verify the expected result and confirm that unrelated state remains unchanged.
5. Run it again when idempotence is part of the contract.
6. Restore the intended clean state.

Do not mutate production just to satisfy this rule. If no safe run exists, state the exact blocker and do not claim the complete workflow was verified.
