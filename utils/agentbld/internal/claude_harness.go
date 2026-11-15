package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

// claudeHandler integrates Claude Code's personal instructions and skills with
// the shared immutable generation lifecycle.
func claudeHandler() harnessHandler {
	return harnessHandler{
		parse:  parseClaudeMode,
		decode: decodeClaudeSnapshot,
		stub: func(value any, _ syncFlags, sections []string, root string) (any, error) {
			mode, ok := value.(claudeMode)
			if !ok {
				return nil, fmt.Errorf("invalid Claude mode")
			}
			return snapshotClaude{Version: 1, Agents: mode.agents, Skills: mode.skills, Plugin: claudePluginType(mode), PluginName: mode.pluginName, RenderDir: root, SectionsDirs: append([]string{}, sections...), FileModes: map[string]uint32{}}, nil
		},
		build: func(value any, _ syncFlags, sections []string, root string) (*renderPlan, any, error) {
			mode, ok := value.(claudeMode)
			if !ok {
				return nil, nil, fmt.Errorf("invalid Claude mode")
			}
			plan, err := buildPlan(sections, !(mode.plugin && !mode.includeGuidance && !mode.skills), "claude")
			if err != nil {
				return nil, nil, err
			}
			plan.instructionsName = "CLAUDE.md"
			namespace, prefix := "local", "/"
			if mode.plugin {
				namespace, prefix = mode.pluginName, "/"+mode.pluginName+":"
			}
			plan.agents = string(resolveSkillReferences([]byte(plan.agents), namespace, prefix))
			if err = bundleResources(plan, plan.agents); err != nil {
				return nil, nil, err
			}
			pluginName := ""
			if mode.plugin {
				pluginName = mode.pluginName
			}
			plan, err = renderSkillReferences(plan, "claude", pluginName, mode.includeGuidance)
			if err != nil {
				return nil, nil, err
			}
			saved := snapshotClaude{Version: 1, Agents: mode.agents, Skills: mode.skills, Plugin: claudePluginType(mode), PluginName: mode.pluginName, RenderDir: root, SectionsDirs: append([]string{}, sections...), LayerCount: len(plan.layers), FragmentCount: plan.fragments, FileModes: map[string]uint32{}}
			for path, source := range plan.files {
				info, err := os.Stat(source)
				if err != nil {
					return nil, nil, err
				}
				saved.FileModes[filepath.ToSlash(path)] = uint32(info.Mode().Perm())
			}
			return plan, saved, nil
		},
		restore: func(dir string, s snapshot) (*renderPlan, error) {
			saved, ok := s["claude"].(snapshotClaude)
			if !ok {
				return nil, fmt.Errorf("generation config has no Claude mode")
			}
			return restorePlan(dir, saved.FileModes, saved.LayerCount, saved.FragmentCount, "CLAUDE.md")
		},
		prepare: prepareClaude,
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
