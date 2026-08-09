# Plan Placement

Read this reference when finding an existing plan or choosing the `<plan-root>` for a new one.

## Choose the plan root

Identify the actual target project and read its applicable instructions. Choose `<plan-root>` in this order:

1. Use the exact plan or location specified by the user or the project's instructions.
2. Otherwise, use the target project's established plan root. When the target is a project inside a repository, prefer its own plan root over the repository's. Use the repository's established plan root when no root specific to the project applies.
3. If there is no applicable local location, use `~/.agents-plans/<project-slug>/`. Derive `<project-slug>` from the verified project or repository name.

The name of a local plan root is configurable. `.agents-plans/`, `planning/`, and other documented project directories are possible conventions, not required names. The directory where a session starts does not determine the target project's plan root.

## Create the required group hierarchy

After resolving `<plan-root>`, create or reuse a group directly under it. Choose its name from the user's instructions, the applicable `AGENTS.md`, and the project's established plan convention:

```text
<plan-root>/
├── _index_.md                 # When owner or project guidance requires it.
└── <plan-group>/
    ├── EXEC_STATE.md
    ├── active/
    │   └── plan-01 <description>.md
    ├── backlog/
    │   └── plan-02 <description>.md
    ├── completed/
    │   └── plan-03 <description>.md
    └── archived/              # Only for actually archived plans.
```

When the applicable instructions choose them, supported group names include:

- `task-<NNN>-<slug>` for a local task.
- `epic-<NNN>-<slug>` for a broader group of related local work.
- `ext-${TRACKER_PROJ}-${NUMBER}-<slug>` for work backed by a verified external issue, such as `ext-DEMO-42-restore-import`.

Choose a short slug that describes the work. For an external issue, preserve its actual tracker project and issue number. Verify that the issue exists. Do not invent an issue or create one just to name a directory.

When applicable instructions require a shared sequence for local tasks and epics, find the highest existing `task-` or `epic-` number and use the next number, written as three digits. Start at `001`. For example, after `task-001-prepare-checkout`, the next local group may be `epic-002-refresh-projects`. External issue numbers belong to their tracker; they do not consume or reset the local sequence. Preserve gaps and older group names.

Keep `EXEC_STATE.md` directly in the group. Place each numbered plan in its `active/`, `backlog/`, or `completed/` directory to match the plan's actual status. Use `archived/` only when a plan has actually been archived.

Map planned or queued work to `backlog/`, and active or blocked work to `active/`. Keep an actual blocker in the plan and use `completed/` only after completion evidence exists. Do not create `planned/` or `blocked/` status directories.

Number plans separately within each group. Choose the next available two-digit plan number across all of its status directories. Preserve a plan's number, content, and history when the user permits a change in status. Update links to its new location. Do not reorganize another owner's plans or older files merely because their layout differs.

When owner or project guidance requires an `_index_.md`, create or update it at `<plan-root>`. Record the owner's purpose, current state, groups, relevant existing plans, and next actions. Verify the required files and directories before claiming the plan is initialized.

## Reuse or create the plan

Read an existing plan before choosing a new file. Preserve its location, scope, and useful history. Reuse a group that already covers the work, even when its directory predates the naming convention. If an existing plan at the root covers the work, link it from the required root index. Do not move, rename, duplicate, or discard an existing plan unless the user asks. Put new substantive plans and their execution state in the appropriate group.

If a focused lookup is needed, inspect only the target project and its applicable plan roots. Create a plan only when the requested work needs one and writing to the chosen location is allowed. If the target project, plan root, or write permission is unclear, ask one focused question.

Follow an explicitly documented different project hierarchy when it applies. Do not restructure an existing document without the user's permission.

## Record an epic's produced artifacts

Create or reuse `EXEC_STATE.md` in each group, task, or epic directory. Include a `Produced External Artifacts` section listing the external outputs actually created for that work:

- Pull requests.
- Linear tickets.
- Documents.
- Slack threads started for the work.
- Other relevant external artifacts.

Include each artifact's verified title or identifier, a direct link, and its known status. Record an unknown link or status as unknown instead of inventing it. If the epic has not produced any external artifacts, say `None produced.`

Keep the list current as verified artifacts are produced or materially change. Recording an artifact does not require creating one, publishing a document, starting a conversation, polling an external service, or changing a project's existing plan layout.

## Read the example correctly

The fictional example demonstrates the required directory structure: each group sits directly under `<plan-root>`, keeps `EXEC_STATE.md` at its root, and files numbered plans in the matching status directory. Its project, names, plan contents, goals, and supporting files do not impose a delivery workflow on a real project.
