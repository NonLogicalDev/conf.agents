---
date: 2026-07-12
goal: settings-modernization
workstream: legacy-config-retirement
status: archived
status_log:
  - "[[2026-07-12]]@12:00 created -> backlog"
  - "[[2026-07-12]]@12:10 backlog -> active"
  - "[[2026-07-12]]@13:30 active -> completed"
  - "[[2026-07-12]]@16:05 completed -> archived"
subject: remove-legacy-config
---

# Purpose

Remove the legacy settings reader after replacement behavior is verified, then
retain the obsoleted result as historical evidence.

# Scope

- **In scope:** the experimental removal of `src/settings/load-v1.ts`, its
  focused version 2 loader proof, recovery from the named backup, and archival
  after additive migration obsoletes the approach.
- **Out of scope:** schema migration, settings import, broader loader redesign,
  and use of this archived plan to direct current implementation.

# Completion Criteria

- [x] `test ! -e src/settings/load-v1.ts` confirms the experimental branch
  removed the legacy reader.
- [x] `npm test -- tests/settings/load-v2.test.ts` reports
  `loads version 2 fixture` passed.
- [x] The plan records `tmp/settings/load-v1.ts.pre-cleanup` as the recovery
  source for a partial cleanup failure.
- [x] The plan is under `archived/`, names the active migration plan that
  obsoletes it, and states that it is historical only.

# Blockers

- None.

# Artifacts

- **Artifact:** [Final legacy-cleanup
  revision](https://example.com/scm/settings-app/commits/dddddddddddddddddddddddddddddddddddddddd)
  - **Type:** Git commit from a removed branch.
  - **State:** Upstream branch `legacy-config-cleanup` was deleted after
    archival.
  - **Contents:** Final historical record of the retired experiment to remove
    the legacy reader.
  - **Location:** [Final legacy-cleanup
    revision](https://example.com/scm/settings-app/commits/dddddddddddddddddddddddddddddddddddddddd).
  - **SCM surface:** full final commit
    `dddddddddddddddddddddddddddddddddddddddd`.
  - **Reconstruct:** clone
    [example/settings-app](https://example.com/scm/settings-app), fetch
    `origin`, then check out
    `dddddddddddddddddddddddddddddddddddddddd`.
- **Artifact:** `src/settings/load-v1.ts`
  - **Type:** Local repository source file.
  - **State:** Removed.
  - **Contents:** Legacy settings reader deleted by the historical experiment.
  - **Location:** `src/settings/load-v1.ts`.
- **Artifact:** `tmp/settings/load-v1.ts.pre-cleanup`
  - **Type:** Backup saved locally.
  - **State:** Retained.
  - **Contents:** Recovery copy used if the historical cleanup is replayed.
  - **Location:** `tmp/settings/load-v1.ts.pre-cleanup`.
# Context

## Fictional Repository

- **Legacy reader:** `src/settings/load-v1.ts` reads schema version 1 directly.
- **Replacement candidate:** `src/settings/load.ts` was expected to read only
  schema version 2 after cleanup.
- **Obsoleting plan:**
  [[epic-002-settings-migration/active/plan-01 2026-07-12 migrate-settings.md]]
  later replaced this cleanup approach with additive migration.

## Original Desired Behavior

- Remove `load-v1.ts` only after version 2 settings loading is proven.
- Keep existing version 2 users unaffected.

# Steering

- [[2026-07-12]]@12:00 - Remove the legacy reader only after replacement proof.
  - **Effect:** cleanup cannot run before version 2 validation.
- [[2026-07-12]]@16:05 - Preserve this completed plan as history after the
  migration plan obsoletes it.
  - **Effect:** move the plan to `archived/` rather than deleting it.

# Decisions

- **Decision:** Keep legacy reader removal separate from schema migration.
  - **Rationale:** At the time, cleanup was expected to happen only after a
    separate replacement path was proven.
- **Decision:** Archive this plan after additive migration replaces the cleanup
  approach.
  - **Rationale:** The historical proof remains useful. Do not use this plan
    to guide new implementation.

# Implementation Steps

1. [x] Add a version 2 fixture so the replacement loader can be checked without
   legacy input.
2. [x] Copy `src/settings/load-v1.ts` to
   `tmp/settings/load-v1.ts.pre-cleanup` before removal.
3. [x] Remove `src/settings/load-v1.ts` from the experimental cleanup branch.
4. [x] Run the version 2 loader test and record its result before marking the
   cleanup complete.
5. [x] Archive the plan after the migration plan adopts additive version 1
   support instead of reader removal.

# Work Log

- [x] [[2026-07-12]]@12:10 - Activated cleanup after identifying the legacy
  reader and replacement test.
- [x] [[2026-07-12]]@12:15 - Saved
  `tmp/settings/load-v1.ts.pre-cleanup` before removing the legacy reader.
- [x] [[2026-07-12]]@13:20 - Ran the version 2 loader test; the replacement
  fixture loaded successfully.
- [x] [[2026-07-12]]@13:30 - Marked cleanup completed after the focused test
  passed.
- [x] [[2026-07-12]]@16:05 - Archived this plan because the additive migration
  plan obsoletes the cleanup approach.

# Validation

- **Legacy reader removal:**
  - **Command or action:** verify the legacy reader path is absent.
  - **Working directory:** fictional repository root.

    ```bash
    test ! -e src/settings/load-v1.ts
    ```

  - **Expected result:** the command exits zero and prints no output.
  - **Failure signals:** the command exits with a nonzero status because the legacy reader
    still exists.
  - **Interpretation:** failure would mean the original cleanup behavior was
    not complete.

- **Version 2 loader proof:**
  - **Command or action:** run the focused version 2 loader test.
  - **Working directory:** fictional repository root.

    ```bash
    npm test -- tests/settings/load-v2.test.ts
    ```

  - **Expected result:** the runner reports `loads version 2 fixture` passed.
  - **Failure signals:** the fixture fails to load or the command exits
    with a nonzero status.
  - **Interpretation:** failure would have blocked the original cleanup plan
    from reaching completed.

# Idempotence And Recovery

- **Repeated execution:** rerunning the focused loader test does not change
  repository files.
- **Partial failure:** if cleanup had failed before validation, restore the
  named backup before retrying.

  ```bash
  cp tmp/settings/load-v1.ts.pre-cleanup src/settings/load-v1.ts
  ```

- **Safe fallback:** the archived plan is historical only; follow the active
  migration plan for current work.

# Outcomes & Retrospective

- The original cleanup experiment proved version 2 fixture loading.
- Additive migration later obsoleted this plan, so the completed work remains
  archived for context rather than directing current implementation.
