---
name: Code_Develop
description: Use when implementing, fixing, refactoring, or reviewing a development goal that needs a small complete solution, current evidence, adversarial review, and a clear stop condition.
---

# Code Develop

## Outcome

Turn a clear development goal into the smallest complete change. Explain it well and check that it works. Treat the first working pass as a draft. Recheck the goal, challenge the design, remove scope that no longer earns its keep, and stop only when current evidence proves the requested outcome.

## Lock The Goal

Before editing, state:

- the exact requested outcome and acceptance criterion;
- constraints and behavior to keep stable;
- excluded work and unrelated changes to preserve;
- the repository instructions, supported workflow, and validation commands; and
- the requested finish point: local implementation, reviewed diff, published change, landed change, or a larger verified outcome.

Do not silently turn a local implementation request into publication, review, integration, deployment, or external communication. A later finish point requires the user's permission and its own evidence.

## Inspect Before Designing

Inspect accepted nearby code, configuration, tests, and recent changes before introducing a new mechanism. Prefer the established pattern when it satisfies the acceptance criterion. Name the exact mismatch before adding a new abstraction, data source, workflow, dependency, or permission requirement.

Orient with the smallest useful reads:

1. Read repository instructions and the relevant module contract.
2. Locate the narrow behavior and its closest accepted example.
3. Read focused tests and the documented validation command.
4. Check current worktree state so unrelated edits remain untouched.

Do not use broad archaeology as a substitute for deciding the smallest change.

## Classify The Failure Before Repairing It

When a request starts from a failure, classify it before designing a repair:

- code defect;
- a missing permission;
- wrong identity or credential source;
- configuration error;
- stale state or obsolete evidence;
- unsupported command or workflow;
- missing tool or dependency;
- transient environment failure; or
- operator action.

One failed route does not prove that shared code needs to change. Check the supported routes and the narrowest configuration, access, or operator action that can satisfy the acceptance criterion before proposing code.

## Prove Necessity And Proportionality

Before a broad, privileged, expensive repair or one that affects several components, record:

- the source acceptance criterion;
- the actual caller, environment, and current state when relevant;
- smaller alternatives checked and why each fails;
- durable beneficiaries beyond the immediate failing run;
- changed components and ownership cost;
- security or authority implications; and
- validation and rollout cost.

If that argument is not defensible, shrink the change or stop. Making a new capability safer does not prove that the capability should exist.

## Build The Smallest Complete First Pass

Choose the narrowest implementation likely to satisfy the real goal.

- Keep the diff inside the stated scope.
- Preserve unrelated user and agent changes.
- Prefer direct code over speculative layers, flags, generic helpers, or compatibility branches.
- Preserve compatibility and the responsibilities of existing owners when current behavior depends on them.
- Add focused tests for the behavior, including a negative case when the failure depends on absence, preservation, or leaving something unchanged.

The first pass should be easy to explain, not merely easy to make compile.

## Run The Review And Repair Loop

After the first working pass, and after every meaningful rewrite:

1. Read the goal, constraints, and excluded work again.
2. Inspect the complete diff, not only the changed line.
3. Ask whether each new helper, branch, type, dependency, config value, and file still earns its keep.
4. Look for a smaller design that preserves the same correctness and explicit constraints.
5. Fix the most important weakness: remove unnecessary surface, collapse accidental indirection, tighten ownership, or add missing rationale.
6. Run focused validation again.
7. Repeat only while the next change materially improves correctness, simplicity, or maintainability.

Do not defend sunk cost. Do not keep polishing after the remaining complexity is tied to real requirements.

## Leave Human Rationale Locally

Comments, names, tests, and documentation are part of the implementation. A future engineer should not need hidden task context or a search across the repository to answer:

- what this code owns and what it deliberately does not own;
- which input or state is authoritative;
- why a guard, ordering constraint, or side effect exists;
- what can change safely and what stays stable;
- why an unusual fixture exists; and
- what a negative assertion proves.

Prefer clear parts and exact names. Add comments to explain the reason, permissions, behavior to preserve, ordering, recovery, and what a negative assertion proves. Do not narrate syntax that the code already states.

## Use Current Evidence Before External Claims

Before declaring a blocker external, a dependency necessary, a diff obsolete, or a change ready for review or integration, refresh the relevant current state. Compare against the current target and effective diff, not only a stale base or an earlier log. Separate pending work from failure, and obsolete runs from current evidence.

Do not broaden the implementation merely to work around an unrelated failure. Record the exact evidence and keep the repair in its owning scope unless the user explicitly expands the task.

## Review The Final Diff With Fresh Eyes

Before handoff, review the result as if another engineer wrote it. Use only the stated goal, repository instructions, current diff, touched files, and validation evidence. Check:

- the diff contains only the intended change;
- the design is the smallest complete solution;
- names, ownership, and control flow are clear;
- comments and tests explain reasons that are not obvious;
- no temporary branches, stale comments, unused imports, dead helpers, or formatting noise remain;
- validation proves the changed behavior rather than only compilation; and
- the requested finish point has actually been reached.

Fix defects that belong to the requested work. Run the affected checks again.

## Check Before Finishing

Complete only when:

- the current solution satisfies the acceptance criterion;
- current evidence still shows the change is needed;
- focused tests and proportional repository validation pass;
- the final diff is narrow, readable, and free of avoidable structure;
- every remaining moving part has a real requirement; and
- the requested finish point, and no later invented condition, is proved.

If a safe test through the complete workflow exists and the goal depends on it, run it through the real entrypoint. If it is unavailable or unsafe, state the exact blocker and do not overclaim completion.

## Red Flags

| Red flag | Correct action |
| --- | --- |
| “The first focused test passes, so stop.” | Run an adversarial review pass before handoff. |
| “The repository already has a draft framework.” | Keep it only when a current requirement needs it. |
| “This permission failure needs a new shared API.” | Classify the failure and exhaust smaller supported routes first. |
| “The old log proves the target branch is broken.” | Refresh current state and effective diff before claiming a blocker. |
| “The code is readable, so the why is obvious.” | Explain important reasons, permissions, and conditions for a safe change where they matter. |
| “The diff is already large, so keep going.” | Cut back to the smallest complete change and preserve unrelated work. |
| “The user asked for implementation, so publish it.” | Stop when the work the user actually requested is complete. |

## Tests

When changing this skill, read [tests/README.md](tests/README.md). Run the relevant scenarios with fresh, independent agents when the available tools permit.
