# Agent configuration

This repository stores layered instructions, skills, resources, and optional marketplaces for coding agents. Its renderer combines source sections into an immutable generation. Activating that generation is a separate, explicit action.

## Source sections

Each section lives under `sections/` and uses the name `{nn}-{ns}-{suffix}`. The numeric prefix controls precedence, the namespace identifies its ownership, and the suffix describes its purpose.

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
    ├── opencode/plugins/
    │   └── {plugin-name}.js
    └── marketplace/
        ├── .agents/
        │   └── plugins/
        │       └── marketplace.json
        └── plugins/
            └── {plugin-name}/
                ├── .codex-plugin/
                │   └── plugin.json
                └── skills/
                    └── {skill-name}/
                        └── SKILL.md
```

A section can omit any component it does not need. Pass one or more directories containing sections with `--sections-dir`. Directories are evaluated in the order supplied; within each directory, sections are evaluated in ascending lexicographic order. Later sections take precedence, even when the sections are in different repositories.

Instruction fragments are ordered by filename across all sections. When two fragments share a filename, only the version from the latest section appears in the generated instructions.

Use `{{%_resources_%}}` for installed resources and `{{%_skills_%}}` for installed skills in instruction fragments. Sync replaces each placeholder with the selected harness destination's corresponding absolute path, keeping machine-specific paths out of the generation.

Each source skill lives in its own directory. Its `SKILL.md.d` contains Markdown fragments, ordered by filename. A later section replaces an earlier fragment with the same filename; different filenames are combined. Any fragment can supply YAML frontmatter. The renderer merges frontmatter in filename order: nested mappings are combined, and later values replace earlier scalars or lists. The final frontmatter must contain a name matching the skill directory and a description. It appears first in the generated `SKILL.md`, followed by fragment bodies in filename order. YAML aliases are not supported in skill fragments. This lets a `00-` fragment add instructions before `10-skill.md`. The fragment directory is not copied into the output.

A later fragment can extend a skill without replacing its earlier supporting files. Supporting files such as references and scripts follow section precedence when they share a path. A source skill can also contain a complete `SKILL.md` for compatibility with existing section directories; a later complete `SKILL.md` replaces the earlier skill before subsequent fragments are applied.

Inside a skill, `{{%_resources_%}}/group/file.md` becomes a path to `resources/group/file.md` within that skill. Paths are relative to the file containing the reference, so a file in `references/` uses `../resources/`. The renderer copies each referenced resource using the selected section's version. It also includes other known resources linked from copied Markdown files. A directory reference copies that directory; the bare placeholder copies all resources. Missing references and conflicts with files already in the skill stop the render. Source files remain unchanged.

`Agent_Guidance` is a generated mirror of the combined instructions. Its source fragments contain one `{{%_agents_%}}` placeholder, which the renderer replaces with the combined instructions from the selected sections. Resource paths point to copies inside the skill. The top-level `AGENTS.md` continues to use installed shared resources. These rules apply to both standalone skills and marketplace packages. Edit the instruction fragments to update the mirror; do not maintain a second copy in the skill. A generated mirror can contain private instructions and must not be exported as public source.

Resources are grouped by purpose and copied without changing their relative paths. A later resource can replace an earlier file or directory. Optional marketplaces remain in their source sections and are not automatically installed or flattened into standalone skills.

## Public sharing

Keep each section within its intended scope. Material intended for public release must be understandable and usable without restricted instructions, credentials, unpublished tools, machine-specific paths, or access to another section.

Export only explicitly approved public source files. Do not export rendered generations, installed links, other sections, repository history, caches, or source symlinks. A public checkout should remain useful without access to the workspace that produced it.

## Rendered generations

By default, generated configuration lives under `${XDG_STATE_HOME:-$HOME/.local/state}/agentbld/`:

```text
${XDG_STATE_HOME:-$HOME/.local/state}/agentbld/
├── latest -> {yyyy}-{mm}-{dd}__{NN}/
└── {yyyy}-{mm}-{dd}__{NN}/
    ├── config.json
    ├── codex/           # If selected
    │   ├── AGENTS.md
    │   ├── skills/
    │   ├── plugin-skills/  # If a Codex plugin is selected
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
        └── resources/
