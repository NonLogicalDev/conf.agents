# Review Scenarios

## 01 Search Before Removing A Wrapper

### Prompt

Use `$Code_Pedantics`.

A review shows a Python function that forwards the same arguments and return
value to another function. The user asks whether to remove it. Choose the next
concrete review plan. Do not modify files.

### Expectations

- Search callers, imports, and tests before deciding the wrapper is removable.
- Keep it when it preserves compatibility, applies a policy, adds
  instrumentation, or supplies a dependency.
- Remove it only when evidence shows it adds no contract.
- Put the search before the verdict.

### Pressure Variant

The wrapper is one line and looks obviously redundant.

- Do not skip the caller and compatibility search.

## 08 Preserve Deliberate Auditability

### Prompt

Use `$Code_Pedantics`.

Two explicit branches look repetitive, but each branch uses a different static
type and has different lifetime and failure behavior. A reviewer asks to
collapse them and delete the nearby comments because the lines look similar.
Choose the review result. Do not modify files.

### Expectations

- Check the actual contracts before treating repetition or comments as
  duplicates.
- Preserve explicit structure when it honestly exposes different type,
  lifetime, permission, side effect, or failure behavior.
- Keep or tighten only a comment that explains an important difference.
- Remove only genuinely duplicate, stale, or misleading comments.
- Do not optimize for line count or unrelated cleanup.

## 30 Separate Process Inputs From Runtime Values

### Prompt

Use `$Code_Pedantics`.

A command that runs for an extended time repeatedly reads a fixed `OUTPUT_ROOT` environment
variable, creates a credential after startup, and refreshes a lease while it
runs. A reviewer proposes moving all three values into constants for the module
to remove repetition. Choose the review result. Do not modify files.

### Expectations

- Load the immutable process input once into a named value and use it
  consistently.
- Preserve the credential as a value created at runtime and pass it explicitly
  to the functions that need it.
- Preserve the lease behind its refresh mechanism instead of freezing it at
  process start.
- Trace each value from definition to use before changing its lifetime or
  owner.
- Do not keep repeated `OUTPUT_ROOT` reads solely because tests replace the
  environment; test the named value set when the process starts instead.
- Do not turn a timestamp, credential, lease, refreshed configuration, or
  other dynamic value into a static constant for cosmetic consistency.

### Adjacent Valid Case

A configuration value is intentionally reloaded while the process runs.

- Keep it behind the explicit reload mechanism and document or name that
  runtime behavior instead of hoisting it into an immutable process input.
