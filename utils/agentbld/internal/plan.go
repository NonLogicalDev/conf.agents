package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"

	"gopkg.in/yaml.v3"
)

const marker = "<!-- generated: agent-config; do not edit -->"
const resourcesPlaceholder = "{{%_resources_%}}"
const skillsPlaceholder = "{{%_skills_%}}"
const agentsPlaceholder = "{{%_agents_%}}"

var layerPattern = regexp.MustCompile(`^\d+-[A-Za-z0-9][A-Za-z0-9._-]*$`)
var namePattern = regexp.MustCompile(`^name\s*:\s*(?:'([^'"\s#]+)'|"([^'"\s#]+)"|([^'"\s#]+))\s*(?:#.*)?$`)
var descPattern = regexp.MustCompile(`^description\s*:\s*\S`)

type renderPlan struct {
	instructionsName string
	layers           []string
	skills           map[string]string
	resources        map[string]string
	agents           string
	fragments        int
	dirs             map[string]bool
	files            map[string]string
	contents         map[string][]byte
	fileModes        map[string]os.FileMode
}

func entries(p string) ([]os.DirEntry, error) { return os.ReadDir(p) }
func isLink(p string) bool                    { i, e := os.Lstat(p); return e == nil && i.Mode()&os.ModeSymlink != 0 }
func exists(p string) bool                    { _, e := os.Lstat(p); return e == nil }
func isDir(p string) bool                     { i, e := os.Stat(p); return e == nil && i.IsDir() }
func isRegular(p string) bool                 { i, e := os.Stat(p); return e == nil && i.Mode().IsRegular() }
func walk(root string, fn func(string, os.DirEntry) error) error {
	return filepath.WalkDir(root, func(p string, d os.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if p == root {
			return nil
		}
		return fn(p, d)
	})
}
func parents(p string) []string {
	var result []string
	for p = filepath.Dir(p); p != "." && p != "/"; p = filepath.Dir(p) {
		result = append(result, p)
	}
	return result
}
func within(p, root string) bool {
	r, e := filepath.Rel(root, p)
	return e == nil && r != ".." && !strings.HasPrefix(r, ".."+string(filepath.Separator))
}
func readSkillName(path string) (string, error) {
	b, e := os.ReadFile(path)
	if e != nil {
		return "", fmt.Errorf("cannot read skill: %s: %w", path, e)
	}
	if !utf8.Valid(b) {
		return "", fmt.Errorf("invalid UTF-8 in skill: %s", path)
	}
	return skillNameFromBytes(b, path, filepath.Base(filepath.Dir(path)))
}

func skillNameFromBytes(b []byte, path, expectedName string) (string, error) {
	if !utf8.Valid(b) {
		return "", fmt.Errorf("invalid UTF-8 in skill: %s", path)
	}
	lines := strings.Split(string(b), "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return "", fmt.Errorf("skill has no YAML frontmatter: %s", path)
	}
	name := ""
	desc := false
	closed := false
	for _, line := range lines[1:] {
		if strings.TrimSpace(line) == "---" || strings.TrimSpace(line) == "..." {
			closed = true
			break
		}
		if m := namePattern.FindStringSubmatch(line); m != nil {
			if name != "" {
				return "", fmt.Errorf("skill has duplicate frontmatter names: %s", path)
			}
			for _, candidate := range m[1:] {
				if candidate != "" {
					name = candidate
					break
				}
			}
		}
		if descPattern.MatchString(line) {
			desc = true
		}
	}
	if !closed {
		return "", fmt.Errorf("skill has unterminated YAML frontmatter: %s", path)
	}
	if name == "" {
		return "", fmt.Errorf("skill has no frontmatter name: %s", path)
	}
	if !desc {
		return "", fmt.Errorf("skill has no frontmatter description: %s", path)
	}
	if name != expectedName {
		return "", fmt.Errorf("skill name %q does not match directory %q: %s", name, expectedName, path)
	}
	return name, nil
}

