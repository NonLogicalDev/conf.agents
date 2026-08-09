---
date: 2026-07-12
task: task-001-settings-import
subject: settings-import
updated: "[[2026-07-12]]@16:20"
---

# Settings Import

## Executive Summary

Users need a supported way to bring external settings into the application
without editing `config/settings.json` by hand. Import work is intentionally
split between independent intake and implementation that waits for its dependency. Intake
may resolve the supported format, conflict behavior, failure behavior, and CLI
owner now. Implementation begins only after the
[[epic-002-settings-migration/active/plan-01 2026-07-12 migrate-settings.md|migration plan]]
records passing tests for repeated runs and interrupted writes, plus loader smoke
output showing that schema version 2 loaded and preserved explicit values.

## Purpose

Deliver supported settings import that is safe to repeat after schema migration
proves the destination format can recover from interruptions.

## Scope

- **In scope:** supported import inputs, validation, writes that are safe to
  repeat, and errors shown to the user.
- **Out of scope:** migrating existing stored settings and publishing the
  cross-epic validation report.
- **Owned surfaces:** import commands, import parsing, and import fixtures.

## Plan Placement

- **Belongs here when:** a plan changes how external settings enter the
  application or how import failures recover.
- **Does not belong here when:** a plan changes stored-schema migration or
  general validation-report publication.
- **New epic or split required when:** the work introduces another durable data
  ingestion product or an independently owned migration program.

## Relationships

- **Goal:** `settings-modernization`
- **Related tasks and epics:** [[epic-002-settings-migration/EXEC_STATE.md]]
  owns migration of existing data; [[task-003-validation-report/EXEC_STATE.md]]
  owns reusable proof publication.

## Produced External Artifacts

- None produced.
