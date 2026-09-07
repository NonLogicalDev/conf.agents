# Owner Home

Scale the owner's setup to the actual work. A simple task may need only one owner working directly, with no separate plan home. When a project needs coordination or lasting context, use one plan home for one owner and let each workstream retain its own execution state and plans.

The following tree is a map of possible locations, not a list of files and folders to create. Add a goal index, owner harness, memory, tools, delegation, or handoffs only as the work shows why each is useful.

```text
<plan-root>/
├── README.md
├── MEMORY.md               # Required read-first map of the project and memory system.
├── OWNERS.md               # Allowed owner threads and authorization basis.
├── AGENTS.md               # The owner's living project instructions.
├── GOALS.md                # Short mission index for /goal.
├── PROSE_STEERING.md       # When writing feedback repeats.
├── CHANGELOG.md            # When owner changes need a short history.
├── _owner/                 # Only when owner support is needed.
│   ├── README.md
│   ├── artifacts/
│   ├── dashboard/          # When an interactive dashboard or site is useful.
│   ├── docs/
│   ├── handoff/
│   ├── memory/
│   ├── state/
│   │   └── worktrees.md    # When the owner has active worktrees.
│   ├── tools/
│   ├── wart.deviation.md
│   ├── wart.guidance.md
│   └── wart.tools.md
├── theme-000-meta/  # Standing task zero when a plan root exists.
│   ├── EXEC_STATE.md
│   └── active/
│       └── plan-01 Meta maintenance.md
├── theme-001-<slug>/
│   ├── EXEC_STATE.md
│   └── active/
│       ├── plan-01 <description>.md
│       └── plan-02 <related description>.md
└── theme-002-<slug>/
    ├── EXEC_STATE.md
    └── backlog/
        └── plan-01 <description>.md
```

`README.md` explains the project, identifies current groups, and links to OWNERS.md and useful project documents. It contains no ownership declaration. Keep changing group status and next actions in `EXEC_STATE.md` instead of copying them into every record.

`OWNERS.md` records active and revoked owner threads, with verified identities and authorization sources. It is the only ownership record; other documents link here without repeating identities or lists. Delegated workers are not automatically allowed owners; follow the skill's membership and migration rules before owner-level changes.

`MEMORY.md` is required whenever a plan root exists. Read it before other owner memory and before substantive work. It is the project's curated map: explain what the project is, where its important instructions, workstreams, source areas, interfaces, checkouts, artifacts, decisions, risks, and durable memory live. Link to the records that hold the details instead of copying their changing contents.

Keep `MEMORY.md` small enough to scan first and rich enough that a clean agent can choose the right next document without searching blindly. It is not a work log, another `EXEC_STATE.md`, or a dump of `_owner/memory/`. Put detailed status and evidence in the owning plan; put durable findings in `_owner/memory/`; summarize and route them from `MEMORY.md`.

Start with a structure like this and omit sections that do not yet help:

```markdown
# Project Memory

Read this file first when starting or resuming owner work. Re-dream it from current evidence instead of trusting it blindly.

## Project Map

- [Owner entry point](README.md): <what it establishes>.
- [Owner instructions](AGENTS.md): <when to read it>.
- [Current execution](theme-001-<slug>/EXEC_STATE.md): <what this workstream owns>.

## Important Areas

- `<path or interface>`: <why it matters and where deeper evidence lives>.

## Durable Decisions and Constraints

- <decision or constraint>: <source or memory link>.

## Worktrees and Runtime

- [Worktree inventory](_owner/state/worktrees.md): <when it matters>.

## Memory Index

- [<memory note>](_owner/memory/<note>.md): <what it answers>.

## Audit

- Last re-dreamed: <time, evidence checked, and any stale areas>.
- Next audit trigger: <material change, handoff, completion, or 30-minute active-work interval>.
```

Audit and re-dream the file when work starts or resumes, after a material change to the project map, before pausing or completing substantial work, and at least every 30 minutes during active long-running work. Compare it with current user direction, `README.md`, `AGENTS.md`, `GOALS.md`, execution state, active plans, worktree inventory, durable memory, source, and observed results. Remove stale claims, repair links, record uncertainty, and update the last audit evidence and next trigger.

`GOALS.md` is an optional short mission index for `/goal`. Keep the outcome, current mission priorities, stable workstreams, and direct links to their existing plans or `EXEC_STATE.md` files. Put status, progress, blockers, review details, and work logs in `EXEC_STATE.md` or the relevant plan.

