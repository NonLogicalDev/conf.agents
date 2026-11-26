package engine

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testTarget(t *testing.T, registry Registry, root, gen, dest string, plan *RenderPlan) string {
	t.Helper()
	target, _, err := registry.MaterializeHarnessTarget(root, gen, "alpha", dest, plan)
	if err != nil {
		t.Fatal(err)
	}
	return target
}

func TestPruneProtectsActiveAndUnverifiableDestinations(t *testing.T) {
	for _, tc := range []struct{ name, want string }{
		{"active", "active alpha generation"},
		{"missing", "cannot verify recorded harness destination"},
		{"replaced", "destination identity changed"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			registry := testGenerationRegistry()
			root := testRenderRoot(t)
			dest := filepath.Join(t.TempDir(), "destination")
			plan := &RenderPlan{InstructionsName: "guide.md", Agents: "original\n", Dirs: map[string]bool{"skills": true}}
			plans := map[string]*RenderPlan{"alpha": plan}
			snapshot := Snapshot{"alpha": testGenerationMode{Enabled: true}}
			oldest := renderTestGeneration(t, registry, root, plans, snapshot)
			target := testTarget(t, registry, root, oldest, dest, plan)
			link := filepath.Join(dest, "guide.md")
			if err := os.Symlink(filepath.Join(target, "guide.md"), link); err != nil {
				t.Fatal(err)
			}

			switch tc.name {
			case "missing", "replaced":
				if err := os.Rename(dest, dest+"-offline"); err != nil {
					t.Fatal(err)
				}
			}
			if tc.name == "replaced" {
				if err := os.Mkdir(dest, 0755); err != nil {
					t.Fatal(err)
				}
			}
			for n := 0; n < 5; n++ {
				plan.Agents = fmt.Sprintf("revision %d\n", n)
				renderTestGeneration(t, registry, root, plans, snapshot)
			}

			if err := registry.Prune(root); err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("prune error = %v, want %q", err, tc.want)
			}
			if _, err := os.Stat(oldest); err != nil {
				t.Fatalf("prune removed the original generation: %v", err)
			}
			if tc.name == "replaced" {
				if err := os.Remove(dest); err != nil {
					t.Fatal(err)
				}
			}
			if tc.name != "active" {
				if err := os.Rename(dest+"-offline", dest); err != nil {
					t.Fatal(err)
				}
			}
			if body, err := os.ReadFile(link); err != nil || string(body) != "original\n" {
				t.Fatalf("linked instructions after refused prune = %q, %v", body, err)
			}

			if err := os.Remove(link); err != nil {
				t.Fatal(err)
			}
			if err := registry.Prune(root); err != nil {
				t.Fatal(err)
			}
			for _, path := range []string{oldest, target} {
				if _, err := os.Lstat(path); !errors.Is(err, fs.ErrNotExist) {
					t.Fatalf("prune left inactive path %s: %v", path, err)
				}
			}
		})
	}
}

func TestTargetPreflightRejectsReplacedDestination(t *testing.T) {
	registry := testGenerationRegistry()
	root := testRenderRoot(t)
	dest := filepath.Join(t.TempDir(), "destination")
	plan := &RenderPlan{InstructionsName: "guide.md", Agents: "original\n"}
	gen := renderTestGeneration(t, registry, root, map[string]*RenderPlan{"alpha": plan}, Snapshot{"alpha": testGenerationMode{}})
	target := testTarget(t, registry, root, gen, dest, plan)

	if err := os.Rename(dest, dest+"-offline"); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(dest, 0755); err != nil {
		t.Fatal(err)
	}
	if err := registry.PreflightHarnessTarget(root, gen, "alpha", dest, plan); err == nil || !strings.Contains(err.Error(), "changed alpha destination") {
		t.Fatalf("preflight did not reject replaced destination: %v", err)
	}
	if _, _, err := registry.MaterializeHarnessTarget(root, gen, "alpha", dest, plan); err == nil {
		t.Fatal("materialization accepted a replaced destination")
	}
	if err := os.Remove(dest); err != nil {
		t.Fatalf("replaced destination changed: %v", err)
	}
	if err := os.Rename(dest+"-offline", dest); err != nil {
		t.Fatal(err)
	}
	if reused := testTarget(t, registry, root, gen, dest, plan); reused != target {
		t.Fatalf("target after restoring destination = %s, want %s", reused, target)
	}
}
