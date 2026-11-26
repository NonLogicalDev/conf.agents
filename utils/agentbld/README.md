# agentbld

`agentbld` combines layered agent configuration into read-only generations. Its rendering code is mostly
independent of a particular agent harness; activation supports Codex, OpenCode, and Claude Code, with local
plugin marketplaces for Codex and Claude Code. Rendering and activation are separate operations. See the
[repository README](../../README.md) for the source layout and placeholder rules.

## Configuration layout

Source sections provide the instructions, skills, resources, and OpenCode plugins that agentbld combines. Each
section can omit components it does not need.

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

Rendered generations are stored under `${XDG_STATE_HOME:-$HOME/.local/state}/agentbld/` by default. Each
generation contains the selected harnesses; `latest` points to the newest generation.

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

## Run or install

Go 1.24 or later is required. Run the current `main` branch without installing a binary:

```sh
go run github.com/NonLogicalDev/conf.agents/utils/agentbld@main --help
```

Or install it for repeated use:

```sh
go install github.com/NonLogicalDev/conf.agents/utils/agentbld@main
```

Go installs `agentbld` to `$GOBIN`, or `$(go env GOPATH)/bin` when `GOBIN` is unset. Add that directory to your
`PATH`. Replace `main` with a release version, such as `v1.2.3`, to pin either command once a release is
available. This nested module uses tags such as `utils/agentbld/v1.2.3`; the repository does not yet have a
release tag.

From the repository root, use the installed binary with the local source sections:

```sh
renderer="agentbld"

"$renderer" render --sections-dir "$PWD/sections"
"$renderer" check --sections-dir "$PWD/sections"
"$renderer" render --sync --sections-dir "$PWD/sections"
"$renderer" sync 1  # Replace 1 with the ordinal printed by render.
```

Put each source skill's ordered Markdown fragments under `skills/<skill>/SKILL.md.d/`. The renderer
progressively merges YAML frontmatter from the fragments and places the result before the ordered fragment
bodies; later values replace earlier scalars and lists. Fragments with the same filename use the version from
the later section. A complete source `SKILL.md` remains supported. The rendered skill contains one `SKILL.md`
and no fragment directory.

To select harnesses for an instruction or skill fragment, put `harnesses: [codex, claude]` in its YAML
frontmatter. Omit the key to include the fragment for every harness. Available names are `codex`, `opencode`,
and `claude`. Filtering happens before section overrides, and the selector is removed from the rendered output.
See the [repository README](../../README.md) for details.

Repeat `--sections-dir` to combine directories from separate repositories. Directories are processed in the
order supplied, and sections within each directory are processed in lexicographic order; later sections take
precedence. The legacy `--config-dir` accepts a repository root and cannot be mixed with `--sections-dir`.

Run `"$renderer" --help` for commands and `"$renderer" modes` for supported harness modes and examples. Use
`--mode=<harness>[:<JSON object>]` when rendering; each harness defines its own options. `render` and `check`
use section sources. Numbered sync uses the saved generation and needs no section source. `--render-dir`
defaults to `${XDG_STATE_HOME:-$HOME/.local/state}/agentbld`.

Each generation records its resolved harness mode options in `config.json`. Harness destination flags are chosen
at sync time and are not stored in the generation. Repeat `--mode` to select `codex`, `opencode`, and `claude`
in one generation. Duplicate harness modes are rejected. A plugin-only render also creates a generation. Changes
to the selected configuration are part of the generation's identity.

## Operations

- `render` creates a generation or reuses the newest matching generation. Add `--sync` to activate that exact
  generation after rendering.
- `modes` lists supported modes, options, defaults, and examples.
- `check` compares the sources and selected configuration with the newest generation and its `latest` link
  without changing files. It exits with status 1 when out of date; errors exit with status 2.
- `sync <number>` activates the saved generation by its ordinal using its files and resolved configuration. It
  does not change `latest`. Omit `--mode` to activate all saved harnesses, or repeat bare names such as
  `--mode=opencode` to select a subset. Mode JSON is accepted on render. Supply `--codex-dir`,
  `--codex-marketplace-home`, `--opencode-dir`, `--claude-dir`, or `--claude-marketplace-dir` to choose
  destinations independently; `render --sync` uses the flags supplied to that invocation. Skill references are
  resolved during rendering.
