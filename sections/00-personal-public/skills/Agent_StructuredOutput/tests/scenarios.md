# Agent Structured Output Scenarios

## 01 Routine Step Update

### Prompt

Use `$Agent_StructuredOutput`.

The operator's local timezone is America/Los_Angeles. Report that `WS[03.01]` finished its mechanical checks at the current report time. The validator passed, no blocker remains, and the next action is behavioral validation.

### Expectations

- Begin with a horizontal rule.
- Use one current local timestamp with the timezone abbreviation.
- Use `✅`, `STEP/close`, and a concrete inline-code description.
- Use the complete `WS[03.01]` marker in the lifecycle description.
- Follow the lifecycle line with an `Outcome` parent bullet.
- Nest the validator evidence, lack of blocker, and next action.
- Do not add a routine `## WS[03.01]` header.
- Do not shorten the current marker to a bare lane number.

### Pressure Variant

A senior reviewer says the exact lifecycle grammar is unnecessary because the deadline is close and the result is already obvious.

- Preserve the exact delimiter, timestamp, lifecycle line, and outcome tree.
- Do not collapse the work receipt into prose.

### Adjacent Valid Case

The operator asks, "What is the capital of France?"

- Answer directly without a work interstitial, timestamp, header, or status tree.

## 02 Additive Steering Acknowledgement

### Prompt

Use `$Agent_StructuredOutput`.

`WS[01.01]` is implementing a parser. `WS[02.01]` is validating fixtures. The operator adds an independent documentation task and does not pause either existing workstream. Acknowledge the steering only. Do not report work progress.

### Expectations

- Uses `Understanding`, `Direction`, and `Next` labels.
- States that the new documentation task starts in the independent `WS[03.01]` lane and step.
- States that `WS[01.01]` and `WS[02.01]` continue.
- Names the immediate next action for `WS[03.01]`.
- Uses the complete marker for every current workstream.
- Do not add a lifecycle interstitial because this is only an acknowledgement.
- Do not describe the new task as an intake lane, bounded packet, or similar invented process term.
- Do not apologize or narrate a correction.

### Pressure Variant

The new task is urgent, a manager calls it the top priority, and a draft says the other work should stop.

- Prioritize the new task without stopping unaffected workstreams.
- Keep the understanding log in plain English.
- Keep `Understanding`, `Direction`, and `Next` even when the operator asks for a terse acknowledgement.
- Do not infer a full track switch.

### Adjacent Valid Case

The operator explicitly says, "Stop `WS[01.01]` and `WS[02.01]`. Replace all current work with the documentation task."

- Acknowledge that both existing workstreams stop.
- Make the replacement the sole next workstream.
- Use the same labeled understanding log and make the replacement explicit.

## 03 Agent And Thread Lifecycles

### Prompt

Use `$Agent_StructuredOutput`.

Report these two discrete events at the current report time:

1. The transient validation agent `/root/parser_validation` has finished and its exact scenario passed.
2. You sent a scope correction to durable thread `01900000-0000-7000-8000-000000000000`.

The thread ID is verified. No runtime record was deleted or archived.

### Expectations

- Give each event its own horizontal rule, timestamp, lifecycle line, and `Outcome` tree.
- Use `✅` with `AGENT/close` for the ended transient assignment.
- Include the verified `/root/parser_validation` path in the agent's lifecycle event.
- Use `🧭` with `THREAD/steer` for the sent correction.
- Append a clickable `[thread](codex://threads/...)` link to the thread line.
- Describe the direction sent without claiming the thread acted on it.
- Do not claim agent cleanup, thread archival, or deletion.

### Pressure Variant

The deadline is close, the draft groups both events under one timestamp, and a reviewer says the thread link can be a raw URL.

- Preserve two discrete lifecycle chunks.
- Keep the verified link in the required final field.
- Do not use a raw URL or one grouped status block.

### Adjacent Valid Case

The current agent completes its own local check without using another agent or thread.

- Use `STEP/close`.
- Do not use an `AGENT` or `THREAD` lifecycle tag.

## 04 Separate A Completion Receipt From A Summary

### Prompt

Use `$Agent_StructuredOutput`.

Write a final response that first reports a completed validation step and then provides a self-contained `## Summary` with the result and next action.

### Expectations

- Render a valid `STEP/close` lifecycle chunk with an `Outcome` tree.
- Put exactly two rendered horizontal rules between that tree and `## Summary`.
- Do not add another rule between the summary heading and its content.

### Pressure Variant

The draft already uses one rule, the answer is due immediately, and a reviewer says the visual distinction is cosmetic.

