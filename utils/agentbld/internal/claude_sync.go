package internal

import (
	"fmt"
	"path/filepath"
)

// prepareClaude checks every selected destination before producing a target or
// changing links. User settings, rules, and unselected components stay intact.
func prepareClaude(gen, root string, s snapshot, plan *renderPlan, d runtimeDestinations) (func() error, error) {
	saved, ok := s["claude"].(snapshotClaude)
	if !ok {
		return nil, fmt.Errorf("generation config has no Claude mode")
	}
	destination := d.Claude
	hasDestination := saved.Agents || saved.Skills
	if hasDestination {
		if err := validDir(destination, "Claude destination"); err != nil {
			return nil, err
		}
		if err := preflightHarnessTarget(root, gen, "claude", destination, plan); err != nil {
			return nil, err
		}
	}
	var marketplace *claudeMarketplacePlan
	if saved.Plugin == "full" || saved.Plugin == "skills" {
		pluginPlan, err := projectPluginSkills(plan, saved.PluginName)
		if err != nil {
			return nil, err
		}
		pluginPlan, err = targetPlan(pluginPlan, destination)
		if err != nil {
			return nil, err
		}
		marketplace, err = inspectClaudeMarketplace(pluginPlan, d.ClaudeMarketplace, saved.PluginName)
		if err != nil {
			return nil, err
		}
	}
	payload := filepath.Join(gen, "claude")
	target := filepath.Join(root, ".targets", "claude", targetHash(destination), filepath.Base(gen))
	var instructions agentsPlan
	var resources, skills linkPlan
	var err error
	if saved.Agents {
		instructions, err = planHarnessInstructions(destination, target, root, "claude", "CLAUDE.md")
		if err != nil {
			return nil, err
		}
		resources, err = planHarnessComponent(destination, payload, target, root, "claude", "resources", true)
		if err != nil {
			return nil, err
		}
	}
	manageSkills := saved.Skills || (saved.Agents && (saved.Plugin == "full" || saved.Plugin == "skills"))
	if manageSkills {
		skills, err = planHarnessComponent(destination, payload, target, root, "claude", "skills", saved.Skills)
		if err != nil {
			return nil, err
		}
	}
	return func() error {
		if hasDestination {
			if _, _, err := materializeHarnessTarget(root, gen, "claude", destination, plan); err != nil {
				return err
			}
		}
		if marketplace != nil {
			if err := populateClaudeMarketplace(marketplace); err != nil {
				return err
			}
		}
		if saved.Agents {
			if err := applyHarnessInstructions(instructions, destination); err != nil {
				return err
			}
			if err := applyLinks(resources); err != nil {
				return err
			}
		}
		if manageSkills {
			if err := applyLinks(skills); err != nil {
				return err
			}
		}
		return nil
	}, nil
}