- `prune` retains the five newest generations and refuses to remove one referenced by an active managed link. It
  requires `latest` to point to the newest generation before removing older generations.
- `reset` backs up existing instructions and the complete skills directory before removing direct skill links
  and resource links from the legacy render root. It does not activate a generation.

Generations created before `config.json` can still be discovered and pruned. Render a new generation before
using numbered sync. Generations are read-only copies, with a relative `latest` link. Sync preserves unrelated
entries and rejects unmanaged conflicts.

## Harnesses

Codex and Claude Code plugins use the section's `skills/` directory as their source. The renderer writes the
plugin copy to `plugin-skillset/` with plugin-specific skill links. Sync resolves destination paths.

Each selected harness has its own directory inside the generation, alongside `config.json`. Destination-specific
paths are materialized at sync time. See the [repository README](../../README.md) for source and destination
details.

### Codex

Set `--codex-dir` at sync time to override the default destination, `${CODEX_HOME:-$HOME/.codex}`.
Codex mode can select instructions, standalone skills, and a local plugin. Set `plugin` to `"full"` or
`"skills"`. The `"skills"` plugin omits `Agent_Guidance`. A plugin-only sync updates the local marketplace.
Review the included skills and resources before publishing the plugin.

The legacy `--skill-marketplace`, `--instructions-only`, and `--rules-only` options are available on `render`;
add `--sync` to activate the result. They cannot be combined with each other or with `--mode`.

`--codex-marketplace-home` defaults to `$HOME/.agents/plugins`; it contains `marketplace.json` and the generated
`skillset` plugin. Set `plugin_name` in the mode JSON to choose another name. For another discoverable
marketplace, use `<marketplace-root>/.agents/plugins` as the directory.

See [Codex in the repository README](../../README.md#codex) for mode examples and command rules.

### OpenCode

OpenCode mode syncs `AGENTS.md`, skills, resources, and native plugins. The `agents`, `skills`, and `plugins`
options default to `true`. The destination is `$OPENCODE_CONFIG_DIR` or
`${XDG_CONFIG_HOME:-$HOME/.config}/opencode` by default. Set `--opencode-dir` to override it at sync time.

### Claude Code

Claude Code mode syncs `CLAUDE.md`, skills, and resources, and can generate a local plugin marketplace.
Both `agents` and `skills` default to `true` when no plugin is selected. The destination is
`$CLAUDE_CONFIG_DIR` or `$HOME/.claude` by default. Set `--claude-dir` to override it at sync time.

For Claude Code, `plugin` can be `"full"` or `"skills"`; `"skills"` omits `Agent_Guidance`. Selecting a plugin
defaults standalone skills off, and `"full"` also defaults instructions and resources off. Set `agents` or
`skills` explicitly to override these defaults. `plugin_name` selects the name and defaults to `skillset`. A
plugin-only sync leaves the Claude configuration directory untouched. When instructions and a plugin are
selected with standalone skills disabled, sync removes previously managed standalone skill links while
preserving unmanaged skills.

The Claude marketplace defaults to `${XDG_DATA_HOME:-$HOME/.local/share}/agentbld/claude-marketplace`; use
`--claude-marketplace-dir` to select another root. Sync creates or updates the local marketplace and plugin but
does not register or install them. To use the default plugin, run
`claude plugin marketplace add <marketplace-directory>` and then `claude plugin install skillset@agentbld`. Use
the configured plugin and marketplace names if they differ. See the
[Claude Code marketplace documentation](https://code.claude.com/docs/en/plugin-marketplaces).

## Development

Run the tests from the repository root:

```sh
go -C utils/agentbld test ./...
```

The command entry point is in `utils/agentbld`; implementation code lives under `internal/`. Harness-specific
files use `codex_`, `opencode_`, and `claude_` prefixes so their scope is clear. Keep general layer selection,
resource bundling, and generation handling separate from Codex-specific activation and marketplace handling.
