package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// managedGeneration identifies only our legacy links or target snapshots for
// the Codex destination containing the link.
func managedGeneration(link, root, component string) (string, bool) {
	if !isLink(link) {
		return "", false
	}
	target, err := os.Readlink(link)
	if err != nil {
		return "", false
	}
	if !filepath.IsAbs(target) {
		target = filepath.Join(filepath.Dir(link), target)
	}
	target, err = filepath.Abs(target)
	if err != nil {
		return "", false
	}
	rel, err := filepath.Rel(root, target)
	if err != nil {
		return "", false
	}
	parts := strings.Split(rel, string(filepath.Separator))
	if len(parts) >= 2 && generationPattern.MatchString(parts[0]) {
		if component == "AGENTS.md" {
			return parts[0], len(parts) == 2 && parts[1] == component
		}
		return parts[0], len(parts) >= 3 && parts[1] == component
	}
	if len(parts) < 4 || parts[0] != ".targets" || !generationPattern.MatchString(parts[2]) {
		return "", false
	}
	codex := filepath.Dir(link)
	if component != "AGENTS.md" {
		codex = filepath.Dir(codex)
	}
	if parts[1] != targetHash(codex) {
		return "", false
	}
	if component == "AGENTS.md" {
		return parts[2], len(parts) == 4 && parts[3] == component
	}
	return parts[2], len(parts) >= 5 && parts[3] == component
}
func managedLink(link, root, component string) bool {
	_, ok := managedGeneration(link, root, component)
	return ok
}

func planAgents(codex, gen, root string) (agentsPlan, error) {
	p := agentsPlan{destination: filepath.Join(codex, "AGENTS.md"), source: filepath.Join(gen, "AGENTS.md")}
	if e := validDir(codex, "Codex destination"); e != nil {
		return p, e
	}
	target, _ := filepath.Rel(codex, p.source)
	if isLink(p.destination) {
		v, e := os.Readlink(p.destination)
		if e != nil {
			return p, e
		}
		if v == target {
			p.unchanged = true
			return p, nil
		}
		if managedLink(p.destination, root, "AGENTS.md") {
			p.remove = true
			return p, nil
		}
		return p, fmt.Errorf("refusing to replace unmanaged Codex instructions: %s", p.destination)
	}
	if !exists(p.destination) {
		return p, nil
	}
	if !isRegular(p.destination) {
		return p, fmt.Errorf("Codex instructions are not a regular file: %s", p.destination)
	}
	b, e := os.ReadFile(p.destination)
	if e != nil {
		return p, e
	}
	if len(b) > 0 && !strings.HasPrefix(string(b), marker) {
		return p, fmt.Errorf("refusing to replace unmanaged Codex instructions: %s", p.destination)
	}
	p.remove = true
	return p, nil
}
func planComponent(codex, gen, root, component string, populate bool) (linkPlan, error) {
	p := linkPlan{directory: filepath.Join(codex, component), component: component}
	if e := validDir(p.directory, component+" target"); e != nil {
		return p, e
	}
	desired := map[string]string{}
	if populate {
		es, e := entries(filepath.Join(gen, component))
		if e != nil {
			return p, e
		}
		for _, entry := range es {
			desired[entry.Name()] = filepath.Join(gen, component, entry.Name())
		}
	}
	remove := map[string]bool{}
	if isDir(p.directory) {
		es, e := entries(p.directory)
		if e != nil {
			return p, e
		}
		for _, entry := range es {
			path := filepath.Join(p.directory, entry.Name())
			if !managedLink(path, root, component) {
				continue
			}
			src, ok := desired[entry.Name()]
			target, _ := filepath.Rel(p.directory, src)
			current, _ := os.Readlink(path)
			if !ok || current != target {
				p.removals = append(p.removals, path)
				remove[path] = true
			}
		}
	}
	for _, name := range sortedKeys(desired) {
		src := desired[name]
		dst := filepath.Join(p.directory, name)
		target, _ := filepath.Rel(p.directory, src)
		current, e := os.Readlink(dst)
		if e == nil && current == target {
			p.unchanged++
			continue
		}
		if !remove[dst] && exists(dst) {
			return p, fmt.Errorf("refusing to replace unmanaged %s entry: %s", component, dst)
		}
		p.links = append(p.links, linkPair{dst, src})
	}
	return p, nil
}
func syncRules(codex, gen, root string) error {
	if e := validDir(codex, "Codex destination"); e != nil {
		return e
	}
	p, e := planComponent(codex, gen, root, "rules", true)
	if e != nil {
		return e
	}
	if e = applyLinks(p); e != nil {
		return e
	}
	fmt.Printf("synced rules: linked=%d unchanged=%d\n", len(p.links), p.unchanged)
	return nil
}

