---
name: Tasker_ThreadState
description: Report verified thread state in chat or a requested file, and keep active owner snapshots current in their configured location.
---

# Tasker Thread State

Use this skill to explain where a thread stands. Report the facts. Do not
change a goal, plan, owner, workstream, outside service, or unrelated file.

## Choose The Thread And Where To Send The Report

Report the current thread unless the user names another thread you can actually
read.

- Follow an explicit user request, including a requested file, another
  destination, or an instruction not to write.
- If applicable user or project instructions configure daily snapshots for
  the current owner thread, use that destination. Those instructions
  authorize the derived file; the user does not need to name it again.
- Otherwise, return the report in the current thread without writing a file.
- Read an existing requested or configured file before changing it. Replace
  the current snapshot unless the user asks to append or preserve content.
- Verify a saved report and return a short confirmation when a reply is
  needed. Post the full saved report only when the user asks to see it.
- Use a writing block only when the app supports it and the user
  has not requested a different presentation.

If the requested thread or destination cannot be accessed, say so. Do not
silently substitute another thread or file.

## Keep Active Owner Snapshots Current

When applicable instructions require owner snapshots, keeping the current
owner's daily snapshot current is part of the work. Follow the configured
naming scheme. Derive the destination from the owner's verified start date,
thread ID, stable slug, and the report's local calendar date.

Create or update the current day's snapshot:

- when substantive owner work begins;
- after a meaningful change to a workstream, decision, blocker, output, or
  next action;
- at least every 30 minutes while the owner is actively working; and
- before pausing, handing off, or completing the work.

Update the current day's file in place. Refresh the verified report time
without inventing progress. Preserve earlier daily snapshots and use the
report format the user or applicable instructions require.

If the configured destination or owner metadata cannot be verified, report
what is missing. Do not guess, update another thread's file, write unrelated
notes, start a timer or background job, or send heartbeat messages.

## Check The Current Facts

Read only the requested thread, its existing plan, its current state, and the
results needed to check the report. Use recent, direct information when an
older summary or checkpoint disagrees.

Record:

- the goal and what would prove it complete;
- each current workstream, its owner, state, next action, and each assigned subagent's verified full `/root/...` path;
- blockers, who can resolve them, and any required user decision;
- each known output, its owning workstream, and verified link;
- the actual sent, updated, or completion time when known; and
- recent results that have actually been completed.

Mark every workstream currently being worked `active`, including workstreams running at the same time. Reserve `pending` for work that has not started, preserve verified blockers, and report `completed` only after checking the result.

Identify each assigned subagent by its verified canonical task path, such as `/root/import_worker` or `/root/import_worker/checks`. If a subagent is known but its path cannot be verified, write `worker identity unknown`; do not invent a path. Work done by the owner alone does not need a subagent entry.

A changed file is an output, not a completed result. Include it in recent
work only when evidence proves it was merged, finished, published, or sent.
A sent request proves only that the message was sent, not that it was
approved or performed.

Say when a goal, owner, required decision, or blocker is unknown. Do not
infer one from a workstream name, sent request, or document title.

Follow `$Agent_StructuredOutput` for workstream labels. Preserve the complete
existing `WS[xx.yy]` markers. Sort by current priority without changing those
markers or inventing a new workstream.

## Follow The Requested Format

Read [references/report-format.md](references/report-format.md) before
writing a report. Use its default format unless the user requests another.

If the user asks for one output section, put the outputs in that section and
show the owning `WS[xx.yy]` beside each output or group. Do not reject the
requested layout or lose the connection between work and its results.

If the user requests a shorter report or limits recent work to one type of
result, follow that request. Do not claim omitted results never happened or
that unfinished work was completed. Keep a real blocker visible whenever the
requested report concerns the blocked work.

Use exact verified names, links, and times in any format. Write `Last updated
unknown`, `Sent unknown`, or `duration unknown` when the matching evidence is
missing. Do not turn `today`, an old checkpoint, or a guessed start time into
a fact.

## Check Before Sending

- Answer the requested question in the requested format.
- List an item needing the user's attention only when the user actually
  needs to act.
- Use `- None.` when no needed user action is known.
- Give each output its actual owner and verified time.
- Mark every running workstream `active` and include each assigned subagent's verified full `/root/...` path.
- List completed results only when you can show the work finished.
- Keep existing workstream names and order work by current priority.
- Do not include secrets, hidden reasoning, raw logs, or invented evidence.
- Verify the requested or configured snapshot after saving it.

## Tests

Read [tests/README.md](tests/README.md). Run the relevant tests with fresh,
independent agents when available.
