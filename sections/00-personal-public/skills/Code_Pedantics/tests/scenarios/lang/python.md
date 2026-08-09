# Python Scenarios

## 22 Preserve Python Runtime Behavior

### Prompt

Use `$Code_Pedantics`.

A Python CLI reads `SERVICE_TOKEN` at import time, then refreshes the token
after login by assigning a new value to that module global. Its default config
path is written inline as `Path(__file__).resolve().parents[2] / "config.yaml"`.
The CLI parses a JSON API response with a broad cast before constructing a
typed result. A cleanup also removes `from __future__ import annotations`
because the local interpreter is new enough. Choose the review result. Do not
modify files.

### Expectations

- Keep the environment default set at startup fixed after import; pass a
  refreshed runtime token explicitly instead of mutating the imported default.
- Name the script directory and derived config path once, and verify the
  directory relationship before replacing the numeric parent traversal.
- Validate the external JSON shape before constructing typed domain values;
  reject malformed fields with errors that identify the violated contract.
- Do not remove postponed annotations from local interpreter evidence alone.
  Check supported runtimes, forward references, runtime annotation inspection,
  and circular imports first.
- Require tests for the default set at startup, explicit runtime override,
  malformed external data, and a safe run through the real CLI entrypoint.

### Pressure Variant

The token is intentionally fixed for the whole process, the named path matches
the verified package layout, and the supported runtime plus imports prove that
no annotation depends on postponed evaluation.

- Accept the fixed default and annotation cleanup once the evidence is local
  and repeatable.
- Still require validation before untrusted JSON enters typed domain code.
