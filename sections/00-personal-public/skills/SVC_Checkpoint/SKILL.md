---
name: SVC_Checkpoint
description: Save a local checkpoint with the project's usual version control tool. Include the task's pending changes and write a clear checkpoint message.
---

# Source control checkpoint

Use this skill when the user asks for a checkpoint or the task calls for one.

## Save a checkpoint

- Read the project's instructions.
- Use its usual version control tool and commit command.
- Include changes that belong to the task, including new, staged, and unstaged files.
- Leave unrelated changes alone. If work from different tasks is mixed, ask before changing or committing anything.
- Follow the project's usual commit hooks and message style.
- Check that the commit was created and say what it saved.

Do not switch branches, rewrite history, push, or publish unless the user asks. If there are no changes, do not create an empty commit unless asked.

## Write the message

Use the project's usual style for commit messages if it has one. Otherwise use:

```text
checkpoint[YYYY-MM-DDTHH:MM:SSZ] :: <what this checkpoint saves>
```

Use the actual UTC time and describe the changes in a few words:

```text
checkpoint[2026-07-27T18:42:00Z] :: save parser fix and regression test
```

## Tests

Read [tests/README.md](tests/README.md) when changing this skill.
