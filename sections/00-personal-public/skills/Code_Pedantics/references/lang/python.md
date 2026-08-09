# Python Pedantics

## Process Inputs And Defaults

Load an environment value once at module import when it is fixed for the
process lifetime:

```python
SERVICE_TOKEN_ENV = os.getenv("SERVICE_TOKEN")
```

Use the value directly. If a credential is created or refreshed after import,
pass it explicitly:

```python
def fetch_resource(service_token: str | None = SERVICE_TOKEN_ENV) -> Resource:
    ...
```

Do not mutate `SERVICE_TOKEN_ENV` after import. A caller that creates a
temporary credential should pass `service_token=token`.

Name a derived default once:

```python
SCRIPT_DIRECTORY = Path(__file__).resolve().parent
DEFAULT_CONFIG_PATH = SCRIPT_DIRECTORY.parent / "config.yaml"
```

Avoid hiding a parent traversal inside an inline default. Verify the directory
relationship before replacing a numeric parent index with named paths.

## Functions And Imports

Remove a wrapper when it forwards the same arguments and return value without
adding a domain operation, policy, public interface, instrumentation, or a
place to provide a dependency.

Before removing a wrapper, search callers, imports, and tests. Update those
references or preserve a required public interface.

Remove `from __future__ import annotations` only when the supported runtime
implements every annotation syntax and the module does not depend on postponed
evaluation. Check forward references, runtime annotation inspection, and
circular imports first.

## Types, Comments, And Tests

- Use a dataclass when validated values travel together and field names improve
  call sites.
- Name repeated external identifiers and protocol values.
- Validate external data before constructing the typed result.
- Keep errors specific enough to identify the invalid field or contract.
- Use a module docstring to state what a script reads, validates, and writes.
- Add function docstrings for behavior that needs explanation and local comments
  where a guard prevents unsafe mutation.

Tests should cover the default set when the process starts, an explicit runtime
override when supported, malformed external data, behavior that leaves state
unchanged, and the exact data the code may modify. Before publication, run the
real Python entrypoint with a safe
representative input rather than only importing functions in tests.
