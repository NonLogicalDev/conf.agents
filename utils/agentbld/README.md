# agentbld

Render layered agent configuration into read-only generations and sync them to your chosen harnesses.

## Run or install

**Prerequisite** — Go 1.24 or later. Add `$GOBIN` to `PATH`, or use `$(go env GOPATH)/bin` if `GOBIN` is unset.

| Use | Version suffix |
| --- | --- |
| Track the main branch | `@main` |
| Pin a tagged release | `@vX.Y.Z`, from tag `utils/agentbld/vX.Y.Z` |

```sh
go run github.com/NonLogicalDev/conf.agents/utils/agentbld@main --help
go install github.com/NonLogicalDev/conf.agents/utils/agentbld@main
```

## Configuration layout

- See [source sections](../../README.md#source-sections) for layering, fragments, and placeholders.

```text
sections/
└── {nn}-{ns}-{suffix}/
    ├── AGENTS.md.d/
    │   └── {order}-{section}.md
    ├── skills/
    │   └── {skill-name}/
    │       ├── SKILL.md.d/
    │       │   ├── 10-skill.md
    │       │   └── {order}-{fragment}.md
    │       ├── references/
    │       └── scripts/
    ├── resources/
    │   └── {resource-group}/
    │       └── {resource-name}
    └── opencode/plugins/
        └── {plugin-name}.js
```

- Override the generation directory with `--render-dir`.
- `config.json` records resolved modes. Destination paths are chosen at sync time.
- The selected source files and resolved options determine whether render reuses the newest generation.

```text
${XDG_STATE_HOME:-$HOME/.local/state}/agentbld/
├── latest -> {yyyy}-{mm}-{dd}__{NN}/
└── {yyyy}-{mm}-{dd}__{NN}/
    ├── config.json
    ├── codex/           # If selected
    │   ├── AGENTS.md
    │   ├── skills/
    │   ├── plugin-skillset/  # If a plugin is selected
    │   ├── resources/
    │   └── rules/
    ├── opencode/        # If selected
    │   ├── AGENTS.md
    │   ├── skills/
    │   ├── resources/
    │   └── plugins/
    └── claude/          # If selected
        ├── CLAUDE.md
        ├── skills/
        ├── plugin-skillset/  # If a plugin is selected
        └── resources/
```

## Operations

| Command | Use it to | Input |
| --- | --- | --- |
| `render` | Create or reuse a generation | Source sections |
| `render --sync` | Render and activate that exact generation | Source sections and destinations |
| `check` | Verify the newest generation and `latest` without writing | Source sections and render options |
| `sync <number>` | Activate a saved generation | Ordinal and destinations |
| `modes` | List harness options and examples | — |
| `prune` | Keep the five newest generations | State root with no active links to older generations |
| `reset` | Back up Codex config and remove stale managed links | Codex destination and source paths |

- **Check exit codes**
  - `0` — current
  - `1` — stale
  - `2` — error
- **Numbered sync**
  - Uses the saved files and options. Omit source flags.
  - Leaves `latest` pointing to the newest generation.
  - Repeat bare `--mode` names to select harnesses.
  - Requires `config.json`. Render again before syncing a generation that predates it.

```sh
# From the repository root.
agentbld render --sections-dir "$PWD/sections"
agentbld check --sections-dir "$PWD/sections"
agentbld render --sync --sections-dir "$PWD/sections"
agentbld sync 1  # Replace 1 with the ordinal printed by render.
```

## Harnesses

- **Render** — pass `--mode=<harness>[:<JSON object>]` once per harness.
- **Sync**
  - Supply destination flags. Sync resolves placeholders for those destinations.
  - Unrelated entries stay in place. Unmanaged conflicts block sync.
- **Plugin packaging** — Codex and Claude Code use source `skills/` and render a separate `plugin-skillset/`.
  Sync writes local packages. Register or install them through the harness.

### Codex

| Sync flag | Default |
| --- | --- |
| `--codex-dir` | `${CODEX_HOME:-$HOME/.codex}` |
| `--codex-marketplace-home` | `$HOME/.agents/plugins` |

See [Codex examples, mode options, legacy flags, and command rules](../../README.md#codex).

### OpenCode

| Mode option | Selects | Default |
| --- | --- | --- |
| `agents` | `AGENTS.md` and resources | On |
| `skills` | Standalone skills | On |
| `plugins` | Native plugins | On |

`--opencode-dir` overrides `$OPENCODE_CONFIG_DIR` or `${XDG_CONFIG_HOME:-$HOME/.config}/opencode`.

### Claude Code

| Plugin selection | `agents` default | `skills` default | Plugin includes `Agent_Guidance` |
| --- | --- | --- | --- |
| No plugin | On | On | — |
| `"skills"` | On | Off | No |
| `"full"` | Off | Off | Yes |

- **`agents`** selects `CLAUDE.md` and resources. **`skills`** selects standalone skills.
- **Overrides**
  - Set `agents` and `skills` explicitly to override the defaults.
  - Set `plugin_name` to choose a name. Defaults to `skillset`.
- **Standalone links**
  - Plugin-only sync leaves the Claude configuration directory untouched.
  - When syncing instructions and a plugin with `skills` off, sync removes previously managed skill links.

| Sync flag | Default |
| --- | --- |
| `--claude-dir` | `${CLAUDE_CONFIG_DIR:-$HOME/.claude}` |
| `--claude-marketplace-dir` | `${XDG_DATA_HOME:-$HOME/.local/share}/agentbld/claude-marketplace` |

**Install the generated plugin after sync**

```sh
claude plugin marketplace add "${XDG_DATA_HOME:-$HOME/.local/share}/agentbld/claude-marketplace"
claude plugin install skillset@agentbld
```

- Use your configured path and names if they differ. See the
  [Claude Code marketplace documentation](https://code.claude.com/docs/en/plugin-marketplaces).

## Development

- Command entry point in `utils/agentbld`. Implementation in `internal/`.
- **Harness-specific code** — `codex_`, `opencode_`, and `claude_` files.
- **Shared code** — source layering, resource bundling, and generation handling.
- Run `go -C utils/agentbld test ./...` from the repository root.
