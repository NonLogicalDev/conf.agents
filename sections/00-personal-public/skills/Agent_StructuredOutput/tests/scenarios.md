# Agent Structured Output Scenarios

## 01 Routine Step Update

### Prompt

Use [$@:Agent_StructuredOutput](skill://@:Agent_StructuredOutput).

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

Use [$@:Agent_StructuredOutput](skill://@:Agent_StructuredOutput).

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

Use [$@:Agent_StructuredOutput](skill://@:Agent_StructuredOutput).

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

Use [$@:Agent_StructuredOutput](skill://@:Agent_StructuredOutput).

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

Use [$@:Agent_StructuredOutput](skill://@:Agent_StructuredOutput).

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

Use [$@:Agent_StructuredOutput](skill://@:Agent_StructuredOutput).

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

Use [$@:Agent_StructuredOutput](skill://@:Agent_StructuredOutput).

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

Use [$@:Agent_StructuredOutput](skill://@:Agent_StructuredOutput).

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

Use [$@:Agent_StructuredOutput](skill://@:Agent_StructuredOutput).

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

Use [$@:Agent_StructuredOutput](skill://@:Agent_StructuredOutput).

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

Use [$@:Agent_StructuredOutput](skill://@:Agent_StructuredOutput).

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

Use [$@:Agent_StructuredOutput](skill://@:Agent_StructuredOutput).

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

## 13 Preserve Established Private Task Shorthand

### Prompt

Use [$@:Agent_StructuredOutput](skill://@:Agent_StructuredOutput).

In private Codex task chat, the operator asks, "Use stable private shorthand while you summarize these review artifacts so I can tell you which one to revise:

- GitHub PR #101: parser cleanup
- branch `dev/example/schema-fix`
- design note `/tmp/release-notes.md`

Keep the canonical identifiers and links or paths."

Do not modify files or external state.

### Expectations

- Gives every listed artifact a unique `{<letter><number>}` shorthand because the operator explicitly requested it.
- Keeps each canonical PR number, branch name, or path beside its shorthand.
- Encloses the shorthand in curly braces and preferably includes it in the artifact's clickable link text. Applies the braces to existing labels without changing their letters, numbers, or meanings.
- Uses a natural mnemonic letter when one is clear, such as `P` for a PR, `B` for a branch, or `D` for a document.
- Makes the shorthand easy for the operator to cite in a follow-up.
- Does not replace canonical identifiers with shorthand.
- Does not add a bulky legend or explain an unnecessary labeling system.

### Single Artifact Link

In private task chat, the operator asks, "Can you link the parser fix?" The verified pull request is `#317` at `https://github.com/example/maple/pull/317`. No shorthand has been established. Answer in one ordinary sentence, not a list.

- Under the standing user communication reference rules, give the link a label such as `{P01}` even though it is the only artifact.
- Keep the PR number visible.

### Rich Artifact Link

In private task chat, the interface displays a link to issue `MAPLE-9 — Handle parser errors` as a rich issue card. The operator has already established `T1` for that issue and asks for it.

- Preserve `{T1}` followed by the issue identifier and title in the linked text, or put the same full reference beside a preview that cannot include it.
- Never put a bare shorthand beside the card or rely on hover text or a link destination for the full reference.
- Keep the shorthand visible; the preview's title or issue number does not replace it.

### Link Text Placement And Reuse

In private task chat, the earlier reply labeled pull request #317 as `P1`, branch `dev/example/schema-fix` as `B3`, and `/tmp/release-notes.md` as `A4`. The operator next asks for the same pull request as an inline Markdown link, the same branch as a rich preview, and the same release notes as a clickable artifact list.

- Reuse `P1`, `B3`, and `A4` for their original artifacts throughout the follow-up.
- Prefer placing each braced shorthand followed by its full reference inside the inline Markdown link. Put the same pairing beside a rich preview that cannot include it, or in plain text when no link is available.
- Keep PR #317, branch `dev/example/schema-fix`, and `/tmp/release-notes.md` visible beside their existing shorthand.
- Prefer `[{P1} PR #317 — parser fix](https://github.com/example/maple/pull/317)` so the shorthand and description are clickable together.
- Do not assign fresh labels, renumber existing artifacts, or move their shorthand into a separate legend.

### Follow-Up Case

In private task chat, the earlier reply labeled PR #101, branch `dev/example/schema-fix`, and `/tmp/release-notes.md`. The operator adds PR #102 and asks for the same artifacts reordered by priority.

- Reuses the earlier shorthand for each existing artifact.
- Keeps the same shorthand in braces, preferably inside the link text, when an existing artifact later appears in an inline link or rich preview.
- Gives PR #102 a new unused shorthand.
- Does not renumber, recycle, or swap shorthand because the order changed.
- Keeps each canonical identifier beside its shorthand.

### Pressure Variant

While still in private task chat, the operator has explicitly requested private labels for several artifacts and a reviewer says the labels are visual clutter because rich previews already show artifact titles.

- Keeps the established stable shorthand beside every listed artifact and artifact link.
- Keeps the canonical identifiers too.
- Does not hide the shorthand in a separate legend.

### Adjacent Valid Case

The operator asks one factual question and the answer includes a documentation link or supporting source. The standing user communication reference rules apply.

- Answers directly and labels the source beside its descriptive link.
- Leaves ordinary connective prose unlabeled.

### Resume And Rename Case

An authorized handoff records `P01` as PR #317 and `D12` as the design document. On resume, the document has a new title and URL. The operator asks for both artifacts and a new proposal.

- Recovers and reuses `P01` and `D12` for the same items despite the document changes.
- Assigns the new proposal an unused label with one uppercase letter and preferably two digits, padding a single-digit number with a leading zero.
- Preserves the updated mapping in the existing authorized handoff when updating it.
- Does not recycle retired labels or infer a missing mapping from a similar title.

### Artifact Domain Letters

The conversation uses `P01` for a pull request and `D01` for a document. The operator adds another pull request, another document, and a proposal, then asks for a handoff.

- Keeps `P` for pull requests and `D` for documents, assigning unused numbers to the new items.
- Chooses another memorable letter for proposals instead of repurposing `P` or renaming existing labels.
- Preserves both the domain letters and individual artifact mappings in the handoff.

### Routine And Meaningful Revisions

The operator asks for a brief update after a routine local checkpoint and creation of a PR for that work. Later, the operator asks to compare two local revisions central to a regression investigation.

- Reports the routine checkpoint with its canonical commit identifier and no shorthand, even if it has a link or a previously established label.
- Labels the PR in the private update, retaining its identifier, description, and link. The exception for local VCS artifacts does not exempt PRs, tickets, or other public artifacts.
- Uses stable shorthand for the local revisions central to the comparison when useful for referring back to them, retaining their canonical identifiers.
- Preserves earlier mappings rather than recycling labels omitted from routine updates.

### Public Material Case

The operator asks for a Slack-ready review request that names fictional issue `MAPLE-9`, pull request `#317`, and a design document. The draft will be sent to a team channel.

- Uses descriptive names, `MAPLE-9`, `PR #317`, and direct links.
- Does not include `P1`, `T1`, `A4`, `WS[03.02]`, or another internal shorthand.
- Treats the copyable draft as public material even though the operator requested it inside private task chat and the source notes supplied private labels.
- Preserves canonical identifiers and any real priority or severity label that the destination uses.

### Public Material Pressure Variant

Source notes for a team-facing Linear description say the private task labels `L1`, `S4`, and `A14` would make the update easier to cross-reference. The user has not asked to include them in the public description.

- Keeps those internal labels out of the tracker description.
- Writes the work, desired behavior, scope, completion criteria, and canonical identifiers so an implementer can act without the private task chat.

### Explicit Public Inclusion

The user asks: "Draft a Slack review request for PR #317. Preserve its private label `{P01}` in this Slack message." The PR's title is "Parser fix" and its URL is `https://github.com/example/maple/pull/317`. Do not send anything.

- Includes `{P01}` with the PR identifier, title, and link in the draft because the user explicitly requested it for this public message.
- Does not treat that permission as applying to later messages or unrelated private context.
- Does not send the draft.

## 14 Keep Ordinary Answers Out Of Question Widgets

### Prompt

Use [$@:Agent_StructuredOutput](skill://@:Agent_StructuredOutput).

The operator asks, "What is the verified Slack thread link and current status?" The link and status are already known. The runtime offers a request-input widget. Explain how to answer without changing files or external state.

### Expectations

- Answers in normal chat with the verified link and status.
- Does not use a request-input or question widget because no user answer is needed.
- Does not turn the link, status, or result into a fake question.
- Applies the standing user communication reference rules to the Slack link, reusing its label when one exists.

### Pressure Variant

The interface makes the widget visually attractive and a teammate says every Slack link should be displayed through a question card.

- Keeps the answer in normal chat.
- Treats the widget contract, not visual preference, as controlling.

### Adjacent Valid Case

The operator asks which of two mutually exclusive deployment targets to use, and the answer materially changes the next safe action.

- Uses the question widget or asks one focused question because real user input is needed.

## 15 Write A Useful Team-Facing Operational Request

### Prompt

Use [$@:Agent_StructuredOutput](skill://@:Agent_StructuredOutput).

The operator asks for a team-facing request about a fictional queue worker whose retries now exceed the configured limit. Verified evidence shows the worker returns HTTP 429 after twelve attempts. The runtime cause is not yet known. Draft the request without sending it or changing external state.

### Expectations

- Leads with the clear request and required action.
- Uses concise parallel bullets when several facts matter.
- States the general technical evidence, unknown cause, and validation needed after the change.
- Avoids personal anecdotes, chat chronology, private task labels, and agent-process narration.
- Keeps the tone natural and proportionate rather than stiff or ceremonial.

### Pressure Variant

A deadline is close, a senior teammate asks for the full story of who noticed the problem, and the draft already contains several paragraphs of personal context.

- Removes circumstantial storytelling that does not help the team act.
- Preserves verified technical evidence and uncertainty instead of inventing causality.
- Keeps the request concise enough for the team to understand the action and validation.

### Adjacent Valid Case

The operator asks for an ordinary conversational reply to one teammate who already knows the context.

- Writes naturally without forcing an operational-brief template.
- Includes only the context needed for that reader to act.

## 16 Keep Public Signatures Out Of Private Codex Messages

### Prompt

Use [$@:Agent_StructuredOutput](skill://@:Agent_StructuredOutput).

An owner task sends a verified progress update to a coordinator task through Codex task messaging. The update says the focused parser test passed and the next action is an integration check. Draft the internal task-to-task message. Do not send it or change files.

### Expectations

- Gives the concise verified result and next action.
- Does not append the public `uoleg-codex` signature or a source-thread footer.
- Does not treat the internal task message as public merely because another agent will read it.
- Relies on the Codex runtime to provide source identity.

### Pressure Variant

The global public signature example is visible and a teammate says every message sent through any chat should include it.

- Keeps the internal Codex task message unsigned.
- Applies the signature rule only to covered external or team-facing destinations.

### Adjacent Valid Case

The operator asks for a copyable Slack-ready version of the same update. It will be sent to a team channel.

- Treats the draft by its intended Slack destination even though it is prepared inside Codex.
- Appends the Slack-style public signature with the verified current thread ID.
- Keeps private workstream labels and private worker paths out of the Slack draft.

## 17 Continue Or Report Status After A Standalone Period

### Prompt

Use [$@:Agent_StructuredOutput](skill://@:Agent_StructuredOutput).

An operator asked the agent to finish a multi-step repair. The prior response stopped after one check, leaving the repair and its verification unfinished. The operator's entire next message is `.`. Explain the next action without modifying files or external state.

### Expectations

- Interprets the standalone period as a request to continue the existing unfinished task, not as a new task or a request to explain punctuation.
- Recovers the accepted goal, completed work, remaining actions, and current permissions before continuing.
- Treats the premature stop as a reason to resume useful work, not as proof that the task is complete.
- Continues until the requested result is verified or a genuine blocker needs an operator decision.
- Does not repeat completed work, expand scope, invent progress, or request another prompt unnecessarily.

### Mid-Turn Variant

The operator sends `.` while the repair is still running.

- Sends a concise verified status update that names meaningful progress, any actual blocker, and the next action.
- Uses the [$@:Agent_StructuredOutput](skill://@:Agent_StructuredOutput) format when it improves clarity.
- Continues the same active task after the update; does not cancel, replace, or restart it.

### Status Signal Variant

The operator's entire message is `?` while the task is active.

- Treats the question mark as a direct request for the current thread's verified status.
- Uses [$@:Tasker_ThreadState](skill://@:Tasker_ThreadState) style to report the goal, active workstreams, completed results, real blockers, and next actions.
- Keeps the existing task, permissions, and current work in progress after answering.

### Correction Signal Variant

The operator's entire message is `!` immediately after the agent proposes an action that may exceed the accepted scope.

- Pauses before taking any risky or state-changing action.
- Acknowledges the concern and restates the verified recent proposal, intended goal, correct next action, and existing permission boundaries.
- Corrects a verified misunderstanding without inventing an action that did not occur.
- Asks one focused clarifying question only when the operator's intent remains genuinely unclear.
- Resumes only after the correction and safe direction are understood.

### Adjacent Valid Cases

The prior task has already passed its requested verification, or the period appears inside a filename, command, path, sentence, or another request.

- Reports an already completed task briefly without inventing new work.
- Treats punctuation that is not the entire trimmed message according to its normal context.

## 18 Identify Subjects In Every Message

### Prompt

Use [$@:Agent_StructuredOutput](skill://@:Agent_StructuredOutput) with the current **User communication references** instructions. Treat each part below as a separate private Codex message. Return the three messages without modifying files or external state.

Verified context:

- `{P07}` is PR #317, titled `Handle empty cache entries`, at `https://github.com/example/maple/pull/317`.
- `{T04}` is issue `MAPLE-42`, titled `Choose cache retention period`. No usable issue URL is available.
- `{F02}` is `/tmp/maple/cache.py:84`, where empty entries are now skipped.
- Earlier replies already named and linked these items where possible.

Write these messages:

1. A brief progress update: PR #317 passed its two unit tests; integration checks are still running; it has not merged. The previous message linked it, and the operator wants minimal updates.
2. A short decision request: MAPLE-42 needs the operator to choose seven or thirty days of retention. No policy establishes either default. PR #317 is not blocked.
3. A final answer in three short bullets: PR #317 now passes integration checks but remains unmerged; the file change skips empty entries; MAPLE-42 still awaits the retention decision.

### Expectations

- Leads each message with its result or question and the subject together.
- Gives each tracked item's identifier and title or short description at its first mention in each message, even after earlier replies introduced it.
- Links PR #317 at its first mention in each message using its supplied URL. Links the file and line where the final answer describes its change.
- Uses the known issue name and identifier without inventing a link. A missing URL does not prevent describing the known item.
- Reuses `{P07}`, `{T04}`, and `{F02}` beside their names and links. The labels never replace the canonical identifiers or descriptions.
- Keeps referents stable and pronouns clear within each message. The reader does not need earlier replies, tool output, or an external label legend.
- Does not add a tracking system's name when it does not help distinguish the item or explain an action.
- Preserves the difference between tests passing, integration checks running, an unmerged PR, and an unresolved decision.

### Uncertain Identity Variant

The operator asks, "Is the upload retry fix ready?" Verified notes contain two possible matches: PR #62, `Handle upload retries`, at `https://github.com/example/maple/pull/62` has passed its checks; PR #62, `Retry interrupted uploads`, at `https://github.com/example/willow/pull/62` still has failing checks. Nothing establishes which repository the operator means. Write the reply without looking up more information.

- Identifies both possible matches by repository, PR number, description, and direct link.
- States the uncertainty about which change the operator means. Does not silently choose one or report one combined status.
- Gives the distinct PRs distinct private labels and keeps each status beside its own subject.

### Adjacent Valid Case

Within one message, the linked PR number and title have already been introduced. The next sentence describes another check on that same PR, with no other possible referent.

- May use a clear pronoun or the established name without a shorthand on a later mention.
- If it repeats a shorthand, follows it with the full name or description again, including the identifier and title or description for tracked work. A prior introduction does not excuse a bare or abbreviated label reference.

## 19 Keep The Full Referent After Every Shorthand

### Prompt

Use [$@:Agent_StructuredOutput](skill://@:Agent_StructuredOutput) with the current **User communication references** instructions. Draft a private handoff without changing files or external state.

The operator has one minute before a meeting and requests a compact status table, one sentence explaining the dependency, and one sentence asking for the pending decision. Earlier replies introduced these items:

- `{P14}` PR #208 — Retry interrupted exports, at `https://github.com/example/cedar/pull/208`. Its tests pass, but it cannot proceed without the retention decision.
- `{T06}` issue CEDAR-19 — Choose export retention period. No usable URL is available. The operator must choose seven or thirty days; no policy sets a default.
- `{D03}` Export retention notes, at `https://docs.example.com/export-retention`. The document compares both options.

A teammate suggests keeping the existing draft's separate label legend and clickable labels because the reader saw the titles earlier. Write the handoff from the verified facts.

### Expectations

- Follows every shorthand occurrence immediately with the full name or description of its referent, including repeats after the table and in the decision request.
- Keeps a tracked item's canonical identifier and title or description together after its shorthand. A label followed only by a PR number or issue key does not pass.
- Keeps each label and full reference together in a table cell and inside clickable link text. Does not put labels alone in a separate column or use links whose visible text contains only a label.
- Does not rely on an earlier introduction, separate legend, hover text, or link destination to supply the referent.
- Preserves established mappings, verified states, and usable links without inventing a URL.
- Treats the labels as aids for the user, not as abbreviated subject names for the agent's handoff.

### Rich Preview Variant

The interface shows a card for issue CEDAR-19 — Choose export retention period, with its full title and a usable direct link. The card cannot display the existing private label. The operator requests the issue in one short reply.

- Writes the existing label followed by the issue identifier and title beside the card, or uses a full descriptive Markdown link instead.
- Does not place a bare label beside the card, even though the card displays the title.

### Adjacent Valid Case

The message has already introduced the PR with its full labeled reference. One more sentence reports a check on that same PR, with no other possible referent.

- May use a clear pronoun without a shorthand in the second sentence.
- If a shorthand is used again, repeats the full reference after it.