// mergeSkillMetadata merges mappings recursively while replacing scalar and sequence values.
func mergeSkillMetadata(dst, src *yaml.Node) {
	for i := 0; i < len(src.Content); i += 2 {
		key, value := src.Content[i], src.Content[i+1]
		index := -1
		for j := 0; j < len(dst.Content); j += 2 {
			if dst.Content[j].Value == key.Value {
				index = j
				break
			}
		}
		if index < 0 {
			dst.Content = append(dst.Content, key, value)
			continue
		}
		previous := dst.Content[index+1]
		if previous.Kind == yaml.MappingNode && value.Kind == yaml.MappingNode {
			mergeSkillMetadata(previous, value)
		} else {
			dst.Content[index+1] = value
		}
	}
}

func validateSkillMetadataNode(node *yaml.Node) error {
	switch node.Kind {
	case yaml.MappingNode:
		keys := map[string]bool{}
		for i := 0; i < len(node.Content); i += 2 {
			key, value := node.Content[i], node.Content[i+1]
			if key.Kind != yaml.ScalarNode || key.Tag != "!!str" {
				return fmt.Errorf("metadata mapping keys must be strings")
			}
			if keys[key.Value] {
				return fmt.Errorf("duplicate metadata key %q", key.Value)
			}
			keys[key.Value] = true
			if err := validateSkillMetadataNode(value); err != nil {
				return err
			}
		}
	case yaml.SequenceNode:
		for _, child := range node.Content {
			if err := validateSkillMetadataNode(child); err != nil {
				return err
			}
		}
	case yaml.AliasNode:
		return fmt.Errorf("metadata aliases are not supported")
	}
	return nil
}

// fragmentForHarness reads fragment metadata and removes the renderer-only selector.
func fragmentForHarness(path, harness string) ([]byte, bool, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, false, err
	}
	if !utf8.Valid(b) {
		return nil, false, fmt.Errorf("invalid UTF-8 in fragment: %s", path)
	}
	lines := strings.SplitAfter(string(b), "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return b, true, nil
	}
	end := 1
	for end < len(lines) && strings.TrimSpace(lines[end]) != "---" && strings.TrimSpace(lines[end]) != "..." {
		end++
	}
	if end == len(lines) {
		return nil, false, fmt.Errorf("unterminated YAML frontmatter: %s", path)
	}
	var doc yaml.Node
	if err := yaml.Unmarshal([]byte(strings.Join(lines[1:end], "")), &doc); err != nil {
		return nil, false, fmt.Errorf("invalid YAML frontmatter: %s: %w", path, err)
	}
	if len(doc.Content) == 0 || doc.Content[0].Kind != yaml.MappingNode {
		return nil, false, fmt.Errorf("YAML frontmatter must be a mapping: %s", path)
	}
	meta := doc.Content[0]
	if err := validateSkillMetadataNode(meta); err != nil {
		return nil, false, fmt.Errorf("invalid YAML frontmatter: %s: %w", path, err)
	}
	for i := 0; i < len(meta.Content); i += 2 {
		if meta.Content[i].Value != "harnesses" {
			continue
		}
		list := meta.Content[i+1]
		if list.Kind != yaml.SequenceNode || len(list.Content) == 0 {
			return nil, false, fmt.Errorf("harnesses must be a non-empty list: %s", path)
		}
		match := false
		seen := map[string]bool{}
		for _, item := range list.Content {
			if item.Kind != yaml.ScalarNode || item.Tag != "!!str" {
				return nil, false, fmt.Errorf("harnesses must contain harness names: %s", path)
			}
			if _, ok := harnessRegistry()[item.Value]; !ok {
				return nil, false, fmt.Errorf("unknown harness %q in %s", item.Value, path)
			}
			if seen[item.Value] {
				return nil, false, fmt.Errorf("duplicate harness %q in %s", item.Value, path)
			}
			seen[item.Value] = true
			match = match || item.Value == harness
		}
		if !match {
			return nil, false, nil
		}
		meta.Content = append(meta.Content[:i], meta.Content[i+2:]...)
		body := strings.Join(lines[end+1:], "")
		if len(meta.Content) == 0 {
			return []byte(body), true, nil
		}
		y, err := yaml.Marshal(meta)
		if err != nil {
			return nil, false, fmt.Errorf("cannot render frontmatter: %s: %w", path, err)
		}
		return []byte("---\n" + string(y) + "---\n" + body), true, nil
	}
	return b, true, nil
}

