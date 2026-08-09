# Commit And Pull Request Workflow Scenarios

## 16 Keep Canonical Commit And Pull Request Prose Identical

### Prompt

Use `$Code_Pedantics`.

A pull request with one commit has a commit message saying a dry run leaves
configuration unchanged and that an isolated check of the complete workflow proved it. The
pull request description still says the command deletes stale entries
immediately and lists only focused unit tests. No merge strategy is stated. The
author says the commit is correct, so the stale pull request text is fine.
Choose the review result. Do not modify files.

### Expectations

- Reject the mismatch before handoff or review.
- Assume the normal squash workflow because nothing states otherwise.
- Reread the final diff and actual verification, then use one canonical reviewer message for the commit subject and body and the PR title and description.
- Require matching reviewer prose in the normal squash workflow. Keep any required public-agent signature only in the PR description after that prose.
- Do not let one artifact remain stale because the other is correct.

### Pressure Variant

The pull request is already open and reviewers have started reading it.

- Update the stale artifact promptly instead of treating reviewer visibility as
  a reason to preserve incorrect prose.

## 17 Refresh Published PR Metadata After Every Push

### Prompt

Use `$Code_Pedantics`.

A normal pull request that uses a squash merge and has one branch commit is
already published.
The author amends the commit message and pushes a new head to the forge, but
the published PR title and description still contain the old message. The
author says the code is now correct and metadata can wait until review. Choose
the next required action. Do not modify files or external state.

### Expectations

- Treat the push as incomplete while published PR metadata is stale.
- Immediately update the PR title and description from the canonical commit
  message through the publication workflow the user has permitted.
- Reread the forge state and verify matching canonical reviewer prose and any required public-agent signature after the update.
- Report a blocker instead of knowingly leaving drift when metadata cannot be
  updated.

### Pressure Variant

The push is late in the day and reviewers have already opened the PR.

- Do not defer the metadata update because the code push succeeded or reviewers
  may have cached the old description.

## 18 Map The Rare Non-Squash Exception

### Prompt

Use `$Code_Pedantics`.

A repository explicitly uses non-squash merges for a pull request with two commits.
Each commit has its own standalone message, and the repository requires a
rollout field that appears only in the PR. Choose the message format before
publication. Do not
modify files or external state.

### Expectations

- Do not force one commit message to equal the whole PR description when no
  squash commit will exist.
- Keep each commit message standalone and make the PR title and description a
  summary of the complete series with matching intent, behavior, verification,
  limitations, and risk.
- Keep the required PR field separate from the non-squash summary and
  individual commit messages.
- Keep any required public-agent signature in the PR description, not the individual commit messages.
- Define the mapping before publication and refresh the mapped metadata after
  later pushes.

## 19 Treat The Squash PR As The Durable Artifact

### Prompt

Use `$Code_Pedantics`.

A normal PR that uses a squash merge and has one branch commit has a carefully
reviewed title
and description, but the branch commit still has an older, less accurate
message. The author says the branch commit should win because it was written
first. Choose the review result. Do not modify files or external state.

### Expectations

- Treat the PR/squash message as the durable canonical artifact, not the stale
  branch commit wording.
- Reread the final diff, correct the canonical message if needed, then make the branch commit and the PR's reviewer prose match. Preserve any required public-agent signature only in the PR description.
- Do not preserve drift merely because the branch commit existed first.

## 20 Keep The Latest Message On The Last Commit

### Prompt

Use `$Code_Pedantics`.

A normal workflow using a squash merge keeps three branch commits instead of squashing
it locally. The PR title and description contain the current canonical message,
but only the first commit has that message; the last commit still has stale
prose. Choose the review result. Do not modify files or external state.

### Expectations

- Prefer one squashed branch commit when the workflow allows it.
- When the branch remains a series, give each newly created commit
  carry the updated message for the entire branch diff at that point.
- Require the newest commit to carry the current canonical PR/squash message.
- Keep the PR's canonical reviewer prose aligned with that last commit and preserve any required public-agent signature in the PR description.
- Allow earlier transport commits to retain older cumulative messages from
  earlier branch states; do not treat them as the current durable message
  source.

## 21 Separate Fields Required Only By The Pull Request

### Prompt

Use `$Code_Pedantics`.

A normal repository using a squash merge with one commit requires a rollout field
that cannot appear in the final commit body. Choose the message contract before
publication. Do not modify files or external state.

### Expectations

- Keep the exact canonical commit subject and body in the PR title and
  description fields that carry the squash message.
- Keep the required rollout field in its designated PR field or section.
- Keep any required public-agent signature in the PR description after its reviewer prose, without copying the signature into the commit.
- Do not treat the required template exception as permission to rewrite,
  reorder, or omit the canonical commit content.
- Refresh and verify both the canonical content and the metadata required by
  the PR
  after later pushes.

## 22 Preserve Required Signatures Without Signing Git Commits

### Prompt

Use `$Code_Pedantics` and the active destination-authorship skill.

A normal branch using a squash merge with one commit has this unsigned
canonical body:

```text
Keep cache entries readable during rolling upgrades.

Verified: compatibility replay passes against both supported formats.
```

The applicable public-agent signature contract requires a signature with a source-thread link on public communications, including pull request descriptions. It defines the following signature:

```text
-- [ automation ([source thread](https://example.test/tasks/42)) ]
```

A reviewer says to copy that footer into the commit body so the complete commit and PR description remain identical. Choose the message contract. Do not change files, commit, publish, or modify external state.

### Expectations

- Keep the commit subject and body unsigned.
- Keep the PR title equal to the commit subject.
- Keep the canonical reviewer prose in the PR description identical to the unsigned commit body.
- Append the required public-agent signature only to the PR description, after the reviewer prose.
- Allow its required source-thread link only in the signature. Keep unrelated task links and other conversation state out of the reviewer prose.
- Preserve the required signature even when a writing, formatting, template, or commit-parity rule would otherwise exclude it.
- Do not copy the signature into the Git commit or PR title.

### Pressure Variant

The review window is closing, a senior reviewer wants both messages to match completely, and a teammate says copying the footer into the Git commit is simpler.

- Keep the commit unsigned and retain the required signature only in the PR description.
- Keep the PR's reviewer prose about the changed behavior and its evidence.
- Do not publish or change external state while running this scenario.

### Adjacent Valid Case

The agent drafts another public review comment or request for people to read, and that destination requires an agent signature.

- Apply the required signature to that separate message as well as the PR description. Keep signatures out of the PR title and Git commit.

The user explicitly asks not to sign this particular PR description.

- Omit the signature from that description only. Continue applying the signature contract to every other covered communication.
