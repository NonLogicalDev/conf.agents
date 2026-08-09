---
name: Skill_Harden
description: Use when creating or updating skills, testing skill behavior with subagents, repairing unclear triggers or observed misbehavior, designing pressure tests, closing loopholes, addressing rationalizations, or validating skill folders before deployment.
---

# Skill Harden

## Overview

Write and improve skills by testing their effects. Observe how an agent behaves without the guidance. Write or revise the skill to address the observed failure. Test it under realistic pressure until the desired behavior holds.

Use the least guidance that helps future agents reliably do the right thing.

## Write In Plain English

A skill's language shapes later sessions. Write short, clear sentences. Use ordinary words, direct instructions, and calm explanations. Avoid legal language and long chains of awkward compound adjectives. Keep familiar terms such as `self-contained` when they make a sentence clearer. Explain what to do, what to check, when to ask, and what to leave unchanged. Preserve exact technical terms, skill names, commands, paths, links, examples, and safeguards.

Write each Markdown prose paragraph or list item on one source line, however wide it needs to be, and let the editor or renderer wrap it. Preserve meaningful Markdown structure; do not hard-wrap prose at 80 characters.

For the complete writing guidance, read `~/.codex/resources/communication-principles/WRITING_STYLEGUIDE.md`.

## Core Workflow

Follow this loop for new skills and for edits to existing skills:

1. Orient to the requested behavior, target runtime, destination folder, and existing skill conventions.
2. Capture a baseline failure before finalizing the skill or edit. For a new skill, test without the new guidance. For an update, test the current skill before the proposed change.
3. Write the smallest useful skill or update that addresses the observed failure.
4. Pressure test with fresh subagents that have empty context windows. First rerun the baseline scenario with the guidance, then add pressure that tempts the agent to misbehave.
5. Close loopholes. Turn observed excuses into clear guidance, a better skill description, signs to watch for, or a clearer section order.
6. Persist reusable behavioral tests when the scenario should protect future changes to the skill. Before adding or planning these artifacts, read `references/test-artifact-templates.md`.
7. Check the folder and test the actual behavior before using the skill.

Read `references/subagent-testing.md` before designing subagent validation, running pressure tests, or repairing a loophole found during testing.

Choose the validation shape from the skill's behavioral purpose. Discipline skills need pressure scenarios. Technique and pattern skills need scenarios that apply the technique, vary the situation, and test a counterexample. Reference skills need checks for retrieval, missing information, and written use. Use `references/subagent-testing.md` for the detailed test shapes.

## Write Effective Skills

Use the local skill format for the target runtime. Create a folder named after the skill, use the runtime's required entrypoint, and add the metadata the runtime's harness expects.

Keep the frontmatter focused on discovery:

- `name` is lowercase, hyphenated, and matches the folder.
- `description` says when to use the skill, including concrete triggers, symptoms, file types, tools, or failure modes.
- Avoid stuffing the description with the workflow. If the description fully summarizes the steps, an agent may follow the summary instead of loading the body.

Keep `SKILL.md` focused on the operating procedure:

- Explain only the knowledge a capable agent would not already know.
- Match specificity to risk. Fragile workflows need exact commands; work that calls for judgment needs clear principles and checks.
- Keep frequently used skills especially compact. Move bulky examples, rare cases, command references, and repeated workflow details into referenced files or existing tools.
- Prefer one strong example over several similar examples. Add another example only when it helps an agent make a different decision.
- Check size when a skill starts to sprawl. If a skill is hard to scan, split reference material out before adding more primary guidance.
- Put rare, bulky, or optional details in reference files directly under the skill.
- Add scripts only when deterministic behavior matters or the same code would otherwise be rewritten often.
- Add assets only when the skill uses them in outputs.
- Do not add README, changelog, installation guide, or process notes unless the runtime explicitly requires them.

## Update Existing Skills

Start from the observed malfunction, not from a broad rewrite impulse.

For each update:

1. Identify the failing behavior, the prompt or scenario that exposed it, and the part of the skill the agent likely relied on.
2. Preserve guidance that still works. Skill edits should close the gap without changing unrelated behavior.
3. Generalize the failure before drafting the patch. When changing guidance that affects an agent's behavior, first name the rule that the agent failed to follow:

   ```markdown
   Invariant:
   <general rule to preserve in future sessions>

   Observed symptom:
   <concrete failure that revealed the problem>

   Adjacent valid case:
   <nearby case the repair should still allow>
   ```

   Put only the general rule in the main guidance. Use observed symptoms and adjacent valid cases in examples, red flags, or retest scenarios instead.
4. Patch the smallest section that can prevent the same failure. Prefer precise wording, reordered emphasis, or a short table over a new framework. Write the primary repair as positive criteria: what future agents should optimize for, include, preserve, or verify. Avoid making the primary repair a negative rule about the observed symptom. Use negative wording only for details that are always invalid across the skill's intended use cases. When quoting patch text, quote only the wording that preserves the general rule. Do not add an alternative that applies only to the observed symptom. If placement depends on the target skill's structure, describe where the same general rule should go.
5. Retest the same scenario. Test a change to prose when that prose guides an agent's behavior.

