# Changelog

Notable changes to this project are recorded here. The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## Unreleased

### Added

- Select instruction and skill fragments by harness using a `harnesses` list in their YAML frontmatter.
- Generate local Claude Code plugins and marketplaces with configurable plugin names, optional guidance, and independent sync destinations.
- Assemble skill instructions from ordered `SKILL.md.d` fragments across sections, with later sections overriding matching filenames and frontmatter merging in order.

- Support OpenCode and Claude Code alongside Codex, including global instructions and skills, and native OpenCode plugins. Select multiple harnesses with repeated `--mode` flags, and activate a saved subset with bare mode names on numbered sync.
- Record resolved harness mode options in each generation and activate a saved generation by ordinal without the source sections. Choose harness destinations at sync time so the same generation can be activated independently in multiple locations.
- Render and activate the resulting generation with `render --sync`; plugin-only renders also create generations.
- Combine section directories from separate repositories with repeated `--sections-dir` flags, in the order supplied.
- Select Codex and local plugin outputs with `render --mode`, including plugin-only output and optional `Agent_Guidance` inclusion.
- Default standalone Codex skills off when a plugin is selected, with an explicit option to use both.
- Default Codex instructions, resources, and rules off for a full plugin; `agents` can explicitly enable them.
- Place the personal marketplace catalog and generated plugin together under `$HOME/.agents/plugins` by default, configurable with `--codex-marketplace-home`; retain `--home` for compatibility.
- Name generated plugins `skillset` by default; `plugin_name` in the mode JSON selects another name.
- Print resolved command configuration before performing an operation.
- Inspect supported modes, options, and examples with `agentbld modes`.

### Changed

- Resolve harness-specific skill links during rendering, with separate Codex plugin and standalone skill copies when their namespaces differ.

- Store this repository’s skill instructions in `SKILL.md.d` directories; generated skills retain a single `SKILL.md`.

- Replaced the Python configuration renderer with the Go `agentbld` command, using Kong for command-line parsing.
- Organized the builder's implementation under `internal/` and prefixed Codex-specific filenames with `codex_`.
- Modeled each builder operation as a Kong command with contextual help. Rendering selects the harness configuration; numbered sync uses the saved generation.
- Default rendered generations to `${XDG_STATE_HOME:-$HOME/.local/state}/agentbld` and the Codex destination to `${CODEX_HOME:-$HOME/.codex}`; both paths can be overridden with flags.
- Check Codex destination conflicts before replacing the local plugin and marketplace catalog.
