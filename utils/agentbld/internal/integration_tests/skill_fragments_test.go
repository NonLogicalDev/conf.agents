package integration_tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
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
