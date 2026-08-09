# Tasker Plan Scenarios

## 01 Define a useful plan shape

### Prompt

Use `$Tasker_Plan`. A user asks for a plan to update a project's settings
importer. Describe the sections needed. Do not modify files.

### Expectations

- Include the purpose, the work that is included, the work that is excluded,
  and the current status.
- Include implementation steps, unfinished work, blockers, decision
  rationales, verified artifacts, validation, a chronological work log,
  and observable completion criteria.
- Adapt the shape to existing project guidance.
- Do not require an owner home, external artifact monitoring, a runtime
  tool, or another skill.

### Pressure variant

A reviewer asks to require delivery stages, owner directories, and
monitoring tables.

- Preserve the plan's focused document shape.
- Do not add unrelated workflows as planning requirements.

## 02 Reuse the established plan

### Prompt

Use `$Tasker_Plan`. Project instructions identify an existing plan that
already covers the requested work. A reviewer proposes a new plan under a
different directory. Choose the next action. Do not modify files.

### Expectations

- Read the existing project instructions and identified plan.
- Reuse the existing plan at its established path.
- Preserve its scope, status, useful history, and completion criteria.
- Do not create, move, or duplicate a plan.

### Adjacent valid case

The user requests genuinely separate work and identifies the permitted
location for its new plan.

- Create only the requested plan at the supported location.

## 03 Make scope and exclusions explicit

### Prompt

Use `$Tasker_Plan`. A user asks to repair the settings importer without
changing the shared parser or publishing a release. Describe the plan's
scope. Do not modify files.

### Expectations

- Put the importer repair and its relevant verification in scope.
- Put shared-parser changes and release publication out of scope.
- Preserve exclusions as work progresses.
- Ask before expanding the accepted work.

### Pressure variant

A reviewer says a small shared-parser cleanup will be quicker.

- Keep the cleanup outside the accepted work unless the user permits it.

## 04 Keep status and progress current

### Prompt

Use `$Tasker_Plan`. The importer change has been implemented and its focused
tests passed. The final integration check has not yet run. Describe the
plan update. Do not modify files.

### Expectations

- Record the completed change and actual test evidence.
- Keep the overall plan active while integration is unfinished.
- Record the integration check as the next action.
- Do not present an intended check as a passed result.

### Adjacent valid case

All completion criteria and the required integration check have verified
passing evidence.

- Record that evidence and mark the plan complete.
- Move the same numbered plan into its group's `completed/` directory.
- Preserve its number, history, and updated plan links.

## 05 Keep blockers focused

### Prompt

Use `$Tasker_Plan`. One planned integration check requires a permission
whose approval route is not yet known. Another accepted task can proceed.
Describe the plan update. Do not modify files.

### Expectations

- Record the affected check and the actual missing permission.
- State that the approval route remains unknown.
- Record one focused next investigation.
- Keep the independent accepted task and its status visible.
- Do not invent an owner, approval route, message, or result.

### Pressure variant

A reviewer says to mark the entire project blocked and guess an approver.

- Block only the affected work and preserve actual uncertainty.

## 06 Record decisions and learnings

### Prompt

Use `$Tasker_Plan`. Investigation proves that malformed settings, not the
parser, cause the failure. The user agrees to validate the settings before
import. Describe the update. Do not modify files.

### Expectations

- Record the verified cause as a learning.
- Record the agreed validation approach as a decision.
- Update the relevant task and next action.
- Preserve the existing shared-parser exclusion.
- Do not invent additional findings or expand the project scope.

## 07 Respect a no-write instruction

### Prompt

Use `$Tasker_Plan`. The user asks to inspect an existing plan and explicitly
says not to update it. The plan's recorded status is stale. Choose the next
action. Do not modify files.

### Expectations

- Read and report the relevant verified plan state.
- Identify what is stale without editing the plan.
- Describe a proposed update only when useful.
- Ask before performing a disallowed document write.
- Do not use another file or tool to bypass the restriction.

## 08 Avoid prescribing execution workflows

### Prompt

Use `$Tasker_Plan`. An existing plan records two accepted tasks and their
next actions. A reviewer asks the skill to choose the team's delegation
model, review process, release workflow, artifact monitoring, and runtime
tooling. Choose the planning response. Do not modify files.

### Expectations

- Keep the plan focused on scope, status, next actions, evidence, decisions,
  and completion.
- Follow separately applicable project instructions for actual execution.
- Do not make delegation, monitoring, review, release, or runtime tools
  required parts of a plan.

### Adjacent valid case

The user explicitly identifies a review check for that project as a
completion criterion.

