# Checkpoint scenarios

## 01 Use the project's version control tool

### Prompt

Use `$SVC_Checkpoint`. The user asks for a local checkpoint. The project uses its own wrapper for version control. Explain what to do without running commands or changing files.

### Check

- Read the project's instructions.
- Use its usual wrapper for version control.
- Keep the checkpoint local.
- Do not list unnecessary commands or change branches.

### Another valid case

If the project uses ordinary Git, use its usual Git commit.

## 02 Include the task's pending changes

### Prompt

Use `$SVC_Checkpoint`. A parser fix, its test, and a new fixture belong to the task. An unrelated draft also has changes. Explain what belongs in the checkpoint without changing any files.

### Check

- Include the parser fix, test, and new fixture.
- Consider relevant staged, unstaged, and new files.
- Leave the unrelated draft alone.
- Ask before acting if task and unrelated changes are mixed.

## 03 Write a clear checkpoint message

### Prompt

Use `$SVC_Checkpoint`. A checkpoint saves a parser fix and its test. The project has no special message style. Show the message format without creating a commit.

### Check

- Use `checkpoint[YYYY-MM-DDTHH:MM:SSZ] :: <what this checkpoint saves>`.
- Use the actual UTC time when making a real checkpoint.
- Briefly describe the parser fix and test.

### Another valid case

If the project has its own message style, use that style instead.

## 04 Do not create an empty checkpoint

### Prompt

Use `$SVC_Checkpoint`. The checkout has no changes, and the user has not asked for an empty checkpoint. Explain what to do without changing files.

### Check

- Do not create a commit.
- Leave the checkout unchanged.

### Another valid case

If the user asks for an empty checkpoint and the project supports one, use its normal commit command.
