package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type linkPair struct{ dst, src string }
type linkPlan struct {
	directory, component string
	removals             []string
	links                []linkPair
	unchanged            int
}
type agentsPlan struct {
	destination, source string
	remove, unchanged   bool
}

func managedLink(link, root, component string) bool {
	if !isLink(link) {
		return false
	}
	target, e := os.Readlink(link)
	if e != nil {
		return false
	}
	if !filepath.IsAbs(target) {
		target = filepath.Join(filepath.Dir(link), target)
	}
	target, _ = filepath.Abs(target)
	rel, e := filepath.Rel(root, target)
	if e != nil {
		return false
	}
	parts := strings.Split(rel, string(filepath.Separator))
	if len(parts) < 2 || !generationPattern.MatchString(parts[0]) {
		return false
	}
	if component == "AGENTS.md" {
		return len(parts) == 2 && parts[1] == component
	}
	return len(parts) >= 3 && parts[1] == component
}
func validDir(path, description string) error {
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
		return fmt.Errorf("%s is not a regular directory: %s", description, path)
	}
	return nil
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
func applyLinks(p linkPlan) error {
	if e := os.MkdirAll(p.directory, 0755); e != nil {
		return e
	}
	for _, path := range p.removals {
		if e := os.Remove(path); e != nil {
			return e
		}
		fmt.Printf("removed %s link %s\n", p.component, path)
	}
	for _, link := range p.links {
		target, e := filepath.Rel(p.directory, link.src)
		if e != nil {
			return e
		}
		if e = os.Symlink(target, link.dst); e != nil {
			return e
		}
		fmt.Printf("linked %s %s -> %s\n", p.component, link.dst, target)
	}
	return nil
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

// Preflight every destination before changing any links, so an unmanaged
// conflict in one component cannot leave the others partially updated.
func syncAll(codex, gen, root string, plan *renderPlan, syncSkills, instructionsOnly bool) error {
	a, e := planAgents(codex, gen, root)
	if e != nil {
		return e
	}
	skills, e := planComponent(codex, gen, root, "skills", syncSkills)
	if e != nil {
		return e
	}
	resources, e := planComponent(codex, gen, root, "resources", true)
	if e != nil {
		return e
	}
	rules, e := planComponent(codex, gen, root, "rules", true)
	if e != nil {
		return e
	}
	if e = os.MkdirAll(codex, 0755); e != nil {
		return e
	}
	if a.remove {
		if e = os.Remove(a.destination); e != nil {
			return e
		}
		fmt.Printf("removed instructions %s\n", a.destination)
	}
	status := "unchanged"
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
	for _, p := range []linkPlan{skills, resources, rules} {
		if e = applyLinks(p); e != nil {
			return e
		}
	}
	mode := "marketplace"
	if instructionsOnly {
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
			if !managedLink(link, root, component) {
				continue
			}
			target, e := os.Readlink(link)
			if e != nil {
				return e
			}
			if !filepath.IsAbs(target) {
				target = filepath.Join(filepath.Dir(link), target)
			}
			target, e = filepath.Abs(target)
			if e != nil {
				return e
			}
			rel, e := filepath.Rel(root, target)
			if e != nil {
				return e
			}
			name := strings.Split(rel, string(filepath.Separator))[0]
			if removable[name] {
				return fmt.Errorf("refusing to prune active Codex generation: %s (%s)", name, link)
			}
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
func reset(codex, legacy string) error {
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
			if within(target, legacy) {
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
	fmt.Printf("reset agent-config: instructions_backed_up=%d skills_backed_up=%d skills_removed=%d resources_removed=%d destination=%s\n", ab, sb, len(skillLinks), len(resourceLinks), codex)
	return nil
}
