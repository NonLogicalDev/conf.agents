# Comments Scenarios

## 03 Explain Permissions Where They Matter

### Prompt

Use `$Code_Pedantics`.

A deletion guard checks whether ownership metadata is missing. The code is
mechanically readable, but a future engineer cannot tell why missing metadata
prevents deletion. Choose the review result. Do not modify files.

### Expectations

- Explain the required permission and conditions for a safe change near the
  code.
- Prefer an exact name, nearby comment, or item documentation over hidden task
  context.
- Do not narrate the condition syntax.
- Add or request a test whose negative evidence proves deletion did not occur.

## 10 Lead Comments With Purpose

### Prompt

Use `$Code_Pedantics`.

A module documentation comment starts by listing helper functions, internal
files, and the order of implementation steps. It is accurate, but a reader
cannot tell what capability the module enables, who benefits from it, or how
to recognize that it is working. Choose the review result. Do not modify
files.

### Expectations

- Ask for the purpose and observable outcome before implementation mechanics.
- Explain the capability, user outcome, or operator outcome the module enables
  and how a reader can recognize successful behavior.
- Keep important responsibilities, constraints, or tradeoffs after that
  purpose.
- Do not replace useful rationale with a longer inventory of files, helpers, or
  control flow.

### Adjacent Valid Case

A short local comment sits beside an important guard and names the risk the
guard prevents.

- Keep the short local rationale; do not require every nearby comment to become
  an overview of the module that spans several paragraphs.

## 31 Explain Long Orchestration Phases

### Prompt

Use `$Code_Pedantics`.

A long coordinator validates ownership, builds a mutation plan, applies it,
and reports the result. Every statement is mechanically readable, but the
phases share state, helpers would need large argument lists, and no local
comment explains why validation precedes planning or why reporting happens
after mutation. Choose the review result. Do not modify files.

### Expectations

- Look for coherent named helpers before accepting one long coordinator.
- Keep phases together when shared state or ordering makes extraction less
  clear, rather than splitting solely to reduce line count.
- Add concise phase comments at transitions when the function remains long.
- Explain each phase's purpose, ordering, required permission, or risk; do not
  narrate assignments, loops, or calls.
- Keep the rationale close enough that a future engineer can change the
  coordinator safely without reconstructing the task conversation.
- Prefer exact names for the owned state, plan, mutation, and result instead
  of vague labels such as “step” or “logic.”

### Adjacent Valid Case

A short coordinator calls clearly named helpers in an obvious order and has no
unusual shared state or permission requirement.

- Leave out comments when the helper names and structure already explain the
  behavior.

## 35 Separate Reusable Documentation From Local Reasons

### Prompt

Use `$Code_Pedantics`.

A Rust extraction that changes only comments introduces a structure for a
rebuild plan and a
verification helper. The struct doc says only that it carries the fields used
during a rebuild. The helper doc says to call it after loading a manifest and
before publishing a replacement. Both statements are accurate for the current
caller, the formatter and tests pass, and the diff is large, but a reader
cannot tell why the type has this shape or why the helper exists. Choose the
review result. Do not modify files.

### Expectations

- Ask the struct doc to explain its architectural purpose and why its fields,
  ownership, or lifetime belong together; do not merely inventory fields or
  restate the type name.
- Ask the function doc to explain the helper's primary purpose,
  responsibility, and stable behavior in a way that remains correct without seeing
  its callers.
- Require enough context for a fresh reviewer to understand the item's
  relevant domain, workflow role, inputs, outputs, relationships, mutation,
  failure, and conditions for a safe change; allow multiple sentences when the
  contract cannot be made legible in one line.
- Move the reasons for the current caller's phase and publication order to a nearby
  local comment; do not present advice for one caller as a reusable
  item contract.
- Apply the distinction across materially equivalent changed items in the
  diff, without expanding into untouched code or changing implementation.
- Keep concise docs only when they fully convey the architectural why; do not
  impose a limit of one line or hide required context in a module overview, and do
  not add boilerplate or syntax narration.

### Pressure Variant

The author says one short doc line per item is sufficient, the diff already
passed CI, and moving the ordering sentence will make the documentation patch
slightly larger.

- Keep reusable documentation separate from the reason for a particular call.
- Do not accept accurate advice in the wrong place just because the patch
  changes only comments, is large, or already passes its checks.

### Adjacent Valid Case

A verification helper always requires a held lease and always runs
before publication, regardless of caller. Its item docs state that
precondition and stable rule; a particular caller also has a separate local
comment explaining why it retries verification after refreshing the lease.

- Keep the universal precondition in the function contract and the retry
  rationale at the call site.
- Keep rules that apply to every caller in the item's documentation.

## 36 Keep Internal Plans Out Of Code Comments

### Prompt

Use `$Code_Pedantics`.

A Rust doc comment says:

```rust
/// Plan-06 carries per-request identity to the request handler, but this
/// durable encoding stays unchanged until a staged mixed-version migration can
/// bind principals without invalidating in-flight sessions.
```

The plan name and sequencing came from a private task conversation. The
compatibility requirement is real. Choose the review result and replacement
direction. Do not modify files.

### Expectations

- Remove the private plan name and task-conversation framing.
- Preserve the concrete request-identity behavior and mixed-version
  compatibility requirement.
- State which durable encoding needs to remain unchanged, using its exact
  name when available.
- Do not replace the private reference with vague phrases such as “this work”
  or “the current plan.”

### Pressure Variant

The author says the plan number is useful provenance, the review deadline is
close, and changing the wording would change the head of the stacked branch.

- Still remove the plan number.
- Keep provenance and temporary sequencing in the internal plan, not the code
  comment.
- Do not preserve the leak merely to avoid amending prose or updating the stack.

### Adjacent Valid Case

A public protocol RFC name or public migration document is the authoritative
source of a rule for the wire format and the link helps a future maintainer verify
it.

- Keep the public source when it is relevant and independently readable.
- Do not treat a private task plan as equivalent to a public specification.
