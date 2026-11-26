package integration_tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestSectionsDirectoryOrderAndPrecedence(t *testing.T) {
	f := newContractFixture(t)
	first := filepath.Join(f.config, "sections")
	second := filepath.Join(f.root, "other", "sections")
	low := filepath.Join(second, "01-lower")
	high := filepath.Join(second, "02-higher")
	for _, item := range []struct{ layer, value string }{{low, "lower"}, {high, "winning"}} {
		f.put(item.layer, "AGENTS.md.d/10-base.md", item.value+" fragment\n")
		f.put(item.layer, "skills/Alpha/SKILL.md", "---\nname: Alpha\ndescription: Example\n---\n"+item.value+" skill\n")
		f.put(item.layer, "resources/guide/first.md", item.value+" resource\n")
		f.put(item.layer, "rules/default.rules", item.value+" rule\n")
	}
	// The later directory wins even when the first directory's layer name sorts later.
	laterInFirst := filepath.Join(first, "99-local")
	f.put(laterInFirst, "AGENTS.md.d/10-base.md", "first source fragment\n")
	args := []string{"--sections-dir", first, "render", "--sections-dir=" + second, "--render-dir", f.render, "--codex-dir", f.codex}
	if out, err := exec.Command(contractBinary, args...).CombinedOutput(); err != nil {
		t.Fatalf("render: %v\n%s", err, out)
	}
	gen := f.latest()
	agents := f.read(filepath.Join(gen, "codex", "AGENTS.md"))
	if !strings.Contains(agents, "winning fragment") || strings.Contains(agents, "first source fragment") || !strings.Contains(agents, "<!-- source: sections[2]/02-higher/AGENTS.md.d/10-base.md -->") || strings.Contains(agents, f.root) {
		t.Fatalf("fragment precedence/provenance: %s", agents)
	}
	for _, c := range []struct{ path, want string }{{"skills/Alpha/SKILL.md", "winning skill"}, {"resources/guide/first.md", "winning resource"}, {"rules/default.rules", "winning rule"}} {
		if got := f.read(filepath.Join(gen, "codex", c.path)); !strings.Contains(got, c.want) {
			t.Errorf("%s: got %q", c.path, got)
		}
	}
}

func TestSectionsRejectInvalidRootsBeforeChanges(t *testing.T) {
	for _, kind := range []string{"missing", "file", "symlink"} {
		t.Run(kind, func(t *testing.T) {
			f := newContractFixture(t)
			bad := filepath.Join(f.root, "bad")
			switch kind {
			case "file":
				if err := os.WriteFile(bad, []byte("file"), 0600); err != nil {
					t.Fatal(err)
				}
			case "symlink":
				if err := os.Symlink(filepath.Join(f.config, "sections"), bad); err != nil {
					t.Fatal(err)
				}
				bad += string(filepath.Separator)
			}
			args := []string{"render", "--sync", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render, "--codex-dir", f.codex}
			args = append(args, "--sections-dir", bad)
			out, err := exec.Command(contractBinary, args...).CombinedOutput()
			exit, ok := err.(*exec.ExitError)
			if !ok || exit.ExitCode() != 2 {
				t.Fatalf("want exit 2: %v\n%s", err, out)
			}
			if _, err := os.Lstat(f.render); !os.IsNotExist(err) {
				t.Fatalf("render root changed: %v", err)
			}
			if _, err := os.Lstat(f.codex); !os.IsNotExist(err) {
				t.Fatalf("Codex destination changed: %v", err)
			}
		})
	}
}
