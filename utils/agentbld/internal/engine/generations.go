package engine

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var generationPattern = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2})__(\d{2,})$`)

type generation struct {
	Ordinal int
	Path    string
}

func (registry Registry) DiscoverGenerations(root string) ([]generation, error) {
	if IsLink(root) {
		return nil, fmt.Errorf("render root must not be a symlink: %s", root)
	}
	if _, err := os.Lstat(root); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	if !IsDir(root) {
		return nil, fmt.Errorf("render root is not a directory: %s", root)
	}
	es, e := Entries(root)
	if e != nil {
		return nil, e
	}
	var result []generation
	seen := map[int]bool{}
	for _, d := range es {
		path := filepath.Join(root, d.Name())
		if strings.HasPrefix(d.Name(), ".") {
			continue
		}
		if d.Name() == "latest" {
			if !IsLink(path) {
				return nil, fmt.Errorf("latest generation alias is not a symlink: %s", path)
			}
			target, e := os.Readlink(path)
			if e != nil {
				return nil, e
			}
			if filepath.IsAbs(target) || filepath.Base(target) != target || !generationPattern.MatchString(target) {
				return nil, fmt.Errorf("latest generation alias must target a local generation: %s", path)
			}
			continue
		}
		m := generationPattern.FindStringSubmatch(d.Name())
		if m == nil {
			return nil, fmt.Errorf("unexpected entry in rendered generation root: %s", path)
		}
		if _, e := time.Parse("2006-01-02", m[1]); e != nil {
			return nil, fmt.Errorf("generation has an invalid date: %s", path)
		}
		if IsLink(path) || !d.IsDir() {
			return nil, fmt.Errorf("generation is not a regular directory: %s", path)
		}
		check := func(p string, d os.DirEntry) error {
			if d.Type()&os.ModeSymlink != 0 {
				return fmt.Errorf("generation must not contain symlinks: %s", p)
			}
			info, e := d.Info()
			if e != nil {
				return e
			}
			if info.Mode().Perm()&0222 != 0 {
				return fmt.Errorf("generation contains a writable path: %s", p)
			}
			return nil
		}
		info, e := os.Stat(path)
		if e != nil {
			return nil, e
		}
		if info.Mode().Perm()&0222 != 0 {
			return nil, fmt.Errorf("generation contains a writable path: %s", path)
		}
		if e = walk(path, check); e != nil {
			return nil, e
		}
		ordinal, e := strconv.Atoi(m[2])
		if e != nil {
			return nil, e
		}
		if seen[ordinal] {
			return nil, fmt.Errorf("duplicate rendered generation ordinal: %02d", ordinal)
		}
		seen[ordinal] = true
		snap, err := registry.ReadSnapshot(path)
		if err != nil {
			return nil, err
		}
		if err := validateGenerationLayout(path, snap); err != nil {
			return nil, err
		}
		result = append(result, generation{ordinal, path})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Ordinal < result[j].Ordinal })
	return result, nil
}
func CopyPlanFile(p *RenderPlan, relative, dst string) error {
	src := p.Files[relative]
	b, ok := p.Contents[relative]
	if !ok {
		var e error
		b, e = os.ReadFile(src)
		if e != nil {
			return e
		}
	}
	info, e := os.Stat(src)
	if e != nil {
		return e
	}
	if e = os.WriteFile(dst, b, info.Mode().Perm()); e != nil {
		return e
	}
	mode := info.Mode().Perm()
	if saved, ok := p.fileModes[relative]; ok {
		mode = saved
	}
	return os.Chmod(dst, mode)
}
func RenderedMatches(root string, p *RenderPlan, config []byte) bool {
	return RenderedPlansMatch(root, map[string]*RenderPlan{"": p}, config)
}
func RenderedPlansMatch(root string, plans map[string]*RenderPlan, config []byte) bool {
	info, e := os.Lstat(root)
	if e != nil || !info.IsDir() || info.Mode().Perm()&0222 != 0 {
		return false
	}
	metadata, e := os.ReadFile(filepath.Join(root, SnapshotFile))
	if e != nil || !bytes.Equal(metadata, config) {
		return false
	}
	expected := map[string]bool{SnapshotFile: true}
	for name, p := range plans {
		if name != "" {
			expected[name] = true
		}
		instruction := p.InstructionsName
		b, err := os.ReadFile(filepath.Join(root, name, instruction))
		if err != nil || string(b) != p.Agents {
			return false
		}
		expected[filepath.Join(name, instruction)] = true
		for d := range p.Dirs {
			expected[filepath.Join(name, d)] = true
		}
		for f := range p.Files {
			expected[filepath.Join(name, f)] = true
		}
	}
	count := 0
	e = walk(root, func(path string, d os.DirEntry) error {
		rel, _ := filepath.Rel(root, path)
		i, e := d.Info()
		if e != nil {
			return e
		}
		if d.Type()&os.ModeSymlink != 0 || i.Mode().Perm()&0222 != 0 || !expected[rel] {
			return fmt.Errorf("unexpected path")
		}
		count++
		return nil
	})
	if e != nil || count != len(expected) {
		return false
	}
	for name, p := range plans {
		for d := range p.Dirs {
			if !IsDir(filepath.Join(root, name, d)) {
				return false
			}
		}
		for rel, src := range p.Files {
			dst := filepath.Join(root, name, rel)
			s, e := os.Stat(src)
			if e != nil {
				return false
			}
			d, e := os.Stat(dst)
			if e != nil || !d.Mode().IsRegular() || d.Mode().Perm() != s.Mode().Perm()&^0222 {
				return false
			}
			expected, ok := p.Contents[rel]
			if !ok {
				expected, e = os.ReadFile(src)
				if e != nil {
					return false
				}
			}
			got, e := os.ReadFile(dst)
			if e != nil || !bytes.Equal(got, expected) {
				return false
			}
		}
	}
	return true
}
func SummaryPlans(plans map[string]*RenderPlan, gen, status string) {
	name, dest := "none", "none"
	if gen != "" {
		name, dest = filepath.Base(gen), gen
	}
	layers, fragments, skills, resources := 0, 0, 0, 0
	for _, p := range plans {
		layers += len(p.Layers)
		fragments += p.Fragments
		skills += len(p.Skills)
		resources += len(p.Resources)
	}
	fmt.Printf("rendered agentbld: layers=%d fragments=%d skills=%d resources=%d generation=%s status=%s destination=%s\n", layers, fragments, skills, resources, name, status, dest)
}
func AliasMatches(root, gen string) bool {
	target, e := os.Readlink(filepath.Join(root, "latest"))
	return e == nil && target == filepath.Base(gen)
}
func updateAlias(root, gen string) error {
	if AliasMatches(root, gen) {
		return nil
	}
	tmp, e := os.MkdirTemp(root, ".latest.")
	if e != nil {
		return e
	}
	defer os.RemoveAll(tmp)
	candidate := filepath.Join(tmp, "latest")
	if e = os.Symlink(filepath.Base(gen), candidate); e != nil {
		return e
	}
	return os.Rename(candidate, filepath.Join(root, "latest"))
}

// PerformRenderPlans reuses an identical immutable generation or stages a new one
// before publishing it and updating the latest alias.
func (registry Registry) PerformRender(root string, p *RenderPlan, config []byte) (string, error) {
	return registry.PerformRenderPlans(root, map[string]*RenderPlan{"": p}, config)
}
func (registry Registry) PerformRenderPlans(root string, plans map[string]*RenderPlan, config []byte) (string, error) {
	gens, e := registry.DiscoverGenerations(root)
	if e != nil {
		return "", e
	}
	ordinal := 1
	if len(gens) > 0 {
		last := gens[len(gens)-1]
		if RenderedPlansMatch(last.Path, plans, config) {
			if e = updateAlias(root, last.Path); e != nil {
				return "", e
			}
			SummaryPlans(plans, last.Path, "up-to-date")
			return last.Path, nil
		}
		ordinal = last.Ordinal + 1
	}
	if e = os.MkdirAll(root, 0755); e != nil {
		return "", e
	}
	name := fmt.Sprintf("%s__%02d", time.Now().Format("2006-01-02"), ordinal)
	dest := filepath.Join(root, name)
	if Exists(dest) {
		return "", fmt.Errorf("rendered generation already exists: %s", dest)
	}
	tmp, e := os.MkdirTemp(root, "."+name+".")
	if e != nil {
		return "", e
	}
	defer os.RemoveAll(tmp)
	candidate := filepath.Join(tmp, name)
	if e = os.Mkdir(candidate, 0755); e != nil {
		return "", e
	}
	if e = os.WriteFile(filepath.Join(candidate, SnapshotFile), config, 0644); e != nil {
		return "", e
	}
	for name, p := range plans {
		base := filepath.Join(candidate, name)
		if e = os.MkdirAll(base, 0755); e != nil {
			return "", e
		}
		instruction := p.InstructionsName
		if e = os.WriteFile(filepath.Join(base, instruction), []byte(p.Agents), 0644); e != nil {
			return "", e
		}
		for _, dir := range SortedKeys(p.Dirs) {
			if e = os.MkdirAll(filepath.Join(base, dir), 0755); e != nil {
				return "", e
			}
		}
		for _, rel := range SortedKeys(p.Files) {
			target := filepath.Join(base, rel)
			if e = os.MkdirAll(filepath.Dir(target), 0755); e != nil {
				return "", e
			}
			if e = CopyPlanFile(p, rel, target); e != nil {
				return "", e
			}
		}
	}
	var paths []string
	if e = walk(candidate, func(path string, _ os.DirEntry) error { paths = append(paths, path); return nil }); e != nil {
		return "", e
	}
	sort.Slice(paths, func(i, j int) bool { return len(paths[i]) > len(paths[j]) })
	for _, path := range paths {
		i, e := os.Stat(path)
		if e != nil {
			return "", e
		}
		if e = os.Chmod(path, i.Mode().Perm()&^0222); e != nil {
			return "", e
		}
	}
	// Keep the candidate root writable until the rename; the children are
	// already read-only. Preserve the umask when locking down the root.
	if e = os.Rename(candidate, dest); e != nil {
		return "", e
	}
	info, e := os.Stat(dest)
	if e != nil {
		return "", e
	}
	if e = os.Chmod(dest, info.Mode().Perm()&^0222); e != nil {
		return "", e
	}
	if e = updateAlias(root, dest); e != nil {
		return "", e
	}
	SummaryPlans(plans, dest, "created")
	return dest, nil
}
