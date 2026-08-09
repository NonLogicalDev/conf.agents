---
date: 2026-07-12
goal: settings-modernization
workstream: settings-import
status: backlog
status_log:
  - "[[2026-07-12]]@13:10 created -> backlog"
depends_on:
  - plan: "[[epic-002-settings-migration/active/plan-01 2026-07-12 migrate-settings.md]]"
    ready_when: "The migration plan Work Log records passing repeat-safe and interrupted-write tests plus loader smoke output containing loaded schema v2 and preserved explicit values: true."
subject: import-settings
---

# Purpose

Add settings import after schema migration proves that version 1 settings can be loaded safely as version 2, so users have a supported alternative to editing `config/settings.json` manually.

# Scope

- **In scope:** resolve the accepted import format, conflict behavior, failure behavior, and CLI owner; implement the supported settings-import command; add focused CLI and loader validation; and preserve the migration plan's readiness dependency.
- **Out of scope:** changing schema migration, inventing an undocumented import format, building settings UI, changing unrelated CLI commands, and broad settings-file cleanup.

# Completion Criteria

- [ ] The dependency plan's `Work Log` records passing tests for repeated runs and interrupted writes, plus loader smoke output containing `loaded schema v2` and `preserved explicit values: true`.
- [ ] `Context` records the accepted import format, conflict behavior, invalid file behavior, and owning CLI module from the exact investigations below.
- [ ] A focused CLI test proves one documented fixture imports successfully and the real loader can read the resulting settings file.
- [ ] A focused failure test proves invalid input returns the repository's existing invalid-file error and exit code without changing the current settings file.
- [ ] `Validation` records the exact commands, expected results, failure signals, and interpretation for the accepted and rejected import paths.

# Blockers

- **Blocked work:** Begin settings-import implementation after its intake unknowns are resolved.
  - **Blocking condition:** The linked schema-migration plan has not recorded passing tests for repeated runs and interrupted writes, plus loader smoke output containing `loaded schema v2` and `preserved explicit values: true`.
  - **Unblock owner:** The schema-migration plan owner produces the evidence in [[epic-002-settings-migration/active/plan-01 2026-07-12 migrate-settings.md]]; the settings-import plan owner verifies it and clears this entry.
  - **Unblock action:** Complete `Implementation Steps` 3–5 in the linked schema-migration plan: write version 2 through a temporary sibling and rename in `src/settings/load.ts`, then add coverage for interrupted writes that stops before rename and proves the original file remains readable. From the fictional repository root, run `npm test -- tests/settings/migrate.test.ts`, then run `npm run settings:smoke -- fixtures/settings/v1-explicit-values.json`; record the passing cases for repeated runs and interrupted writes, plus both expected lines of smoke output in the migration plan's `Work Log`. The settings-import plan owner then verifies those entries and, in the same maintenance pass, records a clearing receipt in this plan's `Work Log`, removes only this dependency blocker, reviews all remaining blockers, and updates the current next action to match `Resume with`.
  - **Unblock evidence:** The schema-migration plan's `Work Log` contains the passing results for repeated runs and interrupted writes, plus `loaded schema v2` and `preserved explicit values: true`, satisfying this plan's `ready_when` condition.
  - **Resume with:** Continue the next unresolved `Next Investigation`. After every intake unknown resolves, add exact `Implementation Steps` and `Validation`; begin implementation only when no blocker prevents it.

- **Blocked work:** Implement the supported import-file contract.
  - **Blocking condition:** The accepted import format and required keys are not yet known.
  - **Unblock owner:** Settings-import plan owner.
  - **Unblock action:** Read `docs/settings-import.md` and compare `fixtures/settings/v1-explicit-values.json` with the other files under `fixtures/settings/`. In the same maintenance pass, record the documented extension, required keys at the top level, and one accepted fixture with exact paths in `Context`; record a clearing receipt in `Work Log`; remove the matching `Unknowns` item and this blocker when that record exists; review the remaining blockers; and update the current next action to match `Resume with`.
  - **Unblock evidence:** `Context` records the documented extension, required keys at the top level, and one accepted fixture with exact paths.
  - **Resume with:** Continue the next unresolved `Next Investigation`. Move this plan to `active/` when that investigation requires several steps. After every intake unknown resolves, add exact `Implementation Steps` and `Validation`; begin implementation only when no blocker prevents it.

