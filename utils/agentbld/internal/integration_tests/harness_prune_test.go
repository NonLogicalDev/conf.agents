package integration_tests

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPruneProtectsActiveHarnessLinksAndRemovesInactiveTargets(t *testing.T) {
	for _, tc := range []struct{ name, harness, flag, active, want string }{
		{"opencode-plugin", "opencode", "--opencode-dir", "plugins/native.js", "export default"},
		{"claude-instructions", "claude", "--claude-dir", "CLAUDE.md", "Base "},
		{"claude-skill", "claude", "--claude-dir", "skills/Alpha", "Example skill"},
		{"claude-resource", "claude", "--claude-dir", "resources/guide", "linked resource"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := modeFixture(t)
			dest := filepath.Join(f.root, "destination")
			if tc.harness == "opencode" {
				f.put(f.layer, "opencode/plugins/native.js", "export default {}\n")
			}
			render := func() {
				lifecycleRun(t, f, 0, "render", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render, "--mode", tc.harness)
			}
			render()
			first := f.latest()
			number := lifecycleOrdinal(t, first)
			sync := func(num string) { lifecycleRun(t, f, 0, "sync", num, "--render-dir", f.render, tc.flag, dest) }
			sync(number)
			active := filepath.Join(dest, tc.active)
			oldTarget := f.resolved(active)
			// Keep just this component linked so it alone must protect the generation.
			for _, pattern := range []string{"AGENTS.md", "CLAUDE.md", "skills/*", "resources/*", "plugins/*"} {
				paths, err := filepath.Glob(filepath.Join(dest, pattern))
				if err != nil {
					t.Fatal(err)
				}
				for _, p := range paths {
					if p == active {
						continue
					}
					if info, err := os.Lstat(p); err == nil && info.Mode()&os.ModeSymlink != 0 {
						if err = os.Remove(p); err != nil {
							t.Fatal(err)
						}
					}
				}
			}
			for n := 0; n < 5; n++ {
				f.put(f.layer, "AGENTS.md.d/10-base.md", fmt.Sprintf("new generation %d\n", n))
				render()
			}
			latest := f.latest()
			pruneArgs := []string{"prune", "--render-dir", f.render, "--codex-dir", filepath.Join(f.root, "unrelated"), "--sections-dir", filepath.Join(f.config, "sections")}
			lifecycleRun(t, f, 2, pruneArgs...)
			if _, err := os.Stat(first); err != nil {
				t.Fatalf("active generation removed: %v", err)
			}
			readPath := active
			if tc.active == "skills/Alpha" {
				readPath = filepath.Join(active, "SKILL.md")
			}
			if tc.active == "resources/guide" {
				readPath = filepath.Join(active, "second.md")
			}
			if got := f.read(readPath); !strings.Contains(got, tc.want) {
				t.Fatalf("active component broken: %q", got)
			}
			sync(lifecycleOrdinal(t, latest))
			lifecycleRun(t, f, 0, pruneArgs...)
			f.absent(first)
			f.absent(oldTarget)
			if _, err := os.Stat(readPath); err != nil {
				t.Fatalf("reactivated component broken: %v", err)
			}
		})
	}
}
