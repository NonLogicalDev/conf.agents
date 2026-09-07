# Worker Contract

The owner is the one person accountable for the project. A subagent owns an assigned piece of work and returns evidence to that owner.

## Choose the work

Use a subagent when its assignment can proceed without constant coordination. Keep small, closely connected work with the owner when delegation would cost more than the work.

Give writing workers distinct files or clearly separate areas. If two pieces of work must edit the same file, sequence them or keep that edit with the owner.

Reuse a worker that already understands the same assignment. Create another worker only when the new work can progress independently.

## Give a complete assignment

```text
Goal: <the result the user actually wants>
Context: <owner AGENTS.md, relevant project facts, and existing plan>
Own: <files or subject this worker may change>
Avoid: <excluded files, actions, or decisions>
Check: <the evidence or test needed>
Return: <outcome, changes, checks, blockers, and next action>
Stop: <the condition that requires the owner to decide>
```

Before starting or resuming an assignment, each worker should read the owner's `<plan-root>/AGENTS.md` when it exists. Re-read it after meaningful steering or instruction changes and before work that depends on its guidance. Share updated instructions with affected workers without interrupting unrelated work.

Do not assign the same shared plan file to several workers. The owner normally updates the root `AGENTS.md`, README, execution state, numbered plans, owner memory, and handoff notes.

Record each worker's verified canonical path, such as `/root/importer` or `/root/documentation`, beside its assignment. If a worker exists but its path cannot be verified, write `worker identity unknown`. Mark every assignment currently being worked `active` in existing execution state and visible runtime steps. Keep work that has not started pending, and mark an assignment complete only after checking its result. Do not invent a worker path or create an otherwise unnecessary record.

When the runtime supports concurrent active steps, mark every running worker's step `in_progress`. Update each assignment when its worker starts, changes tasks, becomes blocked, or finishes.

## Keep the owner available

While workers run, the owner can answer the user, check returned work, prepare integration, record meaningful decisions, and adjust assignments.

When the user changes the goal or scope, tell only the affected workers. Preserve valid work that does not depend on the change.

If a worker discovers a missing decision, permission, fact, or dependency, it should report the issue and the next useful action. The owner decides whether another assignment can continue.

## Return useful evidence

Each worker should report:

- The result and whether its assigned work is complete.
- The files changed or facts established.
- Commands or checks that were actually run and their observed results.
- Useful artifacts, links, or revisions when they exist.
- Remaining blockers, assumptions, risks, and the next action.

An assigned worker is not a completed result. A reported passing check is not proof until the owner has enough evidence to trust or verify it.

## Integrate once

The owner checks worker changes, resolves conflicts, and validates the combined result. Keep one project plan for the actual work; do not create a separate plan only because a worker exists.

Before pausing, handing off, or finishing, record which workers remain active, what each owns, and what should happen next.
