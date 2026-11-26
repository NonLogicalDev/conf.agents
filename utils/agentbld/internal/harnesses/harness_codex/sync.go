package harness_codex

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/NonLogicalDev/conf.agents/utils/agentbld/internal/engine"
)

// syncPlan records all link changes after checking for unmanaged conflicts.
type syncPlan struct {
	agents                     engine.AgentsPlan
	skills, resources, rules   engine.LinkPlan
	manageAgents, manageSkills bool
}

func applySync(codex, gen string, plan *engine.RenderPlan, syncSkills bool, changes syncPlan) error {
	a, skills, resources, rules := changes.agents, changes.skills, changes.resources, changes.rules
	var e error
	if e = os.MkdirAll(codex, 0755); e != nil {
		return e
	}
	status := "skipped"
	if changes.manageAgents {
		if a.Remove {
			if e = os.Remove(a.Destination); e != nil {
				return e
			}
			fmt.Printf("removed instructions %s\n", a.Destination)
		}
		status = "unchanged"
		if !a.Unchanged {
			target, e := filepath.Rel(codex, a.Source)
			if e != nil {
				return e
			}
			if e = os.Symlink(target, a.Destination); e != nil {
				return e
			}
			fmt.Printf("linked instructions %s -> %s\n", a.Destination, target)
			status = "linked"
		}
	}
	if changes.manageSkills {
		if e = engine.ApplyLinks(skills); e != nil {
			return e
		}
	}
	if changes.manageAgents {
		for _, p := range []engine.LinkPlan{resources, rules} {
			if e = engine.ApplyLinks(p); e != nil {
				return e
			}
		}
	}
	mode := "marketplace"
	if !changes.manageSkills {
		mode = "none"
	} else if syncSkills {
		mode = "skills"
	}
	fmt.Printf("synced agentbld: layers=%d fragments=%d skills=%d resources=%d generation=%s agents=%s skills_linked=%d skills_unchanged=%d resources_linked=%d resources_unchanged=%d rules_linked=%d rules_unchanged=%d skill_destination=%s destination=%s\n", len(plan.Layers), plan.Fragments, len(plan.Skills), len(plan.Resources), filepath.Base(gen), status, len(skills.Links), skills.Unchanged, len(resources.Links), resources.Unchanged, len(rules.Links), rules.Unchanged, mode, codex)
	return nil
}

// copyTree preserves symlinks and permissions in reset backups.
func copyTree(src, dst string) error {
	i, e := os.Lstat(src)
	if e != nil {
		return e
	}
	if i.Mode()&os.ModeSymlink != 0 {
		target, e := os.Readlink(src)
		if e != nil {
			return e
		}
		return os.Symlink(target, dst)
	}
	if i.IsDir() {
		if e = os.Mkdir(dst, i.Mode().Perm()|0700); e != nil {
			return e
		}
		es, e := engine.Entries(src)
		if e != nil {
			return e
		}
		for _, d := range es {
			if e = copyTree(filepath.Join(src, d.Name()), filepath.Join(dst, d.Name())); e != nil {
				return e
			}
		}
		return os.Chmod(dst, i.Mode().Perm())
	}
	if !i.Mode().IsRegular() {
		return fmt.Errorf("unsupported backup file: %s", src)
	}
	b, e := os.ReadFile(src)
	if e != nil {
		return e
	}
	if e = os.WriteFile(dst, b, i.Mode().Perm()); e != nil {
		return e
	}
	return os.Chmod(dst, i.Mode().Perm())
}