```

Skill links are resolved for each selected harness during rendering. When Codex selects a plugin, `plugin-skills/` holds its separately resolved copy; `skills/` holds the standalone copy. Sync selects the rendered files and resolves destination paths.

`XDG_STATE_HOME` defaults to `$HOME/.local/state` when it is unset or empty. Each generation uses the local render date and a globally increasing ordinal. Rendering unchanged sources and configuration reuses the latest generation; changes create a new one without modifying existing generations. Each generation records its selected harness modes and their resolved options in `config.json`. Harness destination flags are chosen at sync time and are not stored in the generation.

Generated files and directories are read-only. Executable source files remain executable, and source material is copied instead of linked. The `latest` link is relative and updated atomically.

Changing the state root does not migrate existing generations or retarget active configuration. Existing links to another render root are not managed by the new root, so sync may refuse to replace them. Pass the old location with `--render-dir` to keep using it.

## Render and verify

The configuration builder is a Go program in `utils/agentbld`. Install Go 1.24 or later and provide a writable Go build cache. When the dependencies in `go.mod` are not cached, Go also needs network access to download them. The `scripts/reconcile-config` launcher builds the program in a temporary directory on each invocation and runs it from the caller's working directory. To build a reusable binary instead, run `go -C utils/agentbld build -o /tmp/agentbld .` and set `renderer=/tmp/agentbld` below.

Run these commands from the repository root:

```bash
state_root="${XDG_STATE_HOME:-$HOME/.local/state}/agentbld"
renderer="scripts/reconcile-config"

"$renderer" render \
  --sections-dir "$PWD/sections" \
  --render-dir "$state_root"

"$renderer" check \
  --sections-dir "$PWD/sections" \
  --render-dir "$state_root"
```

Rendering creates or reuses a generation without changing active agent configuration. Checking compares the newest generation with the current sources without changing files.

The examples show explicit paths; `--render-dir` defaults to the state root above, and `--codex-dir` defaults to `${CODEX_HOME:-$HOME/.codex}`. Omit either option to use its default.

To combine this repository with another source, repeat the option in precedence order:

```bash
"$renderer" render \
  --sections-dir "$PWD/sections" \
  --sections-dir /path/to/local-agent-config/sections \
  --render-dir "$state_root"
```

Use the same ordered list for `check` as for `render`. Existing commands using `--config-dir` with a repository root continue to work; do not combine it with `--sections-dir` in one invocation.

## Activate explicitly

Render and activate the exact resulting generation in one command:

```bash
"$renderer" render --sync \
  --sections-dir "$PWD/sections" \
  --render-dir "$state_root"
```

To activate an existing generation, pass its ordinal from the render output. For example, generation `2026-11-14__02` has ordinal `2`:

```bash
"$renderer" sync 2 --render-dir "$state_root" --codex-dir /path/to/first-codex --codex-marketplace-home /path/to/first-home/.agents/plugins
"$renderer" sync 2 --render-dir "$state_root" --codex-dir /path/to/second-codex --codex-marketplace-home /path/to/second-home/.agents/plugins
```

Numbered sync uses the saved harness options and files; it does not require the source sections and does not change `latest`. Each sync resolves instructions for the selected destinations, so the same generation can be activated independently in multiple destinations. `render --sync` uses the destination flags supplied to that invocation. Generations created before `config.json` was introduced can still be discovered and pruned, but cannot be activated by number; render a new generation first.

For a generation containing multiple harnesses, sync them separately with bare `--mode` names, or omit `--mode` to sync all saved harnesses:

```bash
"$renderer" sync 2 --render-dir "$state_root" --mode=opencode --opencode-dir /path/to/opencode
"$renderer" sync 2 --render-dir "$state_root" --mode=claude --claude-dir /path/to/claude
```

Repeat bare `--mode` names to select multiple saved harnesses. Mode JSON belongs on `render`; numbered sync uses the saved options.

Sync updates only links managed by the renderer. It preserves unrelated skills and resources and refuses to overwrite unmanaged configuration. Marketplace installation and publication remain separate actions.

### Modes and skill marketplaces

Run `"$renderer" modes` for supported modes, options, defaults, and examples. Select modes when rendering. Repeat `--mode` for different harnesses: `codex`, `opencode`, and `claude`. Specifying the same harness twice is an error. Omitting `--mode` selects the default Codex configuration.

```bash
# Render and activate instructions, resources, rules, and standalone skills.
"$renderer" render --sync --sections-dir "$PWD/sections"

# Select all three harnesses in the same generation.
"$renderer" render --sync --sections-dir "$PWD/sections" --mode=codex --mode=opencode --mode=claude

# Render a full plugin generation and activate its configured outputs.
"$renderer" render --sync --sections-dir "$PWD/sections" --mode='codex:{"plugin":"full"}'

