# Owner Operating Guide

Use the same adaptable owner system for a simple direct task, a project with several workstreams, or feature delivery that continues for months. Begin with only what the current work needs. Add lasting guidance, memory, tools, delegation, and handoffs as verified needs emerge, while preserving the user's intent.

## Start with the current sources

When starting or resuming owner work, read `<plan-root>/AGENTS.md` when it exists, along with applicable project instructions and only the records that affect the current action. The root `README.md`, `GOALS.md`, group execution state, or numbered plan may help when they exist.

Re-read `AGENTS.md` after meaningful steering or instruction changes and before work that depends on its guidance. Have affected workers read the same instructions before starting or resuming their assignments.

Check the root README's `owner_thread` against the verified current thread before changing owner state. Report a mismatch and stop instead of taking over another owner's work.

Read only memory or supporting documents that help with the current task. After a handoff or lost context, reopen the sources needed to continue instead of trusting an old summary.

Verify a project, owner, checkout, revision, ticket, review, or outside service only when the next action actually depends on it.

## Keep the active plan ahead of project changes

When the work already has or needs an active plan, update that plan before changing project documentation, artifacts, or source code. Record the user's intent, accepted scope, current direction, understanding, decisions, intended changes, and actual state.

Keep the plan current as instructions, findings, or priorities change. It is living documentation of what the owner intends to do, not a report written after the change. A plan does not expand permissions, and a simple direct task does not need one merely to satisfy a template.

## Trust fresh evidence

Current instructions and directly observed state outrank saved memory, cached state, dashboards, and earlier reports.

Keep `GOALS.md` as a short mission index for `/goal`, with the outcome, current priorities, and links to existing workstream records. Put current progress, blockers, checks, and review details in the relevant execution state or plan.

Keep material decisions, work history, blockers, artifacts, and validation in the plan that owns the work. Use support files to reduce repeated investigation, not to create a second version of the plan.

When several workers share a changing file or artifact, make one worker responsible for it. Record its state only when that information will help another action.

## Keep project memory current

Before investigating a question or assigning related work, read the existing owner memory that applies to the current task or workstream.

Proactively save verified observations, project facts, practical learnings, and difficult discoveries in `_owner/memory/` when later work would otherwise repeat a meaningful investigation. Include the relevant task or workstream and supporting evidence when they help another agent reuse the finding.

Correct an existing memory entry whenever verified evidence contradicts it. Update a matching entry instead of creating competing accounts, and share relevant discoveries with workers who can use them.

Keep current status, blockers, and routine work logs in their plans. Add or update memory only when the information can spare a later investigation or help another workstream.

## Keep the owner's living harness useful

When sustained work needs standing owner guidance, keep it in `<plan-root>/AGENTS.md`. A single reference to that file should give another thread enough information to reproduce the owner's intended behavior without access to the original conversation. A simple direct task does not need an instruction file merely to satisfy a template.

Keep `AGENTS.md` current with applicable user instructions, the mission and `GOALS.md`, accepted scope and permissions, project conventions, workstream entry points, delegation, verification expectations, and durable operating decisions. Link existing workstream records rather than copying their status.

When an owner already has `OWNER_PROMPT.md`, preserve its useful guidance in `AGENTS.md`. Keep existing instructions that still apply, replace any that current user direction supersedes, follow the recorded ownership, and do not overwrite another owner's work.

The owner may update its own `AGENTS.md` without waiting for another request as accepted scope, user intent, priorities, workstreams, or project needs change. Improve it as the owner learns, while preserving the user's intent. Do not expand the accepted scope or change permissions. Keep changing progress, blockers, and routine work logs in execution state or plans.

Reading or sharing `AGENTS.md` does not transfer ownership. Check the root README's `owner_thread` before changing owner state, and preserve the current owner's work.

Use a root `CHANGELOG.md` only when a short history of owner changes will actually help.

## Evolve the owner safely

Proactively improve owner guidance, project helpers, or the writable source of this skill when verified project needs show that a small change would better serve the user's intent.

Start with the smallest useful change. Prefer updating the owner's `AGENTS.md`, project memory, or local helper when the problem belongs to one project. Update shared skill guidance only when the lesson applies to other projects and the source can be edited normally.

Keep the existing goal, accepted scope, user preferences, safety checks, and permissions intact. Record the reason or check when that information will matter later.

Record every deliberate departure from an existing owner convention or prior guidance in `_owner/wart.deviation.md`. Explain what changed, why, and how it serves the user's intent. A deviation does not override current instructions or permissions.

If a proposed change would require new permission, affect another owner, remove an important safeguard, publish work, reveal private material, or change a decision the user should make, stop and ask instead.

Edit the writable source, not an installed read-only copy. Do not bypass a filesystem restriction or treat a wart record as permission to do something the user has not allowed.

## Add support only when it helps

- Keep durable project facts in `_owner/memory/`.
- Build and reuse small helpers in `_owner/tools/` when they remove repeated work or common mistakes.
- Keep program data in `_owner/state/` only when a tool actually needs it.
- Add `_owner/dashboard/` when the work calls for an interactive dashboard, small site, or another project view meant to be published.
- Record recurring tool or environment problems in `_owner/wart.tools.md` and problems in skills or instructions in `_owner/wart.guidance.md`.
- Keep the reasons for deliberate departures in `_owner/wart.deviation.md`.
- Use another `wart.<type>.md` only when it describes a separate problem that will matter again.

Keep credentials and other private information out of dashboards or examples intended for other people.

A helper does not change what the user has allowed. Make its purpose clear, keep actions that change files explicit, and test the behavior that matters.

Choose the simplest ownership, coordination, data format, and validation that fits the project. Add or omit support areas as the user and the work require.

## Leave a useful handoff

Before a meaningful pause or handoff, update the records another owner will actually need. Include the current goal, worker assignments, verified outputs, important checks, blockers, or next action when they matter.

Recheck changing information when the handoff depends on it. State what is unknown instead of completing the story from memory.
