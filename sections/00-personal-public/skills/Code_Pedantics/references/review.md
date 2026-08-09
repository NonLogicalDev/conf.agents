# Code Review Pedantics

## Remove Common Review Problems

1. Identify the artifact, audience, intended behavior or message, and stricter
   repository instructions or rules for the destination.
2. Preserve facts, uncertainty, exact technical names, required templates,
   links, signatures, and the author's useful voice.
3. Replace vague nouns, generic praise, and unsupported confidence with exact
   behavior and available evidence.
4. Treat filenames, changed functions, comments, and test files as patch
   geography, not as the change's purpose or validation evidence.
5. Remove repeated ideas, canned contrasts, needless headings, mechanical
   inventories, and chronology the reader does not need.
6. Put the main point first and keep supporting context beside the claim or
   decision it explains.
7. Read the whole artifact again as the next engineer or recipient.

## Code Review Order

1. Identify process inputs, defaults, external behavior, values created at
   startup, and values that may change while the program runs.
2. Trace each changed value from definition to use before changing its lifetime
   or owner.
3. Remove unnecessary wrappers, imports, branches, repeated expressions, and
   speculative configuration.
4. Name important paths, environment variables, timeouts, asset names,
   markers, protocol values, and policy values.
5. Check file and module ownership before adding more responsibility.
6. Make human intent locally available through exact names, decomposition,
   documentation, and rationale comments.
7. Read the complete changed file as a junior engineer would.
8. Run focused tests, formatter checks, and proportional repository validation.

## Values That Stay Fixed While A Process Runs

Load an immutable process input once into a named value when the process does
not expect it to change. Use that value consistently. Do not repeatedly read
the same environment variable, rebuild the same derived path, or hide a stable
default inside an inline expression.

Preserve runtime semantics:

- A credential, timestamp, lease, or configuration created after startup is a
  runtime value. Pass it explicitly to the function that needs it.
- Keep a value intentionally refreshed during a process
  behind its refresh mechanism.
- Tests may replace a value set when the process starts if that is the behavior under
  test; do not keep repeated reads solely for test convenience.

Name a derived default in steps that expose the relationship, such as the
script directory and the final default path. Verify the relationship before
replacing a parent traversal that uses an index.

## Simplicity And Compatibility

Remove a function that only forwards identical arguments and return values
unless it adds a domain operation, policy, a public interface,
instrumentation, or a place to provide a dependency.

Before removing a wrapper, search every caller, import, and test. Update those
references or preserve a required public interface. Put that search
before the verdict in a review answer.

Avoid turning every literal into a constant. Name values that encode an
external contract, repeated value, derived default, important policy, or
configuration. Leave obvious values used once close to where they are needed.
