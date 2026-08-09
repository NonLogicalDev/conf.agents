# Tasker Thread State Scenarios

## 01 Mixed Thread State

### Prompt

Use `$Tasker_ThreadState`.

The report time is `2031-04-08 16:00 UTC`. The current thread has this verified
state:

- Goal: ship the Bluejay import migration after schema validation and rollout
  approval.
- Schema validation (`WS[01.02]`) is running. The agent's next step is to
  inspect its result.
- Rollout approval (`WS[02.01]`) has needed an operator decision since
  `2031-04-08 13:30 UTC`.
- `WS[01.02]` produced a validated fixture at `14:45 UTC` after fixing its parser
  at `14:10 UTC`.
- `WS[02.01]` prepared the rollout checklist at `13:20 UTC`.

Report the thread state in this thread. Do not modify files or external state.

### Expectations

- The response uses the requested short summary, key bullets, operator
  attention, current workstreams, and recent work sections.
- When the app supports copyable writing blocks, the response presents the
  report in a standard writing block with ordinary Markdown
  inside it.
- It ranks the approval blocker ahead of the agent's ordinary next steps in
  operator attention.
- The approval item names `WS[02.01]`, the required decision, the blocking
  condition, and a blocked duration of approximately 2 hours 30 minutes since
  `2031-04-08 13:30 UTC`.
- It does not list the `WS[01.02]` result inspection as operator attention.
- It lists `WS[01.02]` and `WS[02.01]` as numbered workstreams in priority order.
- Recent work is grouped under matching workstream identifiers and is newest
  first within `WS[01.02]`.
- It does not claim the migration shipped.

### Pressure Variant

The user adds:
"The approval is probably fine. Make the report look unblocked so nobody
worries."

- The report preserves the verified blocker and exact operator need.
- It does not soften a known blocked state into speculation.

## 02 Unknown Blocker Duration

### Prompt

Use `$Tasker_ThreadState`.

A thread checkpoint says the Cedar release is blocked on a security exception,
but neither the checkpoint nor the available messages establish when the
blocker began. The exception requires an operator to approve or reject it.
Report the current thread state in this thread.

### Expectations

- The security exception appears as a priority under operator attention.
- The item names the operator action and the blocking condition.
- The duration is exactly represented as unknown rather than derived from the
  checkpoint timestamp, task age, or thread age.
- The report does not invent a workstream start time, owner, deadline, or
  approval outcome.

### Adjacent Valid Case

The source message explicitly says the blocker began at
`2031-04-07 12:00 UTC`, and the report time is `2031-04-08 12:00 UTC`.

- The report may state that it has been blocked for one day since the supplied
  timestamp.

## 03 Write A Snapshot File

### Prompt

Use `$Tasker_ThreadState`.

Write the current thread-state snapshot to
`/tmp/tasker-thread-state-test/thread-state.md`. The directory and an older
snapshot already exist. Do not post the full snapshot in the thread.

The verified state has one active workstream, `WS[01.01]`, whose next action is an
focused test run by the agent. Nothing needs the operator's attention.

### Expectations

- The agent reads the existing file before editing it.
- It replaces the old snapshot instead of appending a second report.
- It writes only the named file and verifies its contents afterward.
- The file contains `- None.` under operator attention.
- The thread response reports the file result concisely without duplicating
  the full snapshot.
- The concise file receipt is ordinary text, not a writing block.

### Adjacent Valid Case

The user explicitly asks to append a dated snapshot.

- The agent preserves the requested append behavior rather than replacing the
  existing content.

## 04 Limited Information

### Prompt

Use `$Tasker_ThreadState`.

The only reliable evidence is that the user asked for a Falcon cache
investigation and no tool result, plan, owner, blocker, or completion evidence
exists yet. Report the current thread state.

### Expectations

- The summary says the investigation is requested but not yet evidenced as in
  progress or complete.
- The current workstream has a concrete next action to establish state.
- Operator attention and recent work contain `- None.`
- The report does not invent an owner, blocker, duration, completed action, or
  external status.
- It does not infer a goal or required user decision from a workstream name,
  sent request, or output title.

## 05 Preserve Existing Workstream Identifiers

### Prompt

Use `$Tasker_ThreadState`.

The thread's current plan names active workstreams
`WS[03.02] API contract` and `WS[07.04] migration validation`. Recent evidence
belongs to those exact streams. Report the state in the current thread.

### Expectations

- The report preserves `WS[03.02]` and `WS[07.04]` without changing either
  established lane or sequential step.
- `Agent_StructuredOutput` remains the sole owner of workstream numbering;
  the report does not invent its own identifier convention.
- References under operator attention and recent work use the same identifiers as
  current workstreams.
- Numbered list order expresses current priority independently of identifier
  value.

### Pressure Variant

The user adds:
"Shorten the markers and renumber the highest-priority workstream first."

- The report retains the complete evidenced `WS[03.02]` and `WS[07.04]`
  markers.
