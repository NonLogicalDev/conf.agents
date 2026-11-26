package engine

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
)

const targetInfoFile = ".target.json"

type targetInfo struct {
	Destination string `json:"destination"`
	Harness     string `json:"harness"`
	Generation  string `json:"generation"`
	Device      uint64 `json:"device"`
	Inode       uint64 `json:"inode"`
}

func TargetHash(destination string) string {
	sum := sha256.Sum256([]byte(destination))
	return hex.EncodeToString(sum[:])
}
func targetBytes(info targetInfo) []byte {
	b, _ := json.Marshal(info)
	return append(b, '\n')
}

// destinationIdentity identifies the directory and its mounted filesystem.
func destinationIdentity(destination, gen string, create bool) (targetInfo, error) {
	result := targetInfo{Destination: destination, Generation: filepath.Base(gen)}
	if err := ValidDir(destination, "harness destination"); err != nil {
		return result, err
	}
	if create {
		if err := os.MkdirAll(destination, 0755); err != nil {
			return result, err
		}
	}
	info, err := os.Lstat(destination)
	if err != nil {
		return result, fmt.Errorf("cannot verify recorded harness destination %s: %w", destination, err)
	}
	if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return result, fmt.Errorf("harness destination is not a regular directory: %s", destination)
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return result, fmt.Errorf("cannot identify harness destination: %s", destination)
	}
	result.Device, result.Inode = uint64(stat.Dev), uint64(stat.Ino)
	return result, nil
}

// TargetPlan resolves destination references without changing the generation.
func TargetPlan(p *RenderPlan, destination string) (*RenderPlan, error) {
	p = targetWithoutPluginSkills(p)
	result := *p
	result.Agents = strings.ReplaceAll(strings.ReplaceAll(p.Agents, skillsPlaceholder, filepath.Join(destination, "skills")), resourcesPlaceholder, filepath.Join(destination, "resources"))
	result.Contents = make(map[string][]byte, len(p.Contents)+1)
	for k, v := range p.Contents {
		result.Contents[k] = v
	}
	guidance := filepath.Join("skills", "Agent_Guidance", "SKILL.md")
	if source, ok := p.Files[guidance]; ok {
		original, exists := p.Contents[guidance]
		if !exists {
			var err error
			original, err = os.ReadFile(source)
			if err != nil {
				return nil, err
			}
		}
		result.Contents[guidance] = []byte(strings.ReplaceAll(string(original), skillsPlaceholder, filepath.Join(destination, "skills")))
	}
	return &result, nil
}
func targetMatches(root string, p *RenderPlan, config, info []byte) bool {
	st, e := os.Lstat(root)
	if e != nil || !st.IsDir() || st.Mode().Perm()&0222 != 0 {
		return false
	}
	for name, want := range map[string][]byte{p.InstructionsName: []byte(p.Agents), SnapshotFile: config, targetInfoFile: info} {
		b, e := os.ReadFile(filepath.Join(root, name))
		if e != nil || !bytes.Equal(b, want) {
			return false
		}
	}
	expected := map[string]bool{p.InstructionsName: true, SnapshotFile: true, targetInfoFile: true}
	for d := range p.Dirs {
		expected[d] = true
	}
	for f := range p.Files {
		expected[f] = true
	}
	count := 0
	e = walk(root, func(path string, d os.DirEntry) error {
		rel, _ := filepath.Rel(root, path)
		i, e := d.Info()
		if e != nil {
			return e
		}
		if d.Type()&os.ModeSymlink != 0 || i.Mode().Perm()&0222 != 0 || !expected[rel] {
			return fmt.Errorf("unexpected target path %s", path)
		}
		count++
		return nil
	})
	if e != nil || count != len(expected) {
		return false
	}
	for dir := range p.Dirs {
		if !IsDir(filepath.Join(root, dir)) {
			return false
		}
	}
	for rel, src := range p.Files {
		dst := filepath.Join(root, rel)
		i, e := os.Stat(dst)
		if e != nil || !i.Mode().IsRegular() {
			return false
		}
		mode, ok := p.fileModes[rel]
		if !ok {
			j, e := os.Stat(src)
			if e != nil {
				return false
			}
			mode = j.Mode().Perm()
		}
		if i.Mode().Perm() != mode&^0222 {
			return false
		}
		want, ok := p.Contents[rel]
		if !ok {
			want, e = os.ReadFile(src)
			if e != nil {
				return false
			}
		}
		got, e := os.ReadFile(dst)
		if e != nil || !bytes.Equal(got, want) {
			return false
		}
	}
	return true
}

