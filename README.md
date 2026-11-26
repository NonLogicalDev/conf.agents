# Agent configuration

Layered instructions, skills, and resources for Codex, OpenCode, and Claude Code.

## Render and verify

**Prerequisites**

- **Go** — 1.24 or later, with a writable build cache. Uncached modules also need network access.
- **PATH** — add `$GOBIN`, or `$(go env GOPATH)/bin` when `GOBIN` is unset.

From the repository root, install, render, and check. See [run or install](utils/agentbld/README.md#run-or-install).

```sh
go install github.com/NonLogicalDev/conf.agents/utils/agentbld@main

state_root="${XDG_STATE_HOME:-$HOME/.local/state}/agentbld"
agentbld render --sections-dir "$PWD/sections" --render-dir "$state_root"
agentbld check --sections-dir "$PWD/sections" --render-dir "$state_root"
```

- **Render** writes a generation. Add `--sync` to activate it.
- **Check** compares the newest generation and `latest` with the current sources and modes. It writes nothing.

## Activate explicitly

| To activate… | Command | Input |
| --- | --- | --- |
| The default Codex setup | `just sync` | The repository's `sections/` |
| The result of a new render | `agentbld render --sync` | Source sections and render options |
| A saved generation | `agentbld sync <number>` | A saved generation with `config.json` |

**Replay a generation**

1. Take the ordinal from the render output. Omit `--sections-dir` and `--config-dir`.
2. Omit `--mode` for all saved harnesses, or repeat bare names such as `--mode=opencode` for a subset.
3. Pass the destination flags for this run. Repeat the same ordinal to sync to other destinations.

- **Saved state** — numbered sync uses the generation's files and options. The `latest` pointer stays put.
- **Older generations** — a generation without `config.json` can be pruned. Render again before numbered sync.

```sh
agentbld render --sync --sections-dir "$PWD/sections" --render-dir "$state_root"

# Replace 2 with the ordinal printed by render.
agentbld sync 2 --render-dir "$state_root"
agentbld sync 2 --render-dir "$state_root" --mode=opencode --opencode-dir /path/to/opencode
agentbld sync 2 --render-dir "$state_root" --mode=claude --claude-dir /path/to/claude
```

## Source sections

`{nn}-{ns}-{suffix}` identifies a section. The prefix sets its order, the namespace identifies its owner,
and the suffix describes its purpose.

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

### Layer sources

1. Pass each source directory with `--sections-dir`, lowest precedence first.
2. Within each directory, sections sort by filename. The later section wins for a matching path.
3. Repeat the same source flags for `render` and `check`.

```sh
agentbld render \
  --sections-dir "$PWD/sections" \
  --sections-dir /path/to/local-agent-config/sections \
  --render-dir "$state_root"
```

- **Optional components** — omit any section component you do not need.
- **Plugin source** — put skills in `skills/`. A section-level `marketplace/` directory is rejected.
- **Legacy source flag** — `--config-dir` takes a repository root. Choose one source flag style per invocation.

### Instructions and skills

| Source | How agentbld assembles it |
| --- | --- |
| `AGENTS.md.d/*.md` | Combine in filename order. Later sections replace matching filenames. |
| `skills/<skill>/SKILL.md.d/*.md` | Combine bodies and merge frontmatter in filename order. |
| `skills/<skill>/SKILL.md` | Use as a new base. Later fragments can extend it. |
| Supporting files and resources | Keep relative paths. Later sections replace matching paths. |

**Skill frontmatter**

- **Required fields**
  - `name` matching the skill directory
  - `description`
- **Merge rules**
  - Nested mappings combine recursively.
  - Later scalars and lists replace earlier values.
  - Use YAML without aliases.
- **Output** — one `SKILL.md` with merged frontmatter followed by the ordered bodies.

### Select harnesses for a fragment

| To include a fragment in… | Set in YAML frontmatter |
| --- | --- |
| Every harness | Omit `harnesses` |
| A subset of harnesses | Set a nonempty list of distinct `codex`, `opencode`, or `claude` names |

```markdown
---
harnesses: [codex, claude]
---
Instructions for Codex and Claude Code.
```

- **Filter order** — select matching fragments, apply filename overrides, then remove the selector.
- **Instruction fragments** — the generated instructions contain the bodies without fragment frontmatter.
- **Complete `SKILL.md` files** — serve as the base for every harness. A harness with no base and no matching
  fragments omits the skill.

### Placeholders

| Location | Placeholder | Result |
| --- | --- | --- |
| Instructions | `{{%_skills_%}}` or `{{%_resources_%}}` | Absolute destination path at sync time |
| Skill | `{{%_resources_%}}/group/file.md` | Copy the file into the skill and make the link relative |
| Skill | Resource directory or bare `{{%_resources_%}}` | Copy that directory or all resources |
| `Agent_Guidance` | One `{{%_agents_%}}` | Mirror the combined instructions |

- **Resources**
  - Resolve relative to the file containing the reference. A file under `references/` uses
    `../resources/` for a sibling resource.
  - Include other known resources linked from copied Markdown files.
  - Stop the render if a referenced resource is missing or conflicts with a skill file.
- **Guidance mirror**
  - Update the source instruction fragments to change the mirror.
  - Bundled resources live inside the skill. The top-level instructions use installed resources.

## Rendered generations

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

| Item | Meaning |
| --- | --- |
| `{yyyy}-{mm}-{dd}__{NN}/` | Local render date and globally increasing ordinal |
| `latest` | Relative link to the newest generation, updated atomically |
| `config.json` | Resolved modes and their options |
| `skills/` | Standalone skills |
| `plugin-skillset/` | Separately resolved plugin skills |

- **Render**
  - Reuse the newest generation if its sources and resolved options match.
  - Write read-only copies and preserve executable permissions.
- **Sync**
  - Resolve destination paths from the flags supplied on this run.
  - Keep using the original `--render-dir` for its managed links. Changing the state root leaves existing
    generations and active links at the old root.

## Harnesses

Repeat `--mode` once per harness at render time. Codex is the default. Run `agentbld modes` for options.

```sh
agentbld render --sync --sections-dir "$PWD/sections" --mode=codex --mode=opencode --mode=claude
```

### Codex

| Mode option | Selects | Default |
| --- | --- | --- |
| `agents` | Instructions, resources, and rules | On except with `"full"` |
| `skills` | Standalone skill links | On when no plugin is selected |
| `plugin` | `"full"` or `"skills"`. The latter omits `Agent_Guidance`. | Off |
| `plugin_name` | Local plugin package name | `skillset` |

| Sync flag | Default |
| --- | --- |
| `--codex-dir` | `${CODEX_HOME:-$HOME/.codex}` |
| `--codex-marketplace-home` | `$HOME/.agents/plugins` |

- **Marketplace**
  - Render creates a generation even for plugin-only mode. Add `--sync` to write the marketplace.
  - Sync writes `marketplace.json` and the plugin side by side after validating frontmatter and the catalog.
  - Sync preserves other catalog entries and existing policy. Local package edits and unsafe symlinks block sync.
  - Register or install the plugin through Codex. Publish it separately.
  - For a separate discoverable marketplace, use `<marketplace-root>/.agents/plugins`.
- **Legacy render flags**
  - `--skill-marketplace`, `--instructions-only`, and `--rules-only`.
  - Choose one and use it without `--mode`.

```sh
# Generate and sync a full plugin.
agentbld render --sync --sections-dir "$PWD/sections" --mode='codex:{"plugin":"full"}'

# Generate and sync only the skills plugin.
agentbld render --sync --sections-dir "$PWD/sections" --mode='codex:{"agents":false,"plugin":"skills"}'

# Sync a full plugin and standalone skill links.
agentbld render --sync --sections-dir "$PWD/sections" --mode='codex:{"plugin":"full","skills":true}'
```

#### Codex command rules

1. Put regular `.rules` files directly in a section's `rules/` directory. Later sections win by filename.
2. Validate rules with `codex execpolicy check`.
3. Sync with `agents` enabled to create managed links under `<codex-dir>/rules/`, then restart Codex.

- **Rule files** — render copies their contents verbatim into the generation.
- **Rule links** — removing a source removes its managed link on the next sync. Unrelated rules stay in place.
  Unmanaged conflicts block sync. Prune protects generations referenced by active rule links.

**Update rules only**

```sh
agentbld render --sync --rules-only --sections-dir "$PWD/sections" --render-dir "$state_root"
```

### OpenCode

| Mode option | Selects | Default |
| --- | --- | --- |
| `agents` | `AGENTS.md` and `resources/` | On |
| `skills` | Standalone `skills/` | On |
| `plugins` | Native `plugins/` | On |

- **Destination** — `$OPENCODE_CONFIG_DIR` or `${XDG_CONFIG_HOME:-$HOME/.config}/opencode`.
  Override with `--opencode-dir`.
- **Native plugins** — regular `.js` or `.ts` files directly in `opencode/plugins/`.
  Symlinks and non-regular files are rejected. Later sections win by filename.
- **Disabled options** — sync skips the corresponding components.
- **Discovery** — OpenCode also searches `~/.claude/skills` and `~/.agents/skills`.
  Sync leaves those directories alone. [Instructions](https://opencode.ai/docs/instructions/) ·
  [Skills](https://opencode.ai/docs/skills/)

### Claude Code

| Mode option | Selects | Default without a plugin |
| --- | --- | --- |
| `agents` | `CLAUDE.md` and `resources/` | On |
| `skills` | Standalone `skills/` | On |

- **Destination** — `${CLAUDE_CONFIG_DIR:-$HOME/.claude}`. Override with `--claude-dir`.
- **Plugin setup** — [Claude Code in the agentbld guide](utils/agentbld/README.md#claude-code).
- **Discovery** — [Instructions](https://code.claude.com/docs/en/memory) ·
  [Skills](https://code.claude.com/docs/en/skills).

## Prune old generations

Prune keeps the five newest generations when both checks pass.

- **Newest generation** — `latest` points to it.
- **Active links** — no managed links point to the generations slated for removal. Prune checks recorded
  destinations for active instructions, skills, resources, rules, and native plugin links.

```sh
agentbld prune --sections-dir "$PWD/sections" --render-dir "$state_root"
```

## Reset explicitly

| Step | What reset does |
| --- | --- |
| Back up | Save instructions as `AGENTS.<epoch>.bak.md` and skills as `skills.<epoch>.bak`. |
| Remove | Remove stale direct skill links and managed legacy resource links. |
| Leave in place | Live skill directories, regular files, unrelated links, nested links, and rule links. |

- **Backups** — both names use the same timestamp when both backups are needed. The skills backup is complete.
- **Legacy resource links**
  - A source named `sections` identifies `agent-config-rendered` beside its repository.
  - A different source directory name leaves those links alone.
- **Activation** — use render and sync separately after reset.

```sh
agentbld reset --sections-dir "$PWD/sections" --codex-dir "${CODEX_HOME:-$HOME/.codex}"
```

## Public sharing

- **Export** — only approved public source files.
- **Public checkout** — keep it usable without private instructions, credentials, unpublished tools,
  machine-specific paths, or access to other sections.
- **Keep private**
  - Rendered generations, installed links, and private sections
  - Repository history, caches, and source symlinks
- **Before distributing a plugin** — review generated guidance and bundled resources.

## Development

- Run `go -C utils/agentbld test ./...`.
- See the [agentbld guide](utils/agentbld/README.md) for command details.
