# Modularity And Ownership

## Keep Ownership Clear

Each module should have a clear owner and responsibility. A file should not keep absorbing unrelated page markup, interaction details, shared state, business mutations, validation, persistence, and side effects just because it already has access to them.

Prefer modules with one reason to change:

- page or screen modules own layout and composition for that page;
- components own reusable or self-contained interactions;
- stores own shared state shape and state transitions;
- domain libraries own data transformations, validation, and rules to preserve; and
- entrypoints and roots own orchestration, routing, dependency wiring, and setup shared across components only.

A new module is useful when it clarifies ownership, not merely when it creates another file.

## When To Act

Act before adding responsibility to a file that is already too large or combines unrelated work.

Typical triggers:

- one component or entrypoint is becoming the whole application;
- new markup adds another large conditional section to an existing file;
- render helpers are replacing named components;
- local state blocks dominate the top of a file;
- prop chains exist only because shared state has no coherent home;
- one feature requires state, UI, persistence, validation, and exports in the same file; or
- a refactor starts with “just one more helper” in a file that is already hard to scan.

Do not wait until the file is impossible to review. Do not extract solely to reach an arbitrary number of lines.

## Choose One Safe Extraction

1. Name what the current file owns.
2. Name which responsibilities are unrelated.
3. Identify the smallest group of related code that can move without changing behavior.
4. Choose one safe extraction slice.
5. Preserve behavior while moving ownership.
6. Check that the extracted code preserves behavior.
7. Review the new shape before choosing another slice.

Prefer these slices, from lower to higher risk:

- static or mostly presentational markup into a component;
- a modal or dialog into a component while keeping submit handlers in the parent;
- repeated interaction mechanics such as menus or pickers into a focused component;
- a page section into a page component;
- shared state shape into a typed store;
- pure data shaping and validation into a domain library; and
- business mutations into store or domain actions after validation covers the behavior.

Do not move UI, state substrate, and business behavior together unless an smaller extraction would be incorrect or unsafe.

## Preserve Behavior While Moving Code

Keep existing handlers, selectors, keyboard interactions, visible behavior, persistence behavior, and external contracts stable during the extraction.

Pass existing rendered nodes or callbacks when that avoids rewriting behavior while refactoring. A temporary thin wrapper is acceptable when it preserves a real call site or public interface during a safe first step. Remove the wrapper when that reason is gone.

After the move, check whether the original file became smaller or materially clearer and whether the new module has a real reason to exist. Remove accidental abstraction introduced only for the extraction.

## State Placement

When state blocks modularity, move state intentionally.

Good patterns:

- introduce a typed store for shared application or editor state;
- move state shape first while preserving existing setter names or callsites when that lowers risk;
- move behavior into actions for the relevant domain after the shared state is stable;
- keep ephemeral DOM refs and purely local widget state inside the component that owns the DOM; and
- keep derived values near the consumer until reuse justifies moving them.

Avoid:

- a giant store that becomes the new giant file;
- moving all mutations and all UI in the same pass;
- passing store state through many layers instead of selecting it where it is used; and
- global state for one component’s temporary UI detail.

## Validate The Extracted Code

Run the relevant build, formatter, type checker, or focused test command. For UI, verify the actual browser flow that uses the extracted component and check runtime errors after a clean reload when hot reload was involved.

Validation should prove that the extracted code preserves behavior. Passing a broad suite without exercising the changed path is supporting evidence, not proof.

## Red Flags

| Red flag | Correct action |
| --- | --- |
| “This file already has all the state, so add it here.” | Stop and identify which module should own the new state. |
| “I’ll extract later after the feature works.” | Extract the smallest safe piece before adding unrelated work. |
| “This helper returns markup, so it is modular enough.” | Turn repeated markup or a self-contained interaction into a named component. |
| “The root can coordinate everything.” | Keep root orchestration thin and move page, component, store, and domain ownership out. |
| “The store can hold everything.” | Put only state shared across components in the store; keep local UI details local. |
| “The refactor touched everything, so keep going.” | Stop and verify the current change before extracting anything else. |

## Stop Condition

Stop a modularity pass when:

- the moved code has a clear owner;
- the original file is smaller or materially easier to scan;
- tests confirm that the extracted code preserves behavior;
- the new module has a real reason to exist; and
- further extraction would mix concerns or rewrite behavior not covered by validation.

Describe any remaining extraction explicitly instead of hiding it in vague TODOs or continuing an unsafe refactor.
