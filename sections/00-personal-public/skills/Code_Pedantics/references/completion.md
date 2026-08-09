# Check Before Finishing

Before handoff:

- search changed files for repeated process inputs and derived paths;
- inspect short wrappers, compatibility imports, and new indirection;
- verify constants did not freeze values that need to change at runtime;
- verify operator documentation names prerequisites;
- verify module ownership is clearer rather than merely moved;
- verify comments and tests explain reasons that are not obvious;
- run focused behavior tests, formatter checks, and proportional validation;
- record a safe check through the actual entrypoint when completion depends
  on it; and
- read the final diff again and check for problems introduced by the cleanup.

For prose, verify the artifact describes the intended change, preserves facts
and required contracts, uses concrete names, avoids diff narration, and makes
only validation claims that the evidence supports.