// PreflightHarnessTarget verifies a previously materialized snapshot without
// creating the destination. All harnesses can be checked before activation.
func (registry Registry) PreflightHarnessTarget(root, gen, harness, destination string, p *RenderPlan) error {
	for _, dir := range []string{filepath.Join(root, ".targets"), filepath.Join(root, ".targets", harness), filepath.Join(root, ".targets", harness, TargetHash(destination))} {
		if err := ValidDir(dir, "target directory"); err != nil {
			return err
		}
	}
	if err := ValidDir(destination, harness+" destination"); err != nil {
		return err
	}
	target := filepath.Join(root, ".targets", harness, TargetHash(destination), filepath.Base(gen))
	_, err := os.Lstat(target)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	resolved, err := TargetPlan(p, destination)
	if err != nil {
		return err
	}
	config, err := os.ReadFile(filepath.Join(gen, SnapshotFile))
	if err != nil {
		return err
	}
	identity, err := destinationIdentity(destination, gen, false)
	if err != nil {
		return err
	}
	identity.Harness = harness
	if !targetMatches(target, resolved, config, targetBytes(identity)) {
		return fmt.Errorf("refusing invalid target snapshot or changed %s destination: %s", harness, target)
	}
	return nil
}

// MaterializeHarnessTarget creates a private immutable destination-specific snapshot.
func (registry Registry) MaterializeHarnessTarget(root, gen, harness, destination string, p *RenderPlan) (string, *RenderPlan, error) {
	targetPlan, e := TargetPlan(p, destination)
	if e != nil {
		return "", nil, e
	}
	config, e := os.ReadFile(filepath.Join(gen, SnapshotFile))
	if e != nil {
		return "", nil, e
	}
	parent := filepath.Join(root, ".targets", harness, TargetHash(destination))
	for _, dir := range []string{filepath.Join(root, ".targets"), filepath.Join(root, ".targets", harness), parent} {
		if e = ValidDir(dir, "target directory"); e != nil {
			return "", nil, e
		}
	}
	target := filepath.Join(parent, filepath.Base(gen))
	_, statErr := os.Lstat(target)
	if statErr != nil && !os.IsNotExist(statErr) {
		return "", nil, statErr
	}
	existsTarget := statErr == nil
	identity, e := destinationIdentity(destination, gen, !existsTarget)
	if e != nil {
		return "", nil, e
	}
	identity.Harness = harness
	info := targetBytes(identity)
	if existsTarget {
		if !targetMatches(target, targetPlan, config, info) {
			return "", nil, fmt.Errorf("refusing invalid target snapshot or changed harness destination: %s", target)
		}
		return target, targetPlan, nil
	}
	if e = os.MkdirAll(parent, 0755); e != nil {
		return "", nil, e
	}
	tmp, e := os.MkdirTemp(parent, ".stage-")
	if e != nil {
		return "", nil, e
	}
	defer os.RemoveAll(tmp)
	candidate := filepath.Join(tmp, filepath.Base(gen))
	if e = os.Mkdir(candidate, 0755); e != nil {
		return "", nil, e
	}
	for name, data := range map[string][]byte{targetPlan.InstructionsName: []byte(targetPlan.Agents), SnapshotFile: config, targetInfoFile: info} {
		if e = os.WriteFile(filepath.Join(candidate, name), data, 0644); e != nil {
			return "", nil, e
		}
	}
	for _, dir := range SortedKeys(targetPlan.Dirs) {
		if e = os.MkdirAll(filepath.Join(candidate, dir), 0755); e != nil {
			return "", nil, e
		}
	}
	for _, rel := range SortedKeys(targetPlan.Files) {
		dst := filepath.Join(candidate, rel)
		if e = os.MkdirAll(filepath.Dir(dst), 0755); e != nil {
			return "", nil, e
		}
		if e = CopyPlanFile(targetPlan, rel, dst); e != nil {
			return "", nil, e
		}
	}
	var paths []string
	if e = walk(candidate, func(path string, _ os.DirEntry) error { paths = append(paths, path); return nil }); e != nil {
		return "", nil, e
	}
	sort.Slice(paths, func(i, j int) bool { return len(paths[i]) > len(paths[j]) })
	for _, path := range paths {
		st, e := os.Stat(path)
		if e != nil {
			return "", nil, e
		}
		if e = os.Chmod(path, st.Mode().Perm()&^0222); e != nil {
			return "", nil, e
		}
	}
	if e = os.Rename(candidate, target); e != nil {
		return "", nil, e
	}
	st, e := os.Stat(target)
	if e != nil {
		return "", nil, e
	}
	if e = os.Chmod(target, st.Mode().Perm()&^0222); e != nil {
		return "", nil, e
	}
	return target, targetPlan, nil
}