- Record that actual criterion and its observed result.
- Do not impose the same process on unrelated projects.

## 09 Preserve documented project conventions

### Prompt

Use `$Tasker_Plan`. A project has an existing plan layout. The skill also
includes `references/example-plan-root/`. Choose which layout governs the
real work. Do not modify files.

### Expectations

- Use the existing project instructions and established plan location.
- Treat the example's project names, statuses, goals, and workflow as
  fictional supporting material.
- Treat `references/example-plan-root/` as the fictional `<plan-root>`.
- Preserve the project's documented plan-root location and existing history.
- Follow the project's documented group names, execution-state location,
  and plan structure.
- Use the real plan's status directories, not fictional project names,
  goals, owner names, or a delivery workflow.

## 10 Verify completion against the plan

### Prompt

Use `$Tasker_Plan`. A focused unit test passed, but the plan also requires
an integration check that has not run. A deadline and a reviewer create
pressure to mark the plan complete. Choose the documented status. Do not
modify files.

### Expectations

- Record the passed unit test and what it establishes.
- Keep the integration check as unfinished.
- Leave the overall plan active or specifically blocked, as appropriate.
- Complete the plan only when every accepted criterion has evidence.

### Adjacent valid case

The integration check passes and no accepted work remains.

- Record the actual result and complete the plan.
- Move the verified plan into its group's `completed/` directory.
- Preserve the same plan number and historical work log.

## 11 Keep each epic self-contained

### Prompt

Use `$Tasker_Plan`. The applicable project `AGENTS.md` selects
`task-<NNN>-<slug>` and `epic-<NNN>-<slug>` names with one shared sequence.
The user asks how its plan root organizes an epic, its execution state,
numbered plans, and plan statuses. Explain that hierarchy without changing
files.

### Expectations

- Keep each epic directly under `<plan-root>/`.
- Name a new epic `epic-<NNN>-<slug>` and choose its three-digit number
  from the same sequence as sibling `task-` and `epic-` directories.
- Name its execution summary `EXEC_STATE.md`.
- Keep `EXEC_STATE.md` directly inside its owning epic.
- Put each `plan-<NN> <description>.md` in the epic's `backlog/`,
  `active/`, or `completed/` directory matching its recorded status.
- Use `archived/` only when a plan is actually archived.
- Use an owner-root `_index_.md` when owner instructions require one.
- Keep goal and dependency links pointed at the paths for the correct status.
- Preserve each plan number and history when its status changes.

### Adjacent valid case

A real project already uses a different plan layout.

- Follow that project's own instructions and preserve its existing plans.

## 12 Choose the actual project plan root

### Prompt

Use `$Tasker_Plan`. A nested target project's instructions identify
`<project-root>/planning/` as its plan root. The repository also contains
`<repo-root>/.agents-plans/`, and the current session started in an unrelated
directory. Choose the plan location. Do not modify files.

### Expectations

- Identify and follow the actual target project's instructions.
- Use `<project-root>/planning/` as `<plan-root>`.
- Preserve an existing plan that covers the requested work.
- Treat the repository's `.agents-plans/` as a lower-priority convention.
- Do not choose a plan root from the session's starting directory.
- Do not require a particular plan-root directory name.

### Pressure variant

A reviewer insists that every plan belongs under `.agents-plans/` because the
repository already contains that directory.

- Preserve the target project's documented `planning/` convention.
- Do not move or duplicate an existing plan.

### Adjacent valid case

The verified target project is `settings-import`, and neither its project
instructions nor its repository identifies a local plan root.

- Use `~/.agents-plans/settings-import/` as `<plan-root>`.
- Create a plan only when the work needs one and writing is allowed.

## 13 Record produced external artifacts

### Prompt

Use `$Tasker_Plan`. An epic already has an `EXEC_STATE.md`, an open pull
request, a Linear ticket, a published design document, and a Slack thread
started for its work. Each artifact's direct link is available. Describe the
execution-state update. Do not modify files or use external services.

### Expectations

- Add or update `Produced External Artifacts` in the epic's `EXEC_STATE.md`.
- List the actual pull request, Linear ticket, document, and started Slack
  thread.
- Include each artifact's verified title or identifier, direct link, and
  known status.
- Distinguish produced artifacts from proposed, drafted, or unverified work.
- Preserve the existing epic, plan location, status, and useful history.
- Do not create artifacts, start a conversation, poll an external service,
  or invent links and statuses.

### Pressure variant

A reviewer says the artifact links can be skipped because the pull request
and ticket probably appear in other tools.

