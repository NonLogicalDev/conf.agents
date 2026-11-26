package harness_claude

import (
	"fmt"
	"path/filepath"

	"github.com/NonLogicalDev/conf.agents/utils/agentbld/internal/engine"
)

// prepareClaude checks every selected destination before producing a target or
// changing links. User settings, rules, and unselected components stay intact.
func prepareClaude(registry engine.Registry, gen, root string, s engine.Snapshot, plan *engine.RenderPlan, d engine.Destination) (func() error, error) {
	saved, ok := s["claude"].(Snapshot)
	if !ok {
		return nil, fmt.Errorf("generation config has no Claude mode")
	}
	destination := d.Config
	hasDestination := saved.Agents || saved.Skills
	if hasDestination {
		if err := engine.ValidDir(destination, "Claude destination"); err != nil {
			return nil, err
		}
		if err := registry.PreflightHarnessTarget(root, gen, "claude", destination, plan); err != nil {
			return nil, err
		}
	}
	var marketplace *claudeMarketplacePlan
	if saved.Plugin == "full" || saved.Plugin == "skills" {
		pluginPlan, err := engine.ProjectPluginSkills(plan)
		if err != nil {
			return nil, err
		}
		pluginPlan, err = engine.TargetPlan(pluginPlan, destination)
		if err != nil {
			return nil, err
		}
		marketplace, err = inspectClaudeMarketplace(pluginPlan, d.Marketplace, saved.PluginName)
		if err != nil {
			return nil, err
		}
	}
	payload := filepath.Join(gen, "claude")
	target := filepath.Join(root, ".targets", "claude", engine.TargetHash(destination), filepath.Base(gen))
	var instructions engine.AgentsPlan
	var resources, skills engine.LinkPlan
	var err error
	if saved.Agents {
		instructions, err = registry.PlanHarnessInstructions(destination, target, root, "claude", "CLAUDE.md")
		if err != nil {
			return nil, err
		}
		resources, err = registry.PlanHarnessComponent(destination, payload, target, root, "claude", "resources", true)
		if err != nil {
			return nil, err
		}
	}
	manageSkills := saved.Skills || (saved.Agents && (saved.Plugin == "full" || saved.Plugin == "skills"))
	if manageSkills {
		skills, err = registry.PlanHarnessComponent(destination, payload, target, root, "claude", "skills", saved.Skills)
		if err != nil {
			return nil, err
		}
	}
	return func() error {
		if hasDestination {
			if _, _, err := registry.MaterializeHarnessTarget(root, gen, "claude", destination, plan); err != nil {
				return err
			}
		}
		if marketplace != nil {
			if err := populateClaudeMarketplace(marketplace); err != nil {
				return err
			}
		}
		if saved.Agents {
			if err := engine.ApplyHarnessInstructions(instructions, destination); err != nil {
				return err
			}
			if err := engine.ApplyLinks(resources); err != nil {
				return err
			}
		}
		if manageSkills {
			if err := engine.ApplyLinks(skills); err != nil {
				return err
			}
		}
		return nil
	}, nil
}
