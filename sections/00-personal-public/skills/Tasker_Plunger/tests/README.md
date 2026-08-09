# Tasker Plunger Behavioral Tests

Run each scenario with a fresh subagent that has an empty context window. Give the subagent `$Tasker_Plunger` and the scenario prompt, but do not show the expectations.

Keep scenarios read-only or use a temporary directory created for the test. Capture the raw response and compare it with the expectations afterward.

Check that the current owner can complete simple tasks directly without a separate owner record, delegation, or unnecessary files. For substantial work, check that one owner controls shared plan state, the root `README.md` identifies the verified owner, and ownership mismatches leave another owner's records unchanged. Keep active plans current with user intent, accepted scope, reasoning, intended changes, and project state before changing documentation, artifacts, or code, or delegating those changes. When a plan would need an almost complete rewrite, check that it is marked `abandoned`, preserved in `archived/`, and replaced with the next numbered active plan. A plan never grants permission. Correct verified mismatches in owner memory, preserve useful findings that were difficult to discover, and keep routine status in its plan. Check that workers receive separate assignments, every running assignment remains visibly active with its verified `/root/...` worker path or `worker identity unknown`, unstarted work remains pending, and reusable helpers justify their existence. Keep `<plan-root>/AGENTS.md` as the living owner instructions, including current user direction, the mission, permissions, workstreams, delegation, verification, and useful decisions. Owners and workers must read it when starting or resuming their work and re-read it after relevant instructions change. Preserve existing `AGENTS.md` content and useful guidance from an older `OWNER_PROMPT.md` during migration. The owner may update its own writable instructions as approved priorities, workstreams, and project needs evolve; it cannot expand the approved scope, change permissions, or take over another owner's project. Keep any root `GOALS.md` as a short `/goal` index of the intended outcome, current priorities, and links to its workstreams; keep changing progress and status in `EXEC_STATE.md` or the numbered plan. Keep recurring tool and environment problems in `_owner/wart.tools.md` and explain intentional departures in `_owner/wart.deviation.md`. Reuse useful owners, plans, workers, memory, and tools across longer projects. Add owner support only when useful and require evidence before claiming completion.

Run the Python helper tests from the skill directory:

```sh
python3 -B -m unittest discover -s tests -p test_plan.py -v
```

Use [scenarios.md](scenarios.md) for the behavioral prompts and pressure cases.
