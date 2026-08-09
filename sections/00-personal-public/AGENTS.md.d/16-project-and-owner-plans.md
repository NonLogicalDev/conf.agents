## Project and owner plans

- Use `<agent-home>` for the user- or project-configured agent home. Unless the user or project specifies otherwise, use the following plan roots. `<plan-root>/_owner/` is the optional owner support folder.

```text
<agent-home>/
├── plans-active/
│   └── <project-or-owner-name>/
│       └── _owner/
└── plans-archive/
    └── <project-or-owner-name>/
        └── _owner/
```

- Use [$Tasker_Plunger]({{%_skills_%}}/Tasker_Plunger/SKILL.md) when one owner coordinates subagents or maintains shared project context. Follow that skill for owner plans, workstreams, structure, and execution.
- Use [$Tasker_Plan]({{%_skills_%}}/Tasker_Plan/SKILL.md) for a standalone plan that does not need owner coordination.
- Name a new owner `<project>__YYYY-qN__<slug>` using the verified project and the year and quarter when its work began. Preserve an existing owner name and history when the quarter changes.
- Base an issue-backed group only on a verified existing ticket's actual team key and number. Do not invent or create a ticket for a group name; ticket numbers do not consume, restart, or change the local group sequence.
- Follow the user's current instructions, the actual project's `AGENTS.md`, and any existing owner or plan before changing shared work.
- Preserve the existing owner, active plans, plan history, and useful standing guidance. Do not take over, overwrite, move, rename, duplicate, archive, or discard another owner's work without the user's permission.
- Keep the owner and affected workers aligned with current project guidance when work begins or resumes, relevant instructions change, or the work depends on that guidance.
- Improve owner guidance or helpers only within the user's intent, accepted scope, and existing permissions. Record and verify meaningful changes.
