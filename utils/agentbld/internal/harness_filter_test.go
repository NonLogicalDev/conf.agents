package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSyncFiltersSavedHarnessesWithoutChangingGeneration(t *testing.T) {
	f := modeFixture(t)
	codex := filepath.Join(f.root, "codex")
	claude := filepath.Join(f.root, "claude")
	opencode := filepath.Join(f.root, "opencode")
	f.put(f.layer, "opencode/plugins/native.js", "export default {}\n")
	lifecycleRun(t, f, 0, "render", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render, "--mode", "codex", "--mode", "claude", "--mode", "opencode")
	gen := f.latest()
	config := f.read(filepath.Join(gen, "config.json"))
	number := lifecycleOrdinal(t, gen)
	args := []string{"sync", number, "--render-dir", f.render, "--codex-dir", codex, "--claude-dir", claude, "--opencode-dir", opencode}
	lifecycleRun(t, f, 0, append(append([]string{}, args...), "--mode", "codex")...)
	f.symlink(filepath.Join(codex, "AGENTS.md"))
	f.absent(claude)
	f.absent(opencode)
	old := f.read(filepath.Join(codex, "AGENTS.md"))
	lifecycleRun(t, f, 0, append(append([]string{}, args...), "--mode", "claude", "--mode", "opencode")...)
	f.symlink(filepath.Join(claude, "CLAUDE.md"))
	f.symlink(filepath.Join(opencode, "AGENTS.md"))
	f.symlink(filepath.Join(opencode, "plugins/native.js"))
	if f.read(filepath.Join(codex, "AGENTS.md")) != old {
		t.Fatal("unselected Codex destination changed")
	}
	if f.read(filepath.Join(gen, "config.json")) != config || f.latest() != gen {
		t.Fatal("sync filter changed saved generation")
	}
}
func TestSyncRejectsInvalidHarnessFiltersBeforeWrites(t *testing.T) {
	for _, tc := range []struct {
		name    string
		filters []string
	}{
		{"unknown", []string{"--mode", "unknown"}},
		{"duplicate", []string{"--mode", "codex", "--mode", "codex"}},
		{"configured", []string{"--mode", `codex:{}`}},
		{"not-in-generation", []string{"--mode", "claude"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := modeFixture(t)
			lifecycleRun(t, f, 0, "render", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render, "--mode", "codex")
			gen := f.latest()
			before := f.read(filepath.Join(gen, "config.json"))
			codex := filepath.Join(f.root, "codex")
			claude := filepath.Join(f.root, "claude")
			opencode := filepath.Join(f.root, "opencode")
			args := []string{"sync", lifecycleOrdinal(t, gen), "--render-dir", f.render, "--codex-dir", codex, "--claude-dir", claude, "--opencode-dir", opencode}
			args = append(args, tc.filters...)
			lifecycleRun(t, f, 2, args...)
			for _, p := range []string{codex, claude, opencode} {
				if _, err := os.Lstat(p); !os.IsNotExist(err) {
					t.Fatalf("destination changed %s: %v", p, err)
				}
			}
			if f.read(filepath.Join(gen, "config.json")) != before {
				t.Fatal("generation config changed")
			}
		})
	}
}

func TestSyncRejectsOverlappingHarnessDestinationsBeforeWrites(t *testing.T) {
	for _, nested := range []bool{false, true} {
		t.Run(map[bool]string{false: "same", true: "nested"}[nested], func(t *testing.T) {
			f := modeFixture(t)
			lifecycleRun(t, f, 0, "render", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render, "--mode", "codex", "--mode", "claude")
			codex := filepath.Join(f.root, "destination")
			claude := codex
			if nested {
				claude = filepath.Join(codex, "child")
			}
			lifecycleRun(t, f, 2, "sync", lifecycleOrdinal(t, f.latest()), "--render-dir", f.render, "--codex-dir", codex, "--claude-dir", claude)
			f.absent(codex)
			f.absent(filepath.Join(f.render, ".targets"))
		})
	}
}
