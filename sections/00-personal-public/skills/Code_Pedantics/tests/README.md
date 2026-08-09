# Code Pedantics Behavioral Tests

Choose the smallest matching file under `scenarios/` for the guidance being changed. Run every affected category when a change crosses reference files.

- Review: [`review.md`](scenarios/review.md) covers 01 wrapper removal, 08 deliberate auditability, and 30 process inputs.
- Modularity: [`modularity.md`](scenarios/modularity.md) covers 02 safe TypeScript extraction.
- Comments: [`comments.md`](scenarios/comments.md) covers 03 reasons for permissions, 10 module documentation that leads with purpose, 31 long orchestration phases, 35 item documentation and local reasons, and 36 exposure of private plans.
- Validation: [`validation.md`](scenarios/validation.md) covers 04 safe updater proof.
- Engineering prose: [`engineering-prose.md`](scenarios/engineering-prose.md) covers 05 external prose, 07 README or user guide, 28 migration/design docs, and 38 natural Markdown wrapping.
- Tests: [`tests.md`](scenarios/tests.md) covers 06 fixtures that prove behavior and 27 regression tests.
- Completion: [`completion.md`](scenarios/completion.md) covers 32 checks before finishing.
- Commit content: [`commit-content.md`](scenarios/comms/commit-content.md) covers 09, 11-15, 29, 36, 37, 39, and 40.
- Workflow mapping: [`workflow-mapping.md`](scenarios/comms/workflow-mapping.md) covers 16-22.
- Python: [`python.md`](scenarios/lang/python.md) covers scenario 22 and checks values that may change while the program runs.
- Dockerfile: [`dockerfile.md`](scenarios/lang/dockerfile.md) covers 23 build ownership.
- Rust: [`rust.md`](scenarios/lang/rust.md) covers 24 intent, 33 module and file layout, and 34 identifier ownership.
- Shell: [`shell.md`](scenarios/lang/shell.md) covers scenario 25 and checks how the script handles data.
- TypeScript/JavaScript: [`typescript-javascript.md`](scenarios/lang/typescript-javascript.md) covers scenario 26 and checks TypeScript ownership and validation.

Run each selected scenario with a fresh isolated agent that has an empty context window. Give it the skill and scenario prompt from the matching category file, but not the expectations or intended answer. Do not let tests change shared files. If a test needs to write, use a temporary directory created for that task.

Capture the raw response and compare it with the expectations afterward. A scenario passes only when every expectation holds and no contrary behavior appears.

When an edit changes behavior, rerun the original scenario plus its pressure variant or adjacent valid case when one is defined.

For guidance about commit messages, ask a fresh reader to check substantial changes when practical: give a fresh agent only the candidate commit text and a checkout or diff for the task, not the task conversation or expected conclusion. The check passes when the agent can reconstruct the intent and independently assess whether the stated verification supports the safety claim.
