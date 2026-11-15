## Skill references

- Write skill links as `[<native skill reference>](skill://<plugin name or local>/<SKILL_NAME>)`.
  - Use the skill's registered name and plugin name, or `local` for standalone skills.
  - Use the harness's native skill reference as the link text. If the harness has no special format, use the skill name.
  - For example, use [$@:Tasker_Stream](skill://@/Tasker_Stream) when that skill is available.
- Treat a skill link as an invocation:
  1. Find the named skill in the available skills for the current harness.
  2. Read the skill at the location supplied by the runtime. Report an unavailable skill instead of substituting another.
- The `skill://` URI is neither a filesystem path nor a web address.
- When referring to skills, omit machine paths, plugin cache paths, and versions. Render links normally; use code formatting only when explaining the syntax.
