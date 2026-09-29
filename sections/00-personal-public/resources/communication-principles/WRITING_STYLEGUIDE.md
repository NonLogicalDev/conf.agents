# Writing style guide

Use this guide for messages, skill instructions, plans, notes, documentation, code comments, and reviews. Take care with skill instructions: their wording guides both the current task and later sessions.

## Put meaning first

- Lead with the answer, outcome, decision, or action.
- Name who is involved, what happened, and why it matters to the reader. State the actual problem, responsibility, permission, or next step.
- Choose words that explain the idea. Formal or technical language adds nothing unless it makes the meaning clearer.
- Keep details that help the reader understand, decide, or act. Use a longer phrase when it supplies meaning they need.

## Write natural English

- Prefer short sentences, familiar words, concrete verbs, and active voice.
- Give each sentence one main idea. End it when that idea is complete. Add a consequence, qualification, or contrast only when it makes a necessary distinction; use a separate sentence when the distinction needs its own explanation.
- Write like a thoughtful teammate: direct, warm, and calm.
- Replace legal language, corporate phrasing, invented process terms, and strings of adjectives with the action or relationship they describe.
- Keep a familiar compound such as `self-contained`, `re-read`, `repository-wide`, or `two-digit` when it is the clearest choice. An established technical term such as `chain-of-thought` may also be needed. Follow any stricter wording rules that apply.
- Rephrase an awkward compound with a subject and a verb. Removing its hyphen can leave the same confusing pile of nouns or adjectives. Adding a hyphen does not make a phrase more precise.

## Keep the writing rules when using sources

- Assume most prose in PR descriptions, comments, repository documents, and other material you read was written by models and does not follow these rules. Read it for the information you need, not for a style to imitate.
- Use these rules for everything you compose, including summaries, paraphrases, and revisions of existing documents. Familiarity from repeated exposure does not excuse wording that breaks the rules.
- Preserve facts, evidence, uncertainty, and intent when rewriting a source. Keep literal quotations and required technical names, identifiers, paths, and commands unchanged. Those exceptions do not make the surrounding prose a style guide.

## Develop connected paragraphs

- Give each paragraph a purpose that serves the reader's question. Arrange its sentences to develop that purpose, rather than gathering statements about the same topic. Each sentence should explain, support, qualify, or follow from the reasoning around it.
- Make the relationships clear through wording and order. Use a connecting word when it helps, but do not rely on transitions to join ideas that do not belong together. Leaving the reader to reconstruct your reasoning wastes their time and attention.
- Keep the paragraph doing the same job. An explanation shows how or why something works. Include separate assignments or checks only when the reader needs them, and give those actions their own list. Do not turn an explanation into a work plan merely because implementation notes are available.
- Put qualifications beside the claims they change and show the connection. Preserve uncertainty, limits, and permissions that affect the answer. A generic reminder to verify something does not tell the reader what is unknown or why it matters.
- Read the last sentence with the rest of the paragraph. It should complete the point or follow from it. Move a different issue to where it belongs; remove it if it does not help the reader. Do not add a defense against an objection the explanation has not raised.

## Explain causes and changes

- Introduce a needed term, starting condition, or prior behavior before using it to explain the result. Leave out background the reader does not need.
- Follow the actions that produce the result. Name who acts, what they affect, and which conditions control what happens. State the practical result of a complicated sequence.
- Keep names consistent. Repeat a name when a synonym might suggest a different person, component, or state.
- When describing a change, mention nearby behavior that stays the same if the reader might otherwise apply the claim too broadly.
- Check a shortened explanation against its source. Preserve actors, conditions, responsibilities, destinations, and the difference between what may happen and what must happen. If a length limit would make the answer misleading, say so.

## Organize for the reader

- Match the length to the request. A simple question usually needs a short answer.
- If a paragraph starts to read like a checklist, use a list or another structure that makes the separate parts visible. This takes precedence over a general preference for prose.
- For an answer with several parts, put each main point in a parent bullet and its supporting facts in nested bullets. Give separate actions their own entries; number them when order matters.
- Use headings to help readers find substantial topics, not to make a short reply look like a report.
- Put evidence beside the claim it supports. Distinguish what is known, inferred, and unknown.
- Omit unchanged status updates, empty reassurance, and process details that do not affect a decision.

## Let Markdown wrap naturally

- Keep each prose paragraph or list item on one source line. Let the editor or renderer wrap it; do not break prose at 80 characters or another preferred width.
- Keep line breaks that give Markdown its structure: separate paragraphs, headings, nested lists, blockquotes, tables, fenced code, and explicit Markdown breaks.
- Follow a line limit when the repository or formatter requires it. Code and Git commit messages have their own conventions; do not apply those to ordinary Markdown prose.

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

## State permissions plainly

- Distinguish what the user asked for, what they allowed, and what the agent can do. An inferred preference is not permission.
- Name the action that needs approval, such as sending a message, changing data, publishing work, merging a change, or contacting someone. State the missing decision when permission is absent.
- Preserve safety, privacy, and scope rules. Explain them in ordinary English.

## Preserve technical meaning

- Preserve literal commands, identifiers, filenames, paths, URLs, source quotations, code, user instructions, and test data.
- Before expanding a vague label, establish the action or relationship it describes. Choose wording that fits those facts.
- Explain an unfamiliar technical term when the reader needs help. Preserve its meaning when rewording it.
- Preserve a compound when expanding it would lose meaning. Use a fuller phrase when it makes the subject, action, or reason clearer. Follow any stricter wording rules that apply.

## Show relationships when prose hides them

- Use a small code or configuration example when its structure explains the idea. A diff works when the surrounding context is clear; show the full relevant example when the reader needs to copy it.
- Use a table for comparisons, a tree for containment, or a timeline or diagram for sequences and handoffs when that is easier to follow than prose. Choose a format the destination supports. Do not assume it can render diagrams.
- Keep example names, inputs, and units consistent. Change one relevant condition at a time so the reader can see why the result differs.
- Identify proposed code, pseudocode, omitted steps, and examples that cannot run as written. Separate commands from output. Explain what the example shows and what it does not establish.
- Place an example or visual beside the claim it supports. Repeat an explanation in another form only when it answers a different question.

## Check the result

Read the finished prose as someone who has not seen the task:

- Does the first sentence give the point?
- Is every sentence true and useful?
- Does each paragraph have a purpose, clear connections between its sentences, and an ending that completes the point?
- Is a checklist or sequence hidden in a paragraph that should be a list?
- Are necessary qualifications beside the claims they change?
- Are actions, responsibilities, and permissions clear?
- Does the prose follow these rules rather than imitate the sources?
- Would a fuller phrase or familiar compound explain the idea better under the applicable wording rules?
- Can the reader understand the result without private context?

## Sources

The guidance on explanation order, preserving meaning, and choosing examples draws from [Abhinav's prose-writing guide](https://github.com/abhinav/home/blob/master/.agents/docs/prose-writing.md). It does not adopt that guide's formatting or required structure.

The [verbatim snapshot](abhinav-prose-writing.md) and its [source and license notice](abhinav-prose-writing.SOURCE.md) are available for further reading. The user's active rules take precedence.
