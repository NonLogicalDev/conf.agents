package internal

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestSkillFragmentsAssembleAcrossSectionsAndHarnesses(t *testing.T) {
	f := newContractFixture(t)
	if err := os.Remove(filepath.Join(f.layer, "skills/Alpha/SKILL.md")); err != nil {
		t.Fatal(err)
	}
	f.put(f.layer, "skills/Alpha/SKILL.md.d/10-frontmatter.md", "---\nname: Alpha\ndescription: Assembled skill\n---\nmetadata body\n")
	f.put(f.layer, "skills/Alpha/SKILL.md.d/20-override.md", "old fragment\n")
	f.put(f.layer, "skills/Alpha/SKILL.md.d/30-last.md", "last fragment\n")
	f.put(f.layer, "skills/Alpha/references/help.md", "reference retained\n")
	f.put(f.layer, "skills/Alpha/scripts/check.sh", "#!/bin/sh\nexit 0\n")
	f.put(f.layer, "skills/Alpha/references/override.md", "old reference\n")
	f.put(f.layer, "skills/Agent_Guidance/SKILL.md.d/00-frontmatter.md", "---\nname: Agent_Guidance\ndescription: Guidance\n---\n")
	f.put(f.layer, "skills/Agent_Guidance/SKILL.md.d/10-body.md", "Guidance: {{%_agents_%}}\n")
	later := filepath.Join(f.root, "second-sections")
	low := filepath.Join(later, "00-low")
	high := filepath.Join(later, "10-high")
	f.put(low, "skills/Alpha/SKILL.md.d/20-override.md", "middle fragment\n")
	f.put(high, "skills/Alpha/SKILL.md.d/00-first.md", "first fragment\n")
	f.put(high, "skills/Alpha/references/override.md", "new reference\n")
	f.put(high, "skills/Alpha/SKILL.md.d/20-override.md", "winning fragment {{%_resources_%}}/guide/first.md\n")
	open, claude := filepath.Join(f.root, "opencode"), filepath.Join(f.root, "claude")
	market := filepath.Join(f.root, "market", ".agents", "plugins")
	lifecycleRun(t, f, 0, "render", "--sync", "--sections-dir", filepath.Join(f.config, "sections"), "--sections-dir", later, "--render-dir", f.render,
		"--mode=codex:{\"agents\":true,\"skills\":true,\"plugin\":\"full\"}", "--mode=opencode", "--mode=claude", "--codex-dir", f.codex, "--opencode-dir", open, "--claude-dir", claude, "--codex-marketplace-home", market)
	gen := f.latest()
	for _, root := range []string{filepath.Join(gen, "codex"), filepath.Join(gen, "opencode"), filepath.Join(gen, "claude"), f.codex, open, claude, filepath.Join(market, "skillset")} {
		skill := filepath.Join(root, "skills/Alpha/SKILL.md")
		got := f.read(skill)
		if !strings.HasPrefix(got, "---\nname: Alpha\n") {
			t.Errorf("%s lost YAML frontmatter: %q", skill, got)
		}
		first, body, winner, last := strings.Index(got, "first fragment"), strings.Index(got, "metadata body"), strings.Index(got, "winning fragment"), strings.Index(got, "last fragment")
		if first < 0 || body <= first || winner <= body || last <= winner || strings.Contains(got, "old fragment") || strings.Contains(got, "middle fragment") {
			t.Errorf("%s wrong fragment order or precedence: %s", skill, got)
		}
		if strings.Contains(got, "{{%_resources_%}}") || !strings.Contains(got, "resources/guide/first.md") {
			t.Errorf("%s resource placeholder: %s", skill, got)
		}
		if got := f.read(filepath.Join(root, "skills/Alpha/references/help.md")); got != "reference retained\n" {
			t.Errorf("reference lost: %q", got)
		}
		if got := f.read(filepath.Join(root, "skills/Alpha/references/override.md")); got != "new reference\n" {
			t.Errorf("reference precedence: %q", got)
		}
		if got := f.read(filepath.Join(root, "skills/Alpha/scripts/check.sh")); !strings.Contains(got, "exit 0") {
			t.Errorf("script lost: %q", got)
		}
		f.absent(filepath.Join(root, "skills/Alpha/SKILL.md.d"))
		guidance := f.read(filepath.Join(root, "skills/Agent_Guidance/SKILL.md"))
		if strings.Contains(guidance, "{{%_agents_%}}") || !strings.Contains(guidance, "Base ") {
			t.Errorf("guidance was not assembled: %s", guidance)
		}
	}
}

