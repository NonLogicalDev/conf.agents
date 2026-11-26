package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type linkPair struct{ dst, src string }
type LinkPlan struct {
	Directory, Component string
	Removals             []string
	Links                []linkPair
	Unchanged            int
}
type AgentsPlan struct {
	Destination, Source string
	Remove, Unchanged   bool
}

func (p *LinkPlan) AddLink(dst, src string) {
	p.Links = append(p.Links, linkPair{dst, src})
}

func ValidDir(path, description string) error {
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
func ApplyLinks(p LinkPlan) error {
	if e := os.MkdirAll(p.Directory, 0755); e != nil {
		return e
	}
	for _, path := range p.Removals {
		if e := os.Remove(path); e != nil {
			return e
		}
		fmt.Printf("removed %s link %s\n", p.Component, path)
	}
	for _, link := range p.Links {
		target, e := filepath.Rel(p.Directory, link.src)
		if e != nil {
			return e
		}
		if e = os.Symlink(target, link.dst); e != nil {
			return e
		}
		fmt.Printf("linked %s %s -> %s\n", p.Component, link.dst, target)
	}
	return nil
}

// ManagedHarnessLink recognizes only links created for this destination and harness.
func ManagedHarnessLink(link, root, harness, component string, instruction bool) bool {
	if !IsLink(link) {
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
	if !instruction {
		destination = filepath.Dir(destination)
	}
	base := filepath.Join(root, ".targets", harness, TargetHash(destination))
	rel, err := filepath.Rel(base, target)
	if err != nil {
		return false
	}
	parts := strings.Split(rel, string(filepath.Separator))
	if len(parts) < 2 || !generationPattern.MatchString(parts[0]) || parts[1] != component {
		return false
	}
	if instruction {
		return len(parts) == 2
	}
	return len(parts) >= 3
}

func (registry Registry) PlanHarnessInstructions(destination, target, root, harness, name string) (AgentsPlan, error) {
	p := AgentsPlan{Destination: filepath.Join(destination, name), Source: filepath.Join(target, name)}
	if err := ValidDir(destination, harness+" destination"); err != nil {
		return p, err
	}
	relative, err := filepath.Rel(destination, p.Source)
	if err != nil {
		return p, err
	}
	if IsLink(p.Destination) {
		current, err := os.Readlink(p.Destination)
		if err != nil {
			return p, err
		}
		if current == relative {
			p.Unchanged = true
			return p, nil
		}
		if ManagedHarnessLink(p.Destination, root, harness, name, true) {
			p.Remove = true
			return p, nil
		}
		return p, fmt.Errorf("refusing to replace unmanaged %s instructions: %s", harness, p.Destination)
	}
	if !Exists(p.Destination) {
		return p, nil
	}
	if !IsRegular(p.Destination) {
		return p, fmt.Errorf("%s instructions are not a regular file: %s", harness, p.Destination)
	}
	b, err := os.ReadFile(p.Destination)
	if err != nil {
		return p, err
	}
	if len(b) != 0 && !strings.HasPrefix(string(b), GeneratedMarker) {
		return p, fmt.Errorf("refusing to replace unmanaged %s instructions: %s", harness, p.Destination)
	}
	p.Remove = true
	return p, nil
}

func (registry Registry) PlanHarnessComponent(destination, payload, target, root, harness, component string, populate bool) (LinkPlan, error) {
	p := LinkPlan{Directory: filepath.Join(destination, component), Component: component}
	if err := ValidDir(p.Directory, component+" target"); err != nil {
		return p, err
	}
	desired := map[string]string{}
	if populate {
		es, err := Entries(filepath.Join(payload, component))
		if err != nil {
			return p, err
		}
		for _, entry := range es {
			desired[entry.Name()] = filepath.Join(target, component, entry.Name())
		}
	}
	remove := map[string]bool{}
	if IsDir(p.Directory) {
		es, err := Entries(p.Directory)
		if err != nil {
			return p, err
		}
		for _, entry := range es {
			path := filepath.Join(p.Directory, entry.Name())
			if !ManagedHarnessLink(path, root, harness, component, false) {
				continue
			}
			src, ok := desired[entry.Name()]
			relative, _ := filepath.Rel(p.Directory, src)
			current, _ := os.Readlink(path)
			if !ok || current != relative {
				p.Removals = append(p.Removals, path)
				remove[path] = true
			}
		}
	}
	for _, name := range SortedKeys(desired) {
		src, dst := desired[name], filepath.Join(p.Directory, name)
		relative, _ := filepath.Rel(p.Directory, src)
		current, err := os.Readlink(dst)
		if err == nil && current == relative {
			p.Unchanged++
			continue
		}
		if !remove[dst] && Exists(dst) {
			return p, fmt.Errorf("refusing to replace unmanaged %s entry: %s", component, dst)
		}
		p.AddLink(dst, src)
	}
	return p, nil
}

func ApplyHarnessInstructions(p AgentsPlan, destination string) error {
	if err := os.MkdirAll(destination, 0755); err != nil {
		return err
	}
	if p.Remove {
		if err := os.Remove(p.Destination); err != nil {
			return err
		}
	}
	if !p.Unchanged {
		relative, err := filepath.Rel(destination, p.Source)
		if err != nil {
			return err
		}
		if err = os.Symlink(relative, p.Destination); err != nil {
			return err
		}
		fmt.Printf("linked instructions %s -> %s\n", p.Destination, relative)
	}
	return nil
}