- **Blocked work:** Implement import conflict behavior.
  - **Blocking condition:** The plan does not yet know whether imported values replace, merge with, or reject existing settings.
  - **Unblock owner:** Settings-import plan owner.
  - **Unblock action:** Inspect `src/settings/write.ts` and `tests/settings/write.test.ts`. In the same maintenance pass, record the verified merge, replace, or rejection behavior with exact source paths in `Context`; update `Completion Criteria` when the behavior changes its proof; record a clearing receipt in `Work Log`; remove the matching `Unknowns` item and this blocker when that record exists; review the remaining blockers; and update the current next action to match `Resume with`.
  - **Unblock evidence:** `Context` records the verified merge, replace, or rejection behavior with exact source paths.
  - **Resume with:** Continue the next unresolved `Next Investigation`. Move this plan to `active/` when that investigation requires several steps. After every intake unknown resolves, add exact `Implementation Steps` and `Validation`; begin implementation only when no blocker prevents it.

- **Blocked work:** Implement invalid-file behavior.
  - **Blocking condition:** The error shown to the user and exit code are not yet known.
  - **Unblock owner:** Settings-import plan owner.
  - **Unblock action:** Inspect `src/cli/errors.ts` and `tests/cli/errors.test.ts`. In the same maintenance pass, record the error shown to the user and exit code with exact source paths in `Context`; record a clearing receipt in `Work Log`; remove the matching `Unknowns` item and this blocker when that record exists; review the remaining blockers; and update the current next action to match `Resume with`.
  - **Unblock evidence:** `Context` records the existing invalid-file error shape and exit code with exact source paths.
  - **Resume with:** Continue the next unresolved `Next Investigation`. Move this plan to `active/` when that investigation requires several steps. After every intake unknown resolves, add exact `Implementation Steps` and `Validation`; begin implementation only when no blocker prevents it.

- **Blocked work:** Implement the settings-import command in its owning CLI module and add focused tests.
  - **Blocking condition:** The owning CLI module and nearest test file are not yet known.
  - **Unblock owner:** Settings-import plan owner.
  - **Unblock action:** From the fictional repository root, run `rg -n "settings|registerCommand" src/cli tests/cli`. In the same maintenance pass, record the exact command module and nearest focused test file in `Context` and `Artifacts`; record a clearing receipt in `Work Log`; remove the matching `Unknowns` item and this blocker when those records exist; review the remaining blockers; and update the current next action to match `Resume with`.
  - **Unblock evidence:** `Context` and `Artifacts` record the exact command module and nearest focused test file.
  - **Resume with:** Continue the next unresolved `Next Investigation`. Move this plan to `active/` when that investigation requires several steps. After every intake unknown resolves, add exact `Implementation Steps` and `Validation`; begin implementation only when no blocker prevents it.

The focused actions in `Next Investigation` remain available while settings-import implementation is blocked.

# Artifacts