func TestSkillFragmentsLegacyWholeSkillReplacement(t *testing.T) {
	f := newContractFixture(t)
	f.put(f.layer, "skills/Alpha/SKILL.md.d/10-earlier.md", "earlier fragment\n")
	later := filepath.Join(f.config, "sections", "50-later")
	f.put(later, "skills/Alpha/SKILL.md", "---\nname: Alpha\ndescription: Replacement\n---\nreplacement base\n")
	final := filepath.Join(f.config, "sections", "90-final")
	f.put(final, "skills/Alpha/SKILL.md.d/20-addition.md", "later fragment\n")
	f.run("render", 0)
	got := f.read(filepath.Join(f.latest(), "codex/skills/Alpha/SKILL.md"))
	if !strings.Contains(got, "replacement base") || !strings.Contains(got, "later fragment") || strings.Contains(got, "earlier fragment") || strings.Contains(got, "Example skill") {
		t.Fatalf("whole skill replacement: %s", got)
	}
}

func TestSkillFragmentsRejectInvalidInputsWithoutWrites(t *testing.T) {
	for _, kind := range []string{"missing-frontmatter", "invalid-yaml", "duplicate-yaml-key", "fragment-symlink", "directory-symlink"} {
		t.Run(kind, func(t *testing.T) {
			f := newContractFixture(t)
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
			cmd := exec.Command(contractBinary, "render", "--sync", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render, "--codex-dir", f.codex)
			if out, err := cmd.CombinedOutput(); err == nil {
				t.Fatalf("accepted invalid %s: %s", kind, out)
			}
			f.absent(filepath.Join(f.render, "latest"))
			f.absent(f.codex)
		})
	}
}

func TestSkillFragmentsSingleFilePreservesSkillContent(t *testing.T) {
	f := newContractFixture(t)
	if err := os.Remove(filepath.Join(f.layer, "skills/Alpha/SKILL.md")); err != nil {
		t.Fatal(err)
	}
	want := "---\nname: Alpha\ndescription: Example\n---\n\nSkill body.\n"
	f.put(f.layer, "skills/Alpha/SKILL.md.d/10-skill.md", want)
	f.run("render", 0)
	if got := f.read(filepath.Join(f.latest(), "codex/skills/Alpha/SKILL.md")); got != want {
		t.Fatalf("single fragment changed: got %q, want %q", got, want)
	}
}

func TestSkillFragmentMetadataMergesRecursively(t *testing.T) {
	f := newContractFixture(t)
	if err := os.Remove(filepath.Join(f.layer, "skills/Alpha/SKILL.md")); err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(f.layer, "skills/Alpha/SKILL.md.d")
	f.put(dir, "00-prefix.md", "---\nmetadata:\n  retained: original\n  shared: old\n  deep:\n    left: keep\n    change: old\ntags: [one, two]\n---\nprefix body\n")
	f.put(dir, "10-skill.md", "---\nname: Alpha\nmetadata:\n  shared: middle\n  deep:\n    change: middle\ntags: [replacement]\n---\nmain body\n")
	f.put(dir, "20-suffix.md", "---\ndescription: Final description\nmetadata:\n  shared: final\n  deep:\n    right: added\n    change: final\n---\nsuffix body\n")
	f.run("render", 0)
	got := f.read(filepath.Join(f.latest(), "codex/skills/Alpha/SKILL.md"))
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
	f := newContractFixture(t)
	f.put(f.layer, "skills/Alpha/SKILL.md.d/10-old.md", "---\nmetadata:\n  obsolete: discard\n---\nold body\n")
	later := filepath.Join(f.config, "sections", "50-later")
	f.put(later, "skills/Alpha/SKILL.md", "---\nname: Alpha\ndescription: Base description\nmetadata:\n  retained: base\n  shared: base\n---\nbase body\n")
	final := filepath.Join(f.config, "sections", "90-final")
	f.put(final, "skills/Alpha/SKILL.md.d/20-addition.md", "---\ndescription: Override description\nmetadata:\n  shared: override\n  added: yes\n---\naddition body\n")
	f.run("render", 0)
	got := f.read(filepath.Join(f.latest(), "codex/skills/Alpha/SKILL.md"))
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
