# Code Behavior Diff Tests

Run each scenario with a fresh subagent that has an empty context window. Give it the skill invocation and scenario prompt, but not the expectations. Keep tests read-only and provide only the source evidence named by the prompt.

Capture the raw response and compare it with the expectations afterward. A scenario passes only when every expectation holds and no contrary behavior appears. When repairing guidance, first run the relevant scenario without the repair, then rerun it after the edit. Also test the adjacent valid case where provided.

Use [scenarios.md](scenarios.md) for the reusable scenarios. These are synthetic examples, not links to real repositories.
