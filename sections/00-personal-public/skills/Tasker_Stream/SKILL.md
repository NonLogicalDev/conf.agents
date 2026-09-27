---
name: Tasker_Stream
description: Adapt task ownership and execution to work of any size, from a simple direct fix to feature delivery spanning months. Use when work benefits from one accountable owner, useful delegation, evolving guidance, living project memory and owner artifacts, repeated iterations, or several workstreams. Invoke $Tasker_Stream migrate to migrate the current owner structure to the latest installed instructions.
---

# Tasker Stream

Tasker Stream is a self-adapting system for work of any size. Keep one owner accountable, apply relevant project memory and existing owner artifacts, delegate when useful, preserve enough context to resume, and check the result. The owner keeps useful project knowledge and artifacts current as part of the work; the user does not manage them. A simple task may need only direct action; a large feature may span months, many iterations, and several workstreams. Evolve the system in the spirit of the user's intent, and let it grow only when the work needs it.

## Scale the core to the work

- Handle a simple task directly when one owner can finish it without helpers, planning files, or durable owner support.
- For work that spans many iterations, workstreams, or months, add useful owner guidance, plans, delegation, memory, artifacts, project tools, and handoffs as the work actually needs them.
- Reuse the project, owner, and plan that already cover the work. The current task can be the owner.
- Keep the owner responsible for user direction, accepted scope, useful project memory and owner artifacts, worker assignments, shared plan state, integration, and completion.
- When the work needs durable records, choose `<plan-root>` from the user, existing project guidance, or an established project convention.
- When an owner has a plan root, keep `<plan-root>/MEMORY.md` as its required read-first project map. It is the first entry point into the owner's memory system: a concise index of the project, active workstreams, important files and interfaces, owned checkouts, durable decisions, artifacts, risks, and the detailed records that support them.
- Record owners only in `<plan-root>/OWNERS.md`. README, MEMORY, AGENTS, plans, and supporting records link to that file instead of repeating owner identities, owner lists, or ownership frontmatter.
- Before resuming or changing owner state, check the verified current thread against `OWNERS.md` as described below. No other document grants owner authority. Preserve unrelated metadata and useful content when removing legacy ownership declarations during migration.
- For work that needs a plan, use a brief root `README.md`, a group `EXEC_STATE.md`, and a numbered plan in the status folder that matches its current state.
- Let each record serve its own purpose: `MEMORY.md` maps the whole project and routes memory, the README introduces the project and links top-level records, execution state summarizes the group, and the plan holds useful work details.
- When a separate goal reference helps, keep `<plan-root>/GOALS.md` as a short mission index suitable for `/goal`. State the intended outcome, current priorities, and workstreams; link each workstream to its actual execution state or plan.
- Preserve existing goal headings and concise `Outcome`, `Priorities`, and `Workstreams` entries. Keep changing status, CI results, commit hashes, blockers, work logs, and checklists in `EXEC_STATE.md` or the relevant plan.
- Follow the workstream naming convention in the user's instructions, applicable `AGENTS.md`, or project. Otherwise, organize related work under broad `theme-<nnn>-<slug>` groups, reserving `theme-000-meta` for task zero and beginning delivery work with `theme-001-<slug>`.
- Track each distinct substantial goal in its own `plan-<nn> <description>.md` record within its existing theme's appropriate status folder. Have delegated workers share the plan that already covers their assignment. Add another theme when the work represents a distinct, lasting area of responsibility; use another group type when the user or project selects it.
- Choose every task, theme, epic, or other group as a meaningful work domain with a coherent responsibility, lifecycle, body of decisions, or set of interfaces. Do not mint a new group merely because a new ask, effort, bug, worker, pull request, or plan appears.
- Before creating a group, compare the new work with existing domains and the full set of observed asks and requirements. Reuse the closest meaningful domain and add or revise numbered plans when the work belongs to the same responsibility.
- When later evidence shows an existing domain was drawn too narrowly, broaden its stated domain and update its execution state, plans, memory index, goals, and owner instructions as needed. Preserve its identity and useful history instead of creating adjacent fragments. This organizational broadening does not expand the user's accepted work scope or permissions.
- Update a record when an important result, decision, blocker, ownership change, or next action changes. Do not create records for their own sake.
- Trust current instructions, actual source, and observed results over an old summary, memory file, or dashboard.

