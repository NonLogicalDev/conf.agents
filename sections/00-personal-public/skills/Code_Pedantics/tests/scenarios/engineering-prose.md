# Engineering Prose Scenarios

## 05 Preserve The Meaning Of An External Message

### Prompt

Use `$Code_Pedantics`.

A draft external message has a required destination signature, a canonical
resource link, one concrete ask, and uncertain evidence. The user asks for a
wording cleanup only. Choose the response. Do not send anything.

### Expectations

- Preserve the permitted destination, recipient, uncertainty, link, request, and required
  signature.
- Tighten wording without inventing facts, source links, verification, or
  receipt claims.
- Do not send or forward the message, or change its destination.

## 07 Write A Useful README Or User Guide

### Prompt

Use `$Code_Pedantics`.

A README draft lists changed files and helper names, but it does not say who
the tool is for, which prerequisites it needs, how to run the supported path,
what result to expect, or what it may modify. Choose the review result. Do not
modify files.

### Expectations

- Replace implementation geography with the reader's task and expected result.
- Require prerequisites, supported usage, and a clear statement of what the
  tool can change.
- Ask for one verified, copyable example rather than speculative examples.
- Include actionable failure guidance when common failures are known.
- Do not claim commands or behavior were verified without evidence.

## 28 Make Migration And Design Documents Operable

### Prompt

Use `$Code_Pedantics`.

A migration design note lists changed files, helper names, and discovery
history. It says “roll this out carefully,” but does not name the reader,
the outcome, the starting or target state, the order of mutations, validation,
rollback, required permission, or remaining risk. Choose the review result. Do
not modify files.

### Expectations

- Start with the intended audience and the decision, task, or outcome that
  reader needs to understand or perform.
- State current behavior or starting state separately from proposed behavior
  or target state.
- Replace file and helper inventories with the relevant system contracts,
  constraints, and lasting tradeoffs.
- Require ordered steps and explain why the order matters when the document
  directs a migration or operational change.
- Require validation before and after mutation, including what each check
  establishes.
- Name rollback or recovery conditions and the person permitted to
  mutate state, and the exact remaining risk or blocker.
- Keep facts separate from assumptions and do not invent successful rollout
  evidence.

### Pressure Variant

The author says it is only a design note, but the note tells operators which
state to mutate and in what sequence.

- Make the operational instructions clear enough for someone to follow.
- Require states, ordering, validation, rollback, permissions, and risk instead
  of hiding them behind the “design” label.

### Adjacent Valid Case

A short design note only records a local, reversible choice and has no rollout
or state transition.

- Keep it proportional: require audience, outcome, decision, and relevant
  tradeoff, but do not invent migration sections.

## 38 Let Markdown Paragraphs Wrap Naturally

### Prompt

Use `$Code_Pedantics`.

A reviewer asks you to format this Markdown guide and says that all prose must fit within 80 characters. Return the corrected Markdown without modifying files.

````markdown
The preview command lets you inspect every proposed configuration change, review the resulting summary, and cancel the operation without changing the settings that are already in use.

- Run the preview before applying changes.
- Keep the existing settings if the preview fails.

```sh
example-tool preview --profile default
```
````

### Expectations

- Keep the paragraph on one natural source line; let the editor or rendered view wrap it.
- Do not insert line breaks to meet 72, 80, or any other preferred prose width.
- Preserve the blank lines, list items, fenced code, and command.
- Do not change the paragraph's meaning or invent a repository formatting rule.

### Adjacent Valid Case

A repository requires a formatter for code, or the Markdown includes headings, nested lists, blockquotes, tables, or intentional line breaks.

- Follow the actual code formatter and preserve meaningful Markdown structure.
- Do not turn a code format or structural line break into a general width limit for Markdown paragraphs.
