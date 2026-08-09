---
date: 2026-07-12
goal: settings-modernization
workstream: import-documentation
status: active
status_log:
  - "[[2026-07-12]]@16:30 created -> active"
tracker: DEMO-42
subject: document-settings-import
---

# Purpose

Document the accepted settings import format for existing fictional Linear issue `DEMO-42`.

# Scope

- **In scope:** the supported file format, accepted keys, and conflict behavior.
- **Out of scope:** creating a Linear issue, changing the importer, and changing the local task and epic numbers.

# Completion Criteria

- [x] The plan identifies existing fictional issue `DEMO-42`.
- [ ] `docs/settings-import.md` describes the accepted file format.
- [ ] The document explains how imported settings handle conflicts.
- [ ] The documented behavior matches the existing importer.

# Blockers

- None.

# Artifacts

- **Artifact:** [DEMO-42](https://example.com/linear/DEMO-42).
  - **Type:** Fictional Linear issue.
  - **State:** Existing; not created by this plan.
  - **Contents:** The request to document settings import.
  - **Location:** `https://example.com/linear/DEMO-42`.
- **Artifact:** `docs/settings-import.md`.
  - **Type:** Project documentation.
  - **State:** Planned.
  - **Contents:** The accepted settings format and conflict behavior.
  - **Location:** `docs/settings-import.md`.

# Decisions

- **Decision:** Name the group after the existing fictional issue.
  - **Rationale:** `DEMO-42` already identifies this work. Its tracker number is separate from the local task and epic sequence.

# Implementation Steps

- [x] Record the existing fictional issue and its requested scope.
- [ ] Inspect the supported import format and conflict behavior.
- [ ] Update `docs/settings-import.md`.
- [ ] Check the documented behavior against the existing importer.

# Work Log

- [x] [[2026-07-12]]@16:30 - Recorded existing fictional issue `DEMO-42` and opened the documentation plan without changing local group numbers.

# Validation

- **Command:** `rg -n 'import|format|conflict' docs/settings-import.md`.
- **Expected result:** The document describes supported settings and conflict behavior.
- **Observed result:** Not run.
- **Interpretation:** Documentation remains unfinished until the check passes and the behavior matches the existing importer.

# Unfinished Work

- [ ] Inspect the existing importer.
- [ ] Document the accepted format and conflict behavior.
- [ ] Run the documentation check.
