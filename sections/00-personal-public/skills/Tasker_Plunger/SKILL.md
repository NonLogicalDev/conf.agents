---
name: Tasker_Plunger
description: Adapt task ownership and execution to work of any size, from a simple direct fix to feature delivery spanning months. Use when work benefits from one accountable owner, useful delegation, evolving guidance, durable context, repeated iterations, or several workstreams.
---

# Tasker Plunger

Tasker Plunger is a self-adapting system for work of any size. Keep one owner accountable, delegate when useful, preserve enough context to resume, and check the result. A simple task may need only direct action; a large feature may span months, many iterations, and several workstreams. Evolve the system in the spirit of the user's intent, and let it grow only when the work needs it.

## Scale the core to the work

- Handle a simple task directly when one owner can finish it without helpers, planning files, or durable owner support.
- For work that spans many iterations, workstreams, or months, add useful owner guidance, plans, delegation, memory, project tools, and handoffs as the work actually needs them.
- Reuse the project, owner, and plan that already cover the work. The current task can be the owner.
- Keep the owner responsible for user direction, accepted scope, worker assignments, shared plan state, integration, and completion.
- When the work needs durable records, choose `<plan-root>` from the user, existing project guidance, or an established project convention.
- Begin every `<plan-root>/README.md` with YAML frontmatter containing `owner_thread: "codex://threads/<thread-id>"`. Use the verified `CODEX_THREAD_ID` for the current owner; never invent a thread ID.
- Before resuming or changing owner state, compare the recorded owner thread with the current thread. If they differ, report the mismatch and stop. Preserve existing frontmatter and README content; do not silently replace another owner.
- For work that needs a plan, use a brief root `README.md`, a group `EXEC_STATE.md`, and a numbered plan in the status folder that matches its current state.
- Let each record serve its own purpose: the README routes the reader, execution state summarizes the group, and the plan holds useful work details.
- When a separate goal reference helps, keep `<plan-root>/GOALS.md` as a short mission index suitable for `/goal`. State the intended outcome, current priorities, and workstreams; link each workstream to its actual execution state or plan.
- Preserve existing goal headings and concise `Outcome`, `Priorities`, and `Workstreams` entries. Keep changing status, CI results, commit hashes, blockers, work logs, and checklists in `EXEC_STATE.md` or the relevant plan.
- Follow the workstream naming convention in the user's instructions, applicable `AGENTS.md`, or project. Otherwise, choose `task-<num>-<slug>` without asking the user; use another `<type>` when the work calls for it.
- Update a record when an important result, decision, blocker, ownership change, or next action changes. Do not create records for their own sake.
- Trust current instructions, actual source, and observed results over an old summary, memory file, or dashboard.

## Update the plan before changing the project

- Treat every active plan as living documentation of the user's intent, accepted scope, direction, current understanding, decisions, and intended changes.
- Before changing project documents, artifacts, or source code, update the relevant active plan so it accurately reflects the intended work and current project state.
- Revise the plan first whenever new steering, evidence, or understanding changes the work. Keep the relevant execution state and plan current as meaningful results, blockers, decisions, or next actions change.
- If a plan would need an almost complete rewrite, mark it `abandoned`, preserve its history in `archived/`, and create the next numbered active plan before changing the project. Keep ordinary revisions in the existing plan.
- An updated plan does not grant new permission or expand the accepted scope. A simple task that does not need a plan does not gain one merely to satisfy this rule.

## Delegate when it helps

- Give independent workers distinct files or subjects, useful context, expected results, and the checks that matter.
- For planned work, update the active plan before asking a worker to change project documents, artifacts, or code.
- Mark every task currently being worked `active` and identify its worker by the verified full `/root/...` path. Write `worker identity unknown` when a worker exists but its path cannot be verified.
- Mark every running runtime step `in_progress` when the UI supports concurrent active steps. Keep all assignments visible, leave work that has not started pending, and mark completion only after checking the result.
- Keep one owner responsible for shared plan files. Assign a shared file to a worker only when that worker is its only writer.
- Reuse workers for related follow-up. Keep closely connected or trivial work with the owner when that is simpler.
- Preserve established workstream names and update only the workers affected by new direction.
- Check meaningful worker results and the combined outcome before claiming completion.

## Keep owner memory and tools useful

- Keep verified observations, project facts, and practical learnings from tasks and workstreams in `_owner/memory/` when they can help later work.
- Read relevant owner memory before investigating a question or assigning related work.
- Always correct relevant owner memory when verified evidence contradicts an existing entry. Save new findings that were difficult to establish and will help avoid repeating the same investigation.
- Update matching entries instead of creating competing accounts, and include the supporting evidence needed to reuse or verify each finding.
- Share useful findings with affected workers. Keep changing status, blockers, and routine work logs in the plans that own them.
- Build small reusable helpers in `_owner/tools/` when they save repeated work, replace awkward commands, or prevent common mistakes. Test their useful behavior and make any action that changes files explicit.

