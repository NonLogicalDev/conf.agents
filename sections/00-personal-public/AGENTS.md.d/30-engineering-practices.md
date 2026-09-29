## Engineering practices

- **File moves:** After moving files, remove old directories only if they are empty. Preserve unrelated files and their directories.
- **Checkpoints:** When editing `agent-config`, use [$@:VCS_Checkpoint](skill://@:VCS_Checkpoint) after meaningful changes and before risky work or a handoff. Include the task's changes. Keep commits local unless the user asks to push.
- **Small, reviewable changes:** When preparing a change, splitting a pull request, or reducing reviewer effort, read `{{%_resources_%}}/engineering-practices/small-changes.md`.
