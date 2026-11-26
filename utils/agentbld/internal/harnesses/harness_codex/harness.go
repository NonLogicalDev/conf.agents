package harness_codex

import (
	"fmt"

	"github.com/NonLogicalDev/conf.agents/utils/agentbld/internal/engine"
)

func Handler(options *Options) engine.Handler {
	return engine.Handler{
		Summary: "Codex instructions, skills, and plugin.",
		ModeHelp: `Codex JSON options:
  agents       Boolean. Sync instructions, resources, and rules. Default: true,
               except plugin "full" defaults to false.
  skills       Boolean. Sync standalone Codex skills. Default: true without a
               plugin; false when a plugin is selected.
  plugin       String: "full" includes Agent_Guidance; "skills" omits it.
               Omit this option to skip plugin output. Skills keeps agents on
               by default; full turns agents off by default.
  plugin_name  String. Name the generated plugin (default: "skillset"). Use
               lowercase letters, digits, hyphens, and underscores; the first
               character must be a letter or digit. Requires a plugin.

Explicit agents and skills values override their defaults. At least one of
agents, skills, or plugin must be selected. Render records the selected mode;
check compares that mode with the latest generation, and sync activates it.`,
		Examples: `  agentbld render --sync --sections-dir ./sections --mode=codex
  agentbld render --sync --sections-dir ./sections --mode='codex:{"plugin":"full"}'
  agentbld render --sync --sections-dir ./sections --mode='codex:{"plugin":"skills"}'
  agentbld render --sync --sections-dir ./sections --mode='codex:{"agents":true,"skills":true,"plugin":"full"}'
  agentbld render --sync --sections-dir ./sections --mode='codex:{"plugin":"skills","plugin_name":"my-skills"}'`,
		ResolveDestination: options.resolveDestination,
		LoadProfile:        options.loadProfile,
		ConfigReport:       configReport,
		Reset:              func(dest engine.Destination, root string) error { return Reset(dest.Config, root) },
		LinkComponents:     []string{"skills", "resources", "rules"},
		Default:            true,
		Parse:              func(value *string) (any, error) { return parseSyncMode(value) },
		Decode:             decodeCodexSnapshot,
		Stub: func(value any, sections []string, root string) (any, error) {
			mode, ok := value.(syncMode)
			if !ok {
				return nil, fmt.Errorf("invalid codex mode")
			}
			c := Snapshot{Version: 1, Agents: mode.agents, Skills: mode.skills, Plugin: "none", PluginName: mode.pluginName, RenderDir: root, SectionsDirs: sections}
			if mode.plugin {
				c.Plugin = "full"
				if !mode.includeGuidance {
					c.Plugin = "skills"
				}
			}
			return c, nil
		},
		Build: func(registry engine.Registry, value any, sections []string, root string) (*engine.RenderPlan, any, error) {
			mode, ok := value.(syncMode)
			if !ok {
				return nil, nil, fmt.Errorf("invalid codex mode")
			}
			namespace, prefix := "local", "$"
			if mode.plugin {
				namespace, prefix = mode.pluginName, "$"+mode.pluginName+":"
			}
			plan, err := buildCodexPlan(registry, sections, !(mode.plugin && !mode.includeGuidance && !mode.skills), !(mode.plugin && !mode.agents), namespace, prefix)
			if err != nil {
				return nil, nil, err
			}
			pluginName := ""
			if mode.plugin {
				pluginName = mode.pluginName
			}
			plan, err = engine.RenderSkillReferences(plan, "$", pluginName, "$"+pluginName+":", mode.includeGuidance)
			if err != nil {
				return nil, nil, err
			}
			saved, err := makeSnapshot(plan, mode, sections, root)
			if err != nil {
				return nil, nil, err
			}
			return plan, saved["codex"], nil
		},
		Restore: func(dir string, s engine.Snapshot) (*engine.RenderPlan, error) {
			c := s["codex"].(Snapshot)
			return engine.RestorePlan(dir, c.FileModes, c.LayerCount, c.FragmentCount, "AGENTS.md")
		},
		Prepare: prepareCodex,
	}
}
