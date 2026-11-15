package internal

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
	CodexDir   string `json:"codex_dir"`
	Harness    string `json:"harness,omitempty"`
	Generation string `json:"generation"`
	Device     uint64 `json:"device"`
	Inode      uint64 `json:"inode"`
}

func targetHash(codex string) string {
	sum := sha256.Sum256([]byte(codex))
	return hex.EncodeToString(sum[:])
}
func targetBytes(info targetInfo) []byte {
	b, _ := json.Marshal(info)
	return append(b, '\n')
}

// destinationIdentity identifies the directory and its mounted filesystem.
func destinationIdentity(codex, gen string, create bool) (targetInfo, error) {
	result := targetInfo{CodexDir: codex, Generation: filepath.Base(gen)}
	if err := validDir(codex, "Codex destination"); err != nil {
		return result, err
	}
	if create {
		if err := os.MkdirAll(codex, 0755); err != nil {
			return result, err
		}
	}
	info, err := os.Lstat(codex)
	if err != nil {
		return result, fmt.Errorf("cannot verify recorded Codex destination %s: %w", codex, err)
	}
	if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return result, fmt.Errorf("Codex destination is not a regular directory: %s", codex)
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return result, fmt.Errorf("cannot identify Codex destination: %s", codex)
	}
	result.Device, result.Inode = uint64(stat.Dev), uint64(stat.Ino)
	return result, nil
}

// targetPlan resolves destination references without changing the generation.
func targetPlan(p *renderPlan, codex string) (*renderPlan, error) {
	p = targetWithoutPluginSkills(p)
	result := *p
	result.agents = strings.ReplaceAll(strings.ReplaceAll(p.agents, skillsPlaceholder, filepath.Join(codex, "skills")), resourcesPlaceholder, filepath.Join(codex, "resources"))
	result.contents = make(map[string][]byte, len(p.contents)+1)
	for k, v := range p.contents {
		result.contents[k] = v
	}
	guidance := filepath.Join("skills", "Agent_Guidance", "SKILL.md")
	if source, ok := p.files[guidance]; ok {
		original, exists := p.contents[guidance]
		if !exists {
			var err error
			original, err = os.ReadFile(source)
			if err != nil {
				return nil, err
			}
		}
		result.contents[guidance] = []byte(strings.ReplaceAll(string(original), skillsPlaceholder, filepath.Join(codex, "skills")))
	}
	return &result, nil
}
func instructionsName(p *renderPlan) string {
	if p.instructionsName != "" {
		return p.instructionsName
	}
	return "AGENTS.md"
}