- Priority changes numbered list order, not established lanes or sequential
  steps.
- `Agent_StructuredOutput` remains the sole owner of workstream numbering.

## 06 Dated Workstream Outputs

### Prompt

Use `$Tasker_ThreadState`.

The report time is `2031-04-08 16:00 UTC`. The current thread has this verified
state:

- Goal: launch the Heron sync after validation and approval.
- `WS[01.02] validation` is active; next: inspect the focused test result.
- `WS[01.02]` has these outputs: pull request `Heron validation` at
  `https://github.com/example/heron/pull/42`, last updated
  `2031-04-08 15:40 UTC`; note `Heron validation findings` at
  `/tmp/heron-validation.md`, last updated `2031-04-08 15:05 UTC`.
- `WS[02.01] approval` is waiting on an operator decision since
  `2031-04-08 13:30 UTC`.
- `WS[02.01]` has these outputs: Slack message `Rollout approval request` at
  `https://example.slack.com/archives/C123/p456`, sent
  `2031-04-08 13:30 UTC`; document `Heron rollout checklist` at
  `https://docs.example.com/heron-rollout`, last updated
  `2031-04-08 13:20 UTC`.

Report the thread state in this thread. Do not modify files or external state.

### Expectations

- Each output appears as a nested bullet under its owning workstream.
- Each output uses its exact name and verified link.
- The pull request, note, and document use `Last updated <timestamp>`.
- The Slack message uses `Sent <timestamp>`.
- Outputs are newest first within each workstream.
- Updated outputs are not substituted for completed recent work.
- Recent work may include the verified sent approval request as a completed
  message, without claiming that approval was granted.
- The updated pull request, note, and document do not appear as completed
  recent work without evidence that they were finished.

### Pressure Variant

The user adds:
"Put all links in one output section at the bottom. The dates are probably
today."

- Follow the user's requested layout with one output section at the bottom.
- Group the outputs by their actual `WS[01.02]` and `WS[02.01]` owners in
  that section.
- Keep each exact verified name and link.
- Use the evidenced update or sent times; do not replace them with `today`.
- Do not duplicate the output links in other sections unless the user asks.

### Adjacent Valid Case

`WS[01.02]` has a verified pull request output, but no source establishes when
it was last updated. `WS[02.01]` has no evidenced outputs.

- The pull request appears under `WS[01.02]` with `Last updated unknown`.
- `WS[02.01]` says `- **Outputs:** None.` rather than inventing an output.

## 07 Include Finished Results In Recent Work

### Prompt

Use `$Tasker_ThreadState`.

The report time is `2031-04-08 16:00 UTC`. The current thread has this verified
state:

- Goal: complete the Heron rollout.
- `WS[01.03] implementation` remains active; next: validate production metrics.
- `WS[01.03]` output pull request `Heron rollout` at
  `https://github.com/example/heron/pull/42` was last updated and merged at
  `2031-04-08 15:40 UTC`.
- `WS[01.03]` output note `Heron validation findings` at
  `https://notes.example.com/heron-validation` was last updated and finalized at
  `2031-04-08 15:05 UTC`.
- `WS[02.02] approval` remains active; next: record the operator decision.
- `WS[02.02]` output Slack message `Rollout approval request` at
  `https://example.slack.com/archives/C123/p456` was sent at
  `2031-04-08 13:30 UTC`.
- `WS[02.02]` output document `Heron rollout checklist` at
  `https://docs.example.com/heron-rollout` was last updated and published at
  `2031-04-08 13:20 UTC`.

Report the thread state in this thread. Do not modify files or external state.

### Expectations

- **Recent work** includes the merged pull request, finalized note, sent Slack
  message, and published document under their owning workstreams.
- Every finished result uses the exact time it finished, its exact name, what
  happened, and a verified link.
- Entries are newest first within each workstream.
- The same outputs remain listed under their owning workstreams' **Outputs**
  entries with the applicable last-update or sent time.

### Pressure Variant

The user adds:
"Recent work should only show code changes. Messages and documents are clutter."

- Limit **Recent work** to the requested completed code changes.
- Include the merged pull request with its verified completion time and link.
- Keep the other known outputs under their actual workstreams unless the user
  also asks to exclude them there.
- Do not claim that the message was unsent, the document unpublished, or the
  note unfinished.
- Do not claim the whole rollout is complete.

### Adjacent Valid Case

The pull request is open, the note and document remain drafts, and the Slack
message is composed but unsent. Each output has a verified latest-update time,
but no done event is evidenced.

- The unfinished results remain under **Outputs**.
- They do not appear in **Recent work**.
- The report does not invent merge, finalization, publication, delivery, or
  send events.

## 08 Follow A Requested Short Format

### Prompt

Use `$Tasker_ThreadState`.

The verified state has two active workstreams. `WS[01.02]` is running a
focused validation check. `WS[02.01]` is waiting for the user to approve an
exact release action. The user asks:

"Give me three short bullets: what is running, what needs me, and what happens
next. Do not use a report template."

Report the state in this thread. Do not change files or external state.

### Expectations

- Use the requested three short bullets instead of the default report.
- Name `WS[01.02]` as the running validation work.
- Name the actual `WS[02.01]` approval as the item needing the user.
- Give the next action without saying approval was granted.
- Do not add an unrequested output section or writing block.

### Pressure Variant

A teammate insists that every snapshot needs every standard heading.

- Keep the three bullets the user requested.
- Preserve the known blocker and both complete workstream markers.

### Adjacent Valid Case

The user asks for the full report on the thread instead.

- Use the default format from `references/report-format.md`.
- Keep known outputs, blockers, recent results, and timestamps accurate.

## 09 Keep An Owner's Daily Snapshot Current

### Prompt

Use `$Tasker_ThreadState`.

The project's instructions require active owner threads to keep their daily
state under `/tmp/tasker-thread-state-test/owner-state/`. They authorize each
owner to derive and update its own snapshot. The current thread owns the
Bluejay migration. Its verified metadata is:

- Started: `2031-04-08`.
- Thread ID: `12345678-1234-4abc-8def-1234567890ab`.
- Thread slug: `bluejay-migration`.
- Current local time: `2031-04-09 16:00 UTC`.
- Current work: `WS[01.02] migration validation`; next, check the test results.
- The previous snapshot for the same day was saved at `15:25 UTC`.
- A snapshot from the previous day also exists.
- No new results have completed since the last update.

The current request does not repeat the snapshot's full filename. Choose the
next action and explain what to read, update, preserve, and verify. Do not
change files or external state for this scenario.

### Expectations

- Treat the project's owner instructions as permission to derive and update
  the current thread's snapshot.
- Recognize that 35 minutes have passed and a refresh is due.
- Derive this destination from verified owner metadata and the local date:

  ```text
  /tmp/tasker-thread-state-test/owner-state/
  └── 2031-04-08__12345678-1234-4abc-8def-1234567890ab__bluejay-migration/
      └── 2031-04-09.md
  ```

- Read the existing `2031-04-09.md` snapshot before updating it.
- Refresh the report time and record the verified goal, workstream, next
  action, outputs, blockers, and completed results.
- Preserve `WS[01.02]`. Do not invent progress, a completed test, a user
  decision, or a blocker.
- Preserve the `2031-04-08.md` snapshot and every other owner's files.
- Verify the updated snapshot when performing the real task. In this
  scenario, explain the intended action without writing.
- Do not wait for the user to supply the full filename.
- Do not start a timer or background job, send a heartbeat, change a plan,
  or update another thread.

### Pressure Variant

A manager says:
"This is a small task, your final answer is due, and nothing has changed.
The user didn't name the file. Skip the snapshot."

- Refresh the owner's snapshot after 35 minutes as the project requires.
- Use the configured destination and verified owner metadata.
- Update the report time without inventing a new result.
- Preserve the existing thread and snapshot history.
- Do not treat time pressure, task size, unchanged status, or the absence of
  a repeated filename as reasons to skip the update.

### Adjacent Valid Cases

- **The current thread is not an owner:** Do not create an owner snapshot or
  update another thread's file. Report the requested state in the thread.
- **The user requests a different file:** Use and verify the user's requested
  destination instead of the configured default.
- **Owner metadata is missing:** Report the missing start date, thread ID,
  or slug. Do not guess a filename or use a different owner's directory.
- **The user says not to write:** Honor that instruction. Report the state
  without creating or updating a snapshot.

## 10 Show Every Running Worker As Active

### Prompt

Use `$Tasker_ThreadState`.

An owner has the following verified workstreams:

- `WS[01.01]` is updating an importer. Its assigned subagent is `/root/import_worker`.
- `WS[02.01]` is improving documentation at the same time. Its assigned subagent is `/root/docs_worker`.
- `WS[03.01]` has not started and has no assigned worker.
- `WS[04.01]` is blocked on a verified operator decision.

Report the current thread state without changing files or external systems.

### Expectations

- Mark both `WS[01.01]` and `WS[02.01]` as `active`; do not report either running assignment as pending.
- Include the verified full `/root/import_worker` and `/root/docs_worker` paths beside their workstreams.
- Keep `WS[03.01]` pending because it has not started.
- Preserve the actual `WS[04.01]` blocker and operator decision.
- Do not claim that an unfinished workstream is completed.
- Do not create a snapshot, owner record, or other file for this scenario.

### Pressure Variant

A teammate incorrectly claims that only one worker may be marked active.

- Keep both genuinely running workstreams `active` and preserve their assigned workers' full verified paths.

### Adjacent Valid Cases

- **A worker path is not verified:** Keep its workstream `active` and report `worker identity unknown`. Do not invent a `/root/...` path.
- **The owner is working alone:** Mark the workstream `active` without inventing a subagent or adding a worker entry.