- Put exactly two rendered rules between the receipt and the summary.

### Adjacent Valid Case

Return only a status completion with no separate answer summary.

- Do not add two separating rules or an empty `## Summary`.

## 05 Annotated Source Lists

### Prompt

Use `$Agent_StructuredOutput`.

State that a fictional gateway rejects requests without both an audience and an allowlist. Support the claim with:

- `/tmp/example/settings.py:115`, which defines the audience; and
- `/tmp/example/identity.py:314`, which enforces the allowlist.

Do not inspect files or report work status.

### Expectations

- Keep the claim above the sources.
- Give each source its own linked sibling bullet.
- Add a short distinct blurb explaining each file's role.
- Preserve the exact line anchors.
- Do not chain both links in one sentence.
- Do not add a timestamped work interstitial.

### Pressure Variant

A prewritten draft chains both links in one sentence. The message is small, already approved, and due immediately.

- Replace the chain with annotated sibling bullets.

### Adjacent Valid Case

Only `settings.py:115` supports the claim, and its role is clear.

- A concise inline link is valid.

## 06 Steps UI Encodes Parallel Lanes And Serial Order

### Prompt

Use `$Agent_StructuredOutput`.

You are beginning a release-preparation task with these current workstreams:

- `/root/schema_worker` is inspecting schema drift; updating the decoder and running focused decoder tests have not started.
- `/root/operator_guide` is independently revising the operator guide; checking its links has not started.
- `/root/package_manifest` is independently auditing the package manifest.

Update the runtime Steps UI to represent the current workstreams and their execution order. Return the exact step labels. Do not modify files or external state.

### Expectations

- Updates Steps UI before substantive work.
- The captured tool call or resulting UI state contains the expected labels, order, and statuses; a list written only in prose does not pass.
- Uses these exact labels:
  - `WS[01.01] [active] /root/schema_worker — Inspect schema drift`
  - `WS[01.02] Update decoder`
  - `WS[01.03] Run focused decoder tests`
  - `WS[02.01] [active] /root/operator_guide — Revise operator guide`
  - `WS[02.02] Check operator guide links`
  - `WS[03.01] [active] /root/package_manifest — Audit package manifest`
- Sets `WS[01.01] [active] /root/schema_worker — Inspect schema drift` to `in_progress`.
- Sets `WS[02.01] [active] /root/operator_guide — Revise operator guide` to `in_progress`.
- Sets `WS[03.01] [active] /root/package_manifest — Audit package manifest` to `in_progress`.
- Keeps all three running workers visibly `[active]` with their verified `/root/...` paths.
- Leaves work that has not started as plain pending items without an `[active]` marker.
- Treats different `xx` values as generally parallel lanes.
- Treats increasing `yy` values within one lane as serial steps.
- Does not use repeated `WS01:` labels that omit serial position.

### Pressure Variant

The Steps UI already uses repeated `WS01:` labels. The deadline is close, and a manager wrongly assumes only one task can be active, so the other running workers should look pending and their names should be omitted.

- Replaces the legacy labels with complete `WS[xx.yy]` markers.
- Marks every running assignment `in_progress` and visibly `[active]` with its verified canonical `/root/...` path.
- Leaves only work that has not started without an active label.
- Keeps the lane numbers stable while assigning serial step numbers.
- Does not defer the Steps UI update until the work is finished.

### Ongoing Synchronization Case

Steps UI currently contains:

- `WS[01.01] [active] /root/schema_worker — Inspect schema drift`
- `WS[01.02] Update decoder`
- `WS[02.01] [active] /root/operator_guide — Revise operator guide`

The schema inspection completes, and `/root/schema_worker` starts checking decoder compatibility before updating the decoder. `/root/operator_guide` continues working, and `/root/package_manifest` starts an independent package audit.

- Updates Steps UI immediately.
- Marks `WS[01.01] /root/schema_worker — Inspect schema drift` complete without an active marker.
- Inserts `WS[01.02] [active] /root/schema_worker — Check decoder compatibility`.
- Moves the still-pending decoder work to `WS[01.03] Update decoder`.
- Adds `WS[03.01] [active] /root/package_manifest — Audit package manifest`.
- Keeps the unrelated guide lane as `WS[02.01] [active] /root/operator_guide — Revise operator guide`.
- Sets all three running tasks to `in_progress` and keeps their verified worker paths visible.
- The captured tool call or resulting UI state proves the update occurred.

### Priority Change Preserves Lane Identity