# Render and activate only a skill plugin, without changing Codex links.
"$renderer" render --sync --sections-dir "$PWD/sections" --mode='codex:{"agents":false,"plugin":"skills"}'
```

`agents` selects Codex instructions, resources, and rules. It defaults to `false` for a `"full"` plugin and `true` otherwise. `skills` selects standalone Codex skill links. It defaults to `false` when `plugin` is set and `true` otherwise; explicitly set `"skills":true` to use both standalone skills and the plugin. Set `plugin` to `"full"` or `"skills"`; `"skills"` omits `Agent_Guidance` from the plugin. A plugin-only render still creates a generation. Rendering without `--sync` does not update the marketplace or Codex links.

OpenCode and Claude Code each accept `agents` and `skills` boolean options, both enabled by default. OpenCode also accepts `plugins`, enabled by default. `agents` selects global instructions and resources; `skills` selects standalone skills. OpenCode's `plugins` selects native `.js` and `.ts` files directly under each section's `opencode/plugins/` directory. Later sections replace an earlier file with the same name. Symlinked or non-regular plugin sources are rejected. Disabled components are left untouched during sync.

Destination defaults can be overridden at sync time:

- `--codex-dir`: `$CODEX_HOME` or `$HOME/.codex`.
- `--opencode-dir`: `$OPENCODE_CONFIG_DIR`, or `${XDG_CONFIG_HOME:-$HOME/.config}/opencode`.
- `--claude-dir`: `$CLAUDE_CONFIG_DIR` or `$HOME/.claude`.
- `--codex-marketplace-home`: `$HOME/.agents/plugins`.

OpenCode uses `AGENTS.md`, `skills/`, `resources/`, and `plugins/` under its destination; Claude Code uses `CLAUDE.md`, `skills/`, and `resources/`. See [OpenCode instructions](https://opencode.ai/docs/instructions/), [OpenCode skills](https://opencode.ai/docs/skills/), [Claude Code instructions](https://code.claude.com/docs/en/memory), and [Claude Code skills](https://code.claude.com/docs/en/skills) for runtime discovery. OpenCode also discovers compatible skills in `~/.claude/skills` and `~/.agents/skills`; syncing its own directory does not remove those entries.

The builder validates skill frontmatter and the marketplace catalog, checks for local package edits and unsafe symlinks, and preserves other catalog entries and existing policy. Building a local plugin does not install or publish it. Omitting `Agent_Guidance` does not establish that other skills or resources are safe for public release.

The legacy `--skill-marketplace`, `--instructions-only`, and `--rules-only` options are available on `render`; add `--sync` to activate the result. They cannot be combined with each other or with `--mode`.

`--codex-marketplace-home` defaults to `$HOME/.agents/plugins`. The catalog is `<codex-marketplace-home>/marketplace.json` and the plugin is `<codex-marketplace-home>/skillset` by default. Set `plugin_name` in the mode JSON to choose another name. For another discoverable marketplace, use `<marketplace-root>/.agents/plugins` as the directory.

## Prune old generations

Remove unused generations while retaining the five newest:

```bash
"$renderer" prune \
  --sections-dir "$PWD/sections" \
  --render-dir "$state_root" \
  --codex-dir "${CODEX_HOME:-$HOME/.codex}"
```

Pruning refuses to remove a generation referenced by active managed instructions, skills, resources, rules, or native plugins across recorded destinations. It also requires `latest` to point to the newest generation before removing older state.

## Reset explicitly

Reset installed configuration only when you intend to remove its current links:

```bash
"$renderer" reset \
  --sections-dir "$PWD/sections" \
  --render-dir "$state_root" \
  --codex-dir "${CODEX_HOME:-$HOME/.codex}"
```

Reset backs up the existing `AGENTS.md` and the complete `skills/` directory as `~/.codex/skills.<epoch>.bak`. When both backups are needed, they use the same timestamp. It then removes stale direct skill links and only resource links managed by the previous generation root. Live skill directories, regular files, unrelated resource links, and nested links remain untouched. Reset does not render or activate another generation.

For legacy resource cleanup, a `--sections-dir` path named `sections` identifies the old `agent-config-rendered` directory beside its parent repository. A sections directory with another name does not identify a legacy render directory, so its resource links are left in place.

Run the Go tests with `go -C utils/agentbld test ./...`.

## Codex command rules

Store Codex `.rules` files in a section's `rules/` directory. Later sections replace earlier files with the same name. Rule sources must be regular files directly inside that directory; symlinks and nested directories are rejected. The private `guardian-usage.rules` source lives in `sections/02-oai-private/rules/`.

Rendering copies rules into each immutable generation's `rules/` directory. A Codex configuration that selects instructions installs managed links under `<codex-dir>/rules/`. Existing unrelated rules are preserved, conflicting unmanaged files cause an error, and deleting a source removes only its managed link on the next sync. Pruning preserves generations referenced by active rule links.

To update only rules while leaving instructions, skills, resources, and marketplace contents unchanged:

```bash
"$renderer" render --sync --rules-only --sections-dir "$PWD/sections" --render-dir "$state_root"
```

Sync installs the rule files; restart Codex to load them. Rule content changes permission behavior, so review and validate them with `codex execpolicy check` before syncing. The renderer copies rule contents without interpreting them. The legacy `reset` operation continues to reset instructions and skills; rule links are reconciled through sync.
