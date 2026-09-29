# Code Design Plan Behavioral Tests

Run each scenario with a fresh subagent that has an empty context window. Give the subagent the named skill invocation (`$Code_DesignPlan`), the skill folder location, and only the scenario prompt and its included source material. Do not give it the expectations or intended answer. Keep tests read-only or confined to a task-local temporary directory.

Capture the raw response and compare it with the expectations afterward. A scenario passes only when every expectation holds and no contrary behavior appears. Judge the resulting design and prose, not a match against fixed wording.

For repair-loop scenarios, first run the relevant scenario against the current guidance. After the edit, rerun the same scenario. Also run a pressure variant or adjacent valid case when the scenario defines one.

Use [scenarios.md](scenarios.md) for the reusable gamut. These are synthetic pattern-application tests, not records of actual incidents. Keep run evidence outside the installed skill unless the project provides a separate location for test results.