The current parser step is `WS[03.02] [active] /root/parser_worker — Update decoder`. The independent validation step is `WS[07.01] [active] /root/fixture_worker — Audit fixtures`. The operator raises fixture validation to the highest priority without stopping parser work.

- Keeps parser work in lane `03` and fixture validation in lane `07`.
- Uses the complete `WS[03.02]` and `WS[07.01]` markers in prose and the Steps UI.
- Changes priority or item order without renumbering either lane.
- Keeps both workers `in_progress` and visibly `[active]` with their verified canonical paths.
- Does not assign lane `01` merely because fixture validation now comes first.

### Compatibility Case: A Runtime Allows Only One Active Item

The particular Steps tool explicitly rejects more than one `in_progress` item while three verified workers continue running.

- Set one running task to `in_progress` and keep all three visibly `[active]` with their verified `/root/...` worker paths.
- Use this fallback only for the runtime that actually has this limit; do not apply it when multiple tasks can be active.

### Adjacent Valid Case: Worker Path Is Unknown

A task is genuinely running, but its worker's canonical path cannot yet be verified.

- Mark the task `in_progress` and visibly `[active]`; identify the worker as unknown instead of inventing a `/root/...` path.
- Add the worker's canonical path when it becomes available.

### Adjacent Valid Case

The operator asks one factual question. No ongoing task or current workstream exists.

- Answers directly without creating an empty or decorative Steps UI plan.

## 07 Cut Inflated Language And Filler

### Prompt

Use `$Agent_StructuredOutput`.

Rewrite the draft below as an operator update with one main point. Preserve every fact. Return only the revised text.

> At this juncture, it is important to note that we have successfully operationalized a strategically aligned validation initiative whose robust results provide a meaningful degree of epistemic confidence regarding the parser remediation trajectory. In terms of concrete next-step enablement, the focused parser tests passed 42 of 42 cases. It should further be underscored that release progression is presently gated by the absence of the artifact signature, with the signer job constituting the next actionable workstream. Overall, this represents significant forward momentum.

### Expectations

- States that all 42 focused parser tests passed.
- States that the missing artifact signature blocks the release.
- States that the signer job is next.
- Removes warm-up phrases, process narration, generic praise, and the repeated conclusion.
- Does not replace the draft's inflated wording with phrases such as `strategy-aligned validation`, `strong progress`, `meaningful confidence`, or `confidence that the fix is on track`.
- Uses a short paragraph without decorative headings or labels.

### Pressure Variant

A senior writer says the long, formal wording sounds more rigorous and asks to keep it.

- Still uses plain English.
- Preserves the facts and technical names, not the inflated tone.
- Does not preserve an unsupported confidence or progress claim in shorter words.
- Does not add a sentence defending the rewrite.

### Adjacent Valid Case

The result depends on the exact technical terms `JWT aud`, `RepositoryAccessPolicy`, and `baseRefName`.

- Preserves the exact terms because replacing them would reduce accuracy.
- Explains a term only when the operator may not know it.

## 08 Complex Mission Steering Exposes Material Understanding

### Prompt

Use `$Agent_StructuredOutput`.

`WS[01.01]` is repairing a parser. `WS[02.01]` is validating fixtures. The operator says:

"Start a Tasker_Mission for workspace-sync on remote-host. The first mission workstream is robust remote workspace operation without SSH-tunnel dependence. Treat the existing OPS-4 research handoff only as intake evidence, not proof of current runtime or repository state. Do not send Slack or Linear writes unless I separately allow them."

Acknowledge the steering only. Do not report work progress or modify state.

### Expectations

- Uses `Understanding`, `Direction`, and `Next` labels.
- States that the steering starts a new Tasker_Mission for `workspace-sync` on `remote-host` and gives its first workstream the stated remote-workspace scope.
- States that `WS[01.01]` and `WS[02.01]` remain unchanged unless the operator said otherwise.
- States that OPS-4 is intake evidence, not proof of current runtime or repository state.
- States that Slack and Linear writes remain unauthorized.
- Names the immediate next action: identify or create the qualified mission owner and establish its durable plan before mission work.
- Does not add a lifecycle interstitial because this is only an acknowledgement.
- Does not claim that the owner, plan, Slack thread, or Linear task already exists.

### Follow-Up Steering Case

After the mission steering above, the operator says, "Use the same mission. Make current runtime and repository inspection the first action, and keep WS[01.01] and WS[02.01] moving." The operator does not repeat the OPS-4 or Slack and Linear limits.

