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
    │       └── SKILL.md
    ├── resources/
    │   └── {resource-group}/
    │       └── {resource-name}
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

A section can omit any component it does not need. Sections are evaluated in ascending lexicographic order, and a later section replaces an earlier skill with the same name or resource at the same relative path.

Instruction fragments are ordered by filename across all sections. When two fragments share a filename, only the version from the latest section appears in the generated `AGENTS.md`.

Use `{{%_resources_%}}` for installed resources and `{{%_skills_%}}` for installed skills in instruction fragments. Rendering replaces each placeholder with the host's corresponding absolute path, keeping machine-specific paths out of the source.

Each standalone skill lives in its own directory. Its `SKILL.md` frontmatter names the skill, describes when to use it, and matches the directory name. The complete winning skill directory is copied into the generated skill tree.

Resources are grouped by purpose and copied without changing their relative paths. A later resource can replace an earlier file or directory. Optional marketplaces remain in their source sections and are not automatically installed or flattened into standalone skills.

## Public sharing

Keep each section within its intended scope. Material intended for public release must be understandable and usable without restricted instructions, credentials, unpublished tools, machine-specific paths, or access to another section.

Export only explicitly approved public source files. Do not export rendered generations, installed links, other sections, repository history, caches, or source symlinks. A public checkout should remain useful without access to the workspace that produced it.

A full source workspace may provide `just export-public` to refresh a separate public checkout. That command is a source-workspace operation, not a requirement or command supplied by a standalone public checkout. Exporting does not commit, push, publish, render, or activate configuration.

## Rendered generations

Generated configuration lives directly under `${XDG_STATE_HOME:-$HOME/.local/state}/agent-config/`:

```text
${XDG_STATE_HOME:-$HOME/.local/state}/agent-config/
├── latest -> {yyyy}-{mm}-{dd}__{NN}/
└── {yyyy}-{mm}-{dd}__{NN}/
    ├── AGENTS.md
    ├── skills/
    │   └── {skill-name}/
    │       └── SKILL.md
    └── resources/
        └── {resource-group}/
            └── {resource-name}
```

`XDG_STATE_HOME` defaults to `$HOME/.local/state` when it is unset or empty. Each generation uses the local render date and a globally increasing ordinal. Rendering unchanged sources reuses the latest generation; changed sources create a new one without modifying existing generations.

Generated files and directories are read-only. Executable source files remain executable, and source material is copied instead of linked. The `latest` link is relative and updated atomically.

Changing the state root does not migrate existing generations or retarget active configuration. Existing agent links remain unchanged until an explicit sync.

## Render and verify

Run these commands from either repository checkout:

```bash
state_root="${XDG_STATE_HOME:-$HOME/.local/state}/agent-config"
renderer="scripts/reconcile-config"

"$renderer" render \
  --config-dir "$PWD" \
  --render-dir "$state_root"

"$renderer" check \
  --config-dir "$PWD" \
  --render-dir "$state_root"
```

Rendering creates or reuses a generation without changing active agent configuration. Checking compares the newest generation with the current sources without changing files.

## Activate explicitly

Activate a generation only when you intend to update the agent's installed configuration:

```bash
"$renderer" sync \
  --config-dir "$PWD" \
  --render-dir "$state_root" \
  --codex-dir "${CODEX_HOME:-$HOME/.codex}"
```

Syncing updates only links managed by the renderer. It preserves unrelated skills and resources and refuses to overwrite unmanaged configuration. Marketplace installation and publication remain separate actions.

## Prune old generations

Remove unused generations while retaining the five newest:

```bash
"$renderer" prune \
  --config-dir "$PWD" \
  --render-dir "$state_root" \
  --codex-dir "${CODEX_HOME:-$HOME/.codex}"
```

Pruning refuses to remove a generation referenced by active instructions, skills, or resources. It also requires `latest` to point to the newest generation before removing older state.

## Reset explicitly

Reset installed configuration only when you intend to remove its current links:

```bash
"$renderer" reset \
  --config-dir "$PWD" \
  --render-dir "$state_root" \
  --codex-dir "${CODEX_HOME:-$HOME/.codex}"
```

Reset backs up the existing `AGENTS.md` and the complete `skills/` directory as `~/.codex/skills.<epoch>.bak`. When both backups are needed, they use the same timestamp. It then removes stale direct skill links and only resource links managed by the previous generation root. Live skill directories, regular files, unrelated resource links, and nested links remain untouched. Reset does not render or activate another generation.

A full source workspace may additionally provide `just agents-render`, `just agents-check`, `just agents-sync`, `just agents-prune`, `just agents-reset`, and `just agents-test`. Those convenience commands depend on that workspace's Justfile; the direct renderer commands above work in either checkout.
