# Prose Steering

Use `<plan-root>/PROSE_STEERING.md` when writing feedback repeats across a project. The file helps an owner recognize patterns, improve future drafts, and avoid making the same mistake again.

This is a record of observed feedback, not a replacement for the user's writing guidance and not permission to change work outside the accepted scope.

## Record themes, not incidents

Add a theme only when it describes a meaningful pattern. When the same issue appears again, update the existing theme and its occurrence count.

```markdown
## PS-001: <Short name for the recurring problem>

- Occurrences: <number of verified examples>.
- Observation: <what the writing actually did>.
- What the owner missed: <the assumption or habit that produced the issue>.
- Evidence: <the relevant file, passage, or verified example>.
- Reader impact: <why the wording confused or misled the reader>.
- Better approach: <the concrete behavior to use next time>.
- Possible guidance improvement: <a useful change, when one is justified>.
```

Keep examples appropriate for the project and its audience. Do not invent feedback, counts, evidence, or a rule the user did not give.

Save useful before, after, or draft examples under `_owner/memory/prose-examples/` when those examples will help future work.

## Apply feedback to the real writing

When a repeated problem appears, inspect the other writing changed by the same task and correct the same pattern where it actually occurs.

Re-read the user's current writing guidance before a substantive prose task. Treat the ledger as project memory, not as a substitute for that guidance.

Check that important claims match the source, test, or other evidence that supports them. Explain the action and result before introducing an abstraction. Give the reader only the background needed to understand the decision.

Match the amount of explanation to the artifact. A comment, plan, design document, and review summary may need different levels of detail.

Keep a substantial draft local until it is ready for the action the user requested. Preserve the precision of dates, measurements, and other evidence; do not add details that were not observed.

## Keep the ledger focused

Use `PROSE_STEERING.md` for repeated writing feedback. Use `_owner/wart.guidance.md` for problems with project instructions, `_owner/wart.tools.md` for tool problems, and other `wart.<type>.md` files only when they describe a different recurring issue.

If the same theme belongs in a shared writing guide, record the proposed improvement without changing that guide unless the user asks.