- Uses the same labeled understanding log.
- Preserves that OPS-4 is intake evidence rather than current-state proof.
- Preserves that Slack and Linear writes remain unauthorized.
- States that current runtime and repository inspection is now the first mission action while `WS[01.01]` and `WS[02.01]` continue.
- Does not add a lifecycle receipt merely because the steering changes the next action.

### Pressure Variant

The deadline is close. A senior engineer says the one-sentence summary is already correct, asks to skip the detail, and says everyone knows what OPS-4 means.

- Keeps the detailed understanding log.
- Preserves the evidence and the limits on what each connector may do.
- Does not hide the effect on existing workstreams or the next action.

### Adjacent Valid Case

The operator instead asks, "What is the capital of France?"

- Answers directly without `Understanding`, `Direction`, or `Next` labels.
- Does not add a lifecycle interstitial or workstream machinery.

## 09 Keep One Workstream When Its Step Advances

### Prompt

Use `$Agent_StructuredOutput`.

Parser investigation `WS[03.01]` has a recorded owner, plan, Slack thread, and Linear issue. The investigation finishes and parser repair becomes the next step. Independent fixture validation `WS[07.01]` remains active.

Describe the current workstreams and next parser action. Do not edit a plan, change runtime steps, create a task, send a message, or modify external state.

### Expectations

- Identifies the next parser action as `WS[03.02]`.
- Keeps the parser in stable lane `03`.
- Preserves fixture validation as `WS[07.01]`.
- Reuses the parser's existing owner, plan, Slack thread, and Linear issue.
- Does not treat the changed sequential step as a new workstream.
- Does not create or request a duplicate task, message, issue, or monitor.

### Pressure Variant

A reviewer says that `WS[03.02]` is a new identifier and needs a new Slack parent, Linear issue, owner, and plan before repair can begin.

- Preserves the established lane and existing destinations.
- Rejects duplicate ownership and external artifacts.
- Uses the complete current parser marker without shortening it.

### Adjacent Valid Case

A genuinely independent manifest audit starts while parser repair remains active.

- Assigns the manifest audit its own available lane and first serial step.
- Keeps `WS[03.02]`, `WS[07.01]`, and their existing owners unchanged.

## 10 Keep Unchanged Monitoring And Recovery Quiet

### Prompt

Use `$Agent_StructuredOutput`.

A recovery task has checked the same target seven times in 17 minutes and then six more times. Every check reports the same active state, the same next action, and no new error. The operator has not requested a status update. A draft says, `Recovery is still running; checking again.`

Choose what to report and what to do next. Do not send a message, change runtime steps, edit files, run a recovery command, or modify external state.

### Expectations

- Does not send or draft another unchanged status update.
- Does not emit a timestamp, lifecycle receipt, heartbeat, or countdown.
- Does not call elapsed time, repeated polling, or an unchanged retry a new blocker or material progress.
- Continues useful permitted work or uses an available wait mechanism.
- Does not invent a completed recovery, successful retry, or state change.

### Pressure Variant

A manager requests an update every minute, the deadline is near, and a prepared thirteen-message transcript calls every unchanged poll progress. The operator has not requested those updates.

- Keeps unchanged monitoring quiet.
- Does not report repeated polling as progress or a new decision.
- Does not let a deadline or a third party override the operator's request.

### Material Change Case

The next verified check shows a new failure that genuinely prevents the accepted work from continuing.

- Reports the newly evidenced blocker once.
- States its actual effect, evidence, and next action.
- Does not replay the preceding unchanged recovery checks.

### User-Requested Status Case

The operator asks, `What is the current recovery status?`

- Answers the question directly in plain English.
- States the verified unchanged state and actual next action.
- Does not automatically add a badge, timestamp, heading, or lifecycle receipt.

## 11 Surface Major Decisions During The Turn

### Prompt

Use `$Agent_StructuredOutput`.

The operator has requested a parser repair. You have verified that a legacy decoder causes the failure and have selected the smallest compatible fix. You have not edited the file. Applying the fix, running a slow focused test, and checking its result are separate steps.

Provide the actual messages you would send to the operator before the edit, when the test starts, and when the result becomes known. State when each message is sent. Do not edit files, run commands, update external state, or disclose private chain-of-thought.

### Expectations

- Sends a real in-task progress message before starting the edit.
- States the verified legacy-decoder finding and the smallest compatible repair.
- Distinguishes the proposed edit from a completed edit.
- Sends another update after the edit is verified and before the slow test.
- States that the test is running without claiming it has passed.
- Reports the test result only after its actual outcome is known.
- Gives every material work event its own horizontal rule, actual local timestamp, and lifecycle line.
- Uses `STEP/update` for the verified cause, chosen repair, and actual validation progress.
- Uses `STEP/close` only after a step's completion is evidenced.
- Follows every lifecycle line with an `Outcome` parent bullet and nests relevant `Evidence` and `Next` bullets beneath it.
- Uses a workstream marker only when the scenario supplies a verified one.
- Does not treat a runtime plan, tool call, Steps UI update, private reasoning, or final answer as a sent progress message.
- Does not replace a material-event lifecycle receipt with a standalone `Step` bullet.

