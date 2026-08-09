# Engineering Prose

## Scope

Use this reference for engineering prose outside code: doc
comments, module and API documentation, READMEs, user guides, operator guides,
runbooks, design notes, migration notes, release notes, handoffs, and external
message drafts.

Commit messages and pull request descriptions have a separate contract in
[comms/commits-and-prs.md](comms/commits-and-prs.md). Explanations beside code also
uses [comments.md](comments.md).

## Let Markdown Wrap Naturally

Keep each Markdown prose paragraph or list item on one source line, however wide it needs to be. Let the editor or renderer wrap it; do not insert line breaks to meet 80 characters or any other preferred width.

Preserve line breaks that define paragraphs, headings, nested lists, blockquotes, tables, fenced code, and deliberate Markdown breaks. Follow actual repository or code formatting rules when they apply.

## Start With Audience And Outcome

Name who will read the artifact and what they need to decide, do, or understand.
Put that outcome first. A reader should not need the task conversation,
implementation chronology, or repository archaeology to recover the point.

Preserve facts, uncertainty, exact technical names, required templates, links,
signatures, and the author's useful voice. Do not make a prose cleanup stronger
than the evidence.

For external artifacts, keep internal coordination private. Code comments,
READMEs, commit messages, pull request descriptions, Slack, and Linear should
not expose private plan names or numbers, task conversation, coordinator
instructions, hidden sequencing rationale, or agent reasoning. Translate the
relevant fact into the behavior, contract, evidence, or risk a new reader needs.
Omit private context that does not change the reader's decision.

## Common Prose Pass

1. State the purpose, problem, behavior, decision, or requested action before
   background.
2. Use exact artifact, command, type, field, service, and behavior names.
3. Replace vague nouns, generic praise, and unsupported confidence with
   concrete behavior and evidence.
4. Keep supporting context beside the claim or decision it explains.
5. Remove repeated ideas, canned contrasts, mechanical inventories, and
   chronology the reader does not need.
6. Use headings only when they help a reader navigate distinct material.
7. Read the whole artifact again as someone without hidden context.

Treat filenames, changed functions, comments, and test files as patch geography,
not as the reason the change matters. Mention a file or symbol when its
contract, ownership, or verified result matters to the reader.

## Doc Comments And API Documentation

Document the contract that callers and maintainers need:

- responsibility and deliberately excluded responsibility;
- inputs, outputs, errors, side effects, and state transitions;
- who can change, delete, save, or access data, or make
  protocol decisions;
- behavior to preserve and operations to run in the correct order;
- conditions for safe changes and compatibility requirements worth explaining; and
- examples when the API is easy to misuse.

Do not narrate syntax, restate a signature, or paraphrase the implementation.
Put local rationale beside the decision in code; use doc comments for the
contract that applies when another module or API calls the code.

## READMEs And User Guides

Write for the person trying to accomplish a task, not for someone who already
knows the implementation.

A useful README or guide usually answers:

- What is this and who is it for?
- What prerequisite tools, permissions, configuration, or inputs are required?
- What is the shortest supported path to a successful first use?
- Which ways of invoking the tool are supported, and what can the tool change?
- What result should the user expect?
- What common failure modes are actionable, and how should the user diagnose
  them?
- What data, files, or external systems may be modified?

Prefer one verified example over several speculative examples. Keep commands
copyable and mark placeholders clearly. Do not claim a command works unless it
was verified or the limitation is stated.

## Design, Migration, And Operational Documents

State the current problem, chosen decision, relevant constraints, alternatives
considered, and evidence. Separate current behavior from proposed behavior and
known facts from assumptions.

For a migration or operational procedure, name:

- the starting state and target state;
- the ordered steps and why their order matters;
- validation before and after mutation;
- rollback or recovery conditions;
- who owns the work and who can change the affected state; and
- the exact remaining risk or blocker.

Do not bury the decision under discovery history. Keep rejected approaches only
when they explain a lasting tradeoff or prevent a likely future mistake.

## External Message Drafts

Preserve the permitted destination, recipient, factual evidence, uncertainty,
canonical links, attribution, and required signature. When asking for action,
make one clear request to one destination.

Use this reference to improve wording only. It does not give permission to
send or forward a message, write to another service, or claim that any action
or check took place.

## Keep Prose Durable

Prefer durable behavior, rationale, contracts, and stable references over
logs from a single run, temporary previews, short-lived check links, or a
history of debugging. Put evidence from the current run where the repository
expects it.

Before handoff, verify that the prose:

- names the real purpose and intended audience;
- uses concrete names and claims supported by evidence;
- preserves required contracts, links, uncertainty, and attribution;
- omits irrelevant patch geography and discovery history; and
- can be understood without the surrounding conversation.
