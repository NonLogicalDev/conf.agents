---
date: 2026-07-12
goal: settings-modernization
workstream: schema-migration
status: active
status_log:
  - "[[2026-07-12]]@13:45 created -> backlog"
  - "[[2026-07-12]]@14:25 backlog -> active"
unblocks:
  - "[[task-001-settings-import/backlog/plan-01 2026-07-12 import-settings.md]]"
obsoletes:
  - "[[epic-004-legacy-config/archived/plan-01 2026-07-12 remove-legacy-config.md]]"
subject: migrate-settings
---

# Purpose

Migrate old settings files to schema version 2 without losing explicit user
values or leaving a partially written file after interruption, so existing users
can adopt the new schema without manual recovery.

# Scope

- **In scope:** version 1 to version 2 conversion in
  `src/settings/migrate.ts`; atomic persistence in `src/settings/load.ts`;
  focused checks of migration, interruption, and the loader smoke test; and release of
  settings-import implementation after its exact readiness evidence exists.
- **Out of scope:** implementation of the settings-import command, removal of
  the version 1 reader, unrelated settings cleanup, and changes outside the
  settings loader and its focused tests.

# Completion Criteria

- [x] `npm test -- tests/settings/migrate.test.ts` proves explicit version 1
  values survive conversion and a second migration produces identical version
  2 JSON.
- [ ] The test for an interrupted write proves that stopping before rename leaves the
  original `config/settings.json` as valid JSON.
- [ ] The loader smoke check prints `loaded schema v2` and
  `preserved explicit values: true` for
  `fixtures/settings/v1-explicit-values.json`.
- [ ] This plan's `Work Log` records passing tests for repeated runs and
  interrupted writes, plus the expected loader smoke output before settings import
  implementation begins.

# Blockers

- None.

# Artifacts

