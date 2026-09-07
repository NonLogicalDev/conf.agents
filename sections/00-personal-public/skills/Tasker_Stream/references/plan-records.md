# Plan Records

Use each document for one clear purpose. Keep the records current without copying the same explanation into every file. Read root `MEMORY.md` first when a plan root exists; it is the project's maintained map to these records.

## Keep the plan ahead of the work

For work that already has or needs an active plan, update that plan before changing project documentation, artifacts, or source code. Keep the user's current intent, accepted scope, direction, understanding, decisions, intended changes, and actual state accurate.

Treat the plan as living documentation. Revise it when the user's direction or the owner's understanding changes and the existing plan still fits the work. If the new direction would require an almost complete rewrite, mark the old plan `abandoned`, briefly explain why, preserve its history in `archived/`, and create the next numbered plan in `active/` before changing the project. Recording an action does not grant permission to take it.

A simple direct task does not need a plan merely to satisfy this rule.

## Root README

`<plan-root>/README.md` introduces the project and routes the next person to useful state:

```markdown
# <Project or owner>

[Owners](OWNERS.md)

## Purpose

<What this owner is trying to finish.>

## Current Themes

- [<Theme name>](theme-001-<slug>/EXEC_STATE.md): <current state and next action>.

## Useful Owner Documents

- [<Design or handoff>](_owner/docs/<document>.md): <why it matters>.
```

Record owners only in root `OWNERS.md`; link to it rather than repeating identities or adding ownership frontmatter. Preserve unrelated metadata and useful content. Include a section only when it helps the current project, keep changing next actions in the relevant `EXEC_STATE.md`, and avoid placeholder sections.

Link `GOALS.md` when the owner has a separate goal reference, `AGENTS.md` when useful owner guidance exists, `_owner/state/worktrees.md` when the owner has active worktrees, and `PROSE_STEERING.md` when repeated writing feedback needs to remain easy to find.

## Project memory index

`<plan-root>/MEMORY.md` is the first entry point into owner memory. Keep it as a concise, periodically re-dreamed index of the project rather than another status document:

- Map the project, important source areas and interfaces, active workstreams, owned checkouts, durable decisions, artifacts, risks, and useful owner records.
- Link to `README.md`, `AGENTS.md`, `GOALS.md`, execution state, active plans, `_owner/state/worktrees.md`, `_owner/memory/`, docs, and artifacts when they matter.
- Point to the detailed record and explain why it matters; do not copy changing status, logs, command output, or full inventories into the index.
- Record the last re-dream time, evidence checked, stale or uncertain areas, and next audit trigger.

Read and audit it before substantive owner work starts or resumes. Re-dream it after a material change to the project map, before pausing or completing substantial work, and at least every 30 minutes while a long-running owner is active. Current instructions, source, and observed evidence win when they disagree with the index.

## Goal reference

`<plan-root>/GOALS.md` is a short mission index with the outcome, current priorities, and stable workstream links:

```markdown
# Goals

## [[2026-07-12]] settings-modernization

- **Outcome:** Users can safely import and migrate their settings.
- **Priorities:** Preserve saved values and support reliable imports.
- **Workstreams:**
  - [settings modernization](theme-001-settings-modernization/EXEC_STATE.md): Import supported settings and preserve existing values.
```

Keep the file short enough to use directly with `/goal`. Update it when the goal, mission priorities, or workstream links change, not when the project's status changes.

Keep progress, CI results, commit hashes, blockers, work logs, checklists, and artifact inventories in the relevant `EXEC_STATE.md` or numbered plan. Create `GOALS.md` only when a separate goal reference helps the work.

## Group execution state

`<group>/EXEC_STATE.md` summarizes the group:

```markdown
# <Group>

## Goal

<The result this group owns.>

## Current State

- Status: <active, backlog, completed, or blocked>.
- Next action: <what happens next>.

## Workers

- /root/<worker-name>: <assignment>; status: <active, pending, blocked, or completed>; next: <next action>.

## Produced External Artifacts

- <Verified title or identifier>: <direct link and actual status>.
```

