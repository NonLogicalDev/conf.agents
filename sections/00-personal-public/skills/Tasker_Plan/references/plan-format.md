# Plan Format

Read this reference when creating a new plan or improving the shape of an
existing one. Preserve the project's documented format and useful history.

## Required Sections For A Numbered Plan

Keep the work log, artifacts, decisions, and validation in the numbered
plan itself. An owner `_index_.md` and epic `EXEC_STATE.md` summarize their
own levels; neither replaces this plan's detailed history. Keep each plan
in its group's `active/`, `backlog/`, or `completed/` directory matching
the actual status recorded in its document.

Preserve useful project frontmatter and use this plan shape:

```markdown
# <Plan title>

## Purpose

<What the work should achieve and why.>

## Scope

### In scope

- <Accepted work.>

### Out of scope

- <Excluded work or constraint.>

## Completion Criteria

- [ ] <Specific outcome and the evidence that will prove it.>
- [x] <Completed outcome and the actual verification.>

## Blockers

- <Blocked action, observed cause, next action, and clearing evidence.>

## Artifacts

- **Artifact:** <Verified name or linked title.>
  - **Type:** <Source file, test, branch, commit, pull request, or document.>
  - **State:** <Planned, created, changed, tested, published, or merged.>
  - **Contents:** <What this artifact contributes.>
  - **Location:** <Verified path or direct link.>
  - **SCM surface:** <Verified repository, branch, or exact revision, when relevant.>

## Decisions

- **Decision:** <Chosen approach.>
  - **Rationale:** <Observed reason and relevant tradeoff.>

## Implementation Steps

- [x] <Completed milestone and relevant evidence.>
- [ ] <Next meaningful implementation step.>

## Work Log

- [x] [[YYYY-MM-DD]]@HH:MM - <Verified material milestone.>
  - **Evidence:** <Observed result, verified file, exact revision, or check.>
  - **Effect:** <Changed state, finding, decision, or cleared blocker.>
  - **Next:** <Actual remaining action, when relevant.>

## Validation

- **Command or action:** <Exact relevant check.>
  - **Expected result:** <What passing would establish.>
  - **Observed result:** <Passed, failed, blocked, or not yet run.>
  - **Failure signals:** <Relevant actionable failure, when useful.>
  - **Interpretation:** <What the actual result does and does not prove.>

## Unfinished Work

- [ ] <Remaining task, dependency, or next action.>
```

Keep the section responsibilities even when an existing project uses an
equivalent established heading. Preserve an existing accurate status or
status frontmatter. Record the actual timestamp only when verified. Use
`None.` when no blocker, artifact, or decision exists; say `unknown` for an
unverified path, link, revision, command, or result.

## Record Meaningful Work

Append one entry for a material investigation, implementation milestone,
artifact change, decision, validation outcome, changed blocker, or user
correction. Group closely related commands under the milestone they prove.
Keep unrelated events distinguishable. Preserve earlier entries and mark
superseded decisions rather than rewriting the past.

A task checkbox alone does not show what was observed, changed, produced,
decided, or verified. Conversely, a command-by-command transcript hides the
meaningful outcome. Each work log entry should let the next agent understand
what happened, what proves it, what changed, and what remains.

## Maintain the artifact inventory

Record both relevant local and external artifacts in the numbered plan.
Update existing entries as their verified state changes. Include an exact
repository, branch, commit, pull request head, or recovery path only when
it exists and matters. Do not imply that a planned file is created, that a
created test passed, or that a draft pull request merged.

Also record actually produced external artifacts in the epic's
`EXEC_STATE.md`. Keep that summary for the epic consistent with the numbered
plan; it is not a substitute for the plan's artifact inventory or work log.

## Keep content accurate

State unknowns directly. Use actual project paths, decisions, owners, and
validation results only when verified and relevant. Distinguish planned,
active, blocked, and completed work.

Update the existing document as understanding changes. Preserve established
useful sections and add new ones only when they clarify the actual work.
