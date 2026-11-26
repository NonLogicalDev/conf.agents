package internal

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFragmentHarnessSelection(t *testing.T) {
	f := newContractFixture(t)
	f.put(f.layer, "AGENTS.md.d/40-choice.md", "---\nharnesses: [codex, opencode]\n---\nEarly instruction\n")
	f.put(f.layer, "skills/Alpha/SKILL.md.d/40-choice.md", "---\nharnesses: [codex, opencode]\nmetadata:\n  selected: early\n---\nEarly skill\n")
	f.put(f.layer, "skills/OnlyClaude/SKILL.md.d/10-skill.md", "---\nharnesses: [claude]\nname: OnlyClaude\ndescription: Claude only\n---\nClaude skill\n")
	later := filepath.Join(f.config, "sections", "90-later")
	f.put(later, "AGENTS.md.d/40-choice.md", "---\r\nharnesses: [claude]\r\nextra: ignored\r\n---\r\nClaude instruction\r\n")
	f.put(later, "skills/Alpha/SKILL.md.d/40-choice.md", "---\nharnesses: [claude]\nmetadata:\n  selected: claude\n---\nClaude addition\n")
	lifecycleRun(t, f, 0, "render", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render, `--mode=codex:{"agents":true,"skills":true,"plugin":"full"}`, "--mode=opencode", `--mode=claude:{"agents":true,"skills":true,"plugin":"skills"}`)
	for _, harness := range []string{"codex", "opencode", "claude"} {
		t.Run(harness, func(t *testing.T) {
			root := filepath.Join(f.latest(), harness)
			instruction := "AGENTS.md"
			want := "Early instruction"
			skillWant := "Early skill"
			if harness == "claude" {
				instruction = "CLAUDE.md"
				want = "Claude instruction"
				skillWant = "Claude addition"
			}
			got := f.read(filepath.Join(root, instruction))
			if !strings.Contains(got, want) || strings.Contains(got, "harnesses:") || strings.Contains(got, "extra: ignored") {
				t.Fatalf("instructions: %s", got)
			}
			skill := f.read(filepath.Join(root, "skills/Alpha/SKILL.md"))
			if !strings.Contains(skill, skillWant) || strings.Contains(skill, "harnesses:") {
				t.Fatalf("skill: %s", skill)
			}
			f.absent(filepath.Join(root, "plugin-skills"))
			if harness != "opencode" {
				pluginSkill := f.read(filepath.Join(root, "plugin-skillset/Alpha/SKILL.md"))
				if !strings.Contains(pluginSkill, skillWant) || strings.Contains(pluginSkill, "harnesses:") {
					t.Fatalf("plugin skill: %s", pluginSkill)
				}
			}
			if harness == "claude" {
				f.read(filepath.Join(root, "skills/OnlyClaude/SKILL.md"))
			} else {
				f.absent(filepath.Join(root, "skills/OnlyClaude"))
			}
		})
	}
}

func TestFragmentHarnessValidation(t *testing.T) {
	for _, header := range []string{"harnesses: codex", "harnesses: []", "harnesses: [typo]", "harnesses: [codex, codex]", "harnesses: [null]", "harnesses: [codex]\nharnesses: [claude]"} {
		t.Run(header, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "fragment.md")
			if err := os.WriteFile(path, []byte("---\n"+header+"\n---\nbody\n"), 0600); err != nil {
				t.Fatal(err)
			}
			if _, _, err := fragmentForHarness(path, "codex"); err == nil {
				t.Fatalf("accepted %q", header)
			}
		})
	}
}
