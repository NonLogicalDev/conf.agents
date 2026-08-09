# Thread State Report Format

Use this reference to write an accurate report. Follow the user's requested
format first. Never invent results or times.

## Default Report

```markdown
# Thread state

_As of <verified report time>_

<One short paragraph explaining the goal, current state, and main need.>

- **Goals:** <goal and a result that shows the work is finished>
- **Now:** <current work or next action>
- **Risk:** <known risk, or `None identified`>

## Operator attention

1. **WS[03.01] — <decision>**
   - **Need:** <action only the operator can take>
   - **Blocked by:** <verified condition>
   - **Blocked for:** <verified duration, or `duration unknown`>

## Current workstreams

1. **WS[01.02] — <name>** — active; worker: `/root/<worker-name>`; next: <next action>.
   - **Outputs:**
     - `Last updated <verified timestamp>` — [<name>](<verified link>).
2. **WS[02.01] — <name>** — active; worker: `/root/<other-worker-name>`; next: <next action>.
   - **Outputs:**
     - `Sent <verified timestamp>` — [<name>](<verified link>).
3. **WS[03.01] — <name>** — blocked; next: <verified operator decision>.

## Recent work

- **WS[01.02] — <name>**
  - `<verified completion time>` — <completed result>.
```

Use `- None.` for an empty section. An ordinary next action does not belong
under operator attention. Include a real user decision there. If you cannot
confirm a requested decision or blocker, use `- None.` instead of guessing.
If no source establishes the goal, write `Goal not verified` instead of
inferring one from a workstream or output.

Show every workstream currently being worked as `active`, not just the first. Include each assigned subagent's verified full `/root/...` path. Keep unstarted work `pending`, preserve verified blockers, and mark a workstream `completed` only when its outcome is verified. If a worker is known but its path is unavailable, write `worker identity unknown`. Omit worker details when the owner is working alone.

## Requested Output Section

When the user requests all outputs in one section, use that section. Preserve
the workstream owner and the known time for every entry:

```markdown
## Outputs

### WS[01.02] — <workstream name>

- `Last updated <verified timestamp>` — [<name>](<verified link>).

### WS[02.01] — <workstream name>

- `Sent <verified timestamp>` — [<name>](<verified link>).
```

Do not repeat the same links in the section for current workstreams unless the
user asks. A separate output section is valid; dropping the owner or guessing
the time is not.

## Requested Short Report

When the user asks for only code changes, one workstream, or another named
part of the work, report that part. Say what the summary covers when leaving
something out might mislead. Do not say that other known work did not happen.

Keep each completed result separate from a merely open, drafted, updated, or
requested output. A sent message is complete as a message; it does not prove
that the requested action was approved or performed.

Do not put an updated pull request, note, or document in `Recent work`
without a verified merge, completion, or publication. A verified sent message
may appear there as a sent message only.

## Times And Links

- Use the actual report time.
- Use `Last updated <timestamp>` for a changed document, note, or change.
- Use `Sent <timestamp>` for a message that was actually sent.
- Use the verified completion time for a completed result.
- If a time is missing, say `Last updated unknown` or `Sent unknown`.
- Calculate a blocker duration only from its verified start time.
- Use a verified link when one exists. Do not invent a link.
- Keep the full existing `WS[xx.yy]` label with its result.
