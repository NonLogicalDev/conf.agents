package integration_tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/NonLogicalDev/conf.agents/utils/agentbld/internal/harnesses/harness_claude"
)

func TestClaudeRenderAndHistoricalSync(t *testing.T) {
	f := modeFixture(t)
	f.put(f.layer, "AGENTS.md.d/10-base.md", "Claude one {{%_skills_%}} {{%_resources_%}}\n")
	mode := `claude:{"agents":true,"skills":true}`
	lifecycleRun(t, f, 0, "render", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render, "--mode", mode)
	gen := f.latest()
	number := lifecycleOrdinal(t, gen)
	if got := f.read(filepath.Join(gen, "claude/CLAUDE.md")); !strings.Contains(got, "Claude one") {
		t.Fatal(got)
	}
	f.absent(f.codex)
	f.absent(filepath.Join(f.root, "claude-a"))
	f.put(f.layer, "AGENTS.md.d/10-base.md", "Claude changed\n")
	if err := os.RemoveAll(f.config); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"claude-a", "claude-b"} {
		dir := filepath.Join(f.root, name)
		lifecycleRun(t, f, 0, "sync", number, "--render-dir", f.render, "--claude-dir", dir)
		text := f.read(filepath.Join(dir, "CLAUDE.md"))
		for _, part := range []string{"Claude one", filepath.Join(dir, "skills"), filepath.Join(dir, "resources")} {
			if !strings.Contains(text, part) {
				t.Fatalf("missing %q in %s", part, text)
			}
		}
		if _, err := os.Stat(filepath.Join(dir, "skills/Alpha/SKILL.md")); err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stat(filepath.Join(dir, "resources/guide/first.md")); err != nil {
			t.Fatal(err)
		}
		f.absent(filepath.Join(dir, "rules"))
	}
}
func TestClaudeSyncPreservesUnselectedAndUnmanagedFiles(t *testing.T) {
	for _, tc := range []struct {
		name, mode     string
		agents, skills bool
	}{
		{"agents-only", `claude:{"agents":true,"skills":false}`, true, false},
		{"skills-only", `claude:{"agents":false,"skills":true}`, false, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := modeFixture(t)
			dest := filepath.Join(f.root, "claude")
			f.put(dest, "settings.json", "private settings")
			f.put(dest, "rules/custom.md", "private rule")
			if !tc.agents {
				f.put(dest, "CLAUDE.md", "private instructions")
				f.put(dest, "resources/guide/local", "private resource")
			}
			if !tc.skills {
				f.put(dest, "skills/Alpha/SKILL.md", "private skill")
			}
			lifecycleRun(t, f, 0, "render", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render, "--mode", tc.mode)
			lifecycleRun(t, f, 0, "sync", lifecycleOrdinal(t, f.latest()), "--render-dir", f.render, "--claude-dir", dest)
			if f.read(filepath.Join(dest, "settings.json")) != "private settings" || f.read(filepath.Join(dest, "rules/custom.md")) != "private rule" {
				t.Fatal("unmanaged Claude configuration changed")
			}
			if tc.agents {
				f.symlink(filepath.Join(dest, "CLAUDE.md"))
				f.symlink(filepath.Join(dest, "resources/guide"))
			} else {
				if f.read(filepath.Join(dest, "CLAUDE.md")) != "private instructions" || f.read(filepath.Join(dest, "resources/guide/local")) != "private resource" {
					t.Fatal("unselected agent components changed")
				}
			}
			if tc.skills {
				f.symlink(filepath.Join(dest, "skills/Alpha"))
			} else if f.read(filepath.Join(dest, "skills/Alpha/SKILL.md")) != "private skill" {
				t.Fatal("unselected skill changed")
			}
		})
	}
}
func TestClaudeSyncRejectsUnmanagedCollisionBeforeWrites(t *testing.T) {
	f := modeFixture(t)
	dest := filepath.Join(f.root, "claude")
	f.put(dest, "CLAUDE.md", "private instructions")
	lifecycleRun(t, f, 0, "render", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render, "--mode", "claude")
	lifecycleRun(t, f, 2, "sync", lifecycleOrdinal(t, f.latest()), "--render-dir", f.render, "--claude-dir", dest)
	if f.read(filepath.Join(dest, "CLAUDE.md")) != "private instructions" {
		t.Fatal("unmanaged instructions changed")
	}
	f.absent(filepath.Join(dest, "skills"))
	f.absent(filepath.Join(dest, "resources"))
}

