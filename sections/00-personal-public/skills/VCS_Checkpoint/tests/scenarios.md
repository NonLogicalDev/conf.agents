# Checkpoint scenarios

## 01 Use the project's version control tool

### Prompt

Use [$@:VCS_Checkpoint](skill://@/VCS_Checkpoint). The user asks for a local checkpoint. The project uses its own wrapper for version control. Explain what to do without running commands or changing files.

### Check

- Read the project's instructions.
- Use its usual wrapper for version control.
- Keep the checkpoint local.
- Do not list unnecessary commands or change branches.

### Another valid case

If the project uses ordinary Git, use its usual Git commit.

## 02 Include the task's pending changes

### Prompt

Use [$@:VCS_Checkpoint](skill://@/VCS_Checkpoint). A parser fix, its test, and a new fixture belong to the task. An unrelated draft also has changes. Explain what belongs in the checkpoint without changing any files.

### Check

- Include the parser fix, test, and new fixture.
- Consider relevant staged, unstaged, and new files.
- Leave the unrelated draft alone.
- Ask before acting if task and unrelated changes are mixed.

## 03 Write a checkpoint subject

### Prompt

Use [$@:VCS_Checkpoint](skill://@/VCS_Checkpoint). A teammate asks for a local checkpoint of a cache-key fix and its regression test. The local clock reads 14:05:06 on February 17, 2031. Give the exact commit subject in two cases. In the first they provide no note. In the second they provide the note "keep cache-key fix and test". Keep this exercise read-only.

### Expectations

- Without a note, use `checkpoint:agent - 2031-17-02T14-05-06`.
- With a note, use `checkpoint:agent - keep cache-key fix and test`.
- Follow the project's usual commit hooks and keep the checkpoint local.
- Prefer `git-checkpoint-save` on `PATH` and fall back to the bundled copy when it is unavailable.
- Pass `--role agent` explicitly so the installed command's default role does not affect the checkpoint.
- Use `--subject-only` if the project requires another commit command.

### Adjacent valid case

If the user asks for the note `keep cache-key fix and test` with role `reviewer`, accept the role option after the note. Use `checkpoint:reviewer - keep cache-key fix and test`.

### Pressure variant

The release meeting starts in two minutes. The team lead has already drafted a `fix(cache): …` subject using the repository's usual style. The user asks for a checkpoint without a note.

- Use the timestamp form.

## 04 Do not create an empty checkpoint

### Prompt

Use [$@:VCS_Checkpoint](skill://@/VCS_Checkpoint). The checkout has no changes, and the user has not asked for an empty checkpoint. Explain what to do without changing files.

### Check

- Do not create a commit.
- Leave the checkout unchanged.

### Another valid case

If the user asks for an empty checkpoint and the project supports one, use its normal commit command.