// Reset backs up instructions and skills before removing managed links.
func Reset(codex, root string) error {
	if e := engine.ValidDir(codex, "Codex path"); e != nil {
		return e
	}
	for _, component := range []string{"skills", "resources"} {
		if e := engine.ValidDir(filepath.Join(codex, component), component+" path"); e != nil {
			return e
		}
	}
	epoch := time.Now().Unix()
	agents := filepath.Join(codex, "AGENTS.md")
	agentBackup := ""
	if engine.Exists(agents) {
		if !engine.IsLink(agents) && !engine.IsRegular(agents) {
			return fmt.Errorf("Codex instructions are not a regular file: %s", agents)
		}
		agentBackup = filepath.Join(codex, fmt.Sprintf("AGENTS.%d.bak.md", epoch))
		if engine.Exists(agentBackup) {
			return fmt.Errorf("backup already exists: %s", agentBackup)
		}
	}
	skills := filepath.Join(codex, "skills")
	skillsBackup := ""
	var skillLinks, resourceLinks []string
	if engine.IsDir(skills) {
		skillsBackup = filepath.Join(codex, fmt.Sprintf("skills.%d.bak", epoch))
		if engine.Exists(skillsBackup) {
			return fmt.Errorf("backup already exists: %s", skillsBackup)
		}
		es, e := engine.Entries(skills)
		if e != nil {
			return e
		}
		for _, d := range es {
			path := filepath.Join(skills, d.Name())
			if engine.ManagedHarnessLink(path, root, Name, "skills", false) {
				skillLinks = append(skillLinks, path)
			}
		}
	}
	resources := filepath.Join(codex, "resources")
	if engine.IsDir(resources) {
		es, e := engine.Entries(resources)
		if e != nil {
			return e
		}
		for _, d := range es {
			path := filepath.Join(resources, d.Name())

			if engine.ManagedHarnessLink(path, root, Name, "resources", false) {
				resourceLinks = append(resourceLinks, path)
			}
		}
	}
	if skillsBackup != "" {
		if e := copyTree(skills, skillsBackup); e != nil {
			os.RemoveAll(skillsBackup)
			return e
		}
		fmt.Printf("backed up %s to %s\n", skills, skillsBackup)
	}
	if agentBackup != "" {
		if e := os.Rename(agents, agentBackup); e != nil {
			return e
		}
		fmt.Printf("backed up %s to %s\n", agents, agentBackup)
	}
	for _, link := range append(skillLinks, resourceLinks...) {
		if e := os.Remove(link); e != nil {
			return e
		}
	}
	var ab, sb int
	if agentBackup != "" {
		ab = 1
	}
	if skillsBackup != "" {
		sb = 1
	}
	fmt.Printf("reset agentbld: instructions_backed_up=%d skills_backed_up=%d skills_removed=%d resources_removed=%d destination=%s\n", ab, sb, len(skillLinks), len(resourceLinks), codex)
	return nil
}

// prepareCodex checks destination conflicts before any harness is activated.
func prepareCodex(registry engine.Registry, gen, root string, snap engine.Snapshot, plan *engine.RenderPlan, dest engine.Destination) (func() error, error) {
	c, ok := snap["codex"].(Snapshot)
	if !ok {
		return nil, fmt.Errorf("generation config has no codex mode")
	}
	manageSkills := c.Skills || (c.Plugin != "none" && c.Agents)
	if c.Agents || manageSkills {
		if err := registry.PreflightHarnessTarget(root, gen, "codex", dest.Config, plan); err != nil {
			return nil, err
		}
	}
	payload := filepath.Join(gen, "codex")
	target := filepath.Join(root, ".targets", "codex", engine.TargetHash(dest.Config), filepath.Base(gen))
	var changes syncPlan
	changes.manageAgents, changes.manageSkills = c.Agents, manageSkills
	var err error
	if c.Agents {
		changes.agents, err = registry.PlanHarnessInstructions(dest.Config, target, root, "codex", "AGENTS.md")
		if err != nil {
			return nil, err
		}
		changes.resources, err = registry.PlanHarnessComponent(dest.Config, payload, target, root, "codex", "resources", true)
		if err != nil {
			return nil, err
		}
		changes.rules, err = registry.PlanHarnessComponent(dest.Config, payload, target, root, "codex", "rules", true)
		if err != nil {
			return nil, err
		}
	}
	if manageSkills {
		changes.skills, err = registry.PlanHarnessComponent(dest.Config, payload, target, root, "codex", "skills", c.Skills)
		if err != nil {
			return nil, err
		}
	}
	var pluginPlan *engine.RenderPlan
	if c.Plugin != "none" {
		pluginPlan, err = engine.ProjectPluginSkills(plan)
		if err != nil {
			return nil, err
		}
		pluginPlan, err = engine.TargetPlan(pluginPlan, dest.Config)
		if err != nil {
			return nil, err
		}
		if _, err = inspectMarketplace(pluginPlan, dest.Marketplace, c.PluginName); err != nil {
			return nil, err
		}
	}
	return func() error {
		var activatedPlan *engine.RenderPlan
		if c.Agents || manageSkills {
			var err error
			target, activatedPlan, err = registry.MaterializeHarnessTarget(root, gen, "codex", dest.Config, plan)
			if err != nil {
				return err
			}
		}
		if pluginPlan != nil {
			result, err := populateMarketplace(pluginPlan, dest.Marketplace, c.PluginName)
			if err != nil {
				return err
			}
			out, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(out))
		}
		if c.Agents || manageSkills {
			return applySync(dest.Config, target, activatedPlan, c.Skills, changes)
		}
		return nil
	}, nil
}
