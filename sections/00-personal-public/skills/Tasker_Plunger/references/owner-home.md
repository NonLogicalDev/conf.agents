# Owner Home

Scale the owner's setup to the actual work. A simple task may need only one owner working directly, with no separate plan home. When a project needs coordination or lasting context, use one plan home for one owner and let each workstream retain its own execution state and plans.

The following tree is a map of possible locations, not a list of files and folders to create. Add a goal index, owner harness, memory, tools, delegation, or handoffs only as the work shows why each is useful.

```text
<plan-root>/
├── README.md
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
│   ├── tools/
│   ├── wart.deviation.md
│   ├── wart.guidance.md
│   └── wart.tools.md
├── <type>-001-<slug>/
│   ├── EXEC_STATE.md
│   └── active/
│       └── plan-01 <description>.md
└── <type>-002-<slug>/
    ├── EXEC_STATE.md
    └── backlog/
        └── plan-01 <description>.md
```

`README.md` records `owner_thread` in YAML frontmatter, explains the project, identifies current groups, and links useful owner documents. Keep changing group status and next actions in `EXEC_STATE.md` instead of copying them into every record.

`GOALS.md` is an optional short mission index for `/goal`. Keep the outcome, current mission priorities, stable workstreams, and direct links to their existing plans or `EXEC_STATE.md` files. Put status, progress, blockers, review details, and work logs in `EXEC_STATE.md` or the relevant plan.

`AGENTS.md` is the owner's living project harness. Make it complete enough for another thread to reproduce the owner's intended behavior from that single reference. Include applicable user instructions, the mission and `GOALS.md`, accepted scope and permissions, project conventions, workstream entry points, delegation, verification, and durable operating decisions. The owner and its workers should read it when starting or resuming work and before work that depends on its guidance. Have affected workers re-read it after meaningful steering or instruction changes.

If an older `OWNER_PROMPT.md` exists, carry its useful guidance into `AGENTS.md` while preserving other instructions that still apply. Replace obsolete instructions when current user direction supersedes them, but do not take over another owner's work. The owner may update its own instructions without another request as accepted scope, user intent, priorities, workstreams, or project needs change. This does not permit the owner to expand scope or change permissions. Keep running status in execution state or plans. Reading `AGENTS.md` does not transfer the owner recorded in `README.md`.

`PROSE_STEERING.md` collects repeated writing feedback. `CHANGELOG.md` can record meaningful improvements to owner guidance, tools, or structure. Keep these files optional and link them from the root README when they help.

Create `_owner/README.md` only when it makes existing owner material easier to understand. Use these folders only when their contents already have a clear purpose:

- `artifacts/` holds collected reports, review material, exports, and other saved outputs.
- `dashboard/` holds an interactive dashboard, small site, or other project view intended for publication when the work calls for one.
- `docs/` holds design summaries, architecture notes, delivery plans, and other documents a person should be able to read directly.
- `handoff/` holds dated notes that explain how another owner can resume or continue the work.
- `memory/` holds verified observations, project facts, and practical learnings that can prevent repeated investigations.
- `state/` holds data meant for tools, such as small JSON files or cached inventories.
- `tools/` holds reusable helpers and jigs that save repeated project work.
- `wart.tools.md` records useful evidence about environment problems, tools, builds, commands, interfaces, and repeated workarounds.
- `wart.guidance.md` records problems in skills, owner instructions, or project guidance that make the agent less effective.
- `wart.deviation.md` records every intentional departure from an existing owner convention or prior guidance and explains its reason.
- Other `wart.<type>.md` files may capture a different recurring problem.

Create wart files only when there is something useful to record. Do not create empty support folders, move unrelated artifacts, or rename existing files unless the user asks.

Keep workstream groups directly under `<plan-root>`. Follow the naming convention in the user's instructions, applicable `AGENTS.md`, or existing project. When no convention is established, use `task-<num>-<slug>` by default and choose another type only when the work calls for it. Do not ask the user to pick a name. Preserve existing names, types, and number widths. When numbered groups share a sequence, continue it without filling gaps. For a verified issue, follow the project's issue naming convention, such as `ext-<project>-<number>-<slug>`.

Keep `EXEC_STATE.md` at the group root. Put each numbered plan in `backlog/`, `active/`, or `completed/` according to its actual state. Preserve abandoned plans in `archived/`. Create a status folder only when a plan belongs there.

Treat `_index_.md` as an older entry point. New owner homes use `README.md`. Preserve an existing index and its links until the user requests a migration.
