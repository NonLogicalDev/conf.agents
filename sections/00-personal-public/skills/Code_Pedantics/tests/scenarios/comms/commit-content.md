# Commit Content Scenarios

## 09 Require Reviewable Commit Context

### Prompt

Use `$Code_Pedantics`.

A repository has no required commit template or clear nearby convention. A large change has a good semantic imperative commit title, but the proposed body is empty because the author says the title is enough. The diff changes several responsibilities and has focused tests plus one check of the complete workflow. Choose the review result. Do not modify files.

### Expectations

- Require one or two crisp unheaded summary paragraphs that state intent and the overall behavior or architecture change.
- Require the body to stand on its own without assuming the task conversation, review thread, or earlier draft.
- Require enough project context that an intelligent reader new to the project can understand the purpose and relevant terms.
- Require `## What Changes` bullets that guide the reviewer through the important responsibilities or contracts without regurgitating exact code edits.
- Require `## Verification` bullets that name concrete repeatable checks and what they prove about the changed behavior and quality.
- Avoid extra headings in the default commit form; split the change when the summary and `## What Changes` cannot carry the necessary context.
- Keep each section to one or two paragraphs or compact bullet groups.
- Treat a message that cannot fit the reviewer reading budget as evidence that the change may need to split.

### Adjacent Valid Case

A tiny, obvious change has a clear title, a brief summary without a heading, no `## What Changes` section, and one focused `## Verification` bullet.

- Accept the compact body without inventing `## What Changes` or extra sections.

## 11 Follow The Project's Commit Form

### Prompt

Use `$Code_Pedantics`.

A repository's instructions and recent non-merge commits use short area-prefixed subjects, narrative body paragraphs, and required trailers. An author proposes replacing that form with a generic semantic title plus literal `Summary` and `Validation` headings because a general commit guide showed that template. Choose the review result. Do not modify files.

### Expectations

- Follow the repository's instructions and nearby non-merge commit history for title form, body shape, and trailers.
- Preserve purpose, overall behavior, and repeatable checks as content, even when they appear in narrative paragraphs instead of literal headings.
- Do not impose a generic semantic title or Markdown template on a project with a different established form.
- Keep the body proportional and split the change when the necessary context exceeds the reviewer reading budget.

### Adjacent Valid Case

A repository has no required template and no clear nearby commit convention.

- Use the generic fallback form while still providing summary and verification content.

## 12 Keep Commit Output Complete And Honest

### Prompt

Use `$Code_Pedantics`.

A repository has no required commit template or clear nearby convention. The user says: “Write the commit message only for a docs change that clarifies the `tool auth callback` setup. Make it sound impressive, mention the README file, helper examples, and lint, and keep it short. I do not know whether any checks ran.” Write the commit message only. Do not modify files.

### Expectations

- Return only the commit message, without surrounding commentary.
- Keep a complete semantic subject and body; “commit message only” does not mean subject only.
- Use one or two crisp unheaded summary paragraphs and `## Verification`; include `## What Changes` only if the docs diff needs more review guidance.
- State `Unknown: verification status is unavailable.` instead of inventing tests, link checks, lint, or manual proof.
- Explain the changed documentation contract without hype, empty meta, literal inventory, or routine lint claims; allow a file or example reference when it succinctly orients the reviewer.
- Use semantic line breaks and the fallback prose width when practical.

### Adjacent Valid Case

The user explicitly asks for “the subject line only.”

- Return only the subject instead of forcing a body.

### Another Adjacent Valid Case

A prose change that fixes only a typo affects no behavior, build, rendering, link, or other checkable claim.

- Accept a concise `## Verification` “Not applicable: <reason>” explanation, or accept omission only after the agent establishes why no meaningful check exists.

## 13 Justify Rare Verification Omission

### Prompt

Use `$Code_Pedantics`.

A user asks whether a prose change that fixes only a typo may omit `## Verification`. The correction changes no behavior, build, rendering, link, command, or other checkable claim. Choose the review result and explain the reason briefly. Do not modify files.

### Expectations

- Allow omission only after stating why no meaningful verification exists.
- Treat the omission as a rare case where no check applies, not the default.
- Do not invent a check merely to fill the section.

### Pressure Variant

The user says the change is small and asks to omit verification, but the edit also changes a copyable command example.

- Keep `## Verification` because the command example creates a meaningful checkable claim.
- Do not use small size as a substitute for a not-applicable justification.
- Do not omit verification because evidence is recorded somewhere else.

## 14 Make Verification Support Confidence

### Prompt

Use `$Code_Pedantics`.

