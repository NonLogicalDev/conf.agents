## Engineering practices

Use this glossary to select the engineering practices relevant to the task.

- **File moves:** Remove old directories after moving their contents, but
  only when they are empty. Preserve unrelated files and their directories.
- **Frequent checkpoints:** When editing `agent-config`, use
  [$SVC_Checkpoint]({{%_skills_%}}/SVC_Checkpoint/SKILL.md) after meaningful changes and before risky work or a
  handoff. Include the task's changes and keep commits local unless the
  user asks to push.
- **Small, reviewable changes:** When preparing a change, splitting a pull
  request, or reducing reviewer effort, read
  `{{%_resources_%}}/engineering-practices/small-changes.md`.
