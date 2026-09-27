# Agent Structured Output Behavioral Tests

Test both format selection and explicit instructions. With the skill available as context, verify that ordinary replies stay plain and consulting the numbering convention does not require lifecycle receipts. When the user, AGENTS.md, or other applicable instructions request the skill or a format, verify compliance for the stated scope and duration. A standing instruction must persist across replies; a request for one report must not become a standing instruction. The format checks below apply when that format is chosen or required.

Run each scenario with a fresh subagent that has an empty context window when the available tools permit. Tell the subagent to use `$Agent_StructuredOutput` and give it the scenario prompt, but not the expectations.

Do not let tests change files or outside systems. Capture the raw response and compare it with the expectations afterward. A scenario passes only when every expectation holds and no contrary behavior appears.

When a scenario requires a Steps UI update, also capture the Steps UI tool call or inspect the resulting UI state. Verify that the update happened before substantive work and that the stored labels, order, and statuses match the expectations. Echoing the expected labels in prose is not sufficient.

For parallel assignments, verify that every running task is `in_progress`, visibly marked `[active]`, and names its verified `/root/...` worker. Only when a particular runtime demonstrably accepts one `in_progress` item may another running task use a different stored status; its visible label must remain `[active]`. Leave work that has not started as plain pending.

For behavior repairs, run the current guidance first, then rerun the same scenario after the edit. Also run the pressure scenario and nearby valid cases.

Verify that private task-chat shorthand stays private: public or team-facing Slack drafts, tracker text, pull request prose, and shared documents use descriptive names, canonical identifiers, and direct links instead, even when the draft is prepared inside private task chat and source notes contain private labels. Verify that ordinary answers, lone links, status, and results stay in normal chat without invented shorthand; a question widget appears only for an actual question that needs user input.

Verify that ordinary private Codex messages, agent-to-agent messages, and task-to-task coordination do not receive the public `uoleg-codex` signature. A copyable external or team-facing draft still follows its destination's required signature contract even when prepared inside Codex.

For progress during a task, verify that each material step or decision is actually sent to the operator before the next action. Verify a separate horizontal rule, actual local timestamp, correct lifecycle kind and action, `Outcome` parent bullet, and relevant nested evidence or next action for each event. A tool call, plan update, Steps UI change, draft, bare `Step` bullet, or final summary does not count as a lifecycle receipt sent in the task.