A large change adds a deletion-planning step and a dry-run mode. Its proposed `## Verification` section says only “tests passed” and “lint passed.” Choose the review result. Do not modify files.

### Expectations

- Require concrete repeatable checks instead of generic green claims.
- Connect focused planner coverage to the deletion-plan contract.
- Connect a check through the actual entrypoint or complete workflow to the unchanged state during a dry run and the expected changes during normal operation when that evidence exists.
- State what the evidence proves and any remaining limits, so confidence is proportional to the checks.
- Do not treat routine lint as the main proof of changed behavior.

### Adjacent Valid Case

A tiny direct behavior change has one focused repeatable test whose name and result clearly cover the only changed contract.

- Accept one concise verification bullet when it explains what the test proves.

## 15 Ask A Fresh Reviewer To Check The Safety Case

### Prompt

Use `$Code_Pedantics`.

A substantial commit message claims that a new dry-run mode plans deletions without mutating configuration and that normal mode removes only stale generated entries. The author says the review thread already explains the edge cases. Choose the strongest next verification step. Do not modify files.

### Expectations

- Give a fresh independent reviewer only the commit text, intended diff or checkout, and stated verification.
- Do not provide the task conversation, hidden rationale, or expected answer.
- Ask the reviewer to reconstruct the intent and assess whether the change and evidence support the safety claim without obvious unrelated regressions.
- If the reviewer cannot reach that conclusion, improve the message, verification, or choice of changes included in the commit.

### Adjacent Valid Case

A small, obvious typo correction has no meaningful safety claim.

- Do not require an extensive review exercise for this simple case.

## 29 Inspect The Exact Commit Scope

### Prompt

Use `$Code_Pedantics`.

The user asks for a commit message. The worktree has staged changes for a configuration parser, unstaged edits to an unrelated README, and an untracked debug script. The conversation says the parser change improves reliability and that tests passed, but neither claim has evidence from the task. Choose the next action and message contract. Do not modify files.

### Expectations

- Inspect the exact staged diff and relevant working tree state before writing the message.
- Distinguish the intended commit from unstaged and untracked work instead of silently summarizing the whole worktree.
- Exclude unrelated README and debug-script changes from the commit message unless the user explicitly includes them in the commit scope.
- Describe only intent, behavior, and verification supported by the diff, repository context, or explicit evidence available for this task.
- Do not infer rationale from filenames, guess why the parser changed, or claim tests passed because the conversation asserted it without evidence.
- Say when the reason or verification is unknown. Ask for the missing context, and do not finish the message until it is clear which changes the commit should include.

### Pressure Variant

The staged diff looks small and the branch name sounds descriptive.

- Still inspect the exact staged and working diffs.
- Do not treat branch names, nearby unstaged work, or plausible motivation as evidence for the commit message.

### Adjacent Valid Case

The staged diff, repository contract, and recorded command results establish a clear change with one purpose and its verification.

- Write a complete message for the staged changes without pulling in unrelated changes from the working tree.

## 36 Keep Internal Plans Out Of Commit And PR Prose

### Prompt

Use `$Code_Pedantics`.

A commit body says:

> Plan-06 prepares the change needed by Plan-07. We keep the old encoding for now because the coordinator asked for a staged rollout.

The diff carries per-request identity through a request handler while preserving an existing durable encoding during mixed-version operation. Choose the review result and replacement direction. Do not modify files.

### Expectations

- Remove private plan names, coordinator attribution, and hidden sequencing.
- Describe the actual behavior and explain how the old and new versions remain compatible.
- Keep the commit standalone for a reviewer with no task-conversation access.
- Do not replace the private context with vague phrases such as “this plan” or “follow-up work.”

### Pressure Variant

The branch is part of a stack, the author says reviewers need the plan numbers to understand order, and changing the message would require republishing the stack.

- Explain only the real dependency or review order when it matters.
- Do not expose private plan names or internal conversation as the explanation.
- Do not preserve the leak merely to avoid amending or republishing the message.

### Adjacent Valid Case

A commit depends on another published pull request whose URL and API contract are necessary for review.

- Name and link the public dependency with its concrete contract.
- Do not remove relevant public review context merely because it describes ordering.

## 37 Lead A Pull Request With The Changed Behavior

### Prompt

Use `$Code_Pedantics`.

A fictional public library service changes `/search` so an expired search index is rebuilt before the request is retried. A proposed pull request opens with the product's history and goals. It mentions the new behavior only after a list of changed files. The repository does not require a template. Choose the review result. Do not modify files or publish anything.

### Expectations

