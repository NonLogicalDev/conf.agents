package integration_tests

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSkillReferencesResolveForTargetsAndPlugin(t *testing.T) {
	f := modeFixture(t)
	link := "[$@:Tasker_Stream](skill://@/Tasker_Stream)"
	f.put(f.layer, "AGENTS.md.d/10-base.md", "Use "+link+"\n")
	f.put(f.layer, "skills/Alpha/SKILL.md", "---\nname: Alpha\ndescription: Example\n---\nUse "+link+"\n")
	f.put(f.layer, "skills/Alpha/notes.md", "Use "+link+"\n")
	f.put(f.layer, "skills/Agent_Guidance/SKILL.md", "---\nname: Agent_Guidance\ndescription: Guidance\n---\n{{%_agents_%}}\n")
	open := filepath.Join(f.root, "opencode")
	claude := filepath.Join(f.root, "claude")
	market := filepath.Join(f.root, "market/.agents/plugins")
	lifecycleRun(t, f, 0, "render", "--sync", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render,
		"--mode=codex:{\"agents\":true,\"skills\":true,\"plugin\":\"full\",\"plugin_name\":\"my-skills\"}",
		"--mode=opencode", "--mode=claude", "--codex-dir", f.codex, "--opencode-dir", open, "--claude-dir", claude, "--codex-marketplace-home", market)
	generation := f.latest()
	for path, want := range map[string]string{
		"codex/AGENTS.md":                               "[$my-skills:Tasker_Stream](skill://my-skills/Tasker_Stream)",
		"codex/skills/Agent_Guidance/SKILL.md":          "[$my-skills:Tasker_Stream](skill://my-skills/Tasker_Stream)",
		"codex/plugin-skillset/Agent_Guidance/SKILL.md": "[$my-skills:Tasker_Stream](skill://my-skills/Tasker_Stream)",
		"opencode/AGENTS.md":                            "[Tasker_Stream](skill://local/Tasker_Stream)",
		"opencode/skills/Agent_Guidance/SKILL.md":       "[Tasker_Stream](skill://local/Tasker_Stream)",
		"claude/CLAUDE.md":                              "[/Tasker_Stream](skill://local/Tasker_Stream)",
		"claude/skills/Agent_Guidance/SKILL.md":         "[/Tasker_Stream](skill://local/Tasker_Stream)",
	} {
		got := f.read(filepath.Join(generation, path))
		if !strings.Contains(got, want) || strings.Contains(got, link) {
			t.Errorf("generation %s: %s", path, got)
		}
	}
	for path, want := range map[string]string{
		"codex/skills/Alpha/SKILL.md":          "[$Tasker_Stream](skill://local/Tasker_Stream)",
		"codex/plugin-skillset/Alpha/SKILL.md": "[$my-skills:Tasker_Stream](skill://my-skills/Tasker_Stream)",
		"codex/plugin-skillset/Alpha/notes.md": "[$my-skills:Tasker_Stream](skill://my-skills/Tasker_Stream)",
		"codex/skills/Alpha/notes.md":          "[$Tasker_Stream](skill://local/Tasker_Stream)",
		"opencode/skills/Alpha/SKILL.md":       "[Tasker_Stream](skill://local/Tasker_Stream)",
		"opencode/skills/Alpha/notes.md":       "[Tasker_Stream](skill://local/Tasker_Stream)",
		"claude/skills/Alpha/SKILL.md":         "[/Tasker_Stream](skill://local/Tasker_Stream)",
		"claude/skills/Alpha/notes.md":         "[/Tasker_Stream](skill://local/Tasker_Stream)",
	} {
		got := f.read(filepath.Join(generation, path))
		if !strings.Contains(got, want) || strings.Contains(got, link) {
			t.Errorf("generation skill %s: %s", path, got)
		}
	}
	assertNoPortableSkillRefs(t, generation)
	originalInstructions := f.read(filepath.Join(generation, "claude/CLAUDE.md"))
	originalSkill := f.read(filepath.Join(generation, "codex/skills/Alpha/SKILL.md"))
	for path, want := range map[string]string{
		filepath.Join(f.codex, "AGENTS.md"):                      "[$my-skills:Tasker_Stream](skill://my-skills/Tasker_Stream)",
		filepath.Join(f.codex, "skills/Alpha/SKILL.md"):          "[$Tasker_Stream](skill://local/Tasker_Stream)",
		filepath.Join(f.codex, "skills/Alpha/notes.md"):          "[$Tasker_Stream](skill://local/Tasker_Stream)",
		filepath.Join(open, "AGENTS.md"):                         "[Tasker_Stream](skill://local/Tasker_Stream)",
		filepath.Join(open, "skills/Alpha/SKILL.md"):             "[Tasker_Stream](skill://local/Tasker_Stream)",
		filepath.Join(claude, "CLAUDE.md"):                       "[/Tasker_Stream](skill://local/Tasker_Stream)",
		filepath.Join(claude, "skills/Alpha/SKILL.md"):           "[/Tasker_Stream](skill://local/Tasker_Stream)",
		filepath.Join(market, "my-skills/skills/Alpha/SKILL.md"): "[$my-skills:Tasker_Stream](skill://my-skills/Tasker_Stream)",
		filepath.Join(market, "my-skills/skills/Alpha/notes.md"): "[$my-skills:Tasker_Stream](skill://my-skills/Tasker_Stream)",
	} {
		if got := f.read(path); !strings.Contains(got, want) {
			t.Errorf("%s lacks %s: %s", path, want, got)
		}
	}
	if err := os.RemoveAll(f.config); err != nil {
		t.Fatal(err)
	}
	// Replay from the immutable generation after the source has disappeared.
	lifecycleRun(t, f, 0, "sync", lifecycleOrdinal(t, generation), "--render-dir", f.render, "--codex-dir", f.codex, "--opencode-dir", open, "--claude-dir", claude, "--codex-marketplace-home", market)
	codex2, open2, claude2 := filepath.Join(f.root, "codex2"), filepath.Join(f.root, "opencode2"), filepath.Join(f.root, "claude2")
	market2 := filepath.Join(f.root, "market2/.agents/plugins")
	lifecycleRun(t, f, 0, "sync", lifecycleOrdinal(t, generation), "--render-dir", f.render, "--codex-dir", codex2, "--opencode-dir", open2, "--claude-dir", claude2, "--codex-marketplace-home", market2)
	if got := f.read(filepath.Join(generation, "claude/CLAUDE.md")); got != originalInstructions {
		t.Errorf("sync changed generation instructions")
	}
	if got := f.read(filepath.Join(generation, "codex/skills/Alpha/SKILL.md")); got != originalSkill {
		t.Errorf("sync changed generation skill")
	}
	for path, want := range map[string]string{
		filepath.Join(codex2, "AGENTS.md"):                        "[$my-skills:Tasker_Stream](skill://my-skills/Tasker_Stream)",
		filepath.Join(codex2, "skills/Alpha/SKILL.md"):            "[$Tasker_Stream](skill://local/Tasker_Stream)",
		filepath.Join(open2, "AGENTS.md"):                         "[Tasker_Stream](skill://local/Tasker_Stream)",
		filepath.Join(claude2, "CLAUDE.md"):                       "[/Tasker_Stream](skill://local/Tasker_Stream)",
		filepath.Join(claude2, "skills/Agent_Guidance/SKILL.md"):  "[/Tasker_Stream](skill://local/Tasker_Stream)",
		filepath.Join(market2, "my-skills/skills/Alpha/SKILL.md"): "[$my-skills:Tasker_Stream](skill://my-skills/Tasker_Stream)",
	} {
		if got := f.read(path); !strings.Contains(got, want) || strings.Contains(got, link) {
			t.Errorf("historical sync %s: %s", path, got)
		}
	}
}

