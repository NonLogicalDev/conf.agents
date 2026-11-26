package engine

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

type skillPlanFixture struct {
	t                     *testing.T
	root, sections, layer string
}

func newSkillPlanFixture(t *testing.T) *skillPlanFixture {
	t.Helper()
	root := t.TempDir()
	sections := filepath.Join(root, "sections")
	f := &skillPlanFixture{t: t, root: root, sections: sections, layer: filepath.Join(sections, "00-base")}
	f.put(f.layer, "AGENTS.md.d/10-base.md", "Instructions\n")
	f.put(f.layer, "skills/Alpha/SKILL.md", "---\nname: Alpha\ndescription: Example skill\n---\nOriginal body\n")
	return f
}

func (f *skillPlanFixture) put(root, name, body string) {
	f.t.Helper()
	path := filepath.Join(root, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		f.t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0644); err != nil {
		f.t.Fatal(err)
	}
}

func (f *skillPlanFixture) plan() *RenderPlan {
	f.t.Helper()
	plan, err := (Registry{"alpha": {}}).BuildPlan([]string{f.sections}, true, "alpha")
	if err != nil {
		f.t.Fatal(err)
	}
	return plan
}
func TestSkillFragmentsWholeSkillReplacement(t *testing.T) {
	f := newSkillPlanFixture(t)
	f.put(f.layer, "skills/Alpha/SKILL.md.d/10-earlier.md", "earlier fragment\n")
	later := filepath.Join(f.sections, "50-later")
	f.put(later, "skills/Alpha/SKILL.md", "---\nname: Alpha\ndescription: Replacement\n---\nreplacement base\n")
	final := filepath.Join(f.sections, "90-final")
	f.put(final, "skills/Alpha/SKILL.md.d/20-addition.md", "later fragment\n")
	got := string(f.plan().Contents["skills/Alpha/SKILL.md"])
	if !strings.Contains(got, "replacement base") || !strings.Contains(got, "later fragment") || strings.Contains(got, "earlier fragment") || strings.Contains(got, "Example skill") {
		t.Fatalf("whole skill replacement: %s", got)
	}
}

func TestSkillFragmentsRejectInvalidInputs(t *testing.T) {
	for _, kind := range []string{"missing-frontmatter", "invalid-yaml", "duplicate-yaml-key", "fragment-symlink", "directory-symlink"} {
		t.Run(kind, func(t *testing.T) {
			f := newSkillPlanFixture(t)
			if err := os.Remove(filepath.Join(f.layer, "skills/Alpha/SKILL.md")); err != nil {
				t.Fatal(err)
			}
			dir := filepath.Join(f.layer, "skills/Alpha/SKILL.md.d")
			switch kind {
			case "missing-frontmatter":
				f.put(dir, "00-first.md", "No metadata\n")
			case "invalid-yaml":
				f.put(dir, "00-first.md", "---\nname: [bad\n---\nbody\n")
			case "duplicate-yaml-key":
				f.put(dir, "00-first.md", "---\nname: Alpha\nname: Other\ndescription: Example\n---\nbody\n")
			case "fragment-symlink":
				f.put(dir, "00-first.md", "---\nname: Alpha\ndescription: Example\n---\n")
				if err := os.Symlink(filepath.Join(f.layer, "AGENTS.md.d/10-base.md"), filepath.Join(dir, "10-link.md")); err != nil {
					t.Fatal(err)
				}
			case "directory-symlink":
				target := filepath.Join(f.root, "external-fragments")
				f.put(target, "00-first.md", "---\nname: Alpha\ndescription: Example\n---\n")
				if err := os.Symlink(target, dir); err != nil {
					t.Fatal(err)
				}
			}
			if _, err := (Registry{"alpha": {}}).BuildPlan([]string{f.sections}, true, "alpha"); err == nil {
				t.Fatalf("accepted invalid %s", kind)
			}
		})
	}
}

func TestSkillFragmentsSingleFilePreservesSkillContent(t *testing.T) {
	f := newSkillPlanFixture(t)
	if err := os.Remove(filepath.Join(f.layer, "skills/Alpha/SKILL.md")); err != nil {
		t.Fatal(err)
	}
	want := "---\nname: Alpha\ndescription: Example\n---\n\nSkill body.\n"
	f.put(f.layer, "skills/Alpha/SKILL.md.d/10-skill.md", want)
	if got := string(f.plan().Contents["skills/Alpha/SKILL.md"]); got != want {
		t.Fatalf("single fragment changed: got %q, want %q", got, want)
	}
}