func TestClaudeAndCodexPreflightBeforeEitherActivation(t *testing.T) {
	for _, collision := range []string{"claude", "codex"} {
		t.Run(collision, func(t *testing.T) {
			f := modeFixture(t)
			claude := filepath.Join(f.root, "claude")
			codex := filepath.Join(f.root, "codex")
			if collision == "claude" {
				f.put(claude, "CLAUDE.md", "unmanaged Claude")
			} else {
				f.put(codex, "AGENTS.md", "unmanaged Codex")
			}
			lifecycleRun(t, f, 0, "render", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render, "--mode", "codex", "--mode", "claude")
			lifecycleRun(t, f, 2, "sync", lifecycleOrdinal(t, f.latest()), "--render-dir", f.render, "--codex-dir", codex, "--claude-dir", claude)
			if collision == "claude" {
				if f.read(filepath.Join(claude, "CLAUDE.md")) != "unmanaged Claude" {
					t.Fatal("Claude instructions changed")
				}
				f.absent(filepath.Join(codex, "AGENTS.md"))
				f.absent(filepath.Join(codex, "skills"))
			} else {
				if f.read(filepath.Join(codex, "AGENTS.md")) != "unmanaged Codex" {
					t.Fatal("Codex instructions changed")
				}
				f.absent(filepath.Join(claude, "CLAUDE.md"))
				f.absent(filepath.Join(claude, "skills"))
			}
		})
	}
}

func TestClaudePathDefaultsAndOverride(t *testing.T) {
	root := t.TempDir()
	t.Setenv("HOME", root)
	t.Setenv("CLAUDE_CONFIG_DIR", "")
	path, err := harness_claude.Destination("")
	if err != nil || path != filepath.Join(root, ".claude") {
		t.Fatalf("default %q, %v", path, err)
	}
	custom := filepath.Join(root, "custom")
	t.Setenv("CLAUDE_CONFIG_DIR", custom)
	path, err = harness_claude.Destination("")
	if err != nil || path != custom {
		t.Fatalf("environment %q, %v", path, err)
	}
	explicit := filepath.Join(root, "explicit")
	path, err = harness_claude.Destination(explicit)
	if err != nil || path != explicit {
		t.Fatalf("explicit %q, %v", path, err)
	}
}

func TestClaudeRenderCheckDoesNotActivate(t *testing.T) {
	f := modeFixture(t)
	mode := `claude:{"agents":true,"skills":true}`
	lifecycleRun(t, f, 0, "render", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render, "--mode", mode)
	out := lifecycleRun(t, f, 0, "check", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render, "--mode", mode)
	if !strings.Contains(out, "status=up-to-date") {
		t.Fatalf("check not current: %s", out)
	}
	f.absent(f.codex)
	f.absent(filepath.Join(f.root, "claude"))
}

func TestThreeHarnessPreflightLastCollisionLeavesOthersUntouched(t *testing.T) {
	f := modeFixture(t)
	claude := filepath.Join(f.root, "claude")
	codex := filepath.Join(f.root, "codex")
	openCode := filepath.Join(f.root, "opencode")
	f.put(openCode, "AGENTS.md", "unmanaged OpenCode instructions")
	lifecycleRun(t, f, 0, "render", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render, "--mode", "claude", "--mode", "codex", "--mode", "opencode")
	lifecycleRun(t, f, 2, "sync", lifecycleOrdinal(t, f.latest()), "--render-dir", f.render, "--claude-dir", claude, "--codex-dir", codex, "--opencode-dir", openCode)
	if f.read(filepath.Join(openCode, "AGENTS.md")) != "unmanaged OpenCode instructions" {
		t.Fatal("unmanaged instructions changed")
	}
	for _, p := range []string{filepath.Join(claude, "CLAUDE.md"), filepath.Join(claude, "skills"), filepath.Join(codex, "AGENTS.md"), filepath.Join(codex, "skills")} {
		f.absent(p)
	}
}

func TestAllThreeHarnessesRenderAndActivateIndependently(t *testing.T) {
	f := modeFixture(t)
	claude := filepath.Join(f.root, "claude")
	codex := filepath.Join(f.root, "codex")
	openCode := filepath.Join(f.root, "opencode")
	f.put(f.layer, "AGENTS.md.d/10-base.md", "Shared {{%_skills_%}} {{%_resources_%}}\n")
	lifecycleRun(t, f, 0, "render", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render, "--mode", "claude", "--mode", "codex", "--mode", "opencode")
	gen := f.latest()
	for _, h := range []string{"claude", "codex", "opencode"} {
		if _, err := os.Stat(filepath.Join(gen, h)); err != nil {
			t.Fatalf("missing %s payload: %v", h, err)
		}
	}
	lifecycleRun(t, f, 0, "sync", lifecycleOrdinal(t, gen), "--render-dir", f.render, "--claude-dir", claude, "--codex-dir", codex, "--opencode-dir", openCode)
	for _, tc := range []struct{ dir, instructions string }{{claude, "CLAUDE.md"}, {codex, "AGENTS.md"}, {openCode, "AGENTS.md"}} {
		content := f.read(filepath.Join(tc.dir, tc.instructions))
		for _, part := range []string{"Shared", filepath.Join(tc.dir, "skills"), filepath.Join(tc.dir, "resources")} {
			if !strings.Contains(content, part) {
				t.Fatalf("%s missing %s: %s", tc.dir, part, content)
			}
		}
		if _, err := os.Stat(filepath.Join(tc.dir, "skills/Alpha/SKILL.md")); err != nil {
			t.Fatal(err)
		}
	}
}
