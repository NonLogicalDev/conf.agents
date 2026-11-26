package integration_tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestClaudePluginGenerationUsesNativeSkillReferences(t *testing.T) {
	for _, tc := range []struct {
		name, mode string
		guidance   bool
	}{
		{"full", `claude:{"plugin":"full","plugin_name":"my-skills","agents":true,"skills":true}`, true},
		{"skills", `claude:{"plugin":"skills","plugin_name":"my-skills","agents":true,"skills":true}`, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := modeFixture(t)
			sourceLink := "[$@:Tasker_Stream](skill://@/Tasker_Stream)"
			f.put(f.layer, "AGENTS.md.d/10-base.md", "Instructions "+sourceLink+"\n")
			f.put(f.layer, "skills/Alpha/SKILL.md", "---\nname: Alpha\ndescription: Example\n---\n"+sourceLink+"\n")
			f.put(f.layer, "skills/Alpha/notes.md", sourceLink+"\n")
			f.put(f.layer, "skills/Agent_Guidance/SKILL.md", "---\nname: Agent_Guidance\ndescription: Guidance\n---\n{{%_agents_%}}\n")
			lifecycleRun(t, f, 0, "render", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render, "--mode", tc.mode)
			gen := f.latest()
			for path, want := range map[string]string{
				"claude/CLAUDE.md":                      "[/my-skills:Tasker_Stream](skill://my-skills/Tasker_Stream)",
				"claude/skills/Alpha/SKILL.md":          "[/Tasker_Stream](skill://local/Tasker_Stream)",
				"claude/skills/Alpha/notes.md":          "[/Tasker_Stream](skill://local/Tasker_Stream)",
				"claude/plugin-skillset/Alpha/SKILL.md": "[/my-skills:Tasker_Stream](skill://my-skills/Tasker_Stream)",
				"claude/plugin-skillset/Alpha/notes.md": "[/my-skills:Tasker_Stream](skill://my-skills/Tasker_Stream)",
			} {
				got := f.read(filepath.Join(gen, path))
				if !strings.Contains(got, want) || strings.Contains(got, sourceLink) {
					t.Errorf("%s: %s", path, got)
				}
			}
			guidance := filepath.Join(gen, "claude/plugin-skillset/Agent_Guidance/SKILL.md")
			if tc.guidance {
				if got := f.read(guidance); !strings.Contains(got, "[/my-skills:Tasker_Stream](skill://my-skills/Tasker_Stream)") {
					t.Errorf("plugin guidance: %s", got)
				}
			} else {
				f.absent(filepath.Dir(guidance))
			}
			assertNoPortableSkillRefs(t, gen)
			f.absent(f.codex)
		})
	}
}

func TestClaudePluginValidationBeforeWrites(t *testing.T) {
	for _, mode := range []string{`claude:{"plugin":"bogus"}`, `claude:{"plugin_name":"my-skills"}`, `claude:{"plugin":"full","plugin_name":"../escape"}`, `claude:{"agents":false,"skills":false}`} {
		t.Run(mode, func(t *testing.T) {
			f := modeFixture(t)
			dest := filepath.Join(f.root, "claude")
			lifecycleRun(t, f, 2, "render", "--sync", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render, "--claude-dir", dest, "--mode", mode)
			f.absent(filepath.Join(f.render, "latest"))
			if _, err := os.Lstat(dest); !os.IsNotExist(err) {
				t.Fatalf("invalid mode changed destination: %v", err)
			}
		})
	}
}

func TestClaudePluginHistoricalSyncToIndependentMarketplaces(t *testing.T) {
	f := modeFixture(t)
	sourceLink := "[$@:Tasker_Stream](skill://@/Tasker_Stream)"
	f.put(f.layer, "AGENTS.md.d/10-base.md", "instructions "+sourceLink+"\n")
	f.put(f.layer, "skills/Alpha/SKILL.md", "---\nname: Alpha\ndescription: Example\n---\nversion one "+sourceLink+"\n")
	mode := `claude:{"plugin":"full","plugin_name":"my-skills","agents":true,"skills":true}`
	lifecycleRun(t, f, 0, "render", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render, "--mode", mode)
	generation := f.latest()
	ordinal := lifecycleOrdinal(t, generation)
	original := f.read(filepath.Join(generation, "claude/plugin-skillset/Alpha/SKILL.md"))
	lifecycleRun(t, f, 0, "check", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render, "--mode", mode)
	f.put(f.layer, "skills/Alpha/SKILL.md", "---\nname: Alpha\ndescription: Example\n---\nversion two\n")
	if err := os.RemoveAll(f.config); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"a", "b"} {
		root := filepath.Join(f.root, name)
		dest := filepath.Join(root, "claude")
		market := filepath.Join(root, "marketplace")
		settings := f.put(dest, "settings.json", "keep existing settings\n")
		lifecycleRun(t, f, 0, "sync", ordinal, "--render-dir", f.render, "--claude-dir", dest, "--claude-marketplace-dir", market)
		manifest := f.json(filepath.Join(market, "plugins/my-skills/.claude-plugin/plugin.json"))
		if manifest["name"] != "my-skills" {
			t.Errorf("wrong plugin manifest: %#v", manifest)
		}
		_ = f.json(filepath.Join(market, ".claude-plugin/marketplace.json"))
		for path, want := range map[string]string{
			filepath.Join(dest, "CLAUDE.md"):                                 "[/my-skills:Tasker_Stream](skill://my-skills/Tasker_Stream)",
			filepath.Join(dest, "skills/Alpha/SKILL.md"):                     "[/Tasker_Stream](skill://local/Tasker_Stream)",
			filepath.Join(market, "plugins/my-skills/skills/Alpha/SKILL.md"): "[/my-skills:Tasker_Stream](skill://my-skills/Tasker_Stream)",
		} {
			got := f.read(path)
			if !strings.Contains(got, want) || strings.Contains(got, "version two") {
				t.Errorf("%s: %s", path, got)
			}
		}
		if got := f.read(settings); got != "keep existing settings\n" {
			t.Errorf("Claude settings changed: %q", got)
		}
	}
	if got := f.read(filepath.Join(generation, "claude/plugin-skillset/Alpha/SKILL.md")); got != original {
		t.Fatal("sync changed immutable generation")
	}
}

