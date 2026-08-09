# Tests

## Test Behavior, Not Geography

Keep unit tests tightly focused on the code they are meant to test. Do not use them to test code that belongs to another unit without an exceptionally good reason. Use integration tests when you need to check how components work together.

Good tests use the public interface and protect its public contract while leaving the internal implementation free to change. They do not repeat the implementation in another form. Protect the invariants the code needs to remain correct, meaningful edge cases, and normal successful use. Check observable results instead of internal steps. Do not lock in details outside the public contract without a good reason. Name tests and assertions so they explain what the code promises; add comments when context or a useful summary makes the test easier to understand.

Start with the smallest test that proves the behavior. Add an integration test or run the real entrypoint when a focused test cannot represent the interaction. Passing tests support review; they do not replace reasoning about correctness, edge cases, security, operations, or ownership.

## Explain Fixtures That Prove Behavior

Explain fixtures and assertions when the proof would otherwise be unclear. Document or comment:

- the scenario and mutation scope represented by the request or fixture;
- why each unusual input, record, state, dependency, or helper exists;
- which workflow or side effect should occur;
- which workflow or side effect to prevent; and
- what an absent call, scan, mutation, log, or message proves.

Document shared helpers and fixtures when their construction encodes a policy, user permission, access rule, or important proof. Keep comments close to the setup or assertion they explain.

Leave comments out of a simple test when its name, fixture names, and assertions already state the behavior.

## Prove The Operation Actually Ran

A zero count, absent call, missing mutation, or preserved record is meaningful only when the test also proves the operation actually ran and could have taken the forbidden path.

Include positive evidence such as a changed managed record, observed normal callback, successful entrypoint result, or other proof that the tested phase executed. Explain why the negative assertion would fail if the regression returned.

## Build Regression Tests Around The Behavior

For a regression:

1. State the behavior or rule that failed.
2. Reproduce the smallest representative failure when practical.
3. Add a focused test that fails for the old behavior and passes for the fix.
4. Add a pressure or counterexample when the shortcut is tempting.
5. Add an adjacent valid case to avoid overcorrecting.
6. Rerun affected validation after the repair.

Do not preserve a real incident transcript, private data, or exact private context in reusable tests. Use invented examples that preserve the decision and the reason an agent might make the wrong choice.

## Validate In Layers

Run focused behavior tests first, then proportional formatter, linter, type-check, integration, build, and repository checks. After a fix, rerun the affected checks before relying on an earlier green result.

When completion depends on real behavior, follow [validation.md](validation.md). Unit tests, mocks, type checks, and builds support proof from the actual entrypoint; they do not replace it.

## Test Behavioral Guidance Independently

When testing a skill, workflow, or other guidance that affects behavior, use an independent evaluator. Give it the guidance and realistic artifacts from the task, not the diagnosis, expected answer, or planned repair. Do not let the test change shared files. Use a temporary directory if the test needs to write.

Capture the raw result before grading. A scenario passes only when every expectation holds and no contrary behavior appears. Report mechanical checks, behavioral evidence, unexecuted tests, and remaining risk separately.
