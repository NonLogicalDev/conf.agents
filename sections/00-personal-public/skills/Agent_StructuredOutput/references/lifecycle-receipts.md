# Lifecycle Receipts

Read this reference before reporting work by an agent, a visible thread, or several actual events. A factual answer, conversational reply, steering acknowledgement, ordinary rewrite, or simple status answer does not acquire a receipt just because this reference exists.

## Keep Unchanged Monitoring Quiet

Do not issue a receipt for another poll, recovery attempt, waiting period, worker check, or update about elapsed time when the verified work state is unchanged. A deadline, repeated retries, a teammate's preference for heartbeats, or a prewritten update does not turn silence into a blocker.

Use the actual waiting mechanism or continue useful permitted work. Report only a verified state change, actual decision, changed blocker, completed result, or operator action. Answer a requested simple status directly without automatically adding a lifecycle line.

## Use One Receipt For Each Real Event

Start each actual work event with a Markdown horizontal rule, a blank line, and one lifecycle line:

```text
<emoji> :: `YYYY-MM-DD HH:MM TZ` :: <kind>/<action> :: `<specific event>`
```

Immediately follow it with an `- **Outcome:**` parent bullet. Nest actual evidence, current blockers, and the next action only when useful.

- `STEP` identifies work by the current agent.
- `AGENT` identifies a transient helper.
- `THREAD` identifies a continuing task the user can see.
- `open` reports a verified start or assignment.
- `update` reports actual progress, a current blocker, or a state change.
- `close` reports an evidenced finish or end; it does not establish runtime deletion, cleanup, or archival.
- `steer` is valid only for an instruction actually sent to an existing helper or visible thread. Do not call a steering send an `update`.
- Use the actual report time in the operator's local timezone, or verified UTC when that timezone is not known. Do not imply the report time proves when the underlying event occurred.
- Use the complete workstream marker defined in the main skill when the event belongs to a current workstream.
- Give every separate event its own rule, timestamp, lifecycle line, and outcome. Never merge unrelated events under one timestamp.

## Choose The Right Emoji

Use `🟢` for an opening event, `🔄` for an update, `✅` for an ended assignment, and `🧭` for direction actually sent. Use `❓` for an actual operator question, `⛔` for an evidenced blocker, and `🔬` for real research. Retain the valid `open`, `update`, or `close` action when using a question, blocker, or research symbol.

Keep only the timestamp and event description in inline code. Do not put the entire line in inline code or a fenced code block.

## Report A Transient Agent

Report a real assignment, progress, direction, or completion with `AGENT/open`, `AGENT/update`, `AGENT/steer`, or `AGENT/close`. Use the exact observed outcome. Include the worker's verified canonical `/root/...` path in the event description when it is available; never invent a worker identity.

```markdown
---

✅ :: `2026-07-25 03:12 PDT` :: AGENT/close :: `WS[03.02] /root/parser_validation — Parser validation`

- **Outcome:** All 42 focused parser tests passed.
  - **Next:** Verify the remaining integration checks.
```

A completed assignment does not prove that the runtime deleted the worker. Do not claim cleanup unless a supported close operation succeeded.

## Report A Visible Thread

Use a verified stable thread ID. End the lifecycle line with a clickable thread link:

```markdown
---

🧭 :: `2026-07-25 03:12 PDT` :: THREAD/steer :: `WS[03.02] Preserve the parser scope` :: ([thread](codex://threads/01900000-0000-7000-8000-000000000000))

- **Outcome:** Sent the scope correction to the existing owner.
  - **Next:** Check the owner's response when it is available.
```

Sending the correction proves delivery only when delivery was verified; it does not prove that the receiving task acted. Do not claim a thread started, finished, changed title, archived, or deleted without evidence.

## Separate A Final Receipt From Its Answer

When a final answer genuinely contains both a completed lifecycle receipt and a separate, self-contained summary, place exactly two horizontal rules between the receipt's outcome and the `## Summary` heading:

```markdown
- **Outcome:** The focused checks passed.

---

---

## Summary

All 42 focused parser tests passed. Integration checks remain.
```

Do not add double rules, an empty summary, or a lifecycle receipt to an ordinary answer, a simple completion status, or sibling sections.
