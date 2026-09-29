## Skill references

- Write skill links as `[$<ns>:<SKILL_NAME>](skill://<ns>:<SKILL_NAME>)`.
  - Use the registered skill name and its plugin namespace or a configured alias.
  - Keep the namespace or alias and skill name identical in the link text and URI.
  - Use `@` for the default namespace only when applicable instructions define it, as in [$@:Tasker_ThreadState](skill://@:Tasker_ThreadState).
- Treat a skill link as an invocation:
  1. Resolve any alias from the applicable instructions. Ask if its mapping is missing; do not guess.
  2. Find the resolved namespace and name among the available skills.
  3. Read the skill at the location supplied by the runtime. Report an unavailable skill instead of substituting another.
- Keep the configured alias in the displayed link and URI. The `skill://` URI is neither a filesystem path nor a web address.
- When referring to skills, omit machine paths, plugin cache paths, and versions. Render links normally; use code formatting only when explaining the syntax.
