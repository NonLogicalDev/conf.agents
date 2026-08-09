# Skill Harden Behavioral Tests

Run each scenario with a fresh subagent that has an empty context window. Tell the subagent to invoke `$Skill_Harden` and give it the scenario prompt. Do not give it the expectations or intended answer. Do not let tests change shared files. If a test needs to write, use a temporary directory created for that task.

Capture the raw response and compare it with the expectations afterward. A scenario passes only when every expectation holds and no contrary behavior appears.

When testing a repair, first run the relevant scenario against the current guidance. After the edit, rerun that exact scenario. Also run a pressure variant or adjacent valid case when the scenario defines one.

Use [scenarios.md](scenarios.md) for the reusable gamut.