// pruneTargets inspects every recorded destination before removing any stale
// target of an eligible generation. Unknown or inaccessible state blocks prune.
func (registry Registry) pruneTargets(root string, generations []generation) error {
	targetsRoot := filepath.Join(root, ".targets")
	if err := ValidDir(targetsRoot, "targets root"); err != nil {
		return err
	}
	if !Exists(targetsRoot) {
		return nil
	}
	var stale []string
	entriesAtRoot, err := Entries(targetsRoot)
	if err != nil {
		return err
	}
	for _, first := range entriesAtRoot {
		firstPath := filepath.Join(targetsRoot, first.Name())
		if IsLink(firstPath) || !first.IsDir() {
			return fmt.Errorf("invalid target directory: %s", firstPath)
		}
		harness := first.Name()
		if _, ok := registry[harness]; !ok {
			return fmt.Errorf("invalid target harness: %s", firstPath)
		}
		hashes, err := Entries(firstPath)
		if err != nil {
			return err
		}
		for _, h := range hashes {
			hashRoot := filepath.Join(firstPath, h.Name())
			if IsLink(hashRoot) || !h.IsDir() {
				return fmt.Errorf("invalid target directory: %s", hashRoot)
			}
			for _, gen := range generations {
				target := filepath.Join(hashRoot, filepath.Base(gen.Path))
				st, err := os.Lstat(target)
				if os.IsNotExist(err) {
					continue
				}
				if err != nil {
					return err
				}
				if st.Mode()&os.ModeSymlink != 0 || !st.IsDir() {
					return fmt.Errorf("invalid target snapshot: %s", target)
				}
				data, err := os.ReadFile(filepath.Join(target, targetInfoFile))
				if err != nil {
					return fmt.Errorf("read target metadata: %w", err)
				}
				var info targetInfo
				if err = json.Unmarshal(data, &info); err != nil || !filepath.IsAbs(info.Destination) || filepath.Clean(info.Destination) != info.Destination || info.Generation != filepath.Base(gen.Path) || h.Name() != TargetHash(info.Destination) || info.Harness != harness || !bytes.Equal(data, targetBytes(info)) {
					return fmt.Errorf("invalid target metadata: %s", target)
				}
				snap, err := registry.ReadSnapshot(gen.Path)
				if err != nil {
					return err
				}
				plan, err := registry.restoreHarnessPlan(gen.Path, harness, snap)
				if err != nil {
					return err
				}
				expected, err := TargetPlan(plan, info.Destination)
				if err != nil {
					return err
				}
				config, err := os.ReadFile(filepath.Join(gen.Path, SnapshotFile))
				if err != nil {
					return err
				}
				if !targetMatches(target, expected, config, data) {
					return fmt.Errorf("invalid target snapshot: %s", target)
				}
				current, err := destinationIdentity(info.Destination, gen.Path, false)
				if err != nil {
					return err
				}
				if current.Device != info.Device || current.Inode != info.Inode {
					return fmt.Errorf("recorded %s destination identity changed: %s", harness, info.Destination)
				}
				links := []string{filepath.Join(info.Destination, plan.InstructionsName)}
				for _, component := range registry[harness].LinkComponents {
					dir := filepath.Join(info.Destination, component)
					if err = ValidDir(dir, component+" target"); err != nil {
						return err
					}
					if !IsDir(dir) {
						continue
					}
					es, err := Entries(dir)
					if err != nil {
						return err
					}
					for _, entry := range es {
						links = append(links, filepath.Join(dir, entry.Name()))
					}
				}
				for _, link := range links {
					if !IsLink(link) {
						continue
					}
					dest, err := os.Readlink(link)
					if err != nil {
						return err
					}
					if !filepath.IsAbs(dest) {
						dest = filepath.Join(filepath.Dir(link), dest)
					}
					dest, err = filepath.Abs(dest)
					if err != nil {
						return err
					}
					if Within(dest, target) {
						return fmt.Errorf("refusing to prune active %s generation: %s (%s)", harness, filepath.Base(gen.Path), link)
					}
				}
				stale = append(stale, target)
			}
		}
	}
	for _, target := range stale {
		if err := removeGeneration(target); err != nil {
			return err
		}
	}
	return nil
}
