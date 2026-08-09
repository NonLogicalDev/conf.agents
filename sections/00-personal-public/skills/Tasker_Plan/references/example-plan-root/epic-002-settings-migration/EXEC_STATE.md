---
date: 2026-07-12
epic: epic-002-settings-migration
subject: settings-migration
updated: "[[2026-07-12]]@16:20"
---

# Settings Migration

## Executive Summary

Move the application's version 1 settings to schema version 2 without losing explicit user values. An earlier attempt to remove the old reader first was replaced by an additive migration that remains safe after an interruption. The active direction writes the new representation before retiring the old reader, uses the [[task-003-validation-report/completed/plan-01 2026-07-12 publish-report.md|completed validation report]] as lasting proof. Conversion is proven safe to repeat. The migration permits [[task-001-settings-import/backlog/plan-01 2026-07-12 import-settings.md|settings import]] implementation only after tests for repeated runs and interrupted writes pass and the loader smoke check shows that schema version 2 loaded and preserved explicit values. Settings import intake can begin before that evidence is ready. Interruption recovery and the loader smoke check are the remaining gaps.

## Purpose

Make schema version 2 the safe stored-settings format while preserving explicit values and recoverability.

## Scope

- **In scope:** schema conversion, additive writes, loader compatibility, interruption recovery, and migration smoke validation.
- **Out of scope:** external settings import and general report tooling.
- **Owned surfaces:** settings loader, schema conversion, migration fixtures, and rollback behavior.

## Plan Placement

- **Belongs here when:** a plan changes the conversion of existing data, compatibility, rollback, or migration proof.
- **Does not belong here when:** a plan only imports new settings or packages evidence already produced.
- **New epic or split required when:** work becomes a separate storage product, independent migration program, or disconnected compatibility domain.

## Relationships

- **Goal:** `settings-modernization`
- **Related tasks and epics:** [[task-001-settings-import/EXEC_STATE.md]] consumes readiness evidence; [[task-003-validation-report/EXEC_STATE.md]] publishes reusable proof; [[epic-004-legacy-config/EXEC_STATE.md]] preserves the earlier approach that removed the old reader first.

## Produced External Artifacts

- **Pull request:** [Add interruption-safe settings migration](https://github.com/example/settings-app/pull/42) — open.
- **Linear ticket:** [SET-42: Migrate stored settings safely](https://linear.example.com/example/issue/SET-42) — in progress.
- **Document:** [Settings migration design](https://docs.example.com/settings-migration-design) — published.
- **Slack thread started:** [Settings migration recovery discussion](https://chat.example.com/archives/C0123456789/p1720801200000000) — started.