func buildPlan(sections []string, includeGuidance bool, harness string) (*renderPlan, error) {
	p := &renderPlan{skills: map[string]string{}, resources: map[string]string{}, dirs: map[string]bool{"skills": true, "resources": true}, files: map[string]string{}, contents: map[string][]byte{}}
	layerSources := map[string]string{}
	for index, root := range sections {
		if isLink(root) {
			return nil, fmt.Errorf("symlinked source sections are not supported: %s", root)
		}
		if !isDir(root) {
			return nil, fmt.Errorf("missing or invalid sections directory: %s", root)
		}
		entries, err := entries(root)
		if err != nil {
			return nil, err
		}
		prefix := filepath.Base(root)
		if len(sections) > 1 {
			prefix = fmt.Sprintf("sections[%d]", index+1)
		}
		for _, d := range entries {
			if !layerPattern.MatchString(d.Name()) {
				continue
			}
			path := filepath.Join(root, d.Name())
			if isLink(path) {
				return nil, fmt.Errorf("symlinked source layers are not supported: %s", path)
			}
			if d.IsDir() {
				p.layers = append(p.layers, path)
				layerSources[path] = prefix
			}
		}
	}
	fragments := map[string]string{}
	fragmentSources := map[string]string{}
	skillFragments := map[string]map[string]string{}
	skillBases := map[string]string{}
	skillAssetLayers := map[string][]string{}
	for _, layer := range p.layers {
		for _, component := range []string{"skills", "resources", "AGENTS.md.d"} {
			dir := filepath.Join(layer, component)
			if isLink(dir) {
				return nil, fmt.Errorf("symlinked source directory is not supported: %s", dir)
			}
			if _, err := os.Lstat(dir); err != nil {
				if os.IsNotExist(err) {
					continue
				}
				return nil, err
			}
			if !isDir(dir) {
				return nil, fmt.Errorf("%s root is not a regular directory: %s", component, dir)
			}
			switch component {
			case "skills":
				ds, e := entries(dir)
				if e != nil {
					return nil, e
				}
				for _, d := range ds {
					// Skills-only plugin builds exclude this skill before inspecting its contents.
					if !includeGuidance && d.Name() == "Agent_Guidance" {
						continue
					}
					if strings.HasPrefix(d.Name(), ".") {
						continue
					}
					path := filepath.Join(dir, d.Name())
					if isLink(path) {
						return nil, fmt.Errorf("symlinked skill directories are not supported: %s", path)
					}
					if !d.IsDir() {
						continue
					}
					skill := filepath.Join(path, "SKILL.md")
					fragmentDir := filepath.Join(path, "SKILL.md.d")
					if isLink(skill) || (exists(skill) && !isRegular(skill)) {
						return nil, fmt.Errorf("skill is not a regular file: %s", skill)
					}
					if isLink(fragmentDir) || (exists(fragmentDir) && !isDir(fragmentDir)) {
						return nil, fmt.Errorf("skill fragments are not a regular directory: %s", fragmentDir)
					}
					if exists(skill) {
						skillBases[d.Name()] = skill
						delete(skillFragments, d.Name())
						skillAssetLayers[d.Name()] = nil
					}
					if exists(fragmentDir) {
						if skillFragments[d.Name()] == nil {
							skillFragments[d.Name()] = map[string]string{}
						}
						fragments, e := entries(fragmentDir)
						if e != nil {
							return nil, e
						}
						for _, fragment := range fragments {
							if filepath.Ext(fragment.Name()) != ".md" {
								continue
							}
							fragmentPath := filepath.Join(fragmentDir, fragment.Name())
							if isLink(fragmentPath) {
								return nil, fmt.Errorf("symlinked skill fragments are not supported: %s", fragmentPath)
							}
							if fragment.Type().IsRegular() {
								_, include, err := fragmentForHarness(fragmentPath, harness)
								if err != nil {
									return nil, err
								}
								if !include {
									continue
								}
								skillFragments[d.Name()][fragment.Name()] = fragmentPath
							}
						}
					} else if !exists(skill) {
						empty := true
						walkErr := walk(path, func(_ string, entry os.DirEntry) error {
							if !entry.IsDir() || entry.Type()&os.ModeSymlink != 0 {
								empty = false
							}
							return nil
						})
						if walkErr != nil {
							return nil, walkErr
						}
						if empty {
							continue
						}
						return nil, fmt.Errorf("skill directory has no regular SKILL.md or SKILL.md.d: %s", path)
					}
					if prev, ok := p.skills[d.Name()]; ok {
						fmt.Printf("skill override %s: %s -> %s\n", d.Name(), prev, path)
					}
					p.skills[d.Name()] = path
					skillAssetLayers[d.Name()] = append(skillAssetLayers[d.Name()], path)
				}
			case "resources":
				e := walk(dir, func(path string, d os.DirEntry) error {
					if d.Type()&os.ModeSymlink != 0 {
						return fmt.Errorf("symlinked resources are not supported: %s", path)
					}
					if d.IsDir() {
						return nil
					}
					if !d.Type().IsRegular() {
						return fmt.Errorf("resource is not a regular file: %s", path)
					}
					rel, _ := filepath.Rel(dir, path)
					if strings.Split(rel, string(filepath.Separator))[0] == "references" {
						return fmt.Errorf("place grouped resources directly under resources, without a references directory: %s", path)
					}
					for _, prev := range sortedKeys(p.resources) {
						if rel == prev || within(rel, prev) || within(prev, rel) {
							fmt.Printf("resource override %s: %s -> %s\n", prev, p.resources[prev], path)
							delete(p.resources, prev)
						}
					}
					p.resources[rel] = path
					return nil
				})
				if e != nil {
					return nil, e
				}
			case "AGENTS.md.d":
				ds, e := entries(dir)
				if e != nil {
					return nil, e
				}
				for _, d := range ds {
					if filepath.Ext(d.Name()) != ".md" {
						continue
					}
					path := filepath.Join(dir, d.Name())
					if isLink(path) {
						return nil, fmt.Errorf("symlinked config fragments are not supported: %s", path)
					}
					if d.Type().IsRegular() {
						_, include, err := fragmentForHarness(path, harness)
						if err != nil {
							return nil, err
						}
						if !include {
							continue
						}
						fragments[d.Name()] = path
						rel, _ := filepath.Rel(layer, path)
						fragmentSources[d.Name()] = filepath.Join(layerSources[layer], filepath.Base(layer), rel)
					}
				}
			}
		}
	}
	parts := []string{marker, "# Global agent instructions"}
	for _, name := range sortedKeys(fragments) {
		path := fragments[name]
		b, _, e := fragmentForHarness(path, harness)
		if e != nil {
			return nil, e
		}
		if !utf8.Valid(b) {
			return nil, fmt.Errorf("invalid UTF-8 in fragment: %s", path)
		}
		text := strings.TrimSpace(string(b))
		lines := strings.Split(text, "\n")
		if len(lines) > 0 && strings.TrimSpace(lines[0]) == "---" {
			for i := 1; i < len(lines); i++ {
				if strings.TrimSpace(lines[i]) == "---" || strings.TrimSpace(lines[i]) == "..." {
					text = strings.TrimSpace(strings.Join(lines[i+1:], "\n"))
					break
				}
			}
		}
		parts = append(parts, "<!-- source: "+filepath.ToSlash(fragmentSources[name])+" -->")
		if text != "" {
			parts = append(parts, text)
		}
		p.fragments++
	}
	p.agents = strings.Join(parts, "\n\n") + "\n"
	for _, name := range sortedKeys(p.skills) {
		if skillBases[name] == "" && len(skillFragments[name]) == 0 {
			delete(p.skills, name)
			continue
		}
		dst := filepath.Join("skills", name)
		p.dirs[dst] = true
		fragmentMap := skillFragments[name]
		{
			paths := []string{}
			if skillBases[name] != "" {
				paths = append(paths, skillBases[name])
			}
			for _, key := range sortedKeys(fragmentMap) {
				paths = append(paths, fragmentMap[key])
			}
			var source string
			var metadata *yaml.Node
			var bodies []string
			var singleOriginal []byte
			for _, path := range paths {
				var b []byte
				var err error
				if path == skillBases[name] {
					b, err = os.ReadFile(path)
				} else {
					b, _, err = fragmentForHarness(path, harness)
				}
				if err != nil {
					return nil, err
				}
				if !utf8.Valid(b) {
					return nil, fmt.Errorf("invalid UTF-8 in skill fragment: %s", path)
				}
				if len(paths) == 1 {
					singleOriginal = b
				}
				text := strings.TrimSpace(string(b))
				if text == "" {
					continue
				}
				lines := strings.SplitAfter(text, "\n")
				if strings.TrimSpace(lines[0]) == "---" {
					end := 1
					for end < len(lines) && strings.TrimSpace(lines[end]) != "---" && strings.TrimSpace(lines[end]) != "..." {
						end++
					}
					if end == len(lines) {
						return nil, fmt.Errorf("skill has unterminated YAML frontmatter: %s", path)
					}
					var document yaml.Node
					raw := strings.Join(lines[1:end], "")
					if err := yaml.Unmarshal([]byte(raw), &document); err != nil {
						return nil, fmt.Errorf("invalid skill YAML frontmatter: %s: %w", path, err)
					}
					if len(document.Content) == 0 || document.Content[0].Kind != yaml.MappingNode {
						return nil, fmt.Errorf("skill YAML frontmatter must be a mapping: %s", path)
					}
					incoming := document.Content[0]
					if err := validateSkillMetadataNode(incoming); err != nil {
						return nil, fmt.Errorf("invalid skill YAML frontmatter: %s: %w", path, err)
					}
					if metadata == nil {
						metadata = incoming
						source = path
					} else {
						mergeSkillMetadata(metadata, incoming)
					}
					text = strings.TrimSpace(strings.Join(lines[end+1:], ""))
				}
				if text != "" {
					bodies = append(bodies, text)
				}
			}
			if metadata == nil {
				return nil, fmt.Errorf("skill has no YAML frontmatter: %s", p.skills[name])
			}
			metadataBytes, err := yaml.Marshal(metadata)
			if err != nil {
				return nil, fmt.Errorf("cannot render skill YAML frontmatter: %s: %w", source, err)
			}
			frontmatter := "---\n" + strings.TrimRight(string(metadataBytes), "\n") + "\n---"
			assembled := []byte(strings.Join(append([]string{frontmatter}, bodies...), "\n\n") + "\n")
			if singleOriginal != nil {
				assembled = singleOriginal
			}
			if _, err := skillNameFromBytes(assembled, source, name); err != nil {
				return nil, err
			}
			var check map[string]any
			if err := metadata.Decode(&check); err != nil {
				return nil, fmt.Errorf("invalid skill YAML frontmatter: %s: %w", source, err)
			}
			if nameValue, ok := check["name"].(string); !ok || nameValue != name {
				return nil, fmt.Errorf("invalid skill name in frontmatter: %s", source)
			}
			if desc, ok := check["description"].(string); !ok || strings.TrimSpace(desc) == "" {
				return nil, fmt.Errorf("invalid skill description in frontmatter: %s", source)
			}
			target := filepath.Join(dst, "SKILL.md")
			p.files[target] = source
			p.contents[target] = assembled
		}
		for _, src := range skillAssetLayers[name] {
			e := walk(src, func(path string, d os.DirEntry) error {
				if d.Type()&os.ModeSymlink != 0 {
					return fmt.Errorf("symlinked skill contents are not supported: %s", path)
				}
				rel, _ := filepath.Rel(src, path)
				if rel == "SKILL.md.d" && d.IsDir() {
					return filepath.SkipDir
				}
				if rel == "SKILL.md" {
					return nil
				}
				target := filepath.Join(dst, rel)
				if d.IsDir() {
					if p.files[target] != "" {
						return fmt.Errorf("skill asset conflicts with a file: %s", path)
					}
					p.dirs[target] = true
					return nil
				}
				if !d.Type().IsRegular() {
					return fmt.Errorf("skill contains a non-regular file: %s", path)
				}
				if p.dirs[target] {
					return fmt.Errorf("skill asset conflicts with a directory: %s", path)
				}
				p.files[target] = path
				return nil
			})
			if e != nil {
				return nil, e
			}
		}
	}
	for rel, src := range p.resources {
		dst := filepath.Join("resources", rel)
		p.files[dst] = src
		for _, parent := range parents(dst) {
			p.dirs[parent] = true
		}
	}
	return p, nil
}
func sortedKeys[V any](m map[string]V) []string {
	r := make([]string, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	sort.Strings(r)
	return r
}
