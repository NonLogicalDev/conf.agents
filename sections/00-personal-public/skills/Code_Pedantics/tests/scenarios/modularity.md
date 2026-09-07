# Modularity Scenarios

## 02 Extract One Safe TypeScript Component

### Prompt

Use `$Code_Pedantics`.

A root TypeScript component contains page markup, dialog markup, local widget state, shared editor state, persistence, and validation. A new dialog would add another large conditional block. Choose the next concrete refactor plan. Do not modify files.

### Expectations

- Name the mixed responsibilities and choose one smallest safe extraction.
- Prefer extracting the self-contained dialog or shared state before moving UI, state, and business rules together.
- Preserve handlers, selectors, keyboard behavior, and visible behavior.
- Validate the actual UI flow and review the result after the first slice.

### Adjacent Valid Case

The new markup is one small condition local to the page, with no reuse or new state.

- Keep it local when extraction would add more indirection than clarity.

## 42 Name Files By Responsibility

### Prompt

Use `$Code_Pedantics`.

A Python benchmark package has authored files named `store.py`, `descriptor.py`, `export.py`, `window.py`, `runner.py`, `telemetry.py`, `render.py`, and `comparison.html`. They split into coherent fixture, fetch, and report responsibilities. The package also contains `__init__.py`, `__main__.py`, `pyproject.toml`, `README.md`, generated historical output, and tests. Explain the naming change you would request without modifying files.

### Expectations

- Use the flat `<namespace>[_<subnamespace>][_<variant>]` basename pattern with a clear responsibility namespace.
- Give coherent families such as `fixture_store.py`, `fixture_descriptor.py`, `fixture_export.py`, `fixture_window.py`, `fetch_runner.py`, `fetch_telemetry.py`, `report_render.py`, and `report_comparison.html`.
- Add a subnamespace or variant only when it names a meaningful distinction rather than filling every slot.
- Keep framework and repository conventions such as `__init__.py`, `__main__.py`, `pyproject.toml`, and `README.md`.
- Mirror renamed subjects in tests with framework-required prefixes, such as `test_fetch_runner.py`.
- Update imports, tests, documentation, and asset references that use a renamed file.
- Leave generated historical output and unrelated files alone.
- Describe the rule as organization and naming, not a new framework.

### Pressure Variant

A reviewer asks to rename every file in the repository into three underscore-separated parts for visual consistency, including generated snapshots and conventional tool files.

- Reject the mechanical repository-wide rename.
- Preserve required conventional names, generated historical artifacts, and unrelated files.
- Use optional subnamespace and variant slots only when they carry meaning.

### Adjacent Valid Case

A repository has a stricter established filename convention or a build tool requires a conventional filename.

- Follow the repository or tool convention instead of forcing the general pattern.
