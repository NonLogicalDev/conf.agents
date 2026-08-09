# TypeScript And JavaScript Pedantics

## Assign Clear Responsibilities

Keep entrypoints and root components thin. They should own orchestration, routing, dependency wiring, and setup shared across the application rather than accumulating page markup, shared state, validation, persistence, and business mutations.

Prefer:

- pages or screens for composition specific to that page;
- components for reusable or self-contained interactions;
- stores for shared state shape and transitions;
- domain modules for pure transformations, validation, and business rules; and
- UI primitives for repeated interaction mechanics.

A helper that returns a large markup tree is usually a component that needs a clear name and responsibility.

## State Placement

Keep ephemeral DOM refs and purely local widget state in the component that owns the DOM. Keep derived values near the consumer until reuse justifies moving them. Move shared state into a typed store only when several components need the same state or transitions.

Avoid:

- a giant store that becomes the new giant file;
- global state for one component's temporary UI detail;
- prop chains that exist only because shared state has no coherent owner; and
- moving all state, mutations, and UI in one pass.

Move the state shape first when that lowers risk. Move domain actions later, after focused tests or a real UI flow cover the behavior.

## Safe Extraction

When a file is mixed or oversized:

1. Name the responsibilities it owns now.
2. Choose one safe group of related code: presentational markup, a dialog, repeated interaction, page section, store substrate, or pure domain operation.
3. Preserve handlers, visible behavior, selectors, keyboard behavior, and external contracts.
4. Pass existing callbacks or rendered nodes when that avoids rewriting behavior during the extraction.
5. Check that the extracted code preserves behavior before extracting more.

Do not extract UI, state substrate, and business rules together unless the current design cannot be made safe incrementally. A temporary thin wrapper is acceptable only when it preserves a real call site or public interface; remove it when that reason is gone.

## Types, Names, And Validation

Use types to show the shape of the state and who may change it. Prefer narrow domain names over generic `data`, `state`, `result`, or `handleThing`. Name protocol values, storage keys, and important timeouts that encode an external contract.

Avoid `any` or broad casts that hide unvalidated input. Parse and validate external data before it enters typed domain code. Keep error messages specific enough to identify the violated contract.

Run the repository formatter, type checker, and focused tests. For UI changes, verify the real browser flow that uses the extracted component, then check for runtime errors after a clean reload when hot reload was involved.
