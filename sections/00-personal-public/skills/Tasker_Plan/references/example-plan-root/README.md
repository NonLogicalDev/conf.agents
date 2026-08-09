# Fictional Example Plan Root

This directory contains fictional examples of project plan documents and supporting files. Its repository names, revisions, URLs, folders, goals, tasks, epics, and statuses are illustrative.

Use the examples when a concrete plan document or plan hierarchy would help. Follow the user's instructions and the real project's guidance for the plan root, group names, and issue references. This example does not require a particular naming convention, repository layout, coordination process, owner, goal, or delivery workflow.

This directory represents `<plan-root>`. The actual root name and location come from the user or target project's instructions.

```text
<plan-root>/
├── _index_.md
├── GOALS.md
├── task-001-settings-import/
│   ├── EXEC_STATE.md
│   └── backlog/
│       └── plan-01 2026-07-12 import-settings.md
├── epic-002-settings-migration/
│   ├── EXEC_STATE.md
│   └── active/
│       └── plan-01 2026-07-12 migrate-settings.md
├── task-003-validation-report/
│   ├── EXEC_STATE.md
│   └── completed/
│       └── plan-01 2026-07-12 publish-report.md
├── epic-004-legacy-config/
│   ├── EXEC_STATE.md
│   └── archived/
│       └── plan-01 2026-07-12 remove-legacy-config.md
└── ext-DEMO-42-import-documentation/
    ├── EXEC_STATE.md
    └── active/
        └── plan-01 2026-07-12 document-settings-import.md
```

This fictional root illustrates one project that uses `task-<NNN>-<slug>` for a focused task and `epic-<NNN>-<slug>` for a broader group. These two kinds share a local number sequence. A group for an existing, verified tracker issue uses `ext-${TRACKER_PROJ}-${NUMBER}-<slug>`. Here, `DEMO-42` is a fictional existing Linear issue. Its number comes from the issue, not the local sequence, so the next local task or epic would use `005`.

The user or project chooses which naming convention applies. Every group keeps `EXEC_STATE.md` in its own directory and places each numbered plan in the folder that matches its recorded status. `_index_.md` records the root's purpose, current state, groups, and next action. `GOALS.md` briefly lists the outcome, current priorities, and workstream plan links; it does not contain status updates.

Each example `EXEC_STATE.md` includes `Produced External Artifacts`. It lists fictional outputs actually produced by that group, or says `None produced.` An existing issue does not become a produced artifact just because a group refers to it.

The example files demonstrate possible scope, blockers, validation, decisions, and evidence. They do not give permission to reorganize an existing project or move another task's plan.
