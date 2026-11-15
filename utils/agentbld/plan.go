package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

const marker = "<!-- generated: agent-config; do not edit -->"
const resourcesPlaceholder = "{{%_resources_%}}"
const skillsPlaceholder = "{{%_skills_%}}"
const agentsPlaceholder = "{{%_agents_%}}"

var layerPattern = regexp.MustCompile(`^\d+-[A-Za-z0-9][A-Za-z0-9._-]*$`)
var namePattern = regexp.MustCompile(`^name\s*:\s*(?:'([^'"\s#]+)'|"([^'"\s#]+)"|([^'"\s#]+))\s*(?:#.*)?$`)
var descPattern = regexp.MustCompile(`^description\s*:\s*\S`)
var schemePattern = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9+.-]*:`)

type renderPlan struct {
	layers    []string
	skills    map[string]string
	resources map[string]string
	agents    string
	fragments int
	dirs      map[string]bool
	files     map[string]string
	contents  map[string][]byte
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
	if name != filepath.Base(filepath.Dir(path)) {
		return "", fmt.Errorf("skill name %q does not match directory %q: %s", name, filepath.Base(filepath.Dir(path)), path)
	}
	return name, nil
}
func buildPlan(config, codex string) (*renderPlan, error) {
	p := &renderPlan{skills: map[string]string{}, resources: map[string]string{}, dirs: map[string]bool{"skills": true, "resources": true, "rules": true}, files: map[string]string{}, contents: map[string][]byte{}}
	if !isDir(config) {
		return nil, fmt.Errorf("missing agent-config directory: %s", config)
	}
	root := filepath.Join(config, "sections")
	if isLink(root) {
		return nil, fmt.Errorf("symlinked source sections are not supported: %s", root)
	}
	if !isDir(root) {
		return nil, fmt.Errorf("missing agent-config sections directory: %s", root)
	}
	es, e := entries(root)
	if e != nil {
		return nil, e
	}
	for _, d := range es {
		if !layerPattern.MatchString(d.Name()) {
			continue
		}
		path := filepath.Join(root, d.Name())
		if isLink(path) {
			return nil, fmt.Errorf("symlinked source layers are not supported: %s", path)
		}
		if d.IsDir() {
			p.layers = append(p.layers, path)
		}
	}
	fragments := map[string]string{}
	rules := map[string]string{}
	for _, layer := range p.layers {
		for _, component := range []string{"skills", "resources", "rules", "AGENTS.md.d"} {
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
					if !isRegular(skill) || isLink(skill) {
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
						if !exists(skill) && empty {
							continue
						}
						return nil, fmt.Errorf("skill directory has no regular SKILL.md: %s", path)
					}
					name, e := readSkillName(skill)
					if e != nil {
						return nil, e
					}
					if prev, ok := p.skills[name]; ok {
						fmt.Printf("skill override %s: %s -> %s\n", name, prev, path)
					}
					p.skills[name] = path
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
			case "rules":
				ds, e := entries(dir)
				if e != nil {
					return nil, e
				}
				for _, d := range ds {
					path := filepath.Join(dir, d.Name())
					if d.Type()&os.ModeSymlink != 0 || !d.Type().IsRegular() || filepath.Ext(path) != ".rules" {
						return nil, fmt.Errorf("rules must be regular .rules files: %s", path)
					}
					rules[d.Name()] = path
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
						fragments[d.Name()] = path
					}
				}
			}
		}
	}
	parts := []string{marker, "# Global agent instructions"}
	for _, name := range sortedKeys(fragments) {
		path := fragments[name]
		b, e := os.ReadFile(path)
		if e != nil {
			return nil, e
		}
		if !utf8.Valid(b) {
			return nil, fmt.Errorf("invalid UTF-8 in fragment: %s", path)
		}
		rel, _ := filepath.Rel(config, path)
		parts = append(parts, "<!-- source: "+filepath.ToSlash(rel)+" -->")
		if text := strings.TrimSpace(string(b)); text != "" {
			parts = append(parts, text)
		}
		p.fragments++
	}
	if codex == "" {
		codex = os.Getenv("CODEX_HOME")
		if codex == "" {
			h, e := os.UserHomeDir()
			if e != nil {
				return nil, e
			}
			codex = filepath.Join(h, ".codex")
		}
		codex, _ = absolute(codex)
	}
	template := strings.ReplaceAll(strings.Join(parts, "\n\n")+"\n", skillsPlaceholder, filepath.Join(codex, "skills"))
	p.agents = strings.ReplaceAll(template, resourcesPlaceholder, filepath.Join(codex, "resources"))
	for name, src := range rules {
		p.files[filepath.Join("rules", name)] = src
	}
	for _, name := range sortedKeys(p.skills) {
		src := p.skills[name]
		dst := filepath.Join("skills", name)
		p.dirs[dst] = true
		e := walk(src, func(path string, d os.DirEntry) error {
			if d.Type()&os.ModeSymlink != 0 {
				return fmt.Errorf("symlinked skill contents are not supported: %s", path)
			}
			rel, _ := filepath.Rel(src, path)
			target := filepath.Join(dst, rel)
			if d.IsDir() {
				p.dirs[target] = true
				return nil
			}
			if !d.Type().IsRegular() {
				return fmt.Errorf("skill contains a non-regular file: %s", path)
			}
			p.files[target] = path
			return nil
		})
		if e != nil {
			return nil, e
		}
	}
	for rel, src := range p.resources {
		dst := filepath.Join("resources", rel)
		p.files[dst] = src
		for _, parent := range parents(dst) {
			p.dirs[parent] = true
		}
	}
	if e := bundleResources(p, template); e != nil {
		return nil, e
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
func unescape(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) && strings.ContainsRune(`!"#$%&'()*+,-./:;<=>?@[\]^_`+"`"+`{|}~`, rune(s[i+1])) {
			i++
		}
		b.WriteByte(s[i])
	}
	return percentDecode(b.String())
}

// percentDecode follows urllib.parse.unquote: invalid escapes remain literal,
// while valid escapes elsewhere still decode.
func percentDecode(s string) string {
	var b strings.Builder
	hex := func(c byte) (byte, bool) {
		switch {
		case c >= '0' && c <= '9':
			return c - '0', true
		case c >= 'a' && c <= 'f':
			return c - 'a' + 10, true
		case c >= 'A' && c <= 'F':
			return c - 'A' + 10, true
		}
		return 0, false
	}
	for i := 0; i < len(s); i++ {
		if s[i] == '%' && i+2 < len(s) {
			a, okA := hex(s[i+1])
			c, okC := hex(s[i+2])
			if okA && okC {
				b.WriteByte(a<<4 | c)
				i += 2
				continue
			}
		}
		b.WriteByte(s[i])
	}
	return strings.ToValidUTF8(b.String(), "\uFFFD")
}

var linkStart = regexp.MustCompile(`\]\(\s*|(?m:^\s*\[[^]\n]+\]:\s*)`)

func markdownLinks(s string) []string {
	var out []string
	for _, m := range linkStart.FindAllStringIndex(s, -1) {
		i := m[1]
		angled := i < len(s) && s[i] == '<'
		if angled {
			i++
		}
		depth := 0
		var b strings.Builder
		for i < len(s) {
			c := s[i]
			i++
			if c == '\\' && i < len(s) && strings.ContainsRune(`!"#$%&'()*+,-./:;<=>?@[\]^_`+"`"+`{|}~`, rune(s[i])) {
				b.WriteByte(s[i])
				i++
				continue
			}
			if c == '\n' || c == '\r' || (angled && c == '>') {
				break
			}
			if !angled {
				if unicode.IsSpace(rune(c)) || (c == ')' && depth == 0) {
					break
				}
				if c == '(' {
					depth++
				}
				if c == ')' {
					depth--
				}
			}
			b.WriteByte(c)
		}
		out = append(out, b.String())
	}
	return out
}

// bundleResources copies each skill's referenced resources and follows links
// between catalog resources, so bundled guides retain their dependencies.
func bundleResources(p *renderPlan, agents string) error {
	targets := map[string]map[string]bool{"": {}}
	for rel := range p.resources {
		targets[""][rel] = true
		targets[rel] = map[string]bool{rel: true}
		for _, parent := range parents(rel) {
			if targets[parent] == nil {
				targets[parent] = map[string]bool{}
			}
			targets[parent][rel] = true
		}
	}
	// Match the most specific catalog path first when names share prefixes.
	paths := sortedKeys(targets)
	sort.SliceStable(paths, func(i, j int) bool { return len(paths[i]) > len(paths[j]) })
	for _, name := range sortedKeys(p.skills) {
		skillRoot := filepath.Join("skills", name)
		resourceRoot := filepath.Join(skillRoot, "resources")
		type item struct{ relative, resource string }
		var pending []item
		for _, rel := range sortedKeys(p.files) {
			if within(rel, skillRoot) {
				pending = append(pending, item{rel, ""})
			}
		}
		copied := map[string]bool{}
		for len(pending) > 0 {
			it := pending[len(pending)-1]
			pending = pending[:len(pending)-1]
			src := p.files[it.relative]
			original, e := os.ReadFile(src)
			if e != nil {
				return e
			}
			guidance := it.relative == filepath.Join("skills", "Agent_Guidance", "SKILL.md")
			if bytes.IndexByte(original, 0) >= 0 && !guidance {
				continue
			}
			if !utf8.Valid(original) {
				if guidance {
					return fmt.Errorf("cannot read guidance template: %s", src)
				}
				continue
			}
			text := string(original)
			if guidance {
				if strings.Count(text, agentsPlaceholder) != 1 {
					return fmt.Errorf("Agent_Guidance must contain one %s placeholder", agentsPlaceholder)
				}
				text = strings.ReplaceAll(text, agentsPlaceholder, strings.TrimRight(agents, "\n"))
			}
			needed := map[string]bool{}
			for start := 0; ; {
				idx := strings.Index(text[start:], resourcesPlaceholder)
				if idx < 0 {
					break
				}
				end := start + idx + len(resourcesPlaceholder)
				suffix := unescape(text[end:])
				start = end
				if !strings.HasPrefix(suffix, "/") {
					for r := range p.resources {
						needed[r] = true
					}
					continue
				}
				suffix = suffix[1:]
				found := false
				for _, path := range paths {
					if !strings.HasPrefix(suffix, path) {
						continue
					}
					rest := suffix[len(path):]
					if _, ok := p.resources[path]; !ok {
						rest = strings.TrimPrefix(rest, "/")
					}
					valid := rest == ""
					if !valid {
						c := rest[0]
						valid = unicode.IsSpace(rune(c)) || strings.ContainsRune("`\"'>)]}?#", rune(c)) || (strings.ContainsRune(".,;:!", rune(c)) && (len(rest) == 1 || unicode.IsSpace(rune(rest[1]))))
					}
					if valid {
						for r := range targets[path] {
							needed[r] = true
						}
						found = true
						break
					}
				}
				if !found {
					return fmt.Errorf("Unknown resource reference in %s: %s", src, strings.Split(suffix, "\n")[0])
				}
			}
			// Only catalog-relative links extend the bundle; external links are left alone.
			if it.resource != "" {
				for _, link := range markdownLinks(text) {
					if schemePattern.MatchString(link) || strings.HasPrefix(link, "/") {
						continue
					}
					path := strings.SplitN(strings.SplitN(link, "#", 2)[0], "?", 2)[0]
					if path == "" {
						continue
					}
					decoded := percentDecode(path)
					target := filepath.Clean(filepath.Join(filepath.Dir(it.resource), decoded))
					for r := range targets[target] {
						needed[r] = true
					}
				}
			}
			if strings.Contains(text, resourcesPlaceholder) {
				if source, ok := p.files[resourceRoot]; ok {
					return fmt.Errorf("Skill resource directory conflicts with a file: %s", source)
				}
				p.dirs[resourceRoot] = true
				local, e := filepath.Rel(filepath.Dir(it.relative), resourceRoot)
				if e != nil {
					return e
				}
				text = strings.ReplaceAll(text, resourcesPlaceholder, filepath.ToSlash(local))
			}
			if text != string(original) {
				p.contents[it.relative] = []byte(text)
			}
			for _, target := range sortedKeys(needed) {
				if copied[target] {
					continue
				}
				dst := filepath.Join(resourceRoot, target)
				if p.files[dst] != "" || p.dirs[dst] {
					return fmt.Errorf("Bundled resource conflicts with skill contents: %s", dst)
				}
				for _, parent := range parents(dst) {
					if p.files[parent] != "" {
						return fmt.Errorf("Bundled resource conflicts with skill contents: %s", dst)
					}
					p.dirs[parent] = true
				}
				p.files[dst] = p.resources[target]
				copied[target] = true
				pending = append(pending, item{dst, target})
			}
		}
	}
	return nil
}