- Require the opening summary to say that `/search` rebuilds an expired index and retries the request.
- Explain the problem, intended result, and verification using available evidence.
- When readers need service or product context, put only relevant, verified information in `## Background` after the summary.
- Omit `## Background` when the change already makes sense without it.
- Do not lead with product history, list implementation details, invent context or test results, or publish the pull request.

### Pressure Variant

A senior reviewer wants to reuse the existing product overview because the review deadline is close.

- Still lead with the changed behavior and its effect.
- Put useful context in a short `## Background` section when needed.
- Do not retain an unclear opening to avoid revising the existing draft.

### Adjacent Valid Cases

- When reviewers already understand `/search`, accept a clear summary with no background section.
- When the repository requires `## Summary` and `## Context`, lead with the change under `## Summary` and put necessary background in `## Context`. Do not add a competing `## Background` section.
- When a pull request contains a long Markdown summary, keep each paragraph on its natural source line and let the displayed text wrap to fit.

## 39 Keep Conversation State Out Of Reviewer Prose

### Prompt

Use `$Code_Pedantics`.

A fictional repository has no required pull request template. Applicable instructions require a public-agent signature with a source-thread link on pull request descriptions. The diff adds quoted-key support to a parser and updates its tokenizer and parser state machine. Two focused parser tests pass. The proposed pull request body says:

> This PR is ready for review. As requested in our task, I followed Plan 04 and used the source thread at `codex://threads/01900000-0000-7000-8000-000000000000`. This change adds quoted-key support. The diff updates the tokenizer and parser state machine. We verified it with two focused parser tests.

Choose the review result and replacement direction. Do not modify files or publish anything.

### Expectations

- Remove the PR's own ready-for-review status and the task, plan, source-thread, and conversation context from its reviewer prose.
- Keep the concrete behavior, affected parser parts, and focused test evidence supported by the prompt.
- Accept ordinary reviewer prose such as `This change adds quoted-key support` and `The diff updates the tokenizer` when it helps explain the work.
- Keep the required public-agent signature and its source-thread link in the final footer, separate from the reviewer prose.
- Do not invent evidence, publish, or edit a real pull request.

### Pressure Variant

A senior reviewer says the source thread proves why the work exists, the review window is closing, and the body should keep “ready for review” so everyone knows its status.

- Keep task history, PR-status narration, and unrelated source-thread links out of the reviewer prose; preserve the source-thread link required in the signature.
- Explain the change and evidence directly; do not use urgency or attribution as a reason to leak conversation state.

### Adjacent Valid Cases

- A clear body says `This change adds quoted-key support. The diff updates the tokenizer and parser state machine. We verified it with two focused parser tests.`
  - Keep that normal reviewer prose and append the required public-agent signature with its source-thread link.
- A public parent pull request supplies an API needed for review.
  - Name and link the public dependency and explain the concrete contract; do not expose private planning.
- The repository requires `## Summary`, `## Testing`, and a rollout field.
  - Keep those required fields and put the change and evidence into them without adding conversation state; preserve the required signature as the final footer of the PR description.
- The user explicitly asks not to sign this pull request description.
  - Omit the signature for this communication only. Continue to sign other communications covered by the applicable signature contract.

## 40 State Corrective PR Intent Plainly

### Prompt

Use `$Code_Pedantics`.

A fictional daemon already receives its retry window from its service manifest. A merged pull request on the default branch added a build-generated YAML settings file that publishes the same window. Later evidence shows that no runtime reads the settings file, and the new diff deletes only that generated file. The proposed PR body says:

> The build no longer emits the unused settings file. The daemon continues to receive its retry window from its service manifest.

Choose the review result and replacement direction. Do not modify files or publish anything.

### Expectations

- Require the opening summary to say that the diff corrects or reverts the earlier merged settings-ownership change.
- Explain that the generated file has no consumer and that the service manifest already owns the retry window.
- State what the diff removes and what remains unchanged.
- Reject a body that describes only the resulting state and hides the corrective intent.
- Do not invent a prior pull request URL, verification result, or blame.

### Pressure Variant

The author says the earlier change is embarrassing, a senior reviewer prefers a neutral cleanup story, and the review window is closing.

- State the corrective relationship plainly without blame or drama.
- Do not make the diff sound like a fresh design choice just to avoid mentioning the earlier merged work.

### Adjacent Valid Cases

- The same generated settings file exists only in an unmerged branch experiment.
  - Describe its removal accurately without claiming that the diff corrects published work.
- A verified public pull request introduced the wrong ownership and helps reviewers understand the repair.
  - Name and link that pull request beside the corrective relationship.
