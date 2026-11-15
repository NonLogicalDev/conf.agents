---
name: Agent_DeliveryPlan
description: Create or update a short Markdown delivery plan that a person can edit for a fixed set of pull requests. Use when the operator asks for a PR delivery plan, delivery order, landing order, merge sequence, dependency tracks, or a document they can edit to choose the order.
---

# Agent Delivery Plan

## Outcome

Create a short document that lets the operator choose and revise the delivery order without rereading the whole task. Record one fixed snapshot of the PR set, its dependency tracks, and its real blockers.

This skill owns the delivery plan document. It does not merge, rebase, push, deploy, request review, send messages, or replace the project's PR workflow or durable project plan.

## Gather the snapshot

- Use the destination, PR set, ordering constraints, and existing document the operator gives you.
- Use live PR facts only when they were already verified or the current request permits the applicable PR workflow to verify them. Do not invent links, head revisions, checks, approvals, dependencies, blockers, or delivery state.
- Record the named base branch and frozen base revision when they are verified. Say that they are not verified when they are unknown.
- Preserve existing human choices, labels, order, tracks, and notes unless the operator asks to change them or verified evidence makes one stale.
- If the destination is unclear, return an editable draft and ask one focused question before writing a file.

## Write the document

Start with one short instruction that tells the operator what to edit. Put any hard ordering rule beside that instruction. Put the frozen base immediately before the delivery list when it is known.

Use this shape:

```markdown
Edit this list to choose the order and tracks. <rule that fixes part of the order, when one exists>

Frozen <base branch>: `<revision>`.

## Delivery order

1. [ ] <stable ID> — [#<number>: <title>](<url>); `<head revision>`; <short verified snapshot note>.

## Tracks

- TR<NN> — **<track name>:** [#<number>: <title>](<url>)
    - [#<number>: <title>](<url>)

## Blockers

- <real unresolved blocker>.
```

- Give every PR one stable `<letter><number>` ID. Choose a letter that helps the operator recognize its track or component. Preserve existing IDs, never reuse them for another PR, and do not renumber them when the order changes.
- Give every track one stable `TR<NN>` ID, such as `TR1` or `TR2`. Keep it with the same track when the track moves, is renamed, or gains or loses PRs. Do not reuse or renumber it. Keep track IDs separate from PR IDs.
- Keep each delivery row compact: stable ID, canonical PR link and title, verified head revision when known, then only the status or proof that helps choose the order.
- Use a numbered checklist for the delivery order. A checked row records the document's delivery state at that snapshot; it does not grant permission or prove current readiness.
- Start each track bullet with its stable track ID. Show independent tracks as sibling bullets. Indent only real descendants in a dependency chain. Reuse the same PR links and stable PR IDs from the delivery order.
- Keep blockers narrow. Include only unresolved approval, ordering, ownership, release, or external dependency facts that affect delivery. Write `- None known.` when no blocker is verified.
- Keep private conversation state, agent reasoning, links to the source thread, and unrelated project status out of the document.

## Check the result

- Every PR appears once in the delivery order and in the right track when a track applies.
- Every PR ID and track ID is unique, unchanged from earlier versions, and attached to the same PR or track everywhere.
- Every PR number, title, link, head revision, status note, dependency, and blocker is verified or clearly marked unknown.
- The frozen base and rules that fix the order are visible before the list.
- The document stays easy for the operator to edit by hand.
- No merge, review, deployment, or communication action happened merely because the document was written.

## Tests

When changing this skill, read [tests/README.md](tests/README.md). Run the relevant scenarios with fresh, independent agents when available.