`AGENTS.md` is the owner's living project harness. Make it complete enough for another thread to reproduce the owner's intended behavior from that single reference. Include applicable user instructions, the mission and `GOALS.md`, accepted scope and permissions, project conventions, workstream entry points, delegation, verification, and durable operating decisions. The owner and its workers should read it when starting or resuming work and before work that depends on its guidance. Have affected workers re-read it after meaningful steering or instruction changes.

If an older `OWNER_PROMPT.md` exists, carry its useful guidance into `AGENTS.md` while preserving other instructions that still apply. Replace obsolete instructions when current user direction supersedes them, but do not take over another owner's work. The owner may update its own instructions without another request as accepted scope, user intent, priorities, workstreams, or project needs change. This does not permit the owner to expand scope or change permissions. Keep running status in execution state or plans. Reading `AGENTS.md` does not grant the authority recorded only in `OWNERS.md`.

`PROSE_STEERING.md` collects repeated writing feedback. `CHANGELOG.md` can record meaningful improvements to owner guidance, tools, or structure. Keep these files optional and link them from the root README when they help.

Create `_owner/README.md` only when it makes existing owner material easier to understand. Use these folders only when their contents already have a clear purpose:

- `artifacts/` holds collected reports, review material, exports, and other saved outputs.
- `dashboard/` holds an interactive dashboard, small site, or other project view intended for publication when the work calls for one.
- `docs/` holds design summaries, architecture notes, delivery plans, and other documents a person should be able to read directly.
- `handoff/` holds dated notes that explain how another owner can resume or continue the work.
- `memory/` holds verified observations, project facts, and practical learnings that can prevent repeated investigations.
- `state/` holds useful owner state, such as a verified worktree recovery inventory, small JSON files, or cached inventories. Create `state/worktrees.md` when the owner has worktrees; record every checkout's active plan, machine, repository, path, branch, full `HEAD`, upstream, base, sparse focus, checkpoint, Git common directory, storage durability, and reconstruction commands.
- `tools/` holds reusable helpers and jigs that save repeated project work.
- `wart.tools.md` records useful evidence about environment problems, tools, builds, commands, interfaces, and repeated workarounds.
- `wart.guidance.md` records problems in skills, owner instructions, or project guidance that make the agent less effective.
- `wart.deviation.md` records every intentional departure from an existing owner convention or prior guidance and explains its reason.
- Other `wart.<type>.md` files may capture a different recurring problem.

Create wart files only when there is something useful to record. Do not create empty support folders, move unrelated artifacts, or rename existing files unless the user asks.

Keep workstream groups directly under `<plan-root>`. Follow the naming convention in the user's instructions, applicable `AGENTS.md`, or existing project. When no convention is established, use broad `theme-<nnn>-<slug>` groups, reserving `theme-000-meta` for the standing maintenance task and beginning delivery work with `theme-001-<slug>`. Follow the task-zero requirements in `SKILL.md`; preserve an equivalent existing maintenance record instead of creating a duplicate. Treat each task, theme, epic, or other group as a meaningful work domain with a coherent responsibility, lifecycle, body of decisions, or set of interfaces; do not create one per ask, effort, bug, worker, pull request, or plan. Track each distinct substantial goal in its own `plan-<nn> <description>.md` record within its existing theme; have delegated workers share the plan covering their assignment. Compare new work with existing domains before creating another group. When observed asks and requirements show that an existing domain is too narrow, broaden its stated domain and update its execution state, plans, memory index, goals, and owner instructions as needed while preserving its identity and useful history. Add another theme only when the work forms a distinct, lasting responsibility. Broadening the organizational domain does not expand the user's accepted work scope or permissions. Preserve existing names, types, and number widths; explicit group types remain available when the user or project chooses them. When numbered groups share a sequence, continue it without filling gaps. When a verified issue forms its own distinct workstream and the project requests issue-based groups, follow its established convention, such as `ext-<project>-<number>-<slug>`.

Keep `EXEC_STATE.md` at the group root. Put each numbered plan in `backlog/`, `active/`, or `completed/` according to its actual state. Preserve abandoned plans in `archived/`. Create a status folder only when a plan belongs there.

Keep the owner directory sufficient for another clean run to read `MEMORY.md`, find the current workstreams, read their work logs, locate every owned checkout, distinguish durable Git objects from disposable worktree files, and rebuild approved work without relying on a previous conversation. Link the worktree inventory from `MEMORY.md` and useful owner records instead of duplicating changing checkout state.

Treat `_index_.md` as an older entry point. New owner homes use `README.md`. Preserve an existing index and its links until the user requests a migration.