## Keep owner instructions current

- Keep the owner's living project instructions in `<plan-root>/AGENTS.md`. When the project needs standing guidance, another thread should be able to reproduce the owner's intended behavior from this one document.
- Have the owner and every affected worker read `AGENTS.md` when starting or resuming work, after meaningful steering or instruction changes, and before work that depends on those instructions.
- Keep `AGENTS.md` current with the user's applicable instructions, mission, `GOALS.md`, accepted scope and permissions, project conventions, workstream entry points, delegation and verification approach, and durable operating decisions.
- Treat an existing `<plan-root>/OWNER_PROMPT.md` only as migration input. Preserve its useful guidance and other instructions that still apply, replace instructions that current user direction supersedes, and never overwrite another owner's files.
- Proactively adapt `AGENTS.md` as the user's direction, accepted project scope, priorities, workstreams, or owner responsibilities change. Remove obsolete instructions and tailor the document to the actual work.
- The owner may update its own writable `AGENTS.md` without waiting for another request when the change stays true to the user's intent, accepted scope, and existing permissions.
- Link to execution records instead of copying running status. Current user instructions and actual project guidance remain authoritative; owner instructions do not replace them or transfer ownership to another thread.

## Add owner support only when useful

- Use `_owner/docs/`, `_owner/artifacts/`, `_owner/handoff/`, `_owner/memory/`, `_owner/tools/`, `_owner/state/`, or `_owner/dashboard/` only when the work actually needs that material.
- Use `_owner/dashboard/` for an interactive project dashboard, small site, or other view intended for publication when that is part of the requested work.
- Add `PROSE_STEERING.md` or `CHANGELOG.md` only when repeated writing feedback or important owner changes justify it.
- Record recurring problems with tools or the work environment in `_owner/wart.tools.md`. Keep problems with skills or instructions in `_owner/wart.guidance.md`. Include useful evidence, impact, and a response; add stable IDs or counts only when they help.
- Keep a current blocker in the active plan. Use a wart file only for a problem or lesson that will matter again.
- Create support folders, status directories, dashboards, checklists, summaries, ledgers, or helper scripts only when the current work needs them.

## Evolve with the user's intent

- Proactively improve the owner harness, project helpers, or this skill's writable source when user direction, accepted scope changes, repeated friction, or verified project needs show that a change would better serve the user's intent.
- Keep the change small, preserve the accepted scope and existing permissions, and check the result.
- Prefer an owner-local improvement when it solves the problem. Change shared skill guidance only when the lesson applies more broadly.
- Record each deliberate departure from the owner's established instructions, conventions, or approach in `_owner/wart.deviation.md`. Explain what changed, why it better serves the user's intent, and which constraints still apply.
- Do not use self-improvement to widen permissions, ignore the user, change another owner's work, publish externally, or bypass a read-only installation.

## Read a reference only when needed

- [Owner home](references/owner-home.md): Where optional project support belongs.
- [Owner operating guide](references/owner-operating-guide.md): How to resume or improve a long-running owner.
- [Worker contract](references/worker-contract.md): How to assign independent helpers and combine their results.
- [Plan records](references/plan-records.md): How to keep a project plan useful without repeating it elsewhere.
- [Tool, guidance, and deviation warts](references/warts.md): How to record recurring friction and explain deliberate departures.
- [Prose steering](references/prose-steering.md): How to keep repeated writing feedback useful.

## Use the plan helper

The public helper in `scripts/plan.py` has no external dependencies:

```sh
python3 scripts/plan.py doctor <plan-root>
python3 scripts/plan.py next-group <plan-root> <slug>
python3 scripts/plan.py next-group <plan-root> <type> <slug>
python3 scripts/plan.py next-plan <group-root> "<description>"
python3 scripts/plan.py create-group <plan-root> <slug> "<description>"
python3 scripts/plan.py create-group <plan-root> <type> <slug> "<description>"
```

Use `doctor` to inspect an existing plan home. The naming commands print the next available group or plan without changing files. Pass `--digits 2` when project instructions require two-digit group numbers. Use the helper only when it matches the project's actual convention, and run `create-group` only when the user or project instructions allow the requested plan files to be created.

`doctor` reports a missing, invalid, or mismatched owner thread. Pass `--owner-thread <thread-id>` to `doctor` or `create-group` when `CODEX_THREAD_ID` is unavailable.

## Tests

When changing this skill, read [tests/README.md](tests/README.md). Run the relevant scenarios with fresh subagents and execute `python3 -B -m unittest discover -s tests -p test_plan.py -v` from the skill directory.