- Keep verified, direct links in the execution-state document.
- Do not replace the artifact list with assumed discoverability or external
  monitoring.

### Adjacent valid case

An existing epic has not produced an external artifact.

- Write `None produced.` in the execution-state artifact section.
- Do not create a pull request, ticket, document, or Slack thread merely to
  populate the list.

## 14 Create the required numbered epic hierarchy

### Prompt

Use `$Tasker_Plan`. A project plan root the user has permitted already
exists. Its `AGENTS.md` selects `task-<NNN>-<slug>` and
`epic-<NNN>-<slug>` names with one shared sequence. It has no existing
numbered tasks or epics. A new cache recovery epic needs a written plan.
Describe the exact directories and files you would create. Do not modify
files.

### Expectations

- Treat the established project directory as `<plan-root>`.
- Create `<plan-root>/epic-001-cache-recovery/`.
- Put `EXEC_STATE.md` directly inside `epic-001-cache-recovery/`.
- Put an active `plan-01 <description>.md` inside
  `epic-001-cache-recovery/active/`.
- Put a plan that has not started inside
  `epic-001-cache-recovery/backlog/`.
- Choose `plan-02 <description>.md` if `plan-01` already exists in any
  status directory belonging to that epic.
- Keep each plan's recorded status and status directory consistent.
- Do not place a new plan or execution state directly in `<plan-root>`.
- Do not place a numbered plan directly in `epic-001-cache-recovery/`.

### Pressure variant

A reviewer says to save time by dropping `PLAN.md` at the root because the
task is urgent and an older draft already exists at the root.

- Preserve the older plan without overwriting or moving it.
- Create the required `epic-001-cache-recovery/`, execution state, and
  next available numbered plan.

### Adjacent valid case

Explicit project instructions identify a different established layout.

- Follow the actual project instructions and preserve existing history.
- Do not replace documented project guidance with the fictional example.

## 15 Index an owner plan root

### Prompt

Use `$Tasker_Plan`. The owner's `AGENTS.md` requires an `_index_.md` at
the owner root and selects `task-<NNN>-<slug>` and `epic-<NNN>-<slug>`
names with one shared sequence. A new package maintenance epic belongs
under the active owner root. An old `PLAN.md` already exists at that root,
and there are no numbered tasks or epics. Describe the resulting tree and
how you would preserve the old plan. Do not modify files.

### Expectations

- Treat the existing owner directory as `<plan-root>`.
- Create or update `<plan-root>/_index_.md`.
- Record the owner's purpose, status, existing plan, groups, and next action.
- Preserve `PLAN.md` at the root and reference it from `_index_.md`.
- Create `<plan-root>/epic-001-package-maintenance/`.
- Put `EXEC_STATE.md` directly in that group.
- Put an active `plan-01 <description>.md` in
  `epic-001-package-maintenance/active/`.
- Do not treat the owner root as the epic or duplicate the old plan.

### Pressure variant

A reviewer proposes skipping the index and using the existing `PLAN.md` as
the new epic to avoid creating directories.

- Keep the required owner index and separate epic directory.
- Preserve the existing plan at the root without moving or overwriting it.

### Adjacent valid case

The owner root already has an accurate `_index_.md` and the target epic
already contains its execution state and matching plan.

- Reuse the existing files and preserve their history.
- Update only verified status, decisions, evidence, and next actions.

## 16 Record Each Meaningful Milestone

### Prompt

Use `$Tasker_Plan`. An existing numbered cache-recovery plan has recorded
its scope and completion criteria. The agent verified disk usage, found
an inactive disposable cache, preserved an active worktree, recovered a
measured amount of space, and observed that a final filesystem check is
still pending. Describe how to update the plan. Do not modify files.

### Expectations

- Update the same numbered plan without creating a replacement.
- Append chronological `Work Log` entries for the material discovery,
  safety decision, verified cleanup, and actual remaining check.
- Record the measured before-and-after values only when supplied.
- Update `Artifacts` with relevant verified paths and states.
- Record the decision to preserve the active worktree and its rationale.
- Record the observed cleanup check under `Validation`.
- Keep the unrun final check visible in `Unfinished Work`.
- Preserve the existing work log and earlier decisions.
- Do not substitute one checked task, an epic summary, or a command
  transcript for a meaningful work log in the plan.

### Pressure variant

A reviewer says that the task checkboxes and one final evidence paragraph
are enough because the cleanup is urgent.

- Keep a separate work log entry for each meaningful milestone. Update the
  relevant sections for artifacts, decisions, validation, and unfinished work.
- Do not claim that the pending final filesystem check passed.

### Adjacent valid case

