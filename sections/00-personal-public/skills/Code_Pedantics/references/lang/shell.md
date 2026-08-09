# Shell Pedantics

## Process Inputs And Paths

Name environment variables and derived paths near the top of the script when
they remain fixed while the process runs. Use `readonly` for values that
stay unchanged:

```bash
readonly SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
readonly CONFIG_PATH="${SCRIPT_DIR}/../config.yaml"
```

Do not repeatedly derive the same directory or read the same environment
fallback in several functions. Preserve intentionally dynamic variables and
values updated by traps, loops, or subprocess results.

## Control Flow

- Prefer a direct command sequence for a short script.
- Use a function when it names a repeated operation, isolates cleanup, or
  allows the behavior to be tested on its own.
- Quote expansions unless intentional word splitting or globbing is documented.
- Use arrays for command arguments instead of constructing command strings.
- Use `set -euo pipefail` only after commands that can return a nonzero
  status and unset
  variables are handled explicitly.

## External Contracts And Validation

Name queues, services, environment variables, timeouts, and file
modes when they encode an external contract or policy. Leave obvious
command flags at the call site.

Run `shellcheck` and the repository formatter after changes. Exercise the
success path and each intentional branch with a nonzero exit. Before publication, invoke
the real script with a safe representative input and its normal environment
contract. Verify filesystem and subprocess effects, then rerun it when
idempotence is part of the contract.