func TestClaudePluginOnlyPreservesClaudeDestination(t *testing.T) {
	f := modeFixture(t)
	dest := filepath.Join(f.root, "claude")
	instructions := f.put(dest, "CLAUDE.md", "private instructions\n")
	skill := f.put(dest, "skills/Alpha/SKILL.md", "private skill\n")
	settings := f.put(dest, "settings.json", "private settings\n")
	market := filepath.Join(f.root, "marketplace")
	lifecycleRun(t, f, 0, "render", "--sync", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render, "--claude-dir", dest, "--claude-marketplace-dir", market, "--mode", `claude:{"plugin":"full"}`)
	for path, want := range map[string]string{instructions: "private instructions\n", skill: "private skill\n", settings: "private settings\n"} {
		if got := f.read(path); got != want {
			t.Errorf("%s changed: %q", path, got)
		}
	}
	_ = f.json(filepath.Join(market, "plugins/skillset/.claude-plugin/plugin.json"))
	f.absent(f.codex)
}

func TestClaudePluginCollisionPreflightProtectsCodex(t *testing.T) {
	f := modeFixture(t)
	dest := filepath.Join(f.root, "claude")
	market := filepath.Join(f.root, "marketplace")
	marker := f.put(market, "plugins/skillset/keep.txt", "unmanaged package\n")
	lifecycleRun(t, f, 0, "render", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render, "--mode", "codex", "--mode", `claude:{"plugin":"full"}`)
	lifecycleRun(t, f, 2, "sync", lifecycleOrdinal(t, f.latest()), "--render-dir", f.render, "--claude-dir", dest, "--claude-marketplace-dir", market, "--codex-dir", f.codex)
	if got := f.read(marker); got != "unmanaged package\n" {
		t.Fatalf("unmanaged plugin changed: %q", got)
	}
	f.absent(filepath.Join(f.codex, "AGENTS.md"))
	f.absent(filepath.Join(dest, "CLAUDE.md"))
	f.absent(filepath.Join(market, ".claude-plugin/marketplace.json"))
}

// Pre-plugin snapshots have no plugin fields and must still activate as standalone Claude.
func TestClaudeRejectsPrePluginSnapshot(t *testing.T) {
	f := modeFixture(t)
	lifecycleRun(t, f, 0, "render", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render, "--mode", "claude")
	gen := f.latest()
	path := filepath.Join(gen, "config.json")
	config := f.read(path)
	old := `,"plugin":"none","plugin_name":"skillset"`
	if !strings.Contains(config, old) {
		t.Fatalf("unexpected snapshot format: %s", config)
	}
	config = strings.Replace(config, old, "", 1)
	if err := os.Chmod(path, 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(config), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(path, 0444); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(f.config); err != nil {
		t.Fatal(err)
	}
	dest := filepath.Join(f.root, "legacy-claude")
	market := filepath.Join(f.root, "legacy-marketplace")
	lifecycleRun(t, f, 2, "sync", lifecycleOrdinal(t, gen), "--render-dir", f.render, "--claude-dir", dest, "--claude-marketplace-dir", market)
	f.absent(dest)
	f.absent(market)
}

func TestClaudePluginTransitionPreservesUnrelatedSkills(t *testing.T) {
	f := modeFixture(t)
	dest := filepath.Join(f.root, "claude")
	market := filepath.Join(f.root, "marketplace")
	sections := filepath.Join(f.config, "sections")
	unmanaged := f.put(dest, "skills/Personal/SKILL.md", "personal skill\n")
	lifecycleRun(t, f, 0, "render", "--sync", "--sections-dir", sections, "--render-dir", f.render, "--claude-dir", dest, "--claude-marketplace-dir", market, "--mode", "claude")
	f.symlink(filepath.Join(dest, "skills/Alpha"))
	lifecycleRun(t, f, 0, "render", "--sync", "--sections-dir", sections, "--render-dir", f.render, "--claude-dir", dest, "--claude-marketplace-dir", market, "--mode", `claude:{"plugin":"full","agents":true}`)
	f.absent(filepath.Join(dest, "skills/Alpha"))
	if got := f.read(unmanaged); got != "personal skill\n" {
		t.Fatalf("unrelated skill changed: %q", got)
	}
	_ = f.read(filepath.Join(market, "plugins/skillset/skills/Alpha/SKILL.md"))
}

func TestClaudeSkillsPluginRejectsEmptyPackage(t *testing.T) {
	f := modeFixture(t)
	if err := os.RemoveAll(filepath.Join(f.layer, "skills/Alpha")); err != nil {
		t.Fatal(err)
	}
	dest := filepath.Join(f.root, "claude")
	market := filepath.Join(f.root, "market")
	lifecycleRun(t, f, 2, "render", "--sync", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render, "--claude-dir", dest, "--claude-marketplace-dir", market, "--mode", `claude:{"plugin":"skills"}`)
	f.absent(filepath.Join(market, "plugins/skillset"))
	f.absent(filepath.Join(dest, "CLAUDE.md"))
}
