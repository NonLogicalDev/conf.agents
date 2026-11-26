package engine

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

var schemePattern = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9+.-]*:`)

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
func BundleResources(p *RenderPlan, agents string) error {
	targets := map[string]map[string]bool{"": {}}
	for rel := range p.Resources {
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
	paths := SortedKeys(targets)
	sort.SliceStable(paths, func(i, j int) bool { return len(paths[i]) > len(paths[j]) })
	for _, name := range SortedKeys(p.Skills) {
		skillRoot := filepath.Join("skills", name)
		resourceRoot := filepath.Join(skillRoot, "resources")
		type item struct{ relative, resource string }
		var pending []item
		for _, rel := range SortedKeys(p.Files) {
			if Within(rel, skillRoot) {
				pending = append(pending, item{rel, ""})
			}
		}
		copied := map[string]bool{}
		for len(pending) > 0 {
			it := pending[len(pending)-1]
			pending = pending[:len(pending)-1]
			src := p.Files[it.relative]
			original, ok := p.Contents[it.relative]
			if !ok {
				var e error
				original, e = os.ReadFile(src)
				if e != nil {
					return e
				}
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
					for r := range p.Resources {
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
					if _, ok := p.Resources[path]; !ok {
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
				if source, ok := p.Files[resourceRoot]; ok {
					return fmt.Errorf("Skill resource directory conflicts with a file: %s", source)
				}
				p.Dirs[resourceRoot] = true
				local, e := filepath.Rel(filepath.Dir(it.relative), resourceRoot)
				if e != nil {
					return e
				}
				text = strings.ReplaceAll(text, resourcesPlaceholder, filepath.ToSlash(local))
			}
			if text != string(original) {
				p.Contents[it.relative] = []byte(text)
			}
			for _, target := range SortedKeys(needed) {
				if copied[target] {
					continue
				}
				dst := filepath.Join(resourceRoot, target)
				if p.Files[dst] != "" || p.Dirs[dst] {
					return fmt.Errorf("Bundled resource conflicts with skill contents: %s", dst)
				}
				for _, parent := range parents(dst) {
					if p.Files[parent] != "" {
						return fmt.Errorf("Bundled resource conflicts with skill contents: %s", dst)
					}
					p.Dirs[parent] = true
				}
				p.Files[dst] = p.Resources[target]
				copied[target] = true
				pending = append(pending, item{dst, target})
			}
		}
	}
	return nil
}
