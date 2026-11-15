package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

// Codex rule files are flat and later layers take precedence.
func discoverCodexRules(dir string, rules map[string]string) error {
	ds, err := entries(dir)
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
func buildCodexPlan(sections []string, includeGuidance, includeRules bool, namespace, prefix string) (*renderPlan, error) {
	plan, err := buildPlan(sections, includeGuidance)
	if err != nil {
		return nil, err
	}
	if includeRules {
		rules := map[string]string{}
		for _, layer := range plan.layers {
			dir := filepath.Join(layer, "rules")
			if isLink(dir) {
				return nil, fmt.Errorf("symlinked source directory is not supported: %s", dir)
			}
			if _, err := os.Lstat(dir); err != nil {
				if os.IsNotExist(err) {
					continue
				}
				return nil, err
			}
			if !isDir(dir) {
				return nil, fmt.Errorf("rules root is not a regular directory: %s", dir)
			}
			if err := discoverCodexRules(dir, rules); err != nil {
				return nil, err
			}
		}
		plan.dirs["rules"] = true
		for name, source := range rules {
			plan.files[filepath.Join("rules", name)] = source
		}
	}
	// Skill links do not depend on the destination; resolve them in the generation.
	plan.agents = string(resolveSkillReferences([]byte(plan.agents), namespace, prefix))
	// Preserve destination placeholders until activation.
	if err := bundleResources(plan, plan.agents); err != nil {
		return nil, err
	}
	return plan, nil
}
