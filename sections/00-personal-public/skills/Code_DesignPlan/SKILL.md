---
name: Code_DesignPlan
description: Write or revise a system design plan or architecture decision record, including plans modeled on a supplied example. Use for system behavior, component ownership, and design tradeoffs, not a PR delivery schedule or task tracker.
---

# Code Design Plan

Produce a plan that explains the intended system behavior, who owns it, and why the design fits the problem. Follow the user's requested format and length.

## Start with the request and evidence

Identify the behavior being designed, the reader, the decision status, and the available evidence. Separate existing behavior from proposed changes and unresolved choices. A design request does not authorize implementation or rollout.

When the user supplies a plan as an example, read it before choosing an outline. Follow both its organization and its prose: what it leads with, how it assigns responsibility, how it explains operations, and how much detail it gives each point. Apply those qualities to the new subject. Keep the new design as the opening subject; cite the example with other references when useful.

For PR order, assignments, or progress tracking, produce the requested delivery artifact instead. Do not turn an approved design into another architecture discussion.

## Organize around the decision

Use a short declarative title and truthful Status and Date. Lead with Decision, or Proposed decision when that makes the status clearer. State what the system will do and who owns the behavior.

Choose the sections that explain the design:

- **Context:** the current behavior, problem, and constraints that justify the change.
- **Ownership:** the responsibilities of components and the work their callers can delegate.
- **Subject-specific sections:** the model, interfaces, operations, and lifecycle. Name headings after the real concepts, such as “Reservation creation” or “Reading saved reports.”
- **Recovery and cleanup:** retry behavior, partial failure, cancellation, retention, or deletion when relevant.
- **Infrastructure and dependencies:** required services, storage, schema, compatibility, or rollout conditions when relevant.
- **Consequences:** concrete benefits, costs, and limitations.
- **Alternatives considered:** plausible alternatives and the reason for the proposed choice.
- **Future work or open choices:** unresolved decisions and work outside the proposal.
- **References:** evidence and relevant prior designs.

Scale the outline to the task. Omit empty or irrelevant sections. A short explicit user format takes precedence over this outline. Keep delivery steps separate from the design when both are requested.

## Write the behavior

Name the component and its action. Explain operation order, valid inputs, eligibility, state changes, visible outcomes, and failure or retry behavior where the reader needs them. Define unfamiliar terms briefly at their first use; do not open with a glossary.

Let the status establish that a proposal is a proposal. Place uncertainty beside the affected decision instead of repeating disclaimers throughout the plan. Do not present a recommendation as an agreed decision or planned behavior as working code.

Preserve supported interfaces, identifiers, facts, and evidence. Mark illustrative models or command shapes as proposed and non-executable when appropriate. Do not invent an existing API, measurement, deployment result, or source.

Follow applicable writing preferences. Use concrete verbs and connected explanations rather than packed noun phrases or commentary about the writing process.

For an optional example of prose and section choice, read [a fictional reading-position plan](references/reading-position-plan.md). A supplied user example takes precedence.

## Check the plan

Read it as someone who has not seen the conversation. Can the reader identify the decision, component responsibilities, operation sequence, failure outcomes, tradeoffs, and open choices? If a user supplied an example, compare both organization and wording with it. Matching headings alone is insufficient.

## Tests

When changing this skill, read [tests/README.md](tests/README.md). Run the relevant scenarios with fresh subagents that have empty context windows.
