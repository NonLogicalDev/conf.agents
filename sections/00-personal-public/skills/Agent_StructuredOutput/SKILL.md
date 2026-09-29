---
name: Agent_StructuredOutput
description: Formats for clear Codex task messages, progress updates, and status reports. Use the formats that help explain the work. Follow this skill when the user, AGENTS.md, or other applicable instructions request it, for the scope they specify.
---

# Agent Structured Output

## Choose the format for the work

Choose the formats that make the message easier to follow. Plain prose suits ordinary replies; structured reports can help explain complex work. Using one format by choice does not commit later messages to the same format.

When the user, AGENTS.md, or other applicable instructions explicitly request this skill or one of its formats, follow that request for its stated scope and duration. A standing instruction remains in effect until changed; do not treat it as a suggestion or limit it to one reply. A request for one report or a reference to the numbering convention applies to that part, not every format in this skill. The sections below describe how to apply each format when chosen or required.

## Outcome

Give the operator a clear, honest answer of the right length. State the result, evidence, actual blocker, and next action when they matter. Use this skill to structure replies, label workstreams, acknowledge steering, and display runtime steps.

Use the appropriate skill for public messages, saved plans, delegation, mission decisions, and outside actions. Formatting a reply does not give you permission to do any of those things.

## Write For The Operator

- Use plain English, short sentences, concrete verbs, and exact technical names. Lead with the answer or outcome.
- Write each Markdown paragraph or list item on one source line, however wide it needs to be. Let the app wrap it instead of inserting newlines at 80 characters or any other preferred width.
- Preserve line breaks that define paragraphs, headings, nested lists, blockquotes, tables, fenced code, or intentional Markdown breaks.
- Match the operator's requested length. Answer a simple question directly. Do not force headers, bullets, status labels, timestamps, lifecycle receipts, or a steering template onto an ordinary answer.
- Remove filler, repetition, reassurance, corporate language, invented process terms, unsupported praise, and claims that evidence does not prove.
- Use a heading or nested bullets only when several substantive facts or independent workstreams need them. Keep evidence, blockers, and next actions under the outcome they support.
- Keep confidential material and private reasoning out of responses. Distinguish a verified fact, a current assumption, and an unfinished action.
- Read [references/operator-writing.md](references/operator-writing.md) only when the answer requires complex sections, annotated source lists, reusable writing blocks, or a separate final receipt and summary.

## Identify Subjects And Keep Labels Private

- Follow the applicable **User communication references** rules in every reply, including plain replies and short updates. At first mention in each message, identify the subject beside its status, result, or question and link it when possible. Include a tracked item's identifier and title or short description. Earlier messages, private shorthand, and later source lists do not replace this introduction.
- Shorthands help the user refer back to items; they are not abbreviations for the agent. Follow every shorthand immediately with the full name or description of its referent. For tracked work, include its canonical identifier and title or description. Keep this pairing on repeat mentions, including in headings, table cells, links, and handoffs. Earlier context or a separate legend never permits a shorthand alone.
- Reuse established labels for the same items across messages. Preserve their mapping in authorized task records and handoffs; do not renumber or recycle labels.
- Keep private artifact and workstream shorthand out of public or team-facing material unless the user explicitly asks to preserve it in that public communication. Use descriptive names, canonical identifiers, and direct links in Slack, Linear, PR prose, and shared documents or sites.
- Classify copyable drafts by intended destination, even inside private chat. Labels supplied in prompts or source notes are context, not permission to publish them; translate or remove them unless the user explicitly requests their inclusion.
- A single item, documentation link, or source still receives a label when the applicable reference rules call for one. Ordinary connective prose does not need labels.

## Use Canonical Workstream Markers

This skill defines the workstream markers. Use the complete `WS[xx.yy]` marker for a current workstream step in private task prose, steering, snapshots, private durable references, lifecycle descriptions, and the Steps UI. Follow every marker with the full workstream or step name. Apply the public communication restriction above when writing for a public or team-facing destination.

