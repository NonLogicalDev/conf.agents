# Tests Scenarios

## 06 Explain Fixtures That Prove Behavior

### Prompt

Use `$Code_Pedantics`.

A regression test creates a selected record and an unrelated record, invokes a mutation, and asserts that one callback count is zero. The test passes, but it does not explain why the unrelated record exists or what the zero count proves. Choose the review result. Do not modify files.

### Expectations

- Require the test to explain the scenario and mutation scope.
- Explain the unrelated record and any other unusual fixtures.
- Require positive evidence that the mutation phase ran so the zero count cannot pass vacuously.
- Explain what forbidden side effect the zero count proves did not occur.
- Avoid narrating assertions that already state their own intent.

### Adjacent Valid Case

A small test has one obvious input, one operation, and one direct assertion.

- Leave out comments when its name, fixture names, and assertion already state the behavior.

## 27 Build A Regression Test Around The Failed Behavior

### Prompt

Use `$Code_Pedantics`.

A bug let a cleanup command delete an unmanaged record when a shortcut matched only its kind. The proposed regression test calls the new helper and asserts that one managed record remains. It does not state the failed rule, show that the old behavior would fail, exercise the tempting shortcut, or cover a nearby record that should still be deleted. Choose the review result. Do not modify files.

### Expectations

- State the failed rule before choosing assertions: do not delete unmanaged records just because their kind matches.
- Require the smallest representative regression test whose assertions fail for the old behavior and pass for the fix.
- Add positive evidence that cleanup ran, so preservation of the unmanaged record cannot pass vacuously.
- Add a pressure scenario or counterexample that exercises the tempting shortcut shortcut.
- Add an adjacent valid case showing that a stale managed record is still deleted, so the fix does not overcorrect.
- Rerun the affected focused validation after the repair rather than relying on an earlier green result.

### Pressure Variant

The author offers only a broad test of the complete workflow that passes both before and after the fix.

- Reject it as the regression proof even if it is useful supporting coverage.
- Require a focused test that fails before the fix and passes afterward.

### Adjacent Valid Case

A focused test already reproduces the old deletion, proves cleanup ran, preserves the unmanaged record after the fix, and has separate pressure and nearby valid cases.

- Accept the regression structure without requiring more tests merely for symmetry.