func TestSyncRejectsOldPluginDirectory(t *testing.T) {
	f := modeFixture(t)
	f.put(f.layer, "skills/Alpha/SKILL.md", "---\nname: Alpha\ndescription: Example\n---\n[$@:Alpha](skill://@/Alpha)\n")
	f.run("render", 0, "--mode", `codex:{"plugin":"full","plugin_name":"my-skills","skills":true}`)
	gen := f.latest()
	codexPayload := filepath.Join(gen, "codex")

	// Recreate the layout and snapshot keys written before the rename.
	if err := os.Chmod(codexPayload, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(filepath.Join(codexPayload, "plugin-skillset"), filepath.Join(codexPayload, "plugin-skills")); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(codexPayload, 0555); err != nil {
		t.Fatal(err)
	}
	config := filepath.Join(gen, "config.json")
	oldConfig := strings.ReplaceAll(f.read(config), "\"plugin-skillset/", "\"plugin-skills/")
	if err := os.Chmod(config, 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(config, []byte(oldConfig), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(config, 0444); err != nil {
		t.Fatal(err)
	}

	market := filepath.Join(f.home, ".agents", "plugins")
	lifecycleRun(t, f, 2, "sync", lifecycleOrdinal(t, gen), "--render-dir", f.render,
		"--codex-dir", f.codex, "--codex-marketplace-home", market)
	f.absent(market)
	f.absent(f.codex)
	f.absent(filepath.Join(gen, "codex", "plugin-skillset"))
}

func TestStandaloneCodexReferencesUseLocal(t *testing.T) {
	f := modeFixture(t)
	link := "[$@:Tasker_Stream](skill://@/Tasker_Stream)"
	f.put(f.layer, "AGENTS.md.d/10-base.md", link+"\n")
	f.put(f.layer, "skills/Alpha/SKILL.md", "---\nname: Alpha\ndescription: Example\n---\n"+link+"\n")
	lifecycleRun(t, f, 0, lifecycleRenderArgs(f, "--sync")...)
	generation := f.latest()
	for _, path := range []string{filepath.Join(generation, "codex/AGENTS.md"), filepath.Join(generation, "codex/skills/Alpha/SKILL.md"), filepath.Join(f.codex, "AGENTS.md"), filepath.Join(f.codex, "skills/Alpha/SKILL.md")} {
		if got := f.read(path); !strings.Contains(got, "[$Tasker_Stream](skill://local/Tasker_Stream)") {
			t.Errorf("unresolved %s: %s", path, got)
		}
	}
}

func assertNoPortableSkillRefs(t *testing.T, root string) {
	t.Helper()
	if err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Type().IsRegular() && strings.HasSuffix(entry.Name(), ".md") {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			if strings.Contains(string(data), "skill://@/") || strings.Contains(string(data), "[$@:") {
				t.Errorf("unresolved reference in %s: %s", path, data)
			}
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

func TestSkillReferencesSkillsPluginOmitsGuidanceFromGeneration(t *testing.T) {
	f := modeFixture(t)
	link := "[$@:Tasker_Stream](skill://@/Tasker_Stream)"
	f.put(f.layer, "skills/Alpha/SKILL.md", "---\nname: Alpha\ndescription: Example\n---\n"+link+"\n")
	f.put(f.layer, "skills/Agent_Guidance/SKILL.md", "---\nname: Agent_Guidance\ndescription: Guidance\n---\n{{%_agents_%}}\n")
	mode := `codex:{"agents":true,"skills":true,"plugin":"skills","plugin_name":"my-skills"}`
	f.run("sync", 0, "--mode", mode)
	gen := f.latest()
	f.absent(filepath.Join(gen, "codex/plugin-skillset/Agent_Guidance"))
	f.absent(filepath.Join(f.home, ".agents/plugins/my-skills/skills/Agent_Guidance"))
	for path, want := range map[string]string{
		filepath.Join(gen, "codex/plugin-skillset/Alpha/SKILL.md"): "[$my-skills:Tasker_Stream](skill://my-skills/Tasker_Stream)",
		filepath.Join(gen, "codex/skills/Alpha/SKILL.md"):          "[$Tasker_Stream](skill://local/Tasker_Stream)",
	} {
		if got := f.read(path); !strings.Contains(got, want) {
			t.Errorf("%s: %s", path, got)
		}
	}
	assertNoPortableSkillRefs(t, gen)
	f.run("check", 0, "--mode", mode)
	if err := os.RemoveAll(f.config); err != nil {
		t.Fatal(err)
	}
	codex2 := filepath.Join(f.root, "skills-codex2")
	market2 := filepath.Join(f.root, "skills-market2/.agents/plugins")
	lifecycleRun(t, f, 0, "sync", lifecycleOrdinal(t, gen), "--render-dir", f.render, "--codex-dir", codex2, "--codex-marketplace-home", market2)
	f.absent(filepath.Join(market2, "my-skills/skills/Agent_Guidance"))
	if got := f.read(filepath.Join(market2, "my-skills/skills/Alpha/SKILL.md")); !strings.Contains(got, "[$my-skills:Tasker_Stream](skill://my-skills/Tasker_Stream)") {
		t.Fatalf("replayed plugin reference: %s", got)
	}
	if got := f.read(filepath.Join(codex2, "skills/Alpha/SKILL.md")); !strings.Contains(got, "[$Tasker_Stream](skill://local/Tasker_Stream)") {
		t.Fatalf("replayed standalone reference: %s", got)
	}
}
