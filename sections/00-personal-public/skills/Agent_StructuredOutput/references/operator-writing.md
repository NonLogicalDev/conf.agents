# Operator Writing And Sources

Read this reference only when a response needs substantial sections,
several sources, an operational brief, or a reusable writing artifact.
A short answer does not need a forced template.

## Keep Language Plain

Write like a direct teammate. Use common words, short sentences, exact
technical names, and concrete facts. Lead with what happened, what is
known, what cannot yet be proved, or what happens next.

Keep a sentence only if it helps the operator understand, decide, or
act. Remove filler, repeated conclusions, process narration, generic
praise, speculative confidence, and inflated alternatives such as
`strategy-aligned validation` or `strong progress`. Never replace a
technical term such as `JWT aud`, `RepositoryAccessPolicy`, or `baseRefName`
with a less precise synonym.

## Structure Only Substantial Information

Use a topic heading for independent substantial workstreams or a real
change in decision, blocker, or direction. Separate major sibling sections
with one Markdown horizontal rule surrounded by blank lines.

Use a parent bullet for the outcome. Nest only relevant evidence,
exceptions, blockers, and next actions. Keep sibling bullets parallel.
Do not add a heading to a routine receipt, an ordinary status, or a
response with one main point.

A final answer separates its sections with two horizontal rules only when it
actually combines a completed lifecycle receipt and a separate
`## Summary`. Read
[lifecycle-receipts.md](lifecycle-receipts.md) for that exact format.

## Annotate Multiple Sources

State the supported conclusion before its sources. Give each source its
own sibling bullet with a descriptive linked filename, exact line number,
and a short explanation of what it establishes:

```markdown
The gateway requires both checks.

- [settings.py](/tmp/example/settings.py:115): Defines the required audience.
- [identity.py](/tmp/example/identity.py:314): Enforces the allowlist.
```

Do not combine two sources into a chain of inline links. One supporting
source may remain inline when its purpose is obvious. Never invent a
path, source, line, or canonical URL.

## Use A Compact Operational Brief

Use this shape only when a status brief improves the answer:

```markdown
**Outcome:** <verified current state>

- **Evidence:** <actual artifact or check>
- **Blocker:** <real blocking condition, when present>
- **Next:** <actual next action>
```

For genuinely independent work, use one parent outcome per workstream
and preserve each complete current marker from the main skill. Omit
empty fields, raw logs, reassurance, and tool narration.

## Use Writing Blocks Only When Supported

Use a writing block only when the current surface supports it and the
operator is likely to reuse a finished snapshot, handoff, guide, or
template. Keep ordinary Markdown inside the block.

Never use a writing block for a simple answer, unfinished work, a file
receipt, or an unsupported runtime. Do not claim that an artifact was
saved merely because a response displayed it.
