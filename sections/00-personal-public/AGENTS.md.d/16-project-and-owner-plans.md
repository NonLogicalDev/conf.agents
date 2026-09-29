## Project and owner plans

Use `<agent-home>` for the agent home configured by the user or project. Use these plan roots unless the user or project specifies others:

```text
<agent-home>/
├── plans-active/
│   └── <project-or-owner-name>/
│       └── _owner/
└── plans-archive/
    └── <project-or-owner-name>/
        └── _owner/
```

`<plan-root>/_owner/` is the optional folder for owner support files.

- Use [$@:Tasker_Stream](skill://@:Tasker_Stream) when one owner coordinates subagents or maintains shared project context. Follow that skill for owner plans, workstreams, structure, and execution.
- Name a new owner `<project>__YYYY-qN__<slug>` using the verified project and the year and quarter when the work began. Keep an existing owner's name and history when the quarter changes.
- When a group represents an issue:
  - Use a verified existing ticket's actual team key and number.
  - Do not invent or create a ticket just to name the group.
  - Keep ticket numbers separate from the local group sequence; they do not consume a number, restart the sequence, or change it.
- Before changing shared work, follow the user's current instructions, the project's `AGENTS.md`, and any existing owner or plan.
- Preserve the existing owner, active plans, plan history, and useful standing guidance. Without the user's permission, do not take over, overwrite, move, rename, duplicate, archive, or discard another owner's work.
- Make sure the owner and affected workers follow current project guidance when:
  - Work begins or resumes.
  - Relevant instructions change.
  - The work depends on that guidance.
- Improve owner guidance or helpers within the user's intent, accepted scope, and existing permissions. Record and verify meaningful changes.
