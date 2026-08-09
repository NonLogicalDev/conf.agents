# Tasker Plan Behavioral Tests

Run each scenario with a fresh, independent agent when the available tools permit.
Provide only the `$Tasker_Plan` skill and the scenario prompt. Do not provide
the expectations or intended answer.

Do not let scenarios change shared files. If a test needs to write, use a
temporary directory the user has permitted. Capture the response and compare
it with the expectations.

Verify that the skill produces a useful plan shape, preserves established
project conventions, tracks explicit scope, records honest status and next
actions, captures decisions and learnings, respects permission to change a plan,
and requires evidence before marking work complete.

Verify that every substantive numbered plan maintains `Artifacts`,
`Decisions`, `Implementation Steps`, `Work Log`, `Validation`, and
`Unfinished Work`. Check chronological entries for each meaningful milestone, verified
artifact state and exact locations, decision rationales, actual check
results, preserved previous entries, and synchronized blocker transitions.
An epic `EXEC_STATE.md` or a task checklist does not replace those records.

Verify that an epic's `EXEC_STATE.md` records only external artifacts actually
produced by that epic, with verified links and known status. Cover pull
requests, Linear tickets, documents, and Slack threads without requiring
artifact creation, external monitoring, or invented outputs.

Verify that `<plan-root>` follows user and project instructions first, the
target project's or repository's established convention second, and the
fallback under the home directory last. Do not require a particular local directory
name.

Verify that each new task, epic, or tracked issue has its own directory
directly under `<plan-root>`. Follow the user's instructions, the applicable
`AGENTS.md`, and the project's established naming convention. Supported
examples include `task-<NNN>-<slug>`, `epic-<NNN>-<slug>`, and
`ext-${TRACKER_PROJ}-${NUMBER}-<slug>`. Do not impose those names on a project
that documents a different convention.

When project instructions choose numbered tasks and epics, check both
kinds together. Choose the next three-digit number after the highest one
already in use. Do not fill a gap, reuse a number, or count tasks and epics
separately. When they choose an external issue format, use an existing,
verified issue key. Do not invent or create an issue, or use its number as
part of the local task and epic sequence.

Keep `EXEC_STATE.md` at the root of each group. Keep two-digit
`plan-<NN> <description>.md` files inside that group's `active/`,
`backlog/`, or `completed/` directory. Verify plan numbering across all
status directories in the same group. Check status changes, preserved
history, and updated links. When owner or project instructions require a
root `_index_.md`, verify that it records the owner and the relevant groups.

Preserve and reuse an existing plan that already covers the work, even when
its directory follows an older naming convention. Do not move, rename,
duplicate, or delete an existing plan or task just to make its name match
the new convention.

Do not turn artifact monitoring, delegation, runtime tools, external
communication, publication, review, or release into planning requirements.

Use [scenarios.md](scenarios.md) for the prompts, pressure variants, and
adjacent valid cases.
