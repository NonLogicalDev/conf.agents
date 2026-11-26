package harness_claude

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/NonLogicalDev/conf.agents/utils/agentbld/internal/engine"
)

// Handler integrates Claude Code's personal instructions and skills with
// the shared immutable generation lifecycle.
func Handler(options *Options) engine.Handler {
	return engine.Handler{
		Summary: "Claude Code instructions, skills, and plugin.",
		ModeHelp: `Claude JSON options:
  agents       Boolean. Sync instructions and resources. Default: true,
               except plugin "full" defaults to false.
  skills       Boolean. Sync standalone skills. Default: true without a
               plugin; false when a plugin is selected.
  plugin       String: "full" includes Agent_Guidance; "skills" omits it.
               Omit to skip plugin output.
  plugin_name  String. Generated plugin name (default: "skillset").
  At least one of agents, skills, or plugin must be selected.`,
		Examples: `  agentbld render --sections-dir ./sections --mode='claude:{"agents":true,"skills":false}'
  agentbld render --sync --sections-dir ./sections --mode='claude:{"plugin":"skills"}'`,
		ResolveDestination: options.resolveDestination,
		LoadProfile:        options.loadProfile,
		ConfigReport: func(_ any, dest engine.Destination) (map[string]any, map[string]string) {
			return nil, dest.Profile
		},
		LinkComponents: []string{"skills", "resources"},

		Parse:  parseClaudeMode,
		Decode: decodeClaudeSnapshot,
		Stub: func(value any, sections []string, root string) (any, error) {
			mode, ok := value.(claudeMode)
			if !ok {
				return nil, fmt.Errorf("invalid Claude mode")
			}
			return Snapshot{Version: 1, Agents: mode.agents, Skills: mode.skills, Plugin: claudePluginType(mode), PluginName: mode.pluginName, RenderDir: root, SectionsDirs: append([]string{}, sections...), FileModes: map[string]uint32{}}, nil
		},
		Build: func(registry engine.Registry, value any, sections []string, root string) (*engine.RenderPlan, any, error) {
			mode, ok := value.(claudeMode)
			if !ok {
				return nil, nil, fmt.Errorf("invalid Claude mode")
			}
			plan, err := registry.BuildPlan(sections, !(mode.plugin && !mode.includeGuidance && !mode.skills), "claude")
			if err != nil {
				return nil, nil, err
			}
			plan.InstructionsName = "CLAUDE.md"
			namespace, prefix := "local", "/"
			if mode.plugin {
				namespace, prefix = mode.pluginName, "/"+mode.pluginName+":"
			}
			plan.Agents = string(engine.ResolveSkillReferences([]byte(plan.Agents), namespace, prefix))
			if err = engine.BundleResources(plan, plan.Agents); err != nil {
				return nil, nil, err
			}
			pluginName := ""
			if mode.plugin {
				pluginName = mode.pluginName
			}
			plan, err = engine.RenderSkillReferences(plan, "/", pluginName, "/"+pluginName+":", mode.includeGuidance)
			if err != nil {
				return nil, nil, err
			}
			saved := Snapshot{Version: 1, Agents: mode.agents, Skills: mode.skills, Plugin: claudePluginType(mode), PluginName: mode.pluginName, RenderDir: root, SectionsDirs: append([]string{}, sections...), LayerCount: len(plan.Layers), FragmentCount: plan.Fragments, FileModes: map[string]uint32{}}
			for path, source := range plan.Files {
				info, err := os.Stat(source)
				if err != nil {
					return nil, nil, err
				}
				saved.FileModes[filepath.ToSlash(path)] = uint32(info.Mode().Perm())
			}
			return plan, saved, nil
		},
		Restore: func(dir string, s engine.Snapshot) (*engine.RenderPlan, error) {
			saved, ok := s["claude"].(Snapshot)
			if !ok {
				return nil, fmt.Errorf("generation config has no Claude mode")
			}
			return engine.RestorePlan(dir, saved.FileModes, saved.LayerCount, saved.FragmentCount, "CLAUDE.md")
		},
		Prepare: prepareClaude,
	}
}

func claudePluginType(mode claudeMode) string {
	if !mode.plugin {
		return "none"
	}
	if mode.includeGuidance {
		return "full"
	}
	return "skills"
}