- **Artifact:** [settings-v2-migration upstream
  branch](https://example.com/scm/settings-app/branches/settings-v2-migration)
  - **Type:** Git branch.
  - **State:** Pushed and active.
  - **Contents:** Schema conversion that is safe to repeat and unfinished
    atomic persistence in the loader.
  - **Location:** [settings-v2-migration upstream
    branch](https://example.com/scm/settings-app/branches/settings-v2-migration).
  - **SCM surface:** repository
    [example/settings-app](https://example.com/scm/settings-app); remote
    `origin`; upstream `refs/remotes/origin/settings-v2-migration`; latest
    verified full commit `aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa`.
  - **Owned paths:** `src/settings/migrate.ts`, `src/settings/load.ts`, and
    `tests/settings/migrate.test.ts`.
  - **Reconstruct:** clone
    [example/settings-app](https://example.com/scm/settings-app), fetch
    `origin`, then check out
    `aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa`.
- **Artifact:** [Draft pull request
  42](https://example.com/scm/settings-app/pulls/42)
  - **Type:** Pull request.
  - **State:** Draft.
  - **Contents:** Migration implementation with proof that it is safe to repeat and
    incomplete interruption validation.
  - **Location:** [Draft pull request
    42](https://example.com/scm/settings-app/pulls/42).
  - **SCM surface:** head repository
    [example/settings-app](https://example.com/scm/settings-app); head ref
    `settings-v2-migration`; head commit
    `aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa`; base ref `main`.
- **Artifact:** `src/settings/migrate.ts`
  - **Type:** Local repository source file.
  - **State:** Verified.
  - **Contents:** Version 1 to version 2 conversion proven safe to repeat by the
    migration unit test.
  - **Location:** `src/settings/migrate.ts`.
- **Artifact:** `src/settings/load.ts`
  - **Type:** Local repository source file.
  - **State:** In progress.
  - **Contents:** Atomic version 2 persistence; interruption behavior and the
    loader smoke check remain unfinished.
  - **Location:** `src/settings/load.ts`.
# Context

## Fictional Repository

- **Loader:** `src/settings/load.ts` reads `config/settings.json`.
- **Migration:** `src/settings/migrate.ts` converts schema version 1 to version
  2.
- **Unit tests:** `tests/settings/migrate.test.ts` covers migration behavior.
- **Smoke command:** `npm run settings:smoke` loads one fixture settings file
  through the real loader.

## Current Behavior

- The loader accepts only schema version 1.
- A draft migration rewrites the settings file in place.
- An interruption between truncating and writing can leave invalid JSON.

## Desired Behavior

- The loader accepts schema version 1, migrates it to version 2, then applies
  defaults.
- Running migration twice produces the same version 2 file.
- An interrupted write leaves either the original file or the complete migrated
  file, never partial JSON.

## Integration Context

- **Existing model:** `load.ts` owns settings loading and default application.
- **New requirement:** keep old settings readable while version 2 becomes
  the persisted format.
- **Integrated model:** keep migration as a pure conversion in `migrate.ts` and
  let `load.ts` call it before applying defaults.
- **Existing pieces that change:**
  - `src/settings/migrate.ts` gains conversion that is safe to repeat.
  - `src/settings/load.ts` gains atomic persistence.
  - `tests/settings/migrate.test.ts` gains coverage for interruptions and repeated runs.
- **Architecture impact:** settings migration stays inside the settings module;
  no separate cleanup command becomes the source of truth.

# Steering

- [[2026-07-12]]@13:45 - Preserve explicit user values during migration.
  - **Effect:** defaults apply only after version 1 values are converted.
- [[2026-07-12]]@14:00 - Keep migration additive until version 2 loading is
  proven.
  - **Effect:** do not remove version 1 reading before checking the complete workflow.
- [[2026-07-12]]@14:25 - Prove interruption recovery and behavior that is safe to repeat
  before releasing settings-import implementation.
  - **Effect:** the dependent implementation remains blocked until the named
    validation evidence exists; independent intake investigation may activate
    and proceed.

# Decisions

- **Decision:** Keep the schema conversion pure and separate from disk writes.
  - **Rationale:** A pure conversion can be tested twice on the same input
    without filesystem state hiding drift.
- **Decision:** Persist through a temporary file and rename.
  - **Rationale:** Rename makes the final replacement atomic on the supported
    filesystem, so interruption cannot leave partial JSON at the target path.
- **Decision:** Keep version 1 reading until the smoke command proves version 2
  loading.
  - **Rationale:** Removing the old reader before proof would strand existing
    settings.

# Implementation Steps

1. [x] Extract version 1 to version 2 conversion into
   `src/settings/migrate.ts` so the same input can be converted repeatedly
   without drift.
2. [x] Add unit cases in `tests/settings/migrate.test.ts` that preserve explicit
   values and prove the second migration run is unchanged.
3. [ ] Change `src/settings/load.ts` to write version 2 through a temporary
   sibling file and rename so interruption cannot leave partial target JSON.
4. [ ] Add interruption coverage that stops before rename and proves the
   original file remains readable.
5. [ ] Run the smoke command through `load.ts` and record proof that version 1
   input becomes readable version 2 output.
6. [ ] Record all three validation results in this plan's `Work Log` so the
   settings-import plan owner can verify its stable `ready_when` condition and
   remove only the dependency blocker. Do not block independent intake
   investigation or activation.

# Work Log

- [x] [[2026-07-12]]@13:45 - Recorded the version 1 to version 2 migration
  goal and the requirement to preserve explicit values.
- [x] [[2026-07-12]]@14:25 - Activated the plan after choosing pure conversion
  plus atomic rename.
- [x] [[2026-07-12]]@14:35 - Added unit coverage for a repeated run; the second
  conversion produced unchanged version 2 JSON.
- [ ] [[2026-07-12]]@14:40 - Add interruption coverage and run the complete
  loader smoke check.

# Validation

- **Repeatable unit test:**
  - **Command or action:** run the migration unit test.
  - **Working directory:** fictional repository root.

    ```bash
    npm test -- tests/settings/migrate.test.ts
    ```

  - **Expected result:** the test runner reports that explicit values are
    preserved and the second run passes.
  - **Failure signals:** a failed assertion, changed JSON after the second run,
    or a nonzero exit code.
  - **Interpretation:** failure means migration is not yet safe to release
    settings-import implementation.

- **Test for an interrupted write:**
  - **Command or action:** run the interruption case only.
  - **Working directory:** fictional repository root.

    ```bash
    npm test -- tests/settings/migrate.test.ts -t "keeps original file before rename"
    ```

  - **Expected result:** the test reports one passed case and confirms the
    original settings file remains valid JSON.
  - **Failure signals:** the target file is missing, invalid JSON, or the test
    exits with a nonzero status.
  - **Interpretation:** failure means the atomic write implementation is not
    safe under interruption.

- **Loader smoke check:**
  - **Command or action:** load a version 1 fixture through the real loader.
  - **Working directory:** fictional repository root.

    ```bash
    npm run settings:smoke -- fixtures/settings/v1-explicit-values.json
    ```

  - **Expected result:** output contains `loaded schema v2` and
    `preserved explicit values: true`.
  - **Failure signals:** output contains `invalid settings`, reports version 1
    after load, or exits with a nonzero status.
  - **Interpretation:** success satisfies the loader smoke part of the import
    plan's `ready_when` condition; the full condition also requires passing
    tests for repeated runs and interrupted writes.

# Idempotence And Recovery

- **Repeated execution:** run the migration unit test twice and check that
  the second run leaves version 2 JSON unchanged.
- **Partial failure:** if interruption occurs before rename, delete only the
  temporary sibling file and rerun; keep the original settings file.
- **Backup:** copy `config/settings.json` to
  `config/settings.json.pre-v2-backup` before the first manual smoke run.
- **Safe fallback:** keep version 1 reading enabled until all three validation
  checks pass.

# Unfinished Work

- [ ] Implement atomic rename in `src/settings/load.ts`.
- [ ] Add and pass interruption coverage.
- [ ] Run the loader smoke check and record its output.
- [ ] Release settings-import implementation only after all three validation
  checks pass.
