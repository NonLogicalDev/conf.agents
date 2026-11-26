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

## 03 Write a clear checkpoint subject

### Prompt

Use [$@:VCS_Checkpoint](skill://@/VCS_Checkpoint). Save a local checkpoint of the cache-key fix and its regression test. The repository typically uses Conventional Commits. Show the exact subject and say what you would save. Keep this exercise read-only.

### Expectations

- Use `agent:checkpoint - <short note>` for the subject, with the exact prefix and separator.
- Briefly describe the cache-key fix and test.
- Follow the project's usual commit hooks and keep the checkpoint local.

### Pressure variant

The release meeting starts in two minutes. The team lead has already drafted a `fix(cache): …` subject in the repository's usual style and wants the fix saved now.

- Keep the checkpoint subject format and use the draft as context for the short note.

## 04 Do not create an empty checkpoint

### Prompt

Use [$@:VCS_Checkpoint](skill://@/VCS_Checkpoint). The checkout has no changes, and the user has not asked for an empty checkpoint. Explain what to do without changing files.

### Check

- Do not create a commit.
- Leave the checkout unchanged.

### Another valid case

If the user asks for an empty checkpoint and the project supports one, use its normal commit command.
