# Writing style guide

Use this guide when writing messages, skill instructions, plans, notes, documentation, code comments, or reviews. A skill's wording matters twice: it guides the current task and influences the language of later sessions.

## Put meaning first

- Lead with the answer, outcome, decision, or action.
- Say what happened, who is involved, and why the reader should care.
- Name the actual responsibility, permission, problem, or next step.
- Choose the wording that explains the idea best, not the wording that sounds the most formal or technical.
- Use a fuller phrase when it gives the reader useful meaning.
- Give the action a subject and a verb. Do not replace a hyphen with a string of nouns or adjectives.
- Omit details that do not help someone understand, decide, or act.

## Write natural English

- Prefer short sentences, familiar words, concrete verbs, and active voice.
- Give each sentence one main idea.
- Write like a thoughtful teammate. Keep the tone direct, warm, and calm.
- Avoid legal language, corporate phrasing, invented process terms, and strings of adjectives that make an ordinary idea harder to understand.
- Keep normal compounds when they are the clearest choice. Examples include `self-contained`, `re-read`, `repository-wide`, and `two-digit`.
- Keep an established technical term such as `chain-of-thought` when it is the exact term the subject calls for.
- Do not remove a hyphen mechanically. Do not add one to make a phrase sound technical.

## Let Markdown wrap naturally

- Keep each Markdown prose paragraph or list item on one source line, however wide it needs to be.
- Let the editor, reading view, or renderer handle wrapping. Do not break prose at 80 characters or any other preferred width.
- Keep intentional line breaks for paragraph boundaries, headings, nested lists, blockquotes, tables, fenced code, and explicit Markdown line breaks.
- Follow a line limit only when the target repository or formatter actually requires one. Keep code and real Git commit conventions separate from ordinary Markdown prose.

## Replace vague labels with their meaning

- Instead of `runtime-created`, say what is created and when:
  - `The program creates the worker when it starts.`
- Instead of `process-start` or `named process-start`, say what happens:
  - `Start the process named worker.`
  - `When the process starts, open the log.`
- Instead of `bounded wait`, say what to do:
  - `Wait up to ten seconds.`
  - `Set a short timeout.`
- Instead of `repeated exact UUID reads`, name the actual operation:
  - `Check the same request again by its UUID.`
- Instead of `authorization boundary`, name the actual permission:
  - `Send the message only if the user has asked you to send it.`
  - `Only the record owner can change this field.`
- Instead of an unexplained `seam` or `boundary`, describe the relationship:
  - `The API validates the input before saving it.`
  - `The payments service owns the refund.`
  - `This module creates the client; the caller supplies its settings.`
- Keep a familiar technical term if it describes something real more clearly than a longer replacement. Explain it if the reader needs the context.

## State permissions plainly

- Say what the user has asked for and what the agent can do.
- Name the action that needs permission: sending a message, changing data, publishing work, merging a change, or contacting someone.
- If permission is missing, say what decision is needed.
- Keep actual safety, privacy, and scope rules. Describe them in ordinary English instead of dressing them up as legal terms.

## Organize for the reader

- Match the length of the response to the question.
- Use a short answer for a simple question.
- When an answer has several parts, put the main point in a parent bullet and the supporting facts in nested bullets.
- Use a heading only when it helps the reader find a substantial topic.
- Put evidence next to the claim it supports.
- State what is known, what is inferred, and what is still unknown.
- Leave out repeated status updates, empty reassurance, and process details that do not change a decision.

## Preserve exact technical details

- Keep commands, identifiers, file names, paths, URLs, source quotations, code, user instructions, and test data exact.
- Keep the actual distinction between what the user allows and what the agent has only inferred.
- Preserve a useful compound when expanding it would lose meaning.
- Expand a compound when the fuller sentence makes the action, subject, or reason easier to understand.

## Check the result

For explanations of how something works or why it changed, also use the guidance below. A short answer does not need a longer structure when it already answers the reader's question.

- Does the first sentence tell the reader the point?
- Does each sentence say something true and useful?
- Is every action, responsibility, and permission clear?
- Would a familiar compound or a fuller phrase express the idea better?
- Can the reader understand the result without the task's hidden context?

## Explain causes and changes

- Introduce a necessary term, starting condition, or prior behavior before reasoning that depends on it. Leave out background the reader does not need.
- Follow the actions that lead to the result. Name who performs each action, what it affects, and any condition that controls it. After a complicated sequence, state its practical consequence.
- Keep names stable across the explanation. Repeating a name is better than a synonym that could suggest a different person, component, or state.
- When describing a change, mention a nearby unchanged behavior if leaving it out would make the reader apply the claim too broadly.
- Recheck a shortened explanation against its source. Do not lose an actor, condition, responsibility, destination, or the difference between something that may happen and something that must happen. If a length limit would force a misleading answer, say so rather than silently changing the meaning.

## Show relationships when prose hides them

- Use a small code or configuration fragment when its structure helps the reader understand the relationship. Use a diff for a change when surrounding context remains clear; show the whole relevant fragment when the reader needs to copy it.
- Use a table for comparisons, a tree for containment, or a timeline or diagram for sequences and handoffs when that makes the explanation easier to follow. Use a format the destination supports; do not assume a diagram renderer is available.
- When an example helps, keep it consistent across the explanation. Preserve its names, inputs, and units, and change one relevant condition at a time so the cause of a different result is clear.
- Label proposed code, pseudocode, omitted steps, and non-runnable examples. Separate commands from output. Explain what an example demonstrates and what it does not establish.
- Keep examples and visuals near the claim they support. Do not repeat the same explanation in several forms unless each answers a different question.

These additions about explanation order, preserving meaning, and choosing examples were informed by [Abhinav's prose-writing guide](https://github.com/abhinav/home/blob/master/.agents/docs/prose-writing.md). They are a selective synthesis, not an adoption of that guide's formatting or mandatory structure rules.

The [verbatim reference snapshot](abhinav-prose-writing.md) is available for further reading. See its [source and license notice](abhinav-prose-writing.SOURCE.md); the user's active rules take precedence over this reference.
