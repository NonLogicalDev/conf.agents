package internal

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
	for _, kind := range []string{"missing", "file", "symlink", "mixed"} {
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
			if kind == "mixed" {
				args = append(args, "--config-dir", f.config)
			} else {
				args = append(args, "--sections-dir", bad)
			}
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

func TestSectionsResetDoesNotInferLegacyRootFromArbitraryName(t *testing.T) {
	f := newContractFixture(t)
	sections := filepath.Join(f.root, "custom")
	if err := os.MkdirAll(sections, 0755); err != nil {
		t.Fatal(err)
	}
	resources := filepath.Join(f.codex, "resources")
	if err := os.MkdirAll(resources, 0755); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(resources, "external")
	target := filepath.Join(filepath.Dir(f.root), "agent-config-rendered", "anything")
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
	out, err := exec.Command(contractBinary, "reset", "--sections-dir", sections, "--render-dir", f.render, "--codex-dir", f.codex).CombinedOutput()
	if err != nil {
		t.Fatalf("reset: %v\n%s", err, out)
	}
	if _, err := os.Lstat(link); err != nil {
		t.Fatalf("unrelated link removed: %v", err)
	}
}

func TestSectionsResourceFileDirectoryConflict(t *testing.T) {
	f := newContractFixture(t)
	f.put(f.layer, "resources/replaced", "old file")
	second := filepath.Join(f.root, "later")
	f.put(filepath.Join(second, "00-later"), "resources/replaced/child", "new child")
	out, err := exec.Command(contractBinary, "render", "--sections-dir", filepath.Join(f.config, "sections"), "--sections-dir", second, "--render-dir", f.render, "--codex-dir", f.codex).CombinedOutput()
	if err != nil {
		t.Fatalf("render: %v\n%s", err, out)
	}
	if got := f.read(filepath.Join(f.latest(), "codex", "resources/replaced/child")); got != "new child" {
		t.Fatalf("resource = %q", got)
	}
}

func TestSectionsRejectMarketplaceInput(t *testing.T) {
	f := newContractFixture(t)
	f.put(f.layer, "marketplace/.agents/plugins/marketplace.json", "{}")
	out := f.run("sync", 2, "--mode", `codex:{"plugin":"skills"}`)
	if !strings.Contains(out, filepath.Join(f.layer, "marketplace")) || !strings.Contains(out, "skills/") {
		t.Fatalf("expected an error with the source path and replacement directory, got %s", out)
	}
	f.absent(f.render)
	f.absent(f.codex)
	f.absent(filepath.Join(f.home, ".agents", "plugins"))
}