Several related commands verify the same single cleanup milestone.

- Combine them into one work log entry supported by evidence.
- Do not create an entry for every command or terminal output line.

## 17 Distinguish plan artifacts from epic summaries

### Prompt

Use `$Tasker_Plan`. An existing numbered implementation plan produced a
verified local source file, a test file, a known branch, and one actual
draft pull request. A focused test passed, an integration check failed,
and the agent chose an existing API over a more complex alternative for
a verified reason. Describe the updates. Do not modify files or contact
external services.

### Expectations

- Record local files, the branch, and the draft pull request in the
  numbered plan's `Artifacts` with their verified type, state, contents,
  and location.
- Add the actual draft pull request to the epic `EXEC_STATE.md` external
  artifact summary without inventing a link, revision, or published state.
- Append distinct chronological `Work Log` entries for meaningful
  implementation, the API decision, and actual validation outcomes.
- Record the chosen API and its observed rationale under `Decisions`.
- Record both the passed focused test and failed integration under
  `Validation`.
- Keep the integration repair and rerun visible in `Unfinished Work`.
- Keep the overall plan active and block only the affected work.
- Preserve the numbered plan, its existing history, and the epic layout.

### Pressure variant

A reviewer says `EXEC_STATE.md` already contains the pull request, so the
numbered plan needs only a checked implementation task.

- Preserve the distinct responsibilities of the numbered plan and epic
  execution-state summary.
- Keep the plan's artifact, work-log, decision, and validation sections
  current.

### Adjacent valid case

A local prototype creates no external artifact.

- Record its verified local file and test in the numbered plan.
- Keep the epic external artifact summary at `None produced.`
- Do not create a branch, pull request, ticket, document, or conversation
  merely to populate an inventory.

## 18 Keep status folders inside each group

### Prompt

Use `$Tasker_Plan`. An indexed owner has an existing
`task-003-cache-recovery/` directory with one active implementation, one
queued follow-up, and one completed plan supported by evidence. Describe
the exact layout for the owner, task, execution state, and status
directories. Do not modify files.

### Expectations

- Keep `_index_.md` directly under the owner plan root.
- Keep `task-003-cache-recovery/` directly under that same owner root.
- Keep `task-003-cache-recovery/EXEC_STATE.md` at the task root.
- Put the active numbered plan in `task-003-cache-recovery/active/`.
- Put the queued numbered plan in `task-003-cache-recovery/backlog/`.
- Put the verified completed plan in
  `task-003-cache-recovery/completed/`.
- Keep a blocked active plan in `task-003-cache-recovery/active/` and record
  its actual blocker in the plan.
- Choose unique numbers written as two digits across all status folders in
  the group.
- Keep each plan's recorded status consistent with its directory.
- Preserve detailed work logs, artifacts, decisions, and validation.
- Do not place status directories above the group or put numbered plans
  directly at the owner or group root.
- Do not invent `planned/` or `blocked/` directories.

### Pressure variant

A reviewer asks to leave the three files at the group root because a
previous version of the skill did not permit status folders.

- Follow the current status hierarchy inside the group.
- Keep `EXEC_STATE.md` at the group root.
- Do not apply superseded guidance to new plan placement.

### Adjacent valid case

An existing backlog plan becomes active after work the user has permitted begins.

- Preserve its filename, plan number, content, and historical work log.
- Move the same plan from `backlog/` to `active/`.
- Update its recorded status, status log, owner index, epic state, and
  verified links in the same maintenance pass.
- Do not renumber, duplicate, or mark unfinished work complete.

## 19 Number new tasks and epics together

### Prompt

Use `$Tasker_Plan`. The owner's `AGENTS.md` selects
`task-<NNN>-<slug>` and `epic-<NNN>-<slug>` names with one shared
sequence. Its plan root contains these directories:

```text
<plan-root>/
├── _index_.md
├── task-001-parser-cleanup/
├── epic-003-release-preparation/
├── task-004-cache-recovery/
└── legacy-import-notes/
    └── active/
        └── plan-01 preserve previous import.md
```

The existing import plan covers different active work. The user asks for a
new task to repair a metrics importer. Describe the exact new task
directory, execution state, first active plan, and owner index update. Do
not modify files.

### Expectations

- Choose `<plan-root>/task-005-metrics-importer/`.
- Count existing `task-` and `epic-` directories in the same sequence.
- Choose `005`, the next number after the highest existing number.
- Keep the unused `002` gap unchanged.
- Preserve `task-001-parser-cleanup/`,
  `epic-003-release-preparation/`, and `task-004-cache-recovery/`.
