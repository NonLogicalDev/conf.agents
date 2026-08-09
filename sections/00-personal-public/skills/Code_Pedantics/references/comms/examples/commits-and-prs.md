# Commit And Pull Request Examples

These are invented examples. Use them to recognize message shape and workflow mapping, not as evidence that any command ran in a real repository. Repository instructions and required templates still override the fallback form.

When applicable instructions require a public-agent signature, append it only to the PR description after the canonical reviewer prose. Omit it only when the user explicitly exempts that communication. Keep actual Git commit messages and PR titles unsigned.

## Squash Workflow With One Commit

Use one canonical message for the single branch commit and the reviewer prose in the PR title and description. Append any required public-agent signature only to the PR description.

```text
feat(config-prune): plan deletions before mutation

`config-prune` previously deleted stale generated entries immediately. Build
and validate a deletion plan first so dry-run can show the intended changes
without mutating configuration, while normal mode removes only stale generated
entries after validation passes.

## What Changes

- Add `--dry-run` so operators can inspect planned deletions without writing.
- Keep unmanaged entries outside the deletion plan.
- Require plan validation before normal mode mutates configuration.

## Verification

- `cargo test -p config-prune planner`
  - Covers plan construction and selective deletion behavior.
- `cargo test -p config-prune cli`
  - Covers dry-run output and command behavior.
- `./config-prune --dry-run --config /tmp/config.toml`
  - Verified the configuration stayed byte-identical while planned deletions
    were printed.
```

## Squash Workflow With A Branch Series

Each new transport commit describes the entire branch at the moment that commit is created. The newest commit and the canonical reviewer prose in the PR match; any required signature belongs only in the PR description.

Commit 1, after adding planning only:

```text
feat(config-prune): plan generated-entry deletions

Build a deletion plan before pruning so the command can identify stale
generated entries without mixing discovery and mutation.

## What Changes

- Separate stale-entry discovery from configuration mutation.

## Verification

- `cargo test -p config-prune planner`
  - Covers selection of stale generated entries and preservation of unmanaged
    entries.
```

Commit 2, after adding a dry run, becomes the newest message for the entire branch. The PR title and canonical reviewer prose copy this message; append any required signature only to the PR description:

```text
feat(config-prune): add safe dry-run pruning

Plan stale generated-entry deletions before mutation so operators can inspect
the same plan that normal mode will apply. Dry-run prints the plan without
writing, and normal mode mutates only after validation passes.

## What Changes

- Separate deletion planning from mutation.
- Add `--dry-run` for inspecting the validated plan without writing.
- Preserve unmanaged entries in both dry-run and normal mode.

## Verification

- `cargo test -p config-prune planner`
  - Covers cumulative planning and selective mutation behavior.
- `cargo test -p config-prune cli`
  - Covers dry-run behavior.
- `./config-prune --dry-run --config /tmp/config.toml`
  - Verified dry-run printed the plan and left configuration unchanged.
```

Do not rewrite Commit 1 merely because the branch later grew. Its message was the cumulative branch message at that earlier point. Commit 2 is now the current canonical PR/squash message.

## Non-Squash Workflow

Each commit stays scoped to its own diff. The PR uses the same message format, but it summarizes the series instead of matching either commit exactly.

Commit 1:

```text
feat(config-prune): build deletion plans

Separate stale-entry discovery from mutation so callers can inspect the exact
set of generated entries proposed for deletion.

## Verification

- `cargo test -p config-prune planner`
  - Covers plan construction and unmanaged-entry preservation.
```

Commit 2:

```text
feat(config-prune): add dry-run output

Expose the deletion plan through `--dry-run` without writing configuration.

## Verification

- `cargo test -p config-prune cli`
  - Covers dry-run output and non-mutation behavior.
```

PR title and description:

```text
feat(config-prune): make pruning inspectable

Make stale generated-entry pruning inspectable before mutation. The series
separates deletion planning from mutation, then exposes that plan through
`--dry-run` so operators can review the proposed deletions safely.

## What Changes

- Build deletion plans that exclude unmanaged entries.
- Add `--dry-run` to print the plan without writing configuration.

## Verification

- `cargo test -p config-prune planner`
  - Covers the planning contract.
- `cargo test -p config-prune cli`
  - Covers the dry-run contract.
```

## Small Change That Needs Little Explanation

Omit `## What Changes` when the summary already carries the complete reviewable story. Keep `## Verification` by default. When there is no meaningful check, briefly explain why verification does not apply.

```text
docs(auth): correct callback spelling

Correct a typo in the callback setup prose without changing a command, link, or
behavior.

## Verification

- Not applicable: this typo-only correction has no meaningful check.
```