- `xx` is the stable workstream lane. Write it as two digits and add a leading zero when necessary.
- `yy` is the current step in that lane. Write it as two digits and add a leading zero when necessary.
- Assign an independent workstream an available lane. Preserve that lane across reprioritization, display changes, owner handoffs, and step changes.
- Advance `yy` as that same lane progresses. Insert a newly discovered prerequisite before pending work when needed. Never renumber another lane or rewrite an already completed step to make the display convenient.
- A workstream's durable owner, plan, private coordination references, and other existing destinations remain associated with its stable lane. Advancing `WS[03.01]` to `WS[03.02]` changes the current step; it does not create another workstream, owner, plan, conversation, or external artifact.
- Use a current marker that you can verify. Do not invent a step or shorten a current workstream marker to `WS03`, `WSNN`, or `WS03:`.
- Other skills consume these markers; they do not define another numbering convention.

## Keep Runtime Steps Current

When substantive work has active workstreams and a Steps UI is available, use the canonical marker as `WS[xx.yy] <step description>`.

- Set up or update the actual Steps UI before substantive work and when a workstream, current step, dependency, priority, or status changes.
- Mark every task that is actually being worked `in_progress` and include `[active]` and its worker's verified canonical path in the visible label. For example: `WS[02.01] [active] /root/operator_guide — Revise operator guide`. Use `/root` when the owner is doing the work. If a worker's path is unknown, say so instead of inventing one.
- Keep all running tasks `in_progress` when the runtime supports multiple active items. Only if the actual tool demonstrably allows one active item, keep one task formally `in_progress` and mark every other running task `[active]` in its visible label. Leave work that has not started as plain pending.
- Update active labels and worker paths when an assignment begins, changes, finishes, or stops. Keep independent lanes distinct and steps within each lane in their actual execution order.
- Mark a step complete only after evidence proves its result. Report a missing tool or failed update. Do not claim an update succeeded when it did not.
- Use `Tasker_Plan` for the saved plan, recorded steering, limits on what you can change, and decisions about completion.
- Do not create a decorative steps list for a factual question, ordinary answer, or task with no active workstream. Do not update runtime steps or a plan when an instruction says not to.

Examples:

```text
WS[01.01] [active] /root/parser_worker — Inspect parser failure
WS[01.02] Repair parser failure
WS[02.01] [active] /root/fixture_worker — Audit fixtures
```

## Acknowledge Actual Steering

When using the steering format, explain what a message that starts, changes, pauses, replaces, or constrains work means before acting:

```markdown
**Understanding:** <the exact instruction and its concrete meaning>
- **Direction:** <what changes and what remains unchanged>
  - <material evidence, scope, permissions, dependency, or workstream effect>
- **Next:** <the immediate action>
```

When this format is selected, keep all three labels. A request for a short acknowledgement alone does not select this format. Use existing workstream markers when relevant. Preserve existing permissions, evidence limits, and unaffected work unless the operator explicitly changes them. Distinguish additive work from an actual stop or replacement.

Do not claim that an owner exists, a change was made, a message was sent, or a check passed without evidence. Do not add a lifecycle receipt to an acknowledgement that only states understanding. An ordinary factual question, requested rewrite, or status answer is not steering and does not require this template.

Use the available runtime delegation tools to route work and change workers. Use `Tasker_Plan` for mission decisions and changes to a saved plan that the user has permitted.

## Keep Unchanged Work Quiet

Monitoring, polling, recovery, retrying, waiting, and elapsed time do not create a material update by themselves. Do not send repeated messages, heartbeats, countdowns, or lifecycle receipts when the verified result, decision, blocker, and next action have not changed.

Continue work that the user has permitted or use the available waiting tool. Report only verified progress, a decision, a new or changed blocker, a cleared blocker, completion, or an exact action the operator needs to take. Do not claim that a retry, proposed recovery, or unfinished check has succeeded.

When the operator asks for status, answer the actual question directly. Do not add a lifecycle badge, timestamp, or heading unless the requested answer itself reports a real material work event.