func TestSkillFragmentMetadataMergesRecursively(t *testing.T) {
	f := newSkillPlanFixture(t)
	if err := os.Remove(filepath.Join(f.layer, "skills/Alpha/SKILL.md")); err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(f.layer, "skills/Alpha/SKILL.md.d")
	f.put(dir, "00-prefix.md", "---\nmetadata:\n  retained: original\n  shared: old\n  deep:\n    left: keep\n    change: old\ntags: [one, two]\n---\nprefix body\n")
	f.put(dir, "10-skill.md", "---\nname: Alpha\nmetadata:\n  shared: middle\n  deep:\n    change: middle\ntags: [replacement]\n---\nmain body\n")
	f.put(dir, "20-suffix.md", "---\ndescription: Final description\nmetadata:\n  shared: final\n  deep:\n    right: added\n    change: final\n---\nsuffix body\n")
	got := string(f.plan().Contents["skills/Alpha/SKILL.md"])
	if !strings.HasPrefix(got, "---\n") {
		t.Fatalf("no leading frontmatter: %q", got)
	}
	rest := strings.TrimPrefix(got, "---\n")
	end := strings.Index(rest, "\n---\n")
	if end < 0 {
		t.Fatalf("no closing frontmatter: %q", got)
	}
	var metadata map[string]any
	if err := yaml.Unmarshal([]byte(rest[:end]), &metadata); err != nil {
		t.Fatalf("invalid merged YAML: %v\n%s", err, got)
	}
	if metadata["name"] != "Alpha" || metadata["description"] != "Final description" {
		t.Fatalf("merged required values: %#v", metadata)
	}
	nested, ok := metadata["metadata"].(map[string]any)
	if !ok {
		t.Fatalf("metadata is not a mapping: %#v", metadata)
	}
	deep, ok := nested["deep"].(map[string]any)
	if !ok {
		t.Fatalf("deep is not a mapping: %#v", nested)
	}
	if nested["retained"] != "original" || nested["shared"] != "final" || deep["left"] != "keep" || deep["change"] != "final" || deep["right"] != "added" {
		t.Fatalf("recursive merge: %#v", nested)
	}
	tags, ok := metadata["tags"].([]any)
	if !ok || len(tags) != 1 || tags[0] != "replacement" {
		t.Fatalf("list replacement: %#v", metadata["tags"])
	}
	body := rest[end+len("\n---\n"):]
	a, b, c := strings.Index(body, "prefix body"), strings.Index(body, "main body"), strings.Index(body, "suffix body")
	if a < 0 || b <= a || c <= b {
		t.Fatalf("body order: %s", body)
	}
}

func TestSkillFragmentMetadataAfterWholeSkillReplacement(t *testing.T) {
	f := newSkillPlanFixture(t)
	f.put(f.layer, "skills/Alpha/SKILL.md.d/10-old.md", "---\nmetadata:\n  obsolete: discard\n---\nold body\n")
	later := filepath.Join(f.sections, "50-later")
	f.put(later, "skills/Alpha/SKILL.md", "---\nname: Alpha\ndescription: Base description\nmetadata:\n  retained: base\n  shared: base\n---\nbase body\n")
	final := filepath.Join(f.sections, "90-final")
	f.put(final, "skills/Alpha/SKILL.md.d/20-addition.md", "---\ndescription: Override description\nmetadata:\n  shared: override\n  added: yes\n---\naddition body\n")
	got := string(f.plan().Contents["skills/Alpha/SKILL.md"])
	rest := strings.TrimPrefix(got, "---\n")
	end := strings.Index(rest, "\n---\n")
	if end < 0 {
		t.Fatalf("invalid frontmatter: %s", got)
	}
	var metadata map[string]any
	if err := yaml.Unmarshal([]byte(rest[:end]), &metadata); err != nil {
		t.Fatal(err)
	}
	nested, ok := metadata["metadata"].(map[string]any)
	if !ok {
		t.Fatalf("metadata is not mapping: %#v", metadata)
	}
	if metadata["description"] != "Override description" || nested["retained"] != "base" || nested["shared"] != "override" || nested["added"] != "yes" || nested["obsolete"] != nil || strings.Contains(got, "old body") || !strings.Contains(got, "base body") || !strings.Contains(got, "addition body") {
		t.Fatalf("whole skill merge: %s", got)
	}
}