## Record which threads may act as owners

Whenever a durable plan root exists, keep its owner-thread register in `<plan-root>/OWNERS.md`, separate from instructions, plans, and the worktree inventory. It is the sole record of owner identities and authorization; link it from other documents without copying its entries. This file lists permission to act as an owner of this plan root; it is not a roster of every worker or every owner elsewhere.

```markdown
# Owner threads

## Owners

| Thread | Status | Authorization |
| --- | --- | --- |
| codex://threads/<verified-thread-uuid> | active | <verified assignment or approval source> |
```

Keep this table machine-readable: use raw thread URLs, `active` or `revoked` status, and a nonempty authorization source. Keep the `## Owners` section limited to this table; put broader explanations before it or under another heading; do not insert `|` characters in cells.

- Use verified thread IDs, not titles or inferred identity. Require at least one active entry. Active entries may act as owners within the accepted scope; coordinate one writer for shared state. Do not keep a separate primary-owner field elsewhere.
- Add or revoke an entry only from the operator's authorization or an established policy that explicitly permits that ownership assignment. Record the basis in the Authorization column; retain revoked entries with `revoked` status and the revocation basis. A worker assignment, matching title, fork, read access, or familiar checkout does not grant owner authority. Never add yourself merely to pass a failed check.
- Check membership on resume, before changing owner state, and after ownership steering. An absent or revoked thread must stop owner-level changes and resolve the assignment. It may still perform a valid delegated task within its worker contract, without changing ownership or taking over shared state.
- If an existing owner lacks OWNERS.md, stop normal owner mutations and use explicit migration. Treat legacy declarations such as README `owner_thread` as migration evidence only: verify the assignment and authority, carry authorized owners into OWNERS.md, then remove those declarations from other files while preserving unrelated metadata and content. A missing register is not permission to enroll yourself. For a genuinely new owner, initialize OWNERS.md from the verified current thread and actual assignment. Do not migrate unrelated owner directories.
- A malformed register, invalid or duplicate identity, or no active owners blocks owner changes until corrected from verified authority. Preserve existing entries and check concurrent edits before writing. Record ownership changes and authorization in OWNERS.md; maintenance plans link to that record rather than duplicating owner identities; keep the register itself at the owner root, not `theme-000-meta/`.

## Migrate an existing owner to the current format

When the operator invokes `$Tasker_Stream migrate`, migrate the current owner's structure using the latest available installed skill instructions. This is an agent workflow, not a `plan.py migrate` command. The invocation authorizes the needed structural migration of the identified owner, not unrelated owner directories or changes to the project's delivery scope.