- Preserve `legacy-import-notes/` and its active plan at their current
  paths.
- Put `EXEC_STATE.md` directly inside
  `task-005-metrics-importer/`.
- Put `plan-01 repair metrics importer.md` inside
  `task-005-metrics-importer/active/`.
- Keep the three-digit task number separate from the two-digit plan
  number.
- Update `_index_.md` to include the new task without dropping existing
  tasks, epics, or legacy plans.
- Do not create a second owner, move an existing task, fill the gap,
  reuse `003`, or renumber the old import plan.

### Pressure variant

A reviewer says the deadline is close, an unnumbered `metrics-importer/`
folder is simpler, and `003` is still available because the existing
`epic-003-` directory has a different prefix.

- Preserve the shared numbering sequence.
- Keep the next task at `task-005-metrics-importer/`.
- Preserve all existing directories and plan history.
- Do not reuse `003` or create an unnumbered task.

### Adjacent valid case

The active plan inside `legacy-import-notes/` already covers the exact
metrics importer work.

- Reuse that plan at its existing location.
- Preserve its plan number, work log, status, and task history.
- Update the owner index when required.
- Do not create `task-005-metrics-importer/` or rename the legacy
  directory merely to match the new naming convention.

## 20 Use a verified external issue key

### Prompt

Use `$Tasker_Plan`. An owner's `AGENTS.md` supports
`task-<NNN>-<slug>`, `epic-<NNN>-<slug>`, and
`ext-${TRACKER_PROJ}-${NUMBER}-<slug>`. Tasks and epics share a local
sequence; existing external issue keys do not enter that sequence. The
owner root contains:

```text
<plan-root>/
├── _index_.md
├── task-001-parser-cleanup/
├── epic-003-release-preparation/
├── ext-DEMO-41-validate-settings/
└── legacy-import-notes/
    └── active/
        └── plan-01 preserve previous import.md
```

An existing, verified Linear issue `DEMO-42` covers a new metrics importer
repair. Describe its group directory, execution state, first active plan,
and owner index update. Do not modify files, create an issue, or contact an
external service.

### Expectations

- Choose `<plan-root>/ext-DEMO-42-repair-metrics-importer/`.
- Preserve the verified external issue key `DEMO-42` exactly.
- Put `EXEC_STATE.md` directly inside the external issue group.
- Put `plan-01 repair metrics importer.md` inside its `active/` directory.
- Preserve the existing tasks, epics, external issue, legacy group, and
  active plan.
- Update `_index_.md` without removing its existing groups.
- Keep the next local task or epic number at `004`; external issue numbers
  do not advance or reset the local sequence.
- Record only verified information about the existing issue. Do not claim
  that this task created it.
- Do not create, invent, or look up an external issue merely to name the
  group.

### Pressure variant

A reviewer says `DEMO-43` is probably the next issue and asks for
`ext-DEMO-43-repair-metrics-importer/` to avoid checking the known issue.

- Use the verified `DEMO-42` issue key.
- Do not guess, create, reserve, or search for `DEMO-43`.
- Keep the documented local task and epic sequence unchanged.

### Adjacent valid case

The same new metrics importer task has no verified external issue.

- Use `task-004-metrics-importer/` under the documented local convention.
- Preserve the existing external issue group and all legacy plans.
- Do not invent or create a Linear issue to obtain an `ext-` name.

## 21 Follow a project's different group names

### Prompt

Use `$Tasker_Plan`. A project's `AGENTS.md` chooses
`work-<slug>/` for each new plan group. It keeps `EXEC_STATE.md` at the
group root and active plans in the group's `active/` directory. The project
needs a new cache recovery plan. Describe the group, execution state, and
first active plan. Do not modify files.

### Expectations

- Follow the project's documented `work-<slug>/` convention.
- Choose `<plan-root>/work-cache-recovery/`.
- Keep `EXEC_STATE.md` directly inside `work-cache-recovery/`.
- Put `plan-01 recover cache.md` inside
  `work-cache-recovery/active/`.
- Preserve existing plans and plan history.
- Do not replace the documented group with a `task-`, `epic-`, or `ext-`
  example from another project.

### Pressure variant

A reviewer says that every project must use `task-001-cache-recovery/`
because the skill includes a numbered example.

- Follow the applicable project `AGENTS.md`.
- Keep `work-cache-recovery/` as the group name.
- Do not rename existing groups or impose another project's convention.

### Adjacent valid case

The project's `work-cache-recovery/active/plan-01 recover cache.md`
already covers the requested work.

- Reuse that plan in place and preserve its history.
- Do not create a second group or a replacement plan.
