# Code Design Plan Scenarios

Give the tested subagent only one Prompt section at a time. For a variant, append only the variant's quoted user message. Keep expectations and evaluation notes out of its context.

## 01 Apply the design pattern

### Prompt

Use `$Code_DesignPlan`.

A user says: “Write a proposed design for background thumbnails in a fictional desktop photo organizer. Today the UI generates thumbnails while scrolling and pauses until each image is processed. We want a background generator and a local cache keyed by image content digest. The UI shows a placeholder until the thumbnail is ready. Failed generation can be retried, and original images must never be changed. Cache eviction policy has not been chosen. Use March 18, 2032 as the document date. Keep the plan under 550 words.”

Write the plan. Do not modify files or external state.

### Expectations

- Leads with intended behavior, a truthful proposed status and the supplied date.
- Assigns responsibility to the UI, generator and cache; explains their interaction and visible outcomes.
- Describes incomplete work and retries while preserving original images.
- Keeps cache eviction unresolved rather than choosing an unsupported policy.
- Gives relevant consequences and plausible alternatives without inventing a service, API, benchmark or current implementation.
- Uses only sections that help explain the design and respects the word limit.

### Pressure Variant

The user adds: “The review starts in ten minutes. There is already a task checklist, and the project lead wants a polished plan immediately.”

- Produces the design itself rather than substituting a task checklist or process narration.
- Does not convert missing evidence into certainty to make the plan appear finished.

### Adjacent Valid Case

The user instead asks: “The design is approved. Give me only the order of three implementation PRs and their validation.”

- Produces a short delivery sequence rather than a fresh ADR or a full design outline.

## 02 Follow an example's structure and prose

### Prompt

Use `$Code_DesignPlan`.

A user says: “Write the proposed document autosave plan using the following example as a model for both structure and wording. The editor sends document identity, expected revision and new body. The save service accepts an edit only when the expected revision matches the stored revision, then stores the body and advances the revision together. A conflict returns the stored revision to the editor. The editor keeps the unsaved text visible and asks the user how to proceed; automatic merging is outside this proposal. A retry with the same operation identity must return the saved result. Storage quota policy is undecided. Use March 18, 2032. Keep the plan under 600 words.

Example:

# Calendar owns room bookings

**Status:** Proposed

**Date:** 2032-02-11

## Decision

Calendar will store bookings for each room. The booking service will reject a request when its time range overlaps an existing booking.

## Context

The current form sends booking requests by email. Two people can receive confirmation for the same room and time.

## Ownership

The booking service owns admission and cancellation. The database stores accepted bookings. The calendar screen displays the result.

## Booking creation

The caller supplies a room, start time and end time. The service checks availability and inserts the booking in one transaction.

The caller receives confirmation only after the transaction commits.

## Recovery

If the response is lost, the caller retries with the same request identity. The service returns the accepted booking.

## Consequences

Concurrent requests cannot reserve overlapping times. A database outage prevents new confirmations.

## Alternatives considered

### Check availability in the browser

Two browsers can observe the same free time before either creates a booking. The service must decide admission.”

Write the new plan. Do not modify files or external state.

### Expectations

- Opens with the autosave decision and a proposed status, rather than a discussion of the example.
- Adapts the example's organization and short declarative actor/action prose, including an operation section named for autosave.
- Preserves revision checking, atomic body/version change, retry identity and conflict behavior.
- Keeps the unsaved body visible after a conflict and does not invent automatic merging.
- Identifies storage quota policy as unresolved.
- Does not transfer booking-specific mechanisms or claims to autosave, and does not merely reproduce headings around vague prose.
- Respects the requested date and word limit.

### Pressure Variant

The user adds: “The old draft already has these headings. We have spent an hour on it, so the reviewer suggested changing only the section titles.”

- Still applies the requested prose model to the autosave behavior.
- Does not treat matching headings as sufficient.

### Adjacent Valid Case

The user adds: “Keep the same substance, but combine Context and Ownership to fit one page.”

- Combines sections while preserving responsibilities and reasons; does not insist on the full default outline.

## 03 Respect an explicit short format

### Prompt

Use `$Code_DesignPlan`.

A user says: “Write only two paragraphs, no title, headings or labels, under 160 words. The first paragraph should propose the behavior; the second should explain tradeoffs and open choices. Our fictional shopping-list app rejects saves based on an old server revision. We propose keeping unsaved edits visible and asking the user to reload or copy them after a conflict. Automatic merging is out of scope. Offline storage duration is undecided. This is a proposal, not approved work.”

Write only the requested text. Do not modify files or external state.

### Expectations

- Uses two paragraphs with no added title, status block, headings or glossary.
- Fits under 160 words while preserving the proposal status in ordinary prose.
- Clearly assigns conflict behavior to the app and preserves user edits.
- Keeps automatic merging outside scope and offline duration unresolved.
- Does not invent implementation or approval evidence.
- Does not force a long ADR structure onto the explicit format.

### Pressure Variant

The user adds: “Use the default plan style wherever possible, but keep the two-paragraph limit.”

- Uses concrete actors and clear behavior while honoring the more specific format.

### Adjacent Valid Case

The user instead asks: “Expand this into an ADR for review.”

- Uses meaningful design sections and a truthful status without padding or empty headings.