// syncPlan records all link changes after checking for unmanaged conflicts.
type syncPlan struct {
	agents                     agentsPlan
	skills, resources, rules   linkPlan
	manageAgents, manageSkills bool
}

// planSync checks every destination before marketplace or link changes begin.
func planSync(codex, gen, root string, manageAgents, manageSkills, populateSkills bool) (syncPlan, error) {
	p := syncPlan{manageAgents: manageAgents, manageSkills: manageSkills}
	var err error
	if err = validDir(codex, "Codex destination"); err != nil {
		return p, err
	}
	if manageAgents {
		if p.agents, err = planAgents(codex, gen, root); err != nil {
			return p, err
		}
		if p.resources, err = planComponent(codex, gen, root, "resources", true); err != nil {
			return p, err
		}
		if p.rules, err = planComponent(codex, gen, root, "rules", true); err != nil {
			return p, err
		}
	}
	if manageSkills {
		if p.skills, err = planComponent(codex, gen, root, "skills", populateSkills); err != nil {
			return p, err
		}
	}
	return p, nil
}

func applySync(codex, gen string, plan *renderPlan, syncSkills, instructionsOnly bool, changes syncPlan) error {
	a, skills, resources, rules := changes.agents, changes.skills, changes.resources, changes.rules
	var e error
	if e = os.MkdirAll(codex, 0755); e != nil {
		return e
	}
	status := "skipped"
	if changes.manageAgents {
		if a.remove {
			if e = os.Remove(a.destination); e != nil {
				return e
			}
			fmt.Printf("removed instructions %s\n", a.destination)
		}
		status = "unchanged"
		if !a.unchanged {
			target, e := filepath.Rel(codex, a.source)
			if e != nil {
				return e
			}
			if e = os.Symlink(target, a.destination); e != nil {
				return e
			}
			fmt.Printf("linked instructions %s -> %s\n", a.destination, target)
			status = "linked"
		}
	}
	if changes.manageSkills {
		if e = applyLinks(skills); e != nil {
			return e
		}
	}
	if changes.manageAgents {
		for _, p := range []linkPlan{resources, rules} {
			if e = applyLinks(p); e != nil {
				return e
			}
		}
	}
	mode := "marketplace"
	if instructionsOnly || !changes.manageSkills {
		mode = "none"
	} else if syncSkills {
		mode = "skills"
	}
	fmt.Printf("synced agent-config: layers=%d fragments=%d skills=%d resources=%d generation=%s agents=%s skills_linked=%d skills_unchanged=%d resources_linked=%d resources_unchanged=%d rules_linked=%d rules_unchanged=%d skill_destination=%s destination=%s\n", len(plan.layers), plan.fragments, len(plan.skills), len(plan.resources), filepath.Base(gen), status, len(skills.links), skills.unchanged, len(resources.links), resources.unchanged, len(rules.links), rules.unchanged, mode, codex)
	return nil
}
func prune(root, codex string) error {
	gens, e := discoverGenerations(root)
	if e != nil {
		return e
	}
	split := len(gens) - 5
	if split < 0 {
		split = 0
	}
	remove, kept := gens[:split], gens[split:]
	if len(remove) > 0 {
		if !aliasMatches(root, gens[len(gens)-1].path) {
			return fmt.Errorf("latest generation alias must target the newest generation before pruning")
		}
		if e = validDir(codex, "Codex destination"); e != nil {
			return e
		}
		// Refuse to remove any generation still referenced by a managed link.
		removable := map[string]bool{}
		for _, g := range remove {
			removable[filepath.Base(g.path)] = true
		}
		links := map[string]string{filepath.Join(codex, "AGENTS.md"): "AGENTS.md"}
		for _, component := range []string{"skills", "resources", "rules"} {
			dir := filepath.Join(codex, component)
			if e = validDir(dir, component+" target"); e != nil {
				return e
			}
			if isDir(dir) {
				es, e := entries(dir)
				if e != nil {
					return e
				}
				for _, d := range es {
					links[filepath.Join(dir, d.Name())] = component
				}
			}
		}
		for link, component := range links {
			name, ok := managedGeneration(link, root, component)
			if ok && removable[name] {
				return fmt.Errorf("refusing to prune active Codex generation: %s (%s)", name, link)
			}
		}
		if err := pruneTargets(root, remove); err != nil {
			return err
		}

		for _, g := range remove {
			if e = removeGeneration(g.path); e != nil {
				return e
			}
		}
	}
	latest := "none"
	if len(kept) > 0 {
		latest = filepath.Base(kept[len(kept)-1].path)
	}
	fmt.Printf("pruned agent-config: kept=%d removed=%d latest=%s destination=%s\n", len(kept), len(remove), latest, root)
	return nil
}

