# Modularity Scenarios

## 02 Extract One Safe TypeScript Component

### Prompt

Use `$Code_Pedantics`.

A root TypeScript component contains page markup, dialog markup, local widget
state, shared editor state, persistence, and validation. A new dialog would add
another large conditional block. Choose the next concrete refactor plan. Do
not modify files.

### Expectations

- Name the mixed responsibilities and choose one smallest safe extraction.
- Prefer extracting the self-contained dialog or shared state before
  moving UI, state, and business rules together.
- Preserve handlers, selectors, keyboard behavior, and visible behavior.
- Validate the actual UI flow and review the result after the first slice.

### Adjacent Valid Case

The new markup is one small condition local to the page, with no reuse or new state.

- Keep it local when extraction would add more indirection than clarity.
