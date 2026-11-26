---
name: VCS_Checkpoint
description: Save a local checkpoint with the project's usual version control tool. Include the task's pending changes and write a clear checkpoint message.
---

# Source control checkpoint

Use this skill when the user asks for a checkpoint or the task calls for one.

## Save a checkpoint

- Read the project's instructions.
- Use its usual version control tool and commit command.
- Include changes that belong to the task, including new, staged, and unstaged files.
- Leave unrelated changes alone. If work from different tasks is mixed, ask before changing or committing anything.
- Follow the project's usual commit hooks.
- Verify the checkpoint's subject and report what it saved.

Do not switch branches, rewrite history, push, or publish unless the user asks. If there are no changes, do not create an empty commit unless asked.

## Write the subject

- Prefer `git-checkpoint-save` on `PATH` when `command -v git-checkpoint-save` finds it. Otherwise, run the [bundled copy](scripts/git-checkpoint-save).
- Pass `--role agent` explicitly with either copy. If the user requests another role, pass that role instead. Pass a supplied note as one quoted argument.
- For ordinary Git, stage the task’s changes and verify the index contains only those changes before running the command. It commits what is staged.
- When the project requires a different commit command, pass `--subject-only` and use its output with the project’s command.
- With a note, the subject is `checkpoint:agent - <short note>`. Without one, it uses local time in `YYYY-DD-MMThh-mm-ss` order, for example `checkpoint:agent - 2031-17-02T14-05-06`.

## Tests

When changing this skill, read [tests/README.md](tests/README.md). Run the relevant scenarios with fresh subagents that have empty context windows.