// Restore permissions on surviving directories if pruning fails partway.
func removeGeneration(root string) (err error) {
	modes := map[string]os.FileMode{}
	defer func() {
		if err != nil {
			for path, mode := range modes {
				_ = os.Chmod(path, mode)
			}
		}
	}()
	info, err := os.Stat(root)
	if err != nil {
		return err
	}
	modes[root] = info.Mode().Perm()
	if err = os.Chmod(root, modes[root]|0200); err != nil {
		return err
	}
	if err = walk(root, func(path string, d os.DirEntry) error {
		if !d.IsDir() {
			return nil
		}
		i, e := d.Info()
		if e != nil {
			return e
		}
		modes[path] = i.Mode().Perm()
		return os.Chmod(path, modes[path]|0200)
	}); err != nil {
		return err
	}
	return os.RemoveAll(root)
}

// copyTree preserves symlinks and permissions in the reset backup. Directories
// remain writable during the copy so read-only source trees can be backed up.
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
		es, e := entries(src)
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

// reset backs up instructions and skills before removing links. Resource links
// are removed only when they point into the legacy render directory.
func reset(codex string, legacyRoots []string) error {
	if e := validDir(codex, "Codex path"); e != nil {
		return e
	}
	for _, component := range []string{"skills", "resources"} {
		if e := validDir(filepath.Join(codex, component), component+" path"); e != nil {
			return e
		}
	}
	epoch := time.Now().Unix()
	agents := filepath.Join(codex, "AGENTS.md")
	agentBackup := ""
	if exists(agents) {
		if !isLink(agents) && !isRegular(agents) {
			return fmt.Errorf("Codex instructions are not a regular file: %s", agents)
		}
		agentBackup = filepath.Join(codex, fmt.Sprintf("AGENTS.%d.bak.md", epoch))
		if exists(agentBackup) {
			return fmt.Errorf("backup already exists: %s", agentBackup)
		}
	}
	skills := filepath.Join(codex, "skills")
	skillsBackup := ""
	var skillLinks, resourceLinks []string
	if isDir(skills) {
		skillsBackup = filepath.Join(codex, fmt.Sprintf("skills.%d.bak", epoch))
		if exists(skillsBackup) {
			return fmt.Errorf("backup already exists: %s", skillsBackup)
		}
		es, e := entries(skills)
		if e != nil {
			return e
		}
		for _, d := range es {
			path := filepath.Join(skills, d.Name())
			if isLink(path) {
				skillLinks = append(skillLinks, path)
			}
		}
	}
	resources := filepath.Join(codex, "resources")
	if isDir(resources) {
		es, e := entries(resources)
		if e != nil {
			return e
		}
		for _, d := range es {
			path := filepath.Join(resources, d.Name())
			if !isLink(path) {
				continue
			}
			target, e := os.Readlink(path)
			if e != nil {
				return e
			}
			if !filepath.IsAbs(target) {
				target = filepath.Join(resources, target)
			}
			target, _ = filepath.Abs(target)
			for _, legacy := range legacyRoots {
				if within(target, legacy) {
					resourceLinks = append(resourceLinks, path)
					break
				}
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
	fmt.Printf("reset agent-config: instructions_backed_up=%d skills_backed_up=%d skills_removed=%d resources_removed=%d destination=%s\n", ab, sb, len(skillLinks), len(resourceLinks), codex)
	return nil
}

// prepareCodex checks destination conflicts before any harness is activated.
func prepareCodex(gen, root string, snap snapshot, plan *renderPlan, dest runtimeDestinations) (func() error, error) {
	c, ok := snap["codex"].(snapshotCodex)
	if !ok {
		return nil, fmt.Errorf("generation config has no codex mode")
	}
	manageSkills := c.Skills || c.Legacy.InstructionsOnly || (c.Plugin != "none" && c.Agents)
	syncSkills := c.Skills && !c.Legacy.InstructionsOnly
	if c.Agents || manageSkills || c.Legacy.RulesOnly {
		if err := preflightHarnessTarget(root, gen, "codex", dest.Codex, plan); err != nil {
			return nil, err
		}
	}
	payload := harnessPayloadDir(gen, "codex")
	target := filepath.Join(root, ".targets", "codex", targetHash(dest.Codex), filepath.Base(gen))
	var changes syncPlan
	changes.manageAgents, changes.manageSkills = c.Agents, manageSkills
	var err error
	if c.Legacy.RulesOnly {
		changes.rules, err = planHarnessComponent(dest.Codex, payload, target, root, "codex", "rules", true)
		if err != nil {
			return nil, err
		}
	} else {
		if c.Agents {
			changes.agents, err = planHarnessInstructions(dest.Codex, target, root, "codex", "AGENTS.md")
			if err != nil {
				return nil, err
			}
			changes.resources, err = planHarnessComponent(dest.Codex, payload, target, root, "codex", "resources", true)
			if err != nil {
				return nil, err
			}
			changes.rules, err = planHarnessComponent(dest.Codex, payload, target, root, "codex", "rules", true)
			if err != nil {
				return nil, err
			}
		}
		if manageSkills {
			changes.skills, err = planHarnessComponent(dest.Codex, payload, target, root, "codex", "skills", syncSkills)
			if err != nil {
				return nil, err
			}
		}
	}
	activatedPlan, err := targetPlan(plan, dest.Codex)
	if err != nil {
		return nil, err
	}
	if c.Plugin != "none" {
		pluginPlan, err := projectPluginSkills(plan, c.PluginName)
		if err != nil {
			return nil, err
		}
		pluginPlan, err = targetPlan(pluginPlan, dest.Codex)
		if err != nil {
			return nil, err
		}
		if c.Plugin == "skills" && c.Skills {
			pluginPlan = pluginWithoutGuidance(activatedPlan)
		}
		if _, err = inspectMarketplace(pluginPlan, dest.Marketplace, c.PluginName); err != nil {
			return nil, err
		}
	}
	return func() error {
		if c.Agents || manageSkills || c.Legacy.RulesOnly {
			var err error
			target, activatedPlan, err = materializeHarnessTarget(root, gen, "codex", dest.Codex, plan)
			if err != nil {
				return err
			}
		}
		if c.Legacy.RulesOnly {
			return applyLinks(changes.rules)
		}
		if c.Plugin != "none" {
			pluginPlan, err := projectPluginSkills(plan, c.PluginName)
			if err != nil {
				return err
			}
			pluginPlan, err = targetPlan(pluginPlan, dest.Codex)
			if err != nil {
				return err
			}
			if c.Plugin == "skills" && c.Skills {
				pluginPlan = pluginWithoutGuidance(pluginPlan)
			}
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
			return applySync(dest.Codex, target, activatedPlan, syncSkills, c.Legacy.InstructionsOnly, changes)
		}
		return nil
	}, nil
}
