package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type harnessHandler struct {
	parse   modeParser
	decode  func(json.RawMessage) (any, error)
	stub    func(any, syncFlags, []string, string) (any, error)
	build   func(any, syncFlags, []string, string) (*renderPlan, any, error)
	restore func(string, snapshot) (*renderPlan, error)
	prepare func(string, string, snapshot, *renderPlan, runtimeDestinations) (func() error, error)
}

type runtimeDestinations struct{ Codex, Marketplace, OpenCode, Claude, ClaudeMarketplace string }

func harnessPayloadDir(gen, harness string) string {
	path := filepath.Join(gen, harness)
	if harness == "codex" && !isDir(path) && isRegular(filepath.Join(gen, "AGENTS.md")) {
		return gen
	}
	return path
}

func restoreHarnessPlan(gen, harness string, s snapshot) (*renderPlan, error) {
	h, ok := harnessRegistry()[harness]
	if !ok || h.restore == nil {
		return nil, fmt.Errorf("unsupported harness in generation: %s", harness)
	}
	return h.restore(harnessPayloadDir(gen, harness), s)
}

func restoreHarnessPlans(gen string, s snapshot) (map[string]*renderPlan, error) {
	plans := map[string]*renderPlan{}
	legacy := len(s) == 1 && harnessPayloadDir(gen, "codex") == gen
	if !legacy {
		entries, err := os.ReadDir(gen)
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			if entry.Name() == snapshotFile && entry.Type().IsRegular() {
				continue
			}
			if _, ok := s[entry.Name()]; !ok || !entry.IsDir() || entry.Type()&os.ModeSymlink != 0 {
				return nil, fmt.Errorf("unexpected generation entry: %s", entry.Name())
			}
		}
	}
	for _, name := range sortedKeys(s) {
		p, err := restoreHarnessPlan(gen, name, s)
		if err != nil {
			return nil, err
		}
		plans[name] = p
	}
	return plans, nil
}

// harnessRegistry contains only implemented harnesses. Each handler owns its
// mode parser, snapshot schema, rendering and activation.
func harnessRegistry() map[string]harnessHandler {
	return map[string]harnessHandler{"opencode": openCodeHandler(), "claude": claudeHandler(), "codex": {
		parse:  func(value *string, flags syncFlags) (any, error) { return parseSyncMode(value, flags) },
		decode: decodeCodexSnapshot,
		stub: func(value any, flags syncFlags, sections []string, root string) (any, error) {
			mode, ok := value.(syncMode)
			if !ok {
				return nil, fmt.Errorf("invalid codex mode")
			}
			c := snapshotCodex{Version: 1, Agents: mode.agents && !flags.RulesOnly, Skills: mode.skills && !flags.InstructionsOnly && !flags.RulesOnly, Plugin: "none", PluginName: mode.pluginName, RenderDir: root, SectionsDirs: sections}
			if mode.plugin {
				c.Plugin = "full"
				if !mode.includeGuidance {
					c.Plugin = "skills"
				}
			}
			c.Legacy.SkillMarketplace = flags.SkillMarketplace
			c.Legacy.InstructionsOnly = flags.InstructionsOnly
			c.Legacy.RulesOnly = flags.RulesOnly
			return c, nil
		},
		build: func(value any, flags syncFlags, sections []string, root string) (*renderPlan, any, error) {
			mode, ok := value.(syncMode)
			if !ok {
				return nil, nil, fmt.Errorf("invalid codex mode")
			}
			namespace, prefix := "local", "$"
			if mode.plugin {
				namespace, prefix = mode.pluginName, "$"+mode.pluginName+":"
			}
			plan, err := buildCodexPlan(sections, !(mode.plugin && !mode.includeGuidance && !mode.skills), !(mode.plugin && !mode.agents), namespace, prefix)
			if err != nil {
				return nil, nil, err
			}
			pluginName := ""
			if mode.plugin {
				pluginName = mode.pluginName
			}
			plan, err = renderSkillReferences(plan, "codex", pluginName, mode.includeGuidance)
			if err != nil {
				return nil, nil, err
			}
			saved, err := makeSnapshot(plan, mode, flags, sections, root)
			if err != nil {
				return nil, nil, err
			}
			return plan, saved["codex"], nil
		},
		restore: func(dir string, s snapshot) (*renderPlan, error) {
			c := s["codex"].(snapshotCodex)
			return restorePlan(dir, c.FileModes, c.LayerCount, c.FragmentCount, "AGENTS.md")
		},
		prepare: prepareCodex,
	}}
}