## Report Real Work And Lifecycles

When using lifecycle receipts, read [references/lifecycle-receipts.md](references/lifecycle-receipts.md). For that scope of work, send a separate, timestamped update in the task as each verified major step or decision occurs:

- Report the chosen approach before a significant action.
- Report an important finding with its evidence and resulting decision.
- Report a substantial edit or workstream step when it actually finishes.
- Report a slow check when it actually starts and say what it will verify.
- Report a changed blocker, risk, direction, or verified validation result.

Do not defer these events to the final answer. A tool call, private reasoning, runtime plan, Steps UI update, or separate `Step` bullet is not a progress update to the operator. Use this format:

```markdown
---

🔄 :: `YYYY-MM-DD HH:MM TZ` :: STEP/update :: `WS[03.02] <specific event>`

- **Outcome:** <verified current result>
  - **Evidence:** <actual source or check>
  - **Blocker:** <real condition, when relevant>
  - **Next:** <actual next action, when relevant>
```

Use `STEP` for the current agent, `AGENT` for a transient helper, and `THREAD` for a persistent visible thread. Use `open`, `update`, or `close` only for an event supported by evidence. Use `steer` only for direction actually sent to an existing `AGENT` or `THREAD`.

Use `🟢` for `open`, `🔄` for `update`, `✅` for `close`, and `🧭` for `steer`. Use `❓`, `⛔`, or `🔬` only when the next action, blocker, or research makes that symbol useful. Use the operator's verified local time and timezone. If they are unavailable, use verified UTC.

Give each actual lifecycle event its own horizontal rule, timestamp, line, and `Outcome`. For a `THREAD` event, append `([thread](codex://threads/<verified-thread-id>))` after the description. Never claim that a sent instruction was acted on or that an agent or thread was closed, deleted, archived, or cleaned up without proof.

## Verify Before Sending

Apply format checks below to the formats chosen or required for this message. The reference checks apply to every message. Plain replies do not need steering labels, markers, or receipts unless applicable instructions require them.

- Answer the actual request in plain English and at the requested length.
- Check each subject's first mention in this message: can the reader identify and locate it from its name or identifier, description when needed, and usable direct link? Keep names stable and pronouns clear. State uncertainty that prevents identification instead of inventing a name, identifier, or link.
- Check every shorthand, not just first mentions: it must be immediately followed by the full name or description, with the canonical identifier and title or description for tracked work. Preserve established meanings. Include private labels in public or team-facing material only when the user explicitly requests it for that communication.
- Before returning a public or team-facing draft, translate private worker paths and any internal labels or workstream markers not authorized for that destination into descriptive names or canonical identifiers.
- Keep ordinary private Codex assistant/user messages, agent-to-agent messages, task-to-task coordination, status updates, and handoffs unsigned because the runtime supplies source identity. A copyable draft for an external or team-facing destination remains covered by that destination's required signature rule even when returned inside private task chat.
- Use the complete, verified `WS[xx.yy]` marker consistently when a current workstream matters.
- Mark every running assignment `in_progress` when the runtime supports it, and keep its `[active]` label and verified `/root/...` worker path visible.
- Preserve all three steering labels for actual steering only.
- Send each major step, decision supported by evidence, material result, and changed blocker to the operator during the turn.
- Format each actual work event as its own timestamped lifecycle receipt with `STEP/open`, `STEP/update`, or `STEP/close` and an `Outcome`.
- Do not replace a lifecycle receipt with a separate `Step` bullet.
- Do not add a lifecycle receipt to ordinary conversation, a factual answer, a steering acknowledgement, or unchanged monitoring.
- Keep unchanged monitoring, recovery, retries, and waiting quiet.
- Preserve evidence, permissions, unaffected work, real next actions, verified links, and accurate timestamps in the operator's timezone.
- Read only the focused reference that the response actually needs.

## Tests

When changing this skill, read [tests/README.md](tests/README.md). Run the relevant scenarios with fresh, independent agents when available.