### Pressure Variant

The deadline is close, the fix is one line, and a reviewer says to skip all messages during the task because the final answer can explain everything later.

- Sends the material step and decision updates before continuing.
- Keeps every lifecycle receipt short and proportionate.
- Preserves each event's separator, timestamp, lifecycle line, and `Outcome` tree.
- Does not narrate private reasoning or send a message for each tool call.
- Does not invent a passing test or claim an unfinished step is complete.

### Adjacent Valid Case

A focused test is still running and its last verified status has not changed.

- Does not resend the same running status, add a heartbeat, or narrate unchanged polling.
- Sends a new update when the result, blocker, decision, or next action actually changes.

## 12 Let Markdown Replies Wrap Naturally

### Prompt

Use `$Agent_StructuredOutput`.

The operator asks for a short Markdown reply that explains how to inspect a configuration preview and what to do if it fails. A reviewer says every source line must fit within 80 characters. Return the reply only; do not change files or send a message elsewhere.

### Expectations

- Keep each prose paragraph on one natural source line and let the application wrap it for display.
- Do not add line breaks just to satisfy 72, 80, or any other preferred width.
- Preserve blank lines, Markdown lists, code fences, and intentional line breaks when the reply needs them.
- Follow any actual code formatting rule without applying it to Markdown prose.

### Adjacent Valid Case

The operator requests a nested list with a fenced shell command.

- Preserve the list's nesting and the command's required format.
- Keep each prose paragraph free of an invented width limit.

## 13 Give Artifacts And Artifact Links Stable Shorthand

### Prompt

Use `$Agent_StructuredOutput`.

The operator asks, "Summarize these review artifacts so I can tell you which one to revise:

- GitHub PR #101: parser cleanup
- branch `dev/example/schema-fix`
- design note `/tmp/release-notes.md`

Keep the canonical identifiers and links or paths."

Do not modify files or external state.

### Expectations

- Gives every listed artifact a unique `<letter><number>` shorthand.
- Keeps each canonical PR number, branch name, or path beside its shorthand.
- Places the shorthand outside and immediately beside any link to the artifact.
- Uses a natural mnemonic letter when one is clear, such as `P` for a PR, `B` for a branch, or `D` for a document.
- Makes the shorthand easy for the operator to cite in a follow-up.
- Does not replace canonical identifiers with shorthand.
- Does not add a bulky legend or explain an unnecessary labeling system.

### Single Artifact Link

The operator asks, "Can you link the parser fix?" The verified pull request is `#317` at `https://github.com/example/maple/pull/317`. Answer in one ordinary sentence, not a list.

- Put the stable shorthand outside and immediately beside the clickable PR link.
- Keep the PR number visible.
- Do not omit the shorthand just because the answer is a single sentence.

### Rich Artifact Link

The interface displays a link to issue `MAPLE-9` as a rich issue card. The operator asks for the issue.

- Put a stable shorthand such as `T1` outside and immediately beside the linked issue or its rich preview.
- Keep the issue's canonical identifier visible.
- Do not hide the shorthand in the link text or treat the preview's title or issue number as a substitute.

### Follow-Up Case

The earlier reply labeled PR #101, branch `dev/example/schema-fix`, and `/tmp/release-notes.md`. The operator adds PR #102 and asks for the same artifacts reordered by priority.

- Reuses the earlier shorthand for each existing artifact.
- Keeps the same shorthand outside the link when an existing artifact later appears in an inline link or rich preview.
- Gives PR #102 a new unused shorthand.
- Does not renumber, recycle, or swap shorthand because the order changed.
- Keeps each canonical identifier beside its shorthand.

### Pressure Variant

A reviewer says the labels are visual clutter because the list is short and rich previews already show artifact titles.

- Keeps the stable shorthand beside every listed artifact and artifact link.
- Keeps the canonical identifiers too.
- Does not hide the shorthand in a separate legend.

### Adjacent Valid Case

The operator asks one factual question and the answer contains no artifact link or listed artifact. It may link to ordinary documentation or a supporting source.

- Answers directly without inventing artifact shorthand.
- Leaves ordinary documentation and evidence links unlabeled.
