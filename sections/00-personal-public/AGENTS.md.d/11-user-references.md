## User communication references

### Identify the subject in every message

- Lead every chat message with the requested answer, outcome, or most important information and its subject. This includes progress updates and final answers.
- At the first mention of each subject in a message:
  - Put its established name or identifier beside the status, result, or decision request.
  - Link directly to it when possible. Use canonical URLs for external resources and file and line links when the location matters.
  - Without a usable link, give the known name or description and state any uncertainty that prevents identification.
  - Add a short description when the identifier alone does not explain what the subject is or why it matters.
- For tracked work, including PRs and tickets, give the identifier and title or short description at first mention in every message. Link the identifier when a usable URL exists; otherwise use plain text. Name the tracking system only to resolve ambiguity or explain an action.
- Reuse established names. Use a pronoun only when the reader can tell what it refers to within the message.
- Make each subject easy to identify and locate without reconstructing earlier messages, tool calls, or agent activity. Final answers must also make sense on their own.
- When several sources support an answer, give each a short bullet explaining what it establishes.

### Stable shorthand in private conversation

- Shorthands help the user refer to things; they are not abbreviations for the agent. Every use must be immediately followed by the full name or description, even on repeat mentions in the same message.
  - For tracked work, include the canonical identifier and title or description.
  - Never leave a shorthand alone in prose, headings, table cells, links, progress updates, or handoffs.
  - Earlier context, separate legends, link destinations, and previews do not replace the full reference.
- Label every item the user may refer back to, even if there is only one: artifacts, links, files, documents, PRs, issues, options, proposals, findings, and questions. Ordinary connective prose needs no label.
- Use `{<Letter><Number>}`, preferably one uppercase letter and two digits, with a leading zero when needed: `{P01}`, `{D12}`. Use more digits when necessary. Add braces to older labels without changing their letters, numbers, or meanings.
- Choose a meaningful letter for each kind of artifact and keep using it. `P` for PRs, `D` for documents, and `C` for commits are examples, not a fixed category system. If another kind of artifact suggests a letter already in use, prefer a different memorable letter and preserve the existing assignments.
- Prefer putting the shorthand and full reference together in the clickable text: `[{P36} PR #1300110 — authentication dependency health](https://flow.apps.openai.org/prs/openai/openai/1300110)`.
  - Without a usable link, keep the same pairing in plain text.
  - If a rich preview cannot contain the shorthand, put the shorthand followed by the full reference beside it.
- Give different items different labels. Keep each label through updates, title or URL changes, and reordering. Never renumber for presentation, swap meanings, or reuse a retired label for another item.
- Preserve the letters and their assignments in existing authorized task records, handoffs, and context summaries. Recover them before assigning new labels, including after resuming work. If you cannot recover an assignment, say so instead of guessing or silently changing it.
- Local VCS commits and revisions have a narrow exception:
  - Label them only when especially meaningful, such as the subjects of a comparison or central evidence in an investigation or decision.
  - Use canonical identifiers without shorthand for routine local checkpoints and incidental revisions, even when linked.
  - Keep existing assignments for later reuse, but omit the shorthand on routine mentions.
- Always label public artifacts, including PRs and tickets, in private chat, even when they relate to commits or revisions. A public artifact does not make the conversation public.

### Messages and documents for other people

- Keep private artifact and workstream shorthand, such as `{P01}` or `WS[03.02]`, out of material for a public audience or team unless the user explicitly asks to preserve it in that communication. This includes Slack, Linear, PR descriptions and comments, operational requests, and shared documents or sites.
- Treat a copyable draft according to its intended destination, even when preparing it inside Codex. Labels supplied as context are not permission to publish them. Replace or remove those labels unless the user has explicitly asked to keep them in the public material.
- Use descriptive names, canonical identifiers, and direct links. Keep tracker keys, PR numbers, severity and priority labels, and other identifiers the audience already uses when they are useful; these are not private shorthand.