func targetMatches(root string, p *renderPlan, config, info []byte) bool {
	st, e := os.Lstat(root)
	if e != nil || !st.IsDir() || st.Mode().Perm()&0222 != 0 {
		return false
	}
	for name, want := range map[string][]byte{instructionsName(p): []byte(p.agents), snapshotFile: config, targetInfoFile: info} {
		b, e := os.ReadFile(filepath.Join(root, name))
		if e != nil || !bytes.Equal(b, want) {
			return false
		}
	}
	expected := map[string]bool{instructionsName(p): true, snapshotFile: true, targetInfoFile: true}
	for d := range p.dirs {
		expected[d] = true
	}
	for f := range p.files {
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
	for dir := range p.dirs {
		if !isDir(filepath.Join(root, dir)) {
			return false
		}
	}
	for rel, src := range p.files {
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
		want, ok := p.contents[rel]
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

// preflightHarnessTarget verifies a previously materialized snapshot without
// creating the destination. All harnesses can be checked before activation.
func preflightHarnessTarget(root, gen, harness, destination string, p *renderPlan) error {
	for _, dir := range []string{filepath.Join(root, ".targets"), filepath.Join(root, ".targets", harness), filepath.Join(root, ".targets", harness, targetHash(destination))} {
		if err := validDir(dir, "target directory"); err != nil {
			return err
		}
	}
	if err := validDir(destination, harness+" destination"); err != nil {
		return err
	}
	target := filepath.Join(root, ".targets", harness, targetHash(destination), filepath.Base(gen))
	_, err := os.Lstat(target)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	snap, err := readSnapshot(gen)
	if err != nil {
		return err
	}
	resolved, err := resolvedHarnessTarget(p, destination, harness, snap)
	if err != nil {
		return err
	}
	config, err := os.ReadFile(filepath.Join(gen, snapshotFile))
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

// materializeTarget creates a private immutable destination-specific snapshot.
func materializeHarnessTarget(root, gen, harness, codex string, p *renderPlan) (string, *renderPlan, error) {
	snap, e := readSnapshot(gen)
	if e != nil {
		return "", nil, e
	}
	targetPlan, e := resolvedHarnessTarget(p, codex, harness, snap)
	if e != nil {
		return "", nil, e
	}
	config, e := os.ReadFile(filepath.Join(gen, snapshotFile))
	if e != nil {
		return "", nil, e
	}
	parent := filepath.Join(root, ".targets", harness, targetHash(codex))
	for _, dir := range []string{filepath.Join(root, ".targets"), filepath.Join(root, ".targets", harness), parent} {
		if e = validDir(dir, "target directory"); e != nil {
			return "", nil, e
		}
	}
	target := filepath.Join(parent, filepath.Base(gen))
	_, statErr := os.Lstat(target)
	if statErr != nil && !os.IsNotExist(statErr) {
		return "", nil, statErr
	}
	existsTarget := statErr == nil
	identity, e := destinationIdentity(codex, gen, !existsTarget)
	if e != nil {
		return "", nil, e
	}
	identity.Harness = harness
	info := targetBytes(identity)
	if existsTarget {
		if !targetMatches(target, targetPlan, config, info) {
			return "", nil, fmt.Errorf("refusing invalid target snapshot or changed Codex destination: %s", target)
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
	for name, data := range map[string][]byte{instructionsName(targetPlan): []byte(targetPlan.agents), snapshotFile: config, targetInfoFile: info} {
		if e = os.WriteFile(filepath.Join(candidate, name), data, 0644); e != nil {
			return "", nil, e
		}
	}
	for _, dir := range sortedKeys(targetPlan.dirs) {
		if e = os.MkdirAll(filepath.Join(candidate, dir), 0755); e != nil {
			return "", nil, e
		}
	}
	for _, rel := range sortedKeys(targetPlan.files) {
		dst := filepath.Join(candidate, rel)
		if e = os.MkdirAll(filepath.Dir(dst), 0755); e != nil {
			return "", nil, e
		}
		if e = copyPlanFile(targetPlan, rel, dst); e != nil {
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

// materializeTarget is the Codex specialization retained for callers.
func materializeTarget(root, gen, codex string, p *renderPlan) (string, *renderPlan, error) {
	return materializeHarnessTarget(root, gen, "codex", codex, p)
}

// pruneTargets inspects every recorded destination before removing any stale
// target of an eligible generation. Unknown or inaccessible state blocks prune.
func pruneTargets(root string, generations []generation) error {
	targetsRoot := filepath.Join(root, ".targets")
	if err := validDir(targetsRoot, "targets root"); err != nil {
		return err
	}
	if !exists(targetsRoot) {
		return nil
	}
	var stale []string
	entriesAtRoot, err := entries(targetsRoot)
	if err != nil {
		return err
	}
	for _, first := range entriesAtRoot {
		firstPath := filepath.Join(targetsRoot, first.Name())
		if isLink(firstPath) || !first.IsDir() {
			return fmt.Errorf("invalid target directory: %s", firstPath)
		}
		harness := first.Name()
		var hashes []os.DirEntry
		legacy := len(harness) == 64
		if legacy {
			hashes = []os.DirEntry{first}
			harness = "codex"
		} else {
			if _, ok := harnessRegistry()[harness]; !ok {
				return fmt.Errorf("invalid target harness: %s", firstPath)
			}
			hashes, err = entries(firstPath)
			if err != nil {
				return err
			}
		}
		for _, h := range hashes {
			hashRoot := filepath.Join(firstPath, h.Name())
			if legacy {
				hashRoot = firstPath
			}
			if isLink(hashRoot) || !h.IsDir() {
				return fmt.Errorf("invalid target directory: %s", hashRoot)
			}
			for _, gen := range generations {
				target := filepath.Join(hashRoot, filepath.Base(gen.path))
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
				expectedHarness := harness
				if legacy {
					expectedHarness = ""
				}
				if err = json.Unmarshal(data, &info); err != nil || !filepath.IsAbs(info.CodexDir) || filepath.Clean(info.CodexDir) != info.CodexDir || info.Generation != filepath.Base(gen.path) || h.Name() != targetHash(info.CodexDir) || info.Harness != expectedHarness || !bytes.Equal(data, targetBytes(info)) {
					return fmt.Errorf("invalid target metadata: %s", target)
				}
				snap, err := readSnapshot(gen.path)
				if err != nil {
					return err
				}
				plan, err := restoreHarnessPlan(gen.path, harness, snap)
				if err != nil {
					return err
				}
				expected, err := resolvedHarnessTarget(plan, info.CodexDir, harness, snap)
				if err != nil {
					return err
				}
				config, err := os.ReadFile(filepath.Join(gen.path, snapshotFile))
				if err != nil {
					return err
				}
				if !targetMatches(target, expected, config, data) {
					// Older targets can contain portable references. Accept their exact
					// original bytes for ownership checks while preserving active links.
					previous, err := targetPlan(plan, info.CodexDir)
					if err != nil {
						return err
					}
					if !targetMatches(target, previous, config, data) {
						namespace := "local"
						if c, ok := snap["codex"].(snapshotCodex); harness == "codex" && ok && c.Plugin != "none" {
							namespace = c.PluginName
						}
						previous, err = resolvePlanReferences(previous, namespace, "$"+namespace+":", "local", "$local:")
						if err != nil || !targetMatches(target, previous, config, data) {
							return fmt.Errorf("invalid target snapshot: %s", target)
						}
					}
				}
				current, err := destinationIdentity(info.CodexDir, gen.path, false)
				if err != nil {
					return err
				}
				if current.Device != info.Device || current.Inode != info.Inode {
					return fmt.Errorf("recorded %s destination identity changed: %s", harness, info.CodexDir)
				}
				links := []string{filepath.Join(info.CodexDir, instructionsName(plan))}
				components := []string{"skills", "resources"}
				if harness == "codex" {
					components = append(components, "rules")
				}
				if harness == "opencode" {
					components = append(components, "plugins")
				}
				for _, component := range components {
					dir := filepath.Join(info.CodexDir, component)
					if err = validDir(dir, component+" target"); err != nil {
						return err
					}
					if !isDir(dir) {
						continue
					}
					es, err := entries(dir)
					if err != nil {
						return err
					}
					for _, entry := range es {
						links = append(links, filepath.Join(dir, entry.Name()))
					}
				}
				for _, link := range links {
					if !isLink(link) {
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
					if within(dest, target) {
						return fmt.Errorf("refusing to prune active %s generation: %s (%s)", harness, filepath.Base(gen.path), link)
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
