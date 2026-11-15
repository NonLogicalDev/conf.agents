package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

// managedHarnessLink recognizes only links created for this destination and harness.
func managedHarnessLink(link, root, harness, component string) bool {
	if harness == "codex" && managedLink(link, root, component) {
		return true
	}
	if !isLink(link) {
		return false
	}
	target, err := os.Readlink(link)
	if err != nil {
		return false
	}
	if !filepath.IsAbs(target) {
		target = filepath.Join(filepath.Dir(link), target)
	}
	target, err = filepath.Abs(target)
	if err != nil {
		return false
	}
	destination := filepath.Dir(link)
	if component != "AGENTS.md" && component != "CLAUDE.md" {
		destination = filepath.Dir(destination)
	}
	base := filepath.Join(root, ".targets", harness, targetHash(destination))
	rel, err := filepath.Rel(base, target)
	if err != nil {
		return false
	}
	parts := strings.Split(rel, string(filepath.Separator))
	if len(parts) < 2 || !generationPattern.MatchString(parts[0]) || parts[1] != component {
		return false
	}
	if component == "AGENTS.md" || component == "CLAUDE.md" {
		return len(parts) == 2
	}
	return len(parts) >= 3
}

func planHarnessInstructions(destination, target, root, harness, name string) (agentsPlan, error) {
	p := agentsPlan{destination: filepath.Join(destination, name), source: filepath.Join(target, name)}
	if err := validDir(destination, harness+" destination"); err != nil {
		return p, err
	}
	relative, err := filepath.Rel(destination, p.source)
	if err != nil {
		return p, err
	}
	if isLink(p.destination) {
		current, err := os.Readlink(p.destination)
		if err != nil {
			return p, err
		}
		if current == relative {
			p.unchanged = true
			return p, nil
		}
		if managedHarnessLink(p.destination, root, harness, name) {
			p.remove = true
			return p, nil
		}
		return p, fmt.Errorf("refusing to replace unmanaged %s instructions: %s", harness, p.destination)
	}
	if !exists(p.destination) {
		return p, nil
	}
	if !isRegular(p.destination) {
		return p, fmt.Errorf("%s instructions are not a regular file: %s", harness, p.destination)
	}
	b, err := os.ReadFile(p.destination)
	if err != nil {
		return p, err
	}
	if len(b) != 0 && !strings.HasPrefix(string(b), marker) {
		return p, fmt.Errorf("refusing to replace unmanaged %s instructions: %s", harness, p.destination)
	}
	p.remove = true
	return p, nil
}

func planHarnessComponent(destination, payload, target, root, harness, component string, populate bool) (linkPlan, error) {
	p := linkPlan{directory: filepath.Join(destination, component), component: component}
	if err := validDir(p.directory, component+" target"); err != nil {
		return p, err
	}
	desired := map[string]string{}
	if populate {
		es, err := entries(filepath.Join(payload, component))
		if err != nil {
			return p, err
		}
		for _, entry := range es {
			desired[entry.Name()] = filepath.Join(target, component, entry.Name())
		}
	}
	remove := map[string]bool{}
	if isDir(p.directory) {
		es, err := entries(p.directory)
		if err != nil {
			return p, err
		}
		for _, entry := range es {
			path := filepath.Join(p.directory, entry.Name())
			if !managedHarnessLink(path, root, harness, component) {
				continue
			}
			src, ok := desired[entry.Name()]
			relative, _ := filepath.Rel(p.directory, src)
			current, _ := os.Readlink(path)
			if !ok || current != relative {
				p.removals = append(p.removals, path)
				remove[path] = true
			}
		}
	}
	for _, name := range sortedKeys(desired) {
		src, dst := desired[name], filepath.Join(p.directory, name)
		relative, _ := filepath.Rel(p.directory, src)
		current, err := os.Readlink(dst)
		if err == nil && current == relative {
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

func applyHarnessInstructions(p agentsPlan, destination string) error {
	if err := os.MkdirAll(destination, 0755); err != nil {
		return err
	}
	if p.remove {
		if err := os.Remove(p.destination); err != nil {
			return err
		}
	}
	if !p.unchanged {
		relative, err := filepath.Rel(destination, p.source)
		if err != nil {
			return err
		}
		if err = os.Symlink(relative, p.destination); err != nil {
			return err
		}
		fmt.Printf("linked instructions %s -> %s\n", p.destination, relative)
	}
	return nil
}
