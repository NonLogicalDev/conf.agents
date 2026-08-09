---
name: Tasker_Plan
description: Create, resume, and maintain clear project plan documents with explicit scope, current status, next actions, blockers, decisions, learnings, and verified completion.
---

# Tasker Plan

Use this skill when substantive work needs a written plan, an existing plan needs to be continued, or project decisions and progress need to remain understandable beyond the current conversation.

This skill owns the plan document. It defines what the plan records and when the plan should be updated. Project instructions own how the underlying work is performed.

Use `$Tasker_Plunger` when one project owner needs to coordinate subagents, maintain shared owner context, and integrate several workstreams.

## Find or create the plan

- Identify the actual target project and read its applicable instructions.
- Choose `<plan-root>` in this order:
  - Use the plan or location specified by the user or project instructions.
  - Otherwise, use the target project's usual plan location. Use the repository's location when the project has no separate convention.
  - If neither exists, use `~/.agents-plans/<project-slug>/`, deriving the slug from the verified target project.
- Treat `.agents-plans/` as one possible local convention, not a required directory name.
- Reuse an existing plan that already covers the requested work.
- Create a plan only when the work needs one and writing to `<plan-root>` is allowed.
- Ask one focused question if the project, plan, scope, location, or write permission cannot be established.
- Do not create a plan for a simple task unless the user asks.

## Use the required plan hierarchy

Keep each substantive task, epic, or tracked issue in its own directory directly under the chosen `<plan-root>`. Follow the user's instructions, the applicable `AGENTS.md`, and the project's existing naming convention:

```text
<plan-root>/
├── _index_.md                 # When owner or project instructions require it.
├── GOALS.md                   # When a short mission index helps.
└── <plan-group>/
    ├── EXEC_STATE.md
    ├── active/
    │   └── plan-01 <description>.md
    ├── backlog/
    │   └── plan-02 <description>.md
    ├── completed/
    │   └── plan-03 <description>.md
    └── archived/              # When an existing plan is actually archived.
```

- Choose a group name from the user's instructions, the applicable `AGENTS.md`, and the project's established plan convention. Supported patterns, when that guidance chooses them, include:
  - `task-<NNN>-<slug>` for a local task.
  - `epic-<NNN>-<slug>` for a broader group of related local work.
  - `ext-${TRACKER_PROJ}-${NUMBER}-<slug>` for work backed by a verified external issue, such as `ext-DEMO-42-restore-import`.
- For an external issue, preserve the actual tracker project and issue number. Verify that the issue exists. Do not invent an issue or create one just to name a directory.
- When local instructions use one sequence for tasks and epics, find the highest existing `task-` or `epic-` number and choose the next number, written as three digits. Start at `001`. Do not fill old gaps, renumber existing directories, or include external issue numbers in that sequence.
- Reuse an existing directory that already covers the work, even if its name predates this convention. Do not rename or move it unless the user asks. Do not treat `<plan-root>` as the group.
- Keep `EXEC_STATE.md` directly inside the group. Put each numbered `plan-<NN> <description>.md` in that group's `active/`, `backlog/`, or `completed/` directory according to its actual status. Use `archived/` only for a plan that is actually archived.
- Put planned or queued work in `backlog/`. Keep unfinished or blocked work in `active/`; record the exact blocker inside the plan. Put only verified finished work in `completed/`.
- Number plans separately within each group. Choose the next available two-digit plan number across all of that directory's status directories. Preserve an existing plan's number and history when its status changes.
- Keep the plan's recorded status and status directory consistent. When work that the user has permitted actually changes status, move the same plan to its group's matching directory and update the affected links. Do not move another owner's plan or restructure unrelated history.
- When owner or project instructions require a root `_index_.md`, create or update it with the owner's purpose, status, groups, and next actions.
- When the project uses `GOALS.md`, keep it as a short `/goal` reference with the intended outcome, current mission priorities, and links to the actual workstreams. Keep changing status, progress, and evidence in execution state or numbered plans.
- Preserve existing plans at the root and their history. Link relevant older plans from `_index_.md`. Do not move, rename, duplicate, or discard them without the user's permission. Put new work in the appropriate group.
- Verify the root, required index, group, `EXEC_STATE.md`, correct status directory, and numbered plan actually exist before reporting setup.

