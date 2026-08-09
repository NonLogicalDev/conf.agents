# TypeScript And JavaScript Scenarios

## 26 Validate TypeScript Data And Ownership

### Prompt

Use `$Code_Pedantics`.

A React entrypoint fetches JSON, casts it directly to `Widget[]`, and passes the result through a helper typed with `any`. The helper catches every error and returns an empty array. During the same extraction, a dialog's open state and temporary input move into a global store even though no other component uses them. Choose the review result. Do not modify files.

### Expectations

- Parse and validate external JSON before it enters typed domain code; do not let `as Widget[]`, `any`, or another broad cast hide unvalidated input.
- Keep errors narrow enough to identify the violated field or contract instead of swallowing every failure into an indistinguishable empty result.
- Keep temporary dialog state and DOM values local to the component that owns them.
- Move state into a typed shared store only when multiple components need the same state or transitions.
- Keep the entrypoint thin and extract one safe component at a time, then run formatter, type checker, focused tests, and the real browser flow.

### Adjacent Valid Case

The data already comes from an internal validated domain function, and two components genuinely share one transition.

- Accept the typed value without validating it again at every consumer.
- Move the shared transition into a typed store while leaving unrelated local widget state in its owning component.
