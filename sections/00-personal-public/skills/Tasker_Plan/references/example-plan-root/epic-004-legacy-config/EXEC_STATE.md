---
date: 2026-07-12
epic: epic-004-legacy-config
subject: legacy-config
updated: "[[2026-07-12]]@16:05"
---

# Legacy Configuration

## Executive Summary

The original approach removed the legacy settings reader after replacement
behavior appeared ready. Migration work later showed that removing the old
reader first weakened interruption recovery, so the
[[epic-002-settings-migration/active/plan-01 2026-07-12 migrate-settings.md|additive schema-migration plan]]
obsoleted this direction. The epic is retained as a contained historical story
that explains why direct cleanup is no longer the active strategy,
but it should not receive new migration implementation plans.

## Purpose

Preserve the retired cleanup of the legacy reader and the reason it was
superseded.

## Scope

- **In scope:** historical explanation of the cleanup that removed the old
  reader first and any
  evidence required to interpret its archived plan.
- **Out of scope:** active schema migration, import, or new compatibility work.
- **Owned surfaces:** the archived plan for removing the legacy reader.

## Plan Placement

- **Belongs here when:** a plan only repairs or clarifies the retained
  historical artifact.
- **Does not belong here when:** a plan advances current migration or cleanup
  implementation.
- **New epic or split required when:** the user permits renewed legacy work as
  a deliverable instead of historical maintenance.

## Relationships

- **Goal:** `settings-modernization`
- **Related tasks and epics:** [[epic-002-settings-migration/EXEC_STATE.md]]
  owns the active replacement strategy.

## Produced External Artifacts

- **Pull request:** [Retire the legacy settings reader](https://github.com/example/settings-app/pull/37)
  — closed without merging.