- **Artifact:** [example/settings-app](https://example.com/scm/settings-app)
  - **Type:** Git repository baseline.
  - **State:** Verified orientation baseline; no feature branch or pull request exists.
  - **Contents:** Repository state used to resolve the settings-import command owner and validation layout.
  - **Location:** [example/settings-app](https://example.com/scm/settings-app).
  - **SCM surface:** remote `origin`; ref `refs/remotes/origin/main`; full commit `bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb`.
  - **Reconstruct:** clone [example/settings-app](https://example.com/scm/settings-app), fetch `origin`, then check out `bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb`.
- **Artifact:** Settings-import command.
  - **Type:** Planned local repository source artifact.
  - **State:** Planned; exact path unresolved.
  - **Contents:** Supported command for importing a documented settings file without manual `config/settings.json` editing.
  - **Location:** Unresolved until `Next Investigation` identifies the owning CLI module.
- **Artifact:** Focused import tests.
  - **Type:** Planned local repository test artifacts.
  - **State:** Planned; exact paths unresolved.
  - **Contents:** Proof of accepted imports, rejected invalid input, and the actual loader.
  - **Location:** Unresolved until `Next Investigation` identifies the existing test layout.
- **Artifact:** Upstream branch or pull request.
  - **Type:** Deliverable tracked in source control.
  - **State:** None exists while this plan remains in backlog.
  - **Contents:** No implementation surface has been created.
  - **Location:** Not available until a branch or pull request is created.

# Context

## Fictional Repository

- **Candidate CLI:** `src/cli/settings.ts` may own a future `settings import` command.
- **Settings loader:** `src/settings/load.ts` is the likely integration point after migration is proven.
- **Fixture directory:** `fixtures/settings/` contains example settings files.

## Current Behavior

- Users can edit `config/settings.json` manually.
- No supported import command exists.
- The accepted import file format and error behavior are not yet decided.

## Dependency

- **Depends on:** [[epic-002-settings-migration/active/plan-01 2026-07-12 migrate-settings.md]]
- **Ready when:** the migration plan's `Work Log` records passing tests for repeated runs and interrupted writes, plus loader smoke output containing `loaded schema v2` and `preserved explicit values: true`.
- **Reason:** do not let import write files that the actual loader cannot read safely.

# Steering

- [[2026-07-12]]@13:10 - Keep settings import behind proven migration behavior.
  - **Effect:** settings-import implementation remains blocked until its `depends_on` condition is satisfied; the plan may activate for independent intake investigation.
- [[2026-07-12]]@13:12 - Do not invent an import format before product intent is known.
  - **Effect:** preserve the missing decision as an explicit unknown.

# Unknowns

- **Accepted format:** whether import accepts JSON only or another documented format.
- **Conflict behavior:** whether imported values replace, merge with, or reject existing settings.
- **Failure behavior:** which error the user sees when the imported file is invalid.
- **CLI ownership:** whether `src/cli/settings.ts` already owns settings subcommands or another module does.

# Next Investigation

- **Resolve accepted format:**
  - **Action:** read `docs/settings-import.md` and compare `fixtures/settings/v1-explicit-values.json` with other files under `fixtures/settings/`.
  - **Expected result:** identify the documented extension, required keys at the top keys, and one accepted fixture.
- **Resolve conflict behavior:**
  - **Action:** inspect `src/settings/write.ts` and `tests/settings/write.test.ts` for current merge or replace behavior.
  - **Expected result:** record whether existing writes preserve unspecified keys or replace the complete settings object.
  - **Fallback:** if the files do not define import conflict behavior, preserve the current `Blocked work` and `Blocking condition`, then replace its normal route fields with these provisional fields:
    - **Route gap:** the product decision owner, request destination, and clearing evidence are unknown.
    - **Route-discovery owner:** Settings-import plan owner.
    - **Route-discovery action:** Ask the user: "Which owner decides settings-import conflict behavior, where should that decision be requested, and what recorded result proves it is accepted?"
    - **Route-discovery evidence:** The plan records the verified decision owner, exact request and destination, observable clearing result, and the source for each fact.
    - **Resume with:** Replace the provisional entry with the verified normal unblock route, then take its `Unblock action`.
- **Resolve failure behavior:**
  - **Action:** inspect `src/cli/errors.ts` and `tests/cli/errors.test.ts`.
  - **Expected result:** identify the existing invalid-file error shape and exit code.
- **Resolve CLI ownership:**
  - **Action:** search command registration and matching tests.
  - **Working directory:** fictional repository root.

    ```bash
    rg -n "settings|registerCommand" src/cli tests/cli
    ```

  - **Expected result:** identify the exact settings command module and its nearest CLI test file.

This intake plan intentionally omits `Implementation Steps` and `Validation` until the intake unknowns are resolved. Add both sections as soon as those unknowns resolve even if the dependency remains unmet; the dependency continues to block implementation, not plan elaboration.

# Work Log

- [ ] [[2026-07-12]]@13:10 - Wait for evidence from the migration plan that covers repeated runs, interrupted writes, and the loader smoke check.
- [ ] [[2026-07-12]]@13:12 - Resolve import format, conflict behavior, failure behavior, and CLI ownership before drafting implementation.

# Unfinished Work

- [ ] Move this plan to `active/` before an investigation that takes several steps.
- [ ] Resolve every intake unknown.
- [ ] Add exact `Implementation Steps` and `Validation` as soon as the unknowns resolve.
- [ ] Verify all three readiness results in the linked migration plan's `Work Log` before implementation.
