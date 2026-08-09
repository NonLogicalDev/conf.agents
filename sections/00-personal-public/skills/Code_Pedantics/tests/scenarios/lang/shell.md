# Shell Scenarios

## 25 Handle Shell Data Safely

### Prompt

Use `$Code_Pedantics`.

A Bash script derives its config directory in three functions, marks a value updated by a retry loop as `readonly`, and builds a command in a string before executing it with unquoted expansion. It enables `set -euo pipefail`, but a missing optional file is checked with a command whose expected nonzero exit is not handled explicitly. The author ran `shellcheck` but did not exercise the real script. Choose the review result. Do not modify files.

### Expectations

- Name fixed environment inputs and derived paths once near the top, and mark only constants that remain fixed for the process `readonly`.
- Preserve intentionally dynamic values for retries, traps, loops, and subprocess results instead of freezing them.
- Quote expansions and use an argument array instead of a command string; document any intentional word splitting or globbing.
- Handle the expected branch with a nonzero exit explicitly before relying on `set -euo pipefail`.
- Require `shellcheck`, the repository formatter, coverage of the success path, coverage of intentional nonzero exits, and a safe representative run through the real script with its normal environment contract.

### Pressure Variant

The script has one fixed derived path, a genuinely dynamic retry counter, and an optional file check whose nonzero status is an intentional branch.

- Keep the derived path readonly, keep the counter mutable, and make the optional-file branch explicit.
- Do not replace a short direct command sequence with functions unless a repeated operation, cleanup task, or behavior needs a name.