1. Resolve the current installed `Tasker_Stream` source from the runtime skill mapping or managed skill link and re-read it and its relevant references. An old generation linked in conversation is context, not the target format. Do not assume a remembered version is current or fetch arbitrary upstream instructions. If the current installation cannot be verified, report that blocker.
2. Identify the current plan root and verified calling thread. Read its MEMORY, README, OWNERS, AGENTS, active execution state and plans, and affected support records. Apply the current ownership checks, including explicit verified migration when OWNERS is absent. Resolve missing or conflicting ownership from verified authorization before changing files; do not bootstrap yourself into another owner's records.
3. Inventory the differences against the current skill and applicable project/user exceptions. Cover ownership, meaningful work domains, group and plan layout, statuses, task zero, memory/indexes, guidance, recovery inventory, and links as applicable. Reuse equivalent existing records and preserve established numbering and explicit conventions where the current instructions permit them. Do not split a theme for each effort or relocate supporting files into theme-000-meta. If the comparison and read-only checks confirm that no migration is needed, report that the structure is already current and stop without changing files or creating checkpoints.
4. Record the intended changes and old-to-new path mapping in the existing meta-maintenance plan before editing. If it is missing, first add the smallest task-zero plan needed to record the migration. Coordinate with active writers, checkpoint the affected records using the project's permitted workflow, and preserve unrelated changes. Where no checkpoint mechanism exists, save recoverable copies of affected files before destructive moves; do not overwrite useful originals.
5. Apply only the necessary changes. Preserve task IDs, history, decisions, active assignments, authorized ownership, and worktree claims. Consolidate owner declarations into OWNERS.md; remove README `owner_thread` and duplicate owner lists or fields elsewhere, replacing them with links. Preserve unrelated frontmatter, project content, and historical work evidence. Add required records from verified evidence, carrying useful older guidance into current locations. Update internal links and entry points when paths change. Do not manufacture missing facts, mark incomplete work done, move live checkouts, delete branches, publish, or alter external systems as part of a structural migration. Remove old directories only when empty.
6. Verify the result against the latest requirements and permitted exceptions. Run `plan.py doctor` with the verified current thread, check changed links and path mappings, confirm active plans and recovery records remain usable, and re-read README, MEMORY, OWNERS, and AGENTS. The helper validates only part of the format; its success alone is not migration completion. Re-running migration on a current structure should require no changes.
7. Update the maintenance plan with the instruction version or resolved source, changes, verification, remaining gaps, and recovery information. Refresh MEMORY to match the resulting structure and checkpoint the completed migration. Report what changed, what was preserved, and any unresolved blocker. If interrupted or partially blocked, leave a precise resumable state rather than claiming completion.

For a direct task with no owner structure, report that there is nothing to migrate; do not create an owner merely for this command. A request to describe or preview migration is read-only. Creating this migration instruction does not itself authorize running it on existing owners.

## Keep task zero for meta maintenance

- Whenever an owner has a plan root, create or reuse a standing **task 0: Meta maintenance**. Under the default layout, use `theme-000-meta/EXEC_STATE.md` and `active/plan-01 Meta maintenance.md` inside that group. Link it from the root README and `MEMORY.md`.
- Put all plans for owner-level meta work in `theme-000-meta/`, using the same execution-state and status-folder conventions as other themes. This groups the plans, not their supporting files: keep memory, recovery state, helpers, handoffs, dashboards, and other owner material in their established locations and link to them from the plans.
- Use task zero to track the directory structure, meaningful work domains, file placement and naming, indexes and links, memory audits, stale or duplicate records, and readiness for a clean agent to resume. Record intended organizational changes before making them, then record what changed and how it was checked. Keep delivery work and its detailed logs in their own plans.
- Keep one standing task zero for the owner's lifetime. Revisit it when starting or resuming owner work and when the directory needs maintenance; finish individual upkeep items without closing the standing task after each audit. Close it when the owner is retired and the final organization and handoff are verified.
- Preserve existing conventions, numbering, and history. If an equivalent maintenance task or an occupied zero slot already exists, reuse the appropriate existing record and identify it as the owner's task-zero responsibility; do not renumber delivery tasks or overwrite an unrelated zero slot. This requirement does not create an owner directory for a simple direct task or authorize deleting records or changing another owner's files.

## Put project memory and owner artifacts to work