Common repairs:

| Observed failure | Repair |
| --- | --- |
| Skill does not trigger | Add concrete symptoms to the description. |
| Skill triggers too broadly | Narrow the description and explain when not to use the skill. |
| Agent skips a required step | Move the step earlier and make the decision point explicit. |
| Agent follows a shortcut | Name the shortcut as a red flag and state the correct action. |
| Agent misses a detail | Move rare detail to a reference and link it at the decision point. |
| Agent treats testing as optional | Require a subagent or scenario check before deployment. |
| Agent overfits a repair to one observed example | State the general behavior to preserve, then test a nearby valid case. |

## Test With Subagents

Use subagents as independent validation surfaces. The point is to learn whether the skill transfers, not whether another agent can infer your intended answer.

When testing:

- Use fresh subagents with empty context windows for independent passes.
- Pass the skill name and artifacts that belong to the task, not your diagnosis, expected answer, or planned fix.
- Keep tests isolated. Ask the subagent for a decision, rationale, diagnostic plan, or patch sketch, not for changes to shared files, commits, deployments, or external systems. If a test requires tool use, allow inspection without changing files, or use a temporary directory created for that task outside the target repository.
- Capture raw output, decisions, skipped steps, and rationalizations.
- Treat subagent failures as diagnostics for the skill or the test setup.

Do not rely only on academic questions such as `What does this skill say?` Good skill tests make the agent respond under realistic constraints, not recite the skill.

## Persist Reusable Tests

When a skill update identifies behavior that future changes should preserve, leave behind reusable tests for future maintainers. Use a `tests/` directory when the target skill does not already define another test location. Any plan to add or update persisted tests is incomplete until it reads `references/test-artifact-templates.md` and uses that file for the test README, scenario template, and skill footer.

The lightweight default layout is:

```text
tests/
  README.md
  scenarios.md
```

`tests/README.md` explains how to run the reusable scenarios. It should say to use fresh subagents with empty context windows, hide expectations from the tested subagent, keep tests from changing files or confine them to a temporary directory created for the task, and compare the raw response with expectations afterward.

`tests/scenarios.md` records the reusable gamut. Store prompts, expectations, pressure variants, adjacent valid cases, and any special harness steps needed to reproduce the behavioral check.

Do not copy an actual failure into a saved scenario. Real failures are diagnostic evidence, not reusable test fixtures. Invent saved scenarios that test the same behavior, temptation, and expected decision without copying the real prompt, private context, names, paths, data, or exact incident shape.

When a real failure produced useful rationalization wording, capture the exact wording in the repair record. For persisted tests, translate that wording into expectations, red flags, or an invented pressure variant unless the quoted text is already generic and safe to reuse.

## Pressure Tests And Loopholes

Pressure tests should make the wrong behavior tempting. For guidance that needs to hold under pressure, combine at least three pressures:

- Time pressure: a deadline, urgent bug, or closing review window.
- Sunk cost: a draft or implementation already exists.
- Authority: a manager, reviewer, or senior engineer asks for the shortcut.
- Exhaustion: end of day, repeated failures, or a long session.
- Social pressure: fear of seeming slow, rigid, or difficult.
- Ambiguity: the prompt leaves room to claim the rule does not apply.
- Temptation to dismiss a small change: the edit looks too minor to test.

When an agent misbehaves, capture the exact reasoning and turn it into skill text.

| Rationalization | Skill repair |
| --- | --- |
| "This is just documentation." | Test prose that affects an agent's behavior. |
| "The skill is clear enough." | Require a realistic subagent pass before deployment. |
| "I'll test after shipping." | Test the skill before deploying it. |
| "This follows the spirit." | State the actual rule that the shortcut violates. |
| "The fix is obvious." | Require retesting the scenario that exposed the gap. |

## Validation

Before considering a skill ready:

- Run the runtime's validator when available.
- Inspect the file list for clutter or missing resources.
- Confirm the frontmatter is valid and says when to use the skill.
- Follow the sufficient testing checklist in `references/subagent-testing.md`.
- Update or create persisted tests when the repair protects behavior that future changes should preserve.
- For updates, rerun the scenario that exposed the original misbehavior.
- Report the behavioral evidence, not only that the Markdown looked correct.

If testing reveals a new loophole, revise the skill and repeat the relevant scenario. Do not treat the first passing run as enough when the skill enforces discipline under pressure.

## Tests

When changing this skill, read [tests/README.md](tests/README.md). Run the relevant scenarios with fresh subagents that have empty context windows. Read [references/test-artifact-templates.md](references/test-artifact-templates.md) when adding or updating persisted test artifacts.
