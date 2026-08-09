---
name: Code_Pedantics
description: Use when writing, reviewing, or simplifying code, tests, engineering documents, commit messages, or review comments. Make the result clear, modular, supported by evidence, and easy for another person to review.
---

# Code Pedantics

## Outcome

Make the next engineer's job easier. Remove defects a careful reviewer should
not need to point out: unclear behavior, missing explanations, vague prose,
misleading validation claims, unnecessary wrappers, mixed ownership,
speculative abstraction, fragile defaults, and avoidable indirection.

Apply only the parts of this skill that fit the work. You do not need to run
code to improve prose. A small code change does not call for changes across
the repository. Follow the repository's instructions and required templates.

When editing a skill, remember that its wording shapes later sessions. Use
simple English, short sentences, direct instructions, and calm explanations.
Avoid legal language and long chains of awkward compound adjectives. Keep
familiar terms such as `self-contained` when they make a sentence clearer.
Preserve exact technical terms, skill names, commands, paths, fixtures, links,
permissions, and safeguards.

When writing or revising a skill, read
`~/.codex/resources/communication-principles/WRITING_STYLEGUIDE.md`. Keep
familiar compounds when they help. Use fuller wording when it explains the
action, reason, or result more clearly.

Write each Markdown prose paragraph or list item on one source line, however wide it needs to be, and let the editor wrap it. Do not introduce an 80-character limit or apply a Git commit convention to other Markdown.

## Choose The Artifact

- Code, tests, scripts, and operator documentation: use the matching review,
  modularity, comments, validation, completion, and language references.
- Commit messages and review descriptions: use the reference on commits and
  pull requests and the reference on finishing the work.
- Messages to other people: keep the permitted destination, recipient, facts,
  uncertainty, links, attribution, and required signature. Use this skill to
  improve wording only. Do not send a message or claim that one was sent.
- Private notes and plans: keep facts, decisions, and evidence intact; remove
  only wording that obscures them.

## Read General References

General guidance lives in the references below. Read only the topics
that match the artifact:

- [references/review.md](references/review.md) for removing avoidable flaws,
  ordering a code review, checking when values change, and keeping code simple.
- [references/modularity.md](references/modularity.md) when a file, module,
  component, store, or entrypoint has mixed ownership.
- [references/comments.md](references/comments.md) for local human rationale,
  documentation, user permissions, and conditions for a safe change.
- [references/tests.md](references/tests.md) for fixtures that prove behavior,
  negative evidence, regression tests, and behavioral evaluators.
- [references/validation.md](references/validation.md) for proportional
  validation and safe checks through the actual entrypoint.
- [references/engineering-prose.md](references/engineering-prose.md) for doc
  comments, documentation, READMEs, guides, plans, and external prose.
- [references/comms/commits-and-prs.md](references/comms/commits-and-prs.md)
  for commit messages and pull request prose.
- [references/completion.md](references/completion.md) before handoff.

## Read Language References

Read every language reference that matches a changed file:

- Dockerfile: [references/lang/dockerfile.md](references/lang/dockerfile.md)
- Python: [references/lang/python.md](references/lang/python.md)
- Rust: [references/lang/rust.md](references/lang/rust.md)
- Shell: [references/lang/shell.md](references/lang/shell.md)
- TypeScript and JavaScript:
  [references/lang/typescript-javascript.md](references/lang/typescript-javascript.md)

When no language reference exists, apply the relevant general references. Add a
language reference only when it offers useful guidance that is not obvious.

## Apply The References

1. Identify the artifact, audience, intended behavior or message, and stricter
   repository instructions or rules for the destination.
2. Load the smallest set of general references that covers
   the work.
3. Load every matching language reference.
4. Preserve facts, uncertainty, exact technical names, required templates,
   links, signatures, and the author's useful voice.
5. Apply the completion reference after the final edit, not to an intermediate
   draft.

## Tests

When changing this skill, read [tests/README.md](tests/README.md). Run the
relevant scenarios with fresh, independent agents when the available tools
permit.