- For any owner with a plan root, read `<plan-root>/MEMORY.md` before substantive investigation, planning, delegation, implementation, or resumption. Use it to find the relevant instructions, plans, artifacts, worktrees, evidence, and deeper memory before acting.
- Before choosing a consequential action, find the existing owner knowledge that could change it. Apply relevant verified instructions, decisions, artifacts, earlier findings, and available tools to the next action instead of merely noting that they exist. Within the current owner's scope, follow its applicable instructions over generic remembered shortcuts.
- Reuse suitable owner records and artifacts instead of starting over. Read only what the current action needs, confirm the relevant information was actually returned, and let current instructions, source, and observed evidence resolve stale or conflicting information.
- Maintain the current owner's useful memory and artifacts as part of normal ownership. Let actual project needs determine what to preserve, update, or create; do not wait for the user to curate records, repeat known context, or request routine upkeep.
- Treat `MEMORY.md` as a curated project index, not a chronological log, status dashboard, or duplicate of every plan. Keep enough verified context and links that a clean agent can orient itself without searching blindly.
- Re-dream it from current evidence instead of appending forever: compare it with current user direction, `README.md`, `AGENTS.md`, `GOALS.md`, execution state, active plans, worktree inventory, durable memory entries, source, and observed results. Remove stale claims, repair broken links, tighten the map, and add newly important concepts or entry points.
- Audit and refresh `MEMORY.md` when owner work starts or resumes, after a material new workstream, decision, artifact, interface, checkout, blocker, or handoff changes the project map, before pausing or completing substantial work, and at least every 30 minutes while actively working in a long-running owner.
- Record when it was last re-dreamed, which evidence was checked, and the next known audit trigger. State uncertainty and stale areas instead of presenting old memory as current fact.
- Keep changing status, detailed work logs, command output, and validation in the owning execution state or plan. Keep durable findings in `_owner/memory/`; make `MEMORY.md` point to them and explain why they matter.
- If `MEMORY.md` is missing or too stale to orient the work, follow existing owner instructions and established indexes such as `_owner/memory/README.md`. After verifying that the current thread owns the project, create or repair its required root index before substantive work continues. Consult another owner's relevant records only as authorized, read-only references. The simple-task exception still applies when no plan root exists.

## Find a plan before starting substantial work

- Before any owner or worker thread starts a substantial piece of work, find the active plan that already tracks it or create the appropriate active plan.
- Treat a multi-step investigation, independent workstream, delegated assignment, or implementation that needs repeated decisions or checks as substantial work.
- Read only the instructions, existing project records, and current state needed to identify the right owner and plan. Do not start substantive investigation, delegate the work, change project files, or produce project artifacts before the plan exists.
- Record the user's goal, accepted scope, link to OWNERS.md, current understanding, and next action before execution begins. Reuse the existing owner's plan instead of creating a competing plan for an assigned worker.
- Keep the simple-task exception: a small direct action that does not need coordination or durable tracking does not require a plan.

## Update the plan before changing the project

- Treat every active plan as living documentation of the user's intent, accepted scope, direction, current understanding, decisions, and intended changes.
- Before changing project documents, artifacts, or source code, update the relevant active plan so it accurately reflects the intended work and current project state.
- Revise the plan first whenever new steering, evidence, or understanding changes the work. Keep the relevant execution state and plan current as meaningful results, blockers, decisions, or next actions change.
- If a plan would need an almost complete rewrite, mark it `abandoned`, preserve its history in `archived/`, and create the next numbered active plan before changing the project. Keep ordinary revisions in the existing plan.
- An updated plan does not grant new permission or expand the accepted scope. A simple task that does not need a plan does not gain one merely to satisfy this rule.

## Keep enough state to rebuild the work

