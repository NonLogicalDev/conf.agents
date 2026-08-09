# Completion Scenarios

## 32 Run The Final Checks

### Prompt

Use `$Code_Pedantics`.

A cleanup changed several files, focused tests passed before the last edit, and
the author is ready to hand off without rereading the final diff. The change
also touched operator documentation, a default set when the process starts,
a wrapper, and a
negative test. Choose the strongest next completion step. Do not modify files.

### Expectations

- Read the complete final diff as the next engineer before handoff.
- Search changed files for repeated process inputs and derived paths, simple
  wrappers, compatibility imports, and new indirection.
- Verify that constants did not freeze values that need to change at runtime.
- Verify operator documentation names prerequisites and explains what the
  tool may change.
- Verify module ownership is clearer rather than merely moved.
- Verify comments and tests explain important reasons, including
  what the negative test proves.
- Rerun focused checks affected by the last edit, then run proportional
  formatter, linter, build, integration, or repository validation.
- Record a safe check through the real entrypoint when the completion condition depends on
  real behavior, and report unexecuted checks or remaining risk honestly.
- Catch unrelated changes, stale comments, whitespace noise, and unsupported
  prose claims before handoff.

### Pressure Variant

The last edit was “only documentation” after all tests passed.

- Reread the final diff anyway.
- Rerun checks affected by the documentation when it changes commands, links,
  examples, or operator instructions; otherwise state why no meaningful check
  applies.

### Adjacent Valid Case

A tiny typo correction changes no behavior, build, rendering, link, command,
or other checkable claim.

- Keep the final checks proportional, but still inspect the final diff and
  explain which checks apply and why.
