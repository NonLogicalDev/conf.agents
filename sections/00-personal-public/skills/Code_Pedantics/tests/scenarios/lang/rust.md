# Rust Scenarios

## 24 Keep Rust Intent Visible

### Prompt

Use `$Code_Pedantics`.

A changed Rust module imports `super::*`, hand-aligns a long async expression
that `rustfmt` collapses, and leaves a joined comparison inline even though it
decides whether an archived record may be mutated. The associated test has
several fixtures and proves only that an unrelated record did not change; its
setup and negative assertion have no rationale. Choose the review result. Do
not modify files.

### Expectations

- Replace the changed glob import with deliberate item imports or
  module-qualified paths unless the file has an explicit established prelude
  exception.
- Let `rustfmt` own layout; name a meaningful intermediate value when the
  formatted expression still hides a fallible or interpretive step.
- Replace the joined comparison with a specific domain predicate and add local
  rationale when the predicate name alone does not explain what may change.
- Use module or item documentation for responsibility, permissions, stable rules,
  side effects, and failure behavior that callers need.
- Give the involved test a short overview of what it proves, explain unusual
  fixtures and preserved state, and add positive evidence that the operation
  actually ran.

### Adjacent Valid Case

A test module uses a prelude glob required by the repository, a short expression stays
clear after `rustfmt`, and a simple test name plus assertions fully state the
proof.

- Keep the established prelude import, inline expression, and simple test
  without unnecessary comments.
- Do not turn a focused change into cleanup of untouched glob imports.

## 33 Prefer Flat Namespaced Rust Module Files

### Prompt

Use `$Code_Pedantics`.

A Rust refactor extracts one implementation unit from `ledger.rs`. The code
can live in either `ledger_events.rs` or `ledger/events.rs`. No other `ledger`
child modules exist or are planned, and the extraction does not introduce a
public hierarchy of child modules. The nested layout is already drafted, the
change is small, `rustfmt` passes, and a senior reviewer says either layout is
idiomatic. Choose the review result. Do not modify files.

### Expectations

- Prefer `ledger_events.rs` because the extraction only splits one module into
  another file.
- Treat `ledger/events.rs` as implying a family of child modules that the refactor
  does not create.
- Request the local rename despite the drafted layout, passing formatter, and
  reviewer preference.
- Use the matching flat module name. Do not add `#[path]` solely to preserve a
  nested private module name.
- Do not expand the change into unrelated module-layout cleanup.

### Adjacent Valid Case

The `ledger` namespace owns several coherent child modules, exposes an
intentional public hierarchy of child modules, or follows a stricter repository
convention that requires a directory.

- Keep `ledger/<name>.rs` for that module family or required convention.
- Do not flatten a real module hierarchy mechanically.

## 34 Keep Rust Identifier Ownership Visible

### Prompt

Use `$Code_Pedantics`.

A Rust file imports long lists of domain types, constants, and functions from
several crates. Most identifiers appear only once or twice, and several names
are generic without their crate. Replacing the lists with qualified uses such
as `archive_graph::NodeIndex` and `archive_store::RecordId` would make each
identifier's owner visible. The existing imports compile, `rustfmt` passes,
and a reviewer says exact item imports are always cleaner than qualified
paths. Choose the review result. Do not modify files.

### Expectations

- Prefer module-qualified paths for the sparse domain identifiers whose module
  names add useful context.
- Avoid preserving long grouped item imports merely because they are exact and
  compile.
- Judge readability at the use sites, not only the import block.
- Treat import count and use count as clues, not automatic thresholds. Choose
  the form that minimizes ambiguity and reader lookup in the actual file.
- Keep deliberate item imports for frequently repeated identifiers, traits
  required in scope, or stricter repository conventions.
- Do not replace the grouped imports with glob imports.

### Adjacent Valid Case

A specific type appears throughout a file, its owner is obvious from the file
contract, and repeating the module prefix would dominate signatures and
expressions. Extension methods also require an imported trait.

- Import the repeated type and required trait explicitly.
- Do not qualify every identifier mechanically.
