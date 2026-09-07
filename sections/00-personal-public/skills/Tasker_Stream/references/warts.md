# Tool, Guidance, and Deviation Warts

Wart files turn recurring project friction into useful memory and explain deliberate departures from existing guidance. They help an owner stop relearning the same problems and adapt to the user's intent.

## Keep a compact summary

When a wart file has several useful entries, give each tool wart a stable `T-<num>` identifier and each guidance wart a stable `G-<num>` identifier. Preserve existing identifiers and add a new number only for a genuinely different issue.

When a ledger has enough entries that a summary helps, keep a short summary at the top:

```markdown
| Wart | Verified occurrences | Status | Current response |
| --- | ---: | --- | --- |
| `T-001` <tool issue> | 2 | Mitigated | <current safe workaround> |
| `T-002` <another issue> | 1 | Active | <the next useful action> |
```

Count separate verified observations when the count helps explain the problem. Do not count guesses, repeated descriptions of one event, or the number of people discussing it. Use the state that the evidence supports:

- `Active`: The issue still affects the work and no sufficient response has been verified.
- `Mitigated`: A safe workaround helps, but the underlying issue remains.
- `Contained`: A harmful effect has stopped, but the cause or full repair remains unresolved.
- `Resolved`: The underlying issue was fixed and the result was verified.

## Record tool problems

Use `_owner/wart.tools.md` for problems in the environment, commands, build steps, tests, interfaces, project helpers, or other tools.

Useful examples include confusing output, unreliable status, missing capabilities, repeated manual steps, expensive setup, or a command that appears to succeed without proving the expected result.

```markdown
## T-001: <Short name for the tool problem>

- Occurrences: <number of separately verified observations>.
- Observed: <what the tool actually did>.
- Evidence: <command, output, file, or verified example>.
- Cause: <the observed cause, or unknown>.
- Impact: <time, uncertainty, repeated work, or blocked progress>.
- Workaround: <the current safe way to continue, when one exists>.
- Improvement: <a practical change to the tool, helper, or workflow>.
- Status: <active, mitigated, contained, or resolved>.
```

Update the existing entry when the same issue recurs. Keep a verified workaround distinct from a verified fix. Leave enough context for the next owner to avoid repeating the same failed action.

Problems in the environment, a build, repository layout, test runner, stale example, changing system state, or tool output belong here even when a better instruction might later help someone avoid them.

## Record skill and instruction problems

Use `_owner/wart.guidance.md` when a skill, owner instructions, project instructions, or other guidance makes the agent less effective.

Look for instructions that conflict, omit useful context, encourage duplicated work, impose the wrong naming scheme, prescribe needless process, or pull the agent away from the user's actual goal.

```markdown
## G-001: <Short name for the guidance problem>

- Occurrences: <number of separately verified observations>.
- Affected guidance: <the skill, owner instructions, or project instructions>.
- User intent: <the outcome or preference the guidance should support>.
- Observed: <what the instruction caused or failed to explain>.
- Evidence: <the relevant instruction, file, or verified outcome>.
- Cause: <the missing, unclear, outdated, or conflicting instruction>.
- Impact: <confusion, repeated work, incorrect behavior, or delay>.
- Workaround: <the current safe interpretation or owner-local fix>.
- Improvement: <the smallest useful clarification or correction>.
- Status: <active, mitigated, contained, or resolved>.
```

When the user has allowed careful self-improvement, apply a small safe correction to the owner's `AGENTS.md`, project helper, or writable skill source if it stays true to the user's intent. Preserve existing permissions and check the result.

Record a guidance wart only when the instruction itself is defective. A broken command, outdated outside example, missing source file, or changing external state is a tool wart unless separate evidence also proves a problem in the guidance.

## Record deliberate deviations

Use `_owner/wart.deviation.md` for every intentional departure from an existing owner convention or prior guidance. Record what changed, why the change better serves the user's intent, and which instructions or permissions still apply.

Update an existing entry when its circumstances change. A deviation does not override current user instructions, accepted scope, safety checks, or permissions.

## Keep each record in the right place

Keep the immediate blocker and next action in the active plan. Keep the reusable tool or guidance lesson in its wart file.

Use `PROSE_STEERING.md` for repeated writing feedback. Use another `wart.<type>.md` only when the problem does not fit tools, guidance, or deviations and remains useful after the current task.

A wart describes a real problem. It does not grant permission to publish, change another owner's work, ignore the user, or skip a required check.
