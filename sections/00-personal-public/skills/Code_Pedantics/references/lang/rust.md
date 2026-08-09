# Rust Pedantics

## Imports

Keep name ownership readable at the use site. Do not replace one mechanical rule with another.

Bad:

```rust
use domain_graph::*;
use domain_store::{Record, RecordId, RecordState, RecordStore, StorageError};
```

Good when the module name explains otherwise generic domain names:

```rust
fn load_record(
    store: &domain_store::RecordStore,
    record_id: domain_store::RecordId,
) -> Result<domain_store::Record, domain_store::StorageError> {
    // ...
}
```

Good when a name is frequent or a trait needs to be in scope:

```rust
use request_runtime::{RequestContext, RequestContextExt};
```

Use these examples to guide judgment, not as a fixed syntax rule. Read the surrounding signatures and expressions. Choose the form that lets a human understand both the operation and each domain name's owner without unnecessary lookup or repeated prefixes.

- Replace changed glob imports with deliberate item imports or qualified paths. Keep a glob only when the repository requires it for a prelude.
- Prefer qualification when the module name supplies useful domain context or disambiguates a generic name.
- Prefer an item import when qualification would dominate the file, the owner is already obvious, or method resolution requires a trait in scope.
- Treat import count and use count as clues, not thresholds. Read actual signatures and expressions as an engineer unfamiliar with the file.
- Make identifiers clear where they are used, not merely in the import block. Do not widen a focused task into cleanup of untouched imports.

## Module Files

Prefer a flat namespaced filename such as `<namespace>_<name>.rs` to `<namespace>/<name>.rs` when an extraction only splits one module into another file. The flat filename keeps the source tree shallow and does not imply a family of child modules that does not exist. Use the matching flat module name; do not add `#[path]` solely to present the flat file as a nested module.

Use a directory when the namespace owns several coherent child modules, when the directory represents an intentional public module hierarchy, or when a stricter repository convention requires it. Do not create a directory for one extracted file merely because Rust permits the layout.

## Names And Expressions

Let `rustfmt` own layout. Restructure an expression when the formatted result is hard to read; do not add indentation that `rustfmt` removes.

Name values by role and kind, such as `config_path`, `request_id`, `response_bytes`, and `loaded_record`. Name intermediate values when one expression performs several operations, an asynchronous call may fail, or a multiline chain hides the meaning.

A short obvious expression may stay inline. Do not name every future or awaited result mechanically.

Replace a complicated `if` condition with a specific domain predicate when joined comparisons force the reader to infer the branch purpose. Add nearby rationale when the predicate name does not explain why the branch matters.

## Documentation And Tests

Use `//!` for a module contract and `///` for an item contract. Item docs should explain primary purpose, architectural role, why the item exists, why its specific shape and responsibilities are useful, and the permissions, rules, side effects, and failure behavior maintainers need. Keep an item contract correct in isolation and across callers. Document every changed function or type deeply enough for a cold reviewer to understand its domain, workflow role, inputs, outputs, relationships, mutation, failure, and conditions for a safe change without repository or task context. Use multiple sentences when needed. Do not impose a limit of one line. Use `//` for local rationale, not syntax narration; do not turn one caller's phase or ordering advice into a function contract unless it is a true universal precondition.

An involved test has several fixtures that prove behavior, logical operations, or negative evidence whose purpose is not evident from names and assertions. Give it a short overview with setup and expected-proof bullets:

```rust
/// Verify that processing the selected record leaves unrelated state intact.
///
/// - Setup:
///   - Store a selected record and an unrelated record.
///   - Record fallback calls.
/// - Expect:
///   - The selected record changes.
///   - The unrelated record remains unchanged.
///   - The fallback-call count remains zero.
```

Give each fixture whose purpose is not obvious a nearby comment. Explain what absent calls or preserved state prove, and include positive evidence that the operation ran so negative assertions cannot pass vacuously. Avoid comments in simple tests when names and assertions already state the behavior.
