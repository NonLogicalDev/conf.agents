package harness_codex

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/NonLogicalDev/conf.agents/utils/agentbld/internal/engine"
)

// Codex rule files are flat and later layers take precedence.
func discoverCodexRules(dir string, rules map[string]string) error {
	ds, err := engine.Entries(dir)
	if err != nil {
		return err
	}
	for _, d := range ds {
		path := filepath.Join(dir, d.Name())
		if d.Type()&os.ModeSymlink != 0 || !d.Type().IsRegular() || filepath.Ext(path) != ".rules" {
			return fmt.Errorf("rules must be regular .rules files: %s", path)
		}
		rules[d.Name()] = path
	}
	return nil
}

// buildCodexPlan adds Codex rules and bundles resources in the source plan.
func buildCodexPlan(registry engine.Registry, sections []string, includeGuidance, includeRules bool, namespace, prefix string) (*engine.RenderPlan, error) {
	plan, err := registry.BuildPlan(sections, includeGuidance, "codex")
	if err != nil {
		return nil, err
	}
	plan.InstructionsName = "AGENTS.md"
	if includeRules {
		rules := map[string]string{}
		for _, layer := range plan.Layers {
			dir := filepath.Join(layer, "rules")
			if engine.IsLink(dir) {
				return nil, fmt.Errorf("symlinked source directory is not supported: %s", dir)
			}
			if _, err := os.Lstat(dir); err != nil {
				if os.IsNotExist(err) {
					continue
				}
				return nil, err
			}
			if !engine.IsDir(dir) {
				return nil, fmt.Errorf("rules root is not a regular directory: %s", dir)
			}
			if err := discoverCodexRules(dir, rules); err != nil {
				return nil, err
			}
		}
		plan.Dirs["rules"] = true
		for name, source := range rules {
			plan.Files[filepath.Join("rules", name)] = source
		}
	}
	// Skill links do not depend on the destination; resolve them in the generation.
	plan.Agents = string(engine.ResolveSkillReferences([]byte(plan.Agents), namespace, prefix))
	// Preserve destination placeholders until activation.
	if err := engine.BundleResources(plan, plan.Agents); err != nil {
		return nil, err
	}
	return plan, nil
}