List every worker currently handling an assignment as `active` and identify it by its verified full `/root/...` path. If a worker exists but its path cannot be verified, write `worker identity unknown` instead of inventing a path. Keep concurrent workers visible, leave work that has not started pending, and mark completion only after checking the result.

Keep detailed decisions, work history, and validation in the numbered plan. Record only artifacts that actually exist. Omit an external artifact section when the group has no artifacts unless an existing project format needs it.

## Numbered plan

A numbered plan records the work itself:

```markdown
# <Plan title>

## Purpose

<The outcome and why it matters.>

## Scope

- In scope: <accepted work>.
- Out of scope: <work the owner should not add>.

## Completion Criteria

- [ ] <Observable outcome and the proof it needs>.

## Current State

- Status: <backlog, active, blocked, completed, or abandoned>.
- Next action: <the next useful step>.

## Work Log

- <Date and time>: <meaningful action, decision, or result>; files: <changed paths>; checkout: <inventory link, branch, and full HEAD>; checkpoint: <commit and verified durability>; verification: <command and observed result>; next: <next action or blocker>.

## Validation

- <Check>: <observed result and what it proves>.
```

Add `Decisions`, `Artifacts`, `Blockers`, `Validation`, or `Implementation Steps` only when the work produces information that belongs there. Skip empty sections and avoid repeating the root README or group execution state.

Keep work logs detailed enough for a fresh run to understand what changed, why it changed, which worktree and commit contain it, how it was checked, and what remains. Record meaningful worker assignments, checkout changes, checkpoints, blockers, and handoffs when they happen; do not invent checks or create entries for unchanged polling.

## Worktree recovery inventory

When the owner has worktrees, keep `_owner/state/worktrees.md` as the complete checkout inventory and link each entry to its numbered plan. Record the machine, link to the owner root and its OWNERS.md, repository identity and credential-free remote, source repository, checkout path, Bean or manual layout, branch, full `HEAD`, upstream and ahead/behind state, verified base and merge base, sparse configuration, Git common directory, and whether its refs and objects remain on persistent storage.

Record the latest authorized task checkpoint, its full commit and timestamp, whether its objects survive host or pod replacement, the clean/dirty/staged state, related stack or pull request when useful, and commands that can recreate the same checkout and sparse focus. State missing upstreams, absent checkpoints, and unrecoverable local changes explicitly. Never treat a disposable `/build` path, a commit hash without surviving objects, or an uncommitted change as durable.

Update the inventory after claims, meaningful commits, branch or sparse changes, handoffs, releases, and approved reconstruction. Keep changing details there and in the owning plan; summarize workstreams in `EXEC_STATE.md` and link to their plans instead of copying every field into the root README or goals.

Keep planned work in `backlog/`. Keep active or blocked work in `active/`; record the actual blocker in the plan. Put a plan in `completed/` only after its completion criteria have been verified. Keep abandoned plans in `archived/`. Create a status folder only when it has a plan.

Choose the next two-digit plan number across every status folder in the same group. Preserve the plan number, history, and useful links when its status changes.

## Handoffs and memory

Put a handoff in `_owner/handoff/` when another person or task needs to continue the work. Include the current goal, verified state, active workers, useful files, blockers, and next action. Re-dream `MEMORY.md` before the handoff so it points to the current records and flags stale areas.

Keep durable project facts in `_owner/memory/`. Correct an existing entry whenever verified evidence contradicts it. Save a difficult, useful discovery when later work would otherwise repeat the same investigation; do not record routine status or create memory only to fill a template. Index useful entries from root `MEMORY.md` so a clean run can find them. Write human-readable design or delivery documents in `_owner/docs/`. Save generated reports or other collected output in `_owner/artifacts/`.

Record tool friction in `_owner/wart.tools.md` and problems in skills or instructions in `_owner/wart.guidance.md`. Explain what happened, why it slowed the work, the current workaround, and what would improve it. Keep other `wart.<type>.md` files for distinct problems that are worth remembering.

Keep repeated writing feedback in `PROSE_STEERING.md`, not in a general wart file. Preserve recurring themes and practical improvements without duplicating the project's writing guide.