- When substantial work already has an owner, make its plan root sufficient for a fresh agent to identify every active workstream, reconstruct every owned checkout, locate its last recoverable checkpoint, and resume without the prior conversation or worker memory.
- Keep the complete owned-checkout inventory in `_owner/state/worktrees.md` when the owner has worktrees. Link each checkout to its active plan and record its machine, repository and source root, full path, link to OWNERS.md, Bean or manual layout, branch, full `HEAD`, upstream and ahead/behind state, verified base and merge base, sparse focus, and reconstruction commands.
- Resolve each checkout's Git common directory and record whether its refs and objects survive machine or pod replacement. Verify the storage lifetime rather than inferring it from directory names; a recorded path or commit hash is not recoverable when its object database disappears with that filesystem.
- Use `$VCS_Checkpoint` after meaningful task-owned changes, before risky or disruptive work, and before pausing, handing off, releasing, or changing a worktree. Record each checkpoint's full commit, time, workstream, storage location, and verified durability. Preserve unrelated staged or unstaged work and honor a user instruction not to commit.
- If no permitted durable checkpoint exists, describe the uncommitted work and recovery gap explicitly. Do not push, rewrite history, publish a bundle, or preserve secrets without separate authorization.
- Keep each numbered plan's work log detailed enough to resume: record dated decisions, changed files, verified commands and results, checkpoint changes, blockers, worker assignments, and the next action. Keep summaries in `EXEC_STATE.md` and link to the inventory instead of copying mutable state into every owner document.
- On a clean resume, read the owner's `MEMORY.md` first, then its README, `AGENTS.md`, execution state, active plans, and worktree inventory; compare their recorded branches, commits, sparse focus, and storage with actual Git state before continuing or proposing an authorized reconstruction.
- Keep the simple-task exception. Do not create an owner directory, worktree inventory, or empty checkpoint merely to satisfy this section.

## Delegate when it helps

- Delegate independent substantial work promptly when available workers can help. Give each worker distinct files or subjects, relevant owner knowledge and artifacts, expected results, and the checks that matter.
- For planned work, update the active plan before asking a worker to change project documents, artifacts, or code.
- Mark every task currently being worked `active` and identify its worker by the verified full `/root/...` path. Write `worker identity unknown` when a worker exists but its path cannot be verified.
- Mark every running runtime step `in_progress` when the UI supports concurrent active steps. Keep all assignments visible, leave work that has not started pending, and mark completion only after checking the result.
- Keep one owner responsible for shared plan files. Give a shared checkout prerequisite or Git mutation one responsible worker, serialize actual conflicting changes, and assign a shared file to a worker only when that worker is its only writer.
- Reuse workers for related follow-up. Keep closely connected or trivial work with the owner when that is simpler.
- Preserve established workstream names and update only the workers affected by new direction.
- Check meaningful worker results and the combined outcome before claiming completion.

## Keep owner memory, artifacts, and tools useful

- Keep verified observations, project facts, and practical learnings from tasks and workstreams in `_owner/memory/` when they can help later work.
- Use root `MEMORY.md` as the first entry point, then read the relevant owner memory before investigating a question or assigning related work.
- Always correct relevant owner memory when verified evidence contradicts an existing entry. Save new findings that were difficult to establish and will help avoid repeating the same investigation.
- Update matching entries instead of creating competing accounts, and include the supporting evidence needed to reuse or verify each finding.
- Reuse and maintain existing owner documents, artifacts, handoffs, and project helpers when they improve the next decision or make work easier to continue. Link useful material from `MEMORY.md`; do not create duplicate records or support files without a clear purpose.
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

Use `doctor` to inspect an existing plan home. The naming commands print the next available group or plan without changing files. A new owner defaults to `theme-001-<slug>`; established group types and number widths remain supported, and `--digits` selects a different width when project instructions require it. Use the helper only when it matches the project's actual convention, and run `create-group` only when the user or project instructions allow the requested plan files to be created.

`doctor` checks OWNERS.md without changing files and reports missing or malformed registers and callers without active membership. It also warns about legacy README owner_thread declarations that migration must remove. README metadata never grants or denies membership. `create-group` initializes OWNERS.md only for a genuinely new owner root; an existing root without the register needs explicit migration. It does not grant additional ownership. The initial authorization records the verified bootstrap basis; add the actual assignment source in OWNERS.md when maintaining the register. Pass `--owner-thread <thread-id>` to `doctor` or `create-group` when `CODEX_THREAD_ID` is unavailable.

The helper allocates ordinary delivery groups from one onward; it does not create the reserved task zero. Create or reuse that standing maintenance record explicitly when establishing or resuming an owner with a plan root.

## Tests

When changing this skill, read [tests/README.md](tests/README.md). Run the relevant scenarios with fresh subagents and execute `python3 -B -m unittest discover -s tests -p test_plan.py -v` from the skill directory.
