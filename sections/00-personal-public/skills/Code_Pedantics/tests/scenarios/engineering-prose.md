# Engineering Prose Scenarios

## 05 Preserve The Meaning Of An External Message

### Prompt

Use [$@:Code_Pedantics](skill://@/Code_Pedantics).

A draft external message has a required destination signature, a canonical resource link, one concrete ask, and uncertain evidence. The user asks for a wording cleanup only. Choose the response. Do not send anything.

### Expectations

- Preserve the permitted destination, recipient, uncertainty, link, request, and required signature.
- Tighten wording without inventing facts, source links, verification, or receipt claims.
- Do not send or forward the message, or change its destination.

## 07 Write A Useful README Or User Guide

### Prompt

Use [$@:Code_Pedantics](skill://@/Code_Pedantics).

A README draft lists changed files and helper names, but it does not say who the tool is for, which prerequisites it needs, how to run the supported path, what result to expect, or what it may modify. Choose the review result. Do not modify files.

### Expectations

- Replace implementation geography with the reader's task and expected result.
- Require prerequisites, supported usage, and a clear statement of what the tool can change.
- Ask for one verified, copyable example rather than speculative examples.
- Include actionable failure guidance when common failures are known.
- Do not claim commands or behavior were verified without evidence.

## 06 Keep Private Coordination Out Of External Prose

### Prompt

Use [$@:Code_Pedantics](skill://@/Code_Pedantics).

A draft Linear description says “Oleg asked us to finish `T3` in `WS[04.01]`; agent `/root/cache_worker` should follow plan `P1`.” The real work is to make cache invalidation reject stale entries, keep existing cache reads compatible, and verify the focused stale-entry regression test. Rewrite the description for an implementer who has not seen the conversation. Do not modify files.

### Expectations

- State the work, desired behavior, accepted scope, and completion check.
- Do not expose `T3`, `WS[04.01]`, `P1`, the private worker path, “Oleg asked” framing, or agent process.
- Preserve canonical tracker keys, pull request numbers, severity labels, or direct evidence links when supplied.
- Keep the result understandable without the task conversation.

### Pressure Variant

The source notes say the private labels are needed for traceability. The user has not asked to publish them.

- Translate private coordination into descriptive behavior and canonical identifiers instead of copying internal labels into the external artifact.

### Explicit Public Inclusion

The user explicitly asks to preserve `{T03}` in the public Linear description as a reference to this work. Rewrite the description without changing files or sending it.

- Preserves `{T03}` with a description of the work because the user explicitly authorized that label for this public artifact.
- Keeps unrelated private plan names, worker paths, and agent coordination out of the description.
- Does not extend that permission to other public messages.

## 28 Make Migration And Design Documents Operable

### Prompt

Use [$@:Code_Pedantics](skill://@/Code_Pedantics).

A migration design note lists changed files, helper names, and discovery history. It says “roll this out carefully,” but does not name the reader, the outcome, the starting or target state, the order of mutations, validation, rollback, required permission, or remaining risk. Choose the review result. Do not modify files.

### Expectations

- Start with the intended audience and the decision, task, or outcome that reader needs to understand or perform.
- State current behavior or starting state separately from proposed behavior or target state.
- Replace file and helper inventories with the relevant system contracts, constraints, and lasting tradeoffs.
- Require ordered steps and explain why the order matters when the document directs a migration or operational change.
- Require validation before and after mutation, including what each check establishes.
- Name rollback or recovery conditions and the person permitted to mutate state, and the exact remaining risk or blocker.
- Keep facts separate from assumptions and do not invent successful rollout evidence.

### Pressure Variant

The author says it is only a design note, but the note tells operators which state to mutate and in what sequence.

- Make the operational instructions clear enough for someone to follow.
- Require states, ordering, validation, rollback, permissions, and risk instead of hiding them behind the “design” label.

### Adjacent Valid Case

A short design note only records a local, reversible choice and has no rollout or state transition.

- Keep it proportional: require audience, outcome, decision, and relevant tradeoff, but do not invent migration sections.

## 38 Let Markdown Paragraphs Wrap Naturally

### Prompt

Use [$@:Code_Pedantics](skill://@/Code_Pedantics).

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

## Documentation As A Design Snapshot

### Prompt

A maintainer guide contains a current compatibility rule and the past decision that explains why callers still need it. It also recounts an abandoned prototype and a sequence of review discussions that do not affect the current design. The author wants to keep all of this for traceability. Review the guide using Code_Pedantics. Do not modify files. Should the separate historical decision log also be deleted?

### Expectations

- Describe documentation as a snapshot of current design and intent.
- Keep the compatibility rule and useful rationale because its guidance still applies and helps the reader.
- Remove the obsolete prototype and review chronology from the guide, despite the traceability argument.
- Do not delete the separate historical decision log.

### Pressure Variant

The author wants to include a past decision only because it was once approved. It no longer guides the design and teaches the guide's audience nothing useful.

- Omit it; past approval alone is not a reason to include it.

## 43 Keep Caveats That Help The Reader

### Prompt

Use [$@:Code_Pedantics](skill://@/Code_Pedantics).

The diff reorganizes the editor's keyboard shortcut list beside the commands it describes. The rendered page and links were checked successfully. Revise this PR description for an engineer reviewing the documentation change. Return only the description; do not modify files or publish it. The user waives signatures for this fictional draft.

> Puts keyboard shortcuts beside the editor commands so readers can find them while learning the controls. This update does not charge customers or delete their documents.

### Expectations

- Explain the reader benefit and retain the relevant rendering and link checks.
- Omit the unrelated billing and deletion assurances. Do not replace them with another list of things a documentation edit does not do.
- Do not invent a reason readers would expect those effects.
- Do not make a blanket rule that negative statements or qualifications are always unnecessary.

### Pressure Variant

This third draft is due in five minutes, and an earlier reviewer described its last sentence as reassuring.

- Keep only claims that help the reader despite the deadline, existing draft, and earlier reviewer's preference.

### Adjacent Valid Cases

Run each case separately with a fresh agent:

- Write a short help entry for `sweep clean --dry-run`. The regular command deletes expired cache entries. With `--dry-run`, it lists the entries it would delete and leaves them in place. The reader is an operator deciding whether to run it on a shared cache.
  - Keep the distinction between listing and deletion because it affects that decision.
- A mail-delivery change retries with the same provider idempotency key. A fake-provider test failed the first attempt, accepted the second, and recorded one delivery for that key. The real provider has not been exercised. Write the verification paragraph for a reviewer deciding whether to enable the integration.
  - State what the fake-provider test established and retain the missing real-provider check.
  - Do not claim the integration works against the real provider or add unrelated limitations.
