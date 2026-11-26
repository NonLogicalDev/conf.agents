package integration_tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestOpenCodeGenerationSyncAcrossDestinations(t *testing.T) {
	f := modeFixture(t)
	f.put(f.layer, "opencode/plugins/native.js", "export default {}\n")
	f.put(f.layer, "AGENTS.md.d/10-base.md", "Uses {{%_skills_%}} and {{%_resources_%}}\n")
	lifecycleRun(t, f, 0, "render", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render, "--mode=opencode")
	gen := f.latest()
	if !strings.Contains(f.read(filepath.Join(gen, "config.json")), `"opencode"`) {
		t.Fatal("generation omitted OpenCode mode")
	}
	if err := os.RemoveAll(f.config); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"open-a", "open-b"} {
		dest := filepath.Join(f.root, name)
		lifecycleRun(t, f, 0, "sync", lifecycleOrdinal(t, gen), "--render-dir", f.render, "--opencode-dir", dest)
		instructions := f.read(filepath.Join(dest, "AGENTS.md"))
		if !strings.Contains(instructions, filepath.Join(dest, "skills")) || !strings.Contains(instructions, filepath.Join(dest, "resources")) {
			t.Fatalf("incorrect instructions for %s: %s", dest, instructions)
		}
		f.symlink(filepath.Join(dest, "skills/Alpha"))
		f.symlink(filepath.Join(dest, "plugins/native.js"))
		if got := f.read(filepath.Join(dest, "plugins/native.js")); got != "export default {}\n" {
			t.Fatalf("plugin %q", got)
		}
	}
}

func TestOpenCodeUnmanagedConflictDoesNotChangeLinks(t *testing.T) {
	f := modeFixture(t)
	dest := filepath.Join(f.root, "opencode")
	f.put(dest, "skills/Alpha/SKILL.md", "unmanaged")
	lifecycleRun(t, f, 2, "render", "--sync", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render, "--mode=opencode", "--opencode-dir", dest)
	if got := f.read(filepath.Join(dest, "skills/Alpha/SKILL.md")); got != "unmanaged" {
		t.Fatal("unmanaged skill changed")
	}
	f.absent(filepath.Join(dest, "AGENTS.md"))
}

func TestOpenCodeConflictPreventsOtherHarnessActivation(t *testing.T) {
	f := modeFixture(t)
	dest := filepath.Join(f.root, "opencode")
	f.put(dest, "skills/Alpha/SKILL.md", "unmanaged")
	lifecycleRun(t, f, 2, "render", "--sync", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render, "--mode=codex", "--mode=opencode", "--codex-dir", f.codex, "--opencode-dir", dest)
	f.absent(filepath.Join(f.codex, "AGENTS.md"))
	f.absent(filepath.Join(dest, "AGENTS.md"))
}

func TestOpenCodeAndCodexActivateTogether(t *testing.T) {
	f := modeFixture(t)
	dest := filepath.Join(f.root, "opencode")
	f.put(f.layer, "opencode/plugins/native.ts", "export default {}\n")
	lifecycleRun(t, f, 0, "render", "--sync", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render, "--mode=codex", "--mode=opencode", "--codex-dir", f.codex, "--opencode-dir", dest)
	f.symlink(filepath.Join(f.codex, "AGENTS.md"))
	f.symlink(filepath.Join(dest, "AGENTS.md"))
	f.symlink(filepath.Join(dest, "plugins/native.ts"))
	if f.latest() == "" {
		t.Fatal("generation was not created")
	}
}

func TestOpenCodePluginOnlyPreservesExistingInstructionsAndSkills(t *testing.T) {
	f := modeFixture(t)
	dest := filepath.Join(f.root, "opencode")
	f.put(f.layer, "opencode/plugins/native.js", "plugin\n")
	f.put(dest, "AGENTS.md", "personal instructions\n")
	f.put(dest, "skills/Alpha/SKILL.md", "personal skill\n")
	lifecycleRun(t, f, 0, "render", "--sync", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render, "--mode=opencode:{\"agents\":false,\"skills\":false}", "--opencode-dir", dest)
	if got := f.read(filepath.Join(dest, "AGENTS.md")); got != "personal instructions\n" {
		t.Fatal("instructions changed")
	}
	if got := f.read(filepath.Join(dest, "skills/Alpha/SKILL.md")); got != "personal skill\n" {
		t.Fatal("skill changed")
	}
	f.symlink(filepath.Join(dest, "plugins/native.js"))
}