## Shape and maintain the plan

Each numbered plan is the durable record of its own work. Preserve useful existing content and keep these sections current:

- **Purpose:** The outcome and why the work matters.
- **Scope:** Say which work is included and which work is not.
- **Completion Criteria:** Specific outcomes with their actual proof.
- **Blockers:** What is blocked, the observed cause, the next action, and the evidence that will clear it. Write `None.` when none exist.
- **Artifacts:** Each relevant local or external work product, its type, verified state, contents, exact location, and applicable revision.
- **Decisions:** Material decisions, their rationale, and whether a later verified decision superseded them.
- **Implementation Steps:** Meaningful checked and unchecked work.
- **Work Log:** A chronological, history-preserving record of material work, observed results, changes in direction, and next actions.
- **Validation:** The relevant check or command, expected result, observed result, sign of failure, and what the outcome actually proves.
- **Unfinished Work:** Current remaining tasks, dependencies, and next steps.

Keep an existing accurate status or status frontmatter. Add context, steering, recovery, or other sections only when they help explain the real work. Read [Plan format](references/plan-format.md) before creating or substantially updating a plan.

## Keep the plan current

- Update status when work begins, changes, becomes blocked, or finishes.
- Append a dated `Work Log` entry after each meaningful investigation, implementation milestone, decision, validation result, changed blocker, or material user correction. Record what happened, supporting evidence, its effect, and the actual next action.
- Choose one entry per meaningful milestone. Combine closely related work; keep independent decisions, state changes, and checks distinguishable. Do not replace the log with a coarse checklist or narrate every command.
- Update the plan's `Artifacts` when a source file, test, branch, commit, pull request, document, or other relevant output is created or changes. Record its verified type, state, contents, location, and exact revision when applicable. Label unknown or merely planned details explicitly.
- Record each material decision and its reason under `Decisions`. Preserve the original history and label a changed decision `Superseded` instead of silently rewriting it.
- Update `Validation` with checks that passed, failed, were blocked, or have not yet run. Record the actual result. Say what each result proves and what still needs to be checked.
- Add useful decisions and learnings as soon as they materially affect the plan or remaining work.
- State which work is included and which is not. Apply the user's direction to the affected parts without silently expanding the accepted task.
- Record a blocker only for the work it actually prevents. Keep independent planned work visible and note its actual progress.
- Record real evidence. Do not present an intended action or partial result as a verified outcome.
- Update an existing `EXEC_STATE.md` when an epic produces or materially changes an external artifact. Record only verified outputs; do not invent links, imply publication, or start external monitoring.
- Keep the plan's detailed work log, artifacts, decisions, and validation in the numbered plan. `EXEC_STATE.md` summarizes the whole epic and its produced external artifacts; it does not replace the plan's own records.
- Mark a plan complete only when its recorded completion criteria are met.
- Respect an instruction not to create or update the document.

The plan owns its document hierarchy and tracks execution. It does not prescribe the target project's source layout, delivery process, delegation strategy, runtime tool, review procedure, release workflow, or another skill.

## Read only the relevant reference

- [Plan format](references/plan-format.md): Document shape and example fields.
- [Plan placement](references/placement-and-epics.md): Choosing `<plan-root>`, creating the required group and execution-state files, preserving existing plans, and recording verified external artifacts.
- [Goals and workstreams](references/goals-and-workstreams.md): Representing multiple accepted goals, steps, and dependencies in the plan.
- [Evidence and blockers](references/evidence-blockers-and-completion.md): Recording progress, actual blockers, evidence, and completion.
- [Status and steering](references/runtime-and-steering.md): Updating the document after meaningful progress or user direction.
- [Fictional example plan tree](references/example-plan-root/README.md): A concrete example of the required plan hierarchy; its project names, statuses, and workflow remain illustrative.

## Tests

When changing this skill, read [tests/README.md](tests/README.md). Run relevant scenarios with fresh, independent agents when they are available. Have them inspect the work without changing files or outside systems.
