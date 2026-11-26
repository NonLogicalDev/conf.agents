package engine

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

type testGenerationMode struct {
	Enabled bool `json:"enabled"`
}

func testGenerationRegistry() Registry {
	handler := Handler{
		Decode: func(data json.RawMessage) (any, error) {
			var mode testGenerationMode
			err := json.Unmarshal(data, &mode)
			return mode, err
		},
		Restore: func(path string, _ Snapshot) (*RenderPlan, error) {
			return RestorePlan(path, nil, 0, 0, "guide.md")
		},
	}
	return Registry{"alpha": handler, "beta": handler}
}

func renderTestGeneration(t *testing.T, registry Registry, root string, plans map[string]*RenderPlan, snapshot Snapshot) string {
	t.Helper()
	metadata, err := SnapshotBytes(snapshot)
	if err != nil {
		t.Fatal(err)
	}
	gen, err := registry.PerformRenderPlans(root, plans, metadata)
	if err != nil {
		t.Fatal(err)
	}
	return gen
}

func testRenderRoot(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	t.Cleanup(func() {
		if err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() {
				return os.Chmod(path, 0700)
			}
			return nil
		}); err != nil {
			t.Error(err)
		}
	})
	return root
}

func TestGenerationPersistsModesAndReusesMatchingSnapshot(t *testing.T) {
	registry := testGenerationRegistry()
	root := testRenderRoot(t)
	plan := &RenderPlan{InstructionsName: "guide.md", Agents: "Instructions\n", Dirs: map[string]bool{"skills": true}}
	plans := map[string]*RenderPlan{"alpha": plan, "beta": plan}
	snapshot := Snapshot{"alpha": testGenerationMode{Enabled: true}, "beta": testGenerationMode{}}

	first := renderTestGeneration(t, registry, root, plans, snapshot)
	saved, err := registry.ReadSnapshot(first)
	if err != nil || !reflect.DeepEqual(saved, snapshot) {
		t.Fatalf("saved snapshot = %v, %v, want %v", saved, err, snapshot)
	}
	metadata, err := os.Stat(filepath.Join(first, SnapshotFile))
	if err != nil {
		t.Fatal(err)
	}
	if metadata.Mode().Perm()&0222 != 0 {
		t.Fatalf("writable generation metadata: %v", metadata.Mode())
	}
	if reused := renderTestGeneration(t, registry, root, plans, snapshot); reused != first {
		t.Fatalf("identical render created %s, want %s", reused, first)
	}

	snapshot["beta"] = testGenerationMode{Enabled: true}
	if second := renderTestGeneration(t, registry, root, plans, snapshot); second == first || !AliasMatches(root, second) {
		t.Fatalf("changed mode did not publish a new latest generation: %s", second)
	}
	delete(snapshot, "beta")
	metadataBytes, err := SnapshotBytes(snapshot)
	if err != nil {
		t.Fatal(err)
	}
	if RenderedPlansMatch(first, map[string]*RenderPlan{"alpha": plan}, metadataBytes) {
		t.Fatal("generation with a removed harness still matched")
	}
}

func TestDecodeSnapshotRejectsUnknownAndDuplicateHarnesses(t *testing.T) {
	registry := testGenerationRegistry()
	for _, data := range []string{
		`{"alpha":{"enabled":true},"gamma":{}}`,
		`{"alpha":{"enabled":true},"alpha":{"enabled":true}}`,
	} {
		if _, err := DecodeSnapshot([]byte(data+"\n"), registry); err == nil {
			t.Fatalf("accepted invalid snapshot: %s", data)
		}
	}
}

func TestDiscoveryAndPruneRejectInvalidGeneration(t *testing.T) {
	for _, tc := range []struct{ kind, metadata, want string }{
		{"missing", "", "no config.json"},
		{"malformed", "{bad", "invalid generation config"},
		{"flat payload", "{\"alpha\":{\"enabled\":true}}\n", "unexpected generation entry"},
	} {
		t.Run(tc.kind, func(t *testing.T) {
			registry := testGenerationRegistry()
			root := testRenderRoot(t)
			gen := filepath.Join(root, "2025-01-01__01")
			if err := os.Mkdir(gen, 0755); err != nil {
				t.Fatal(err)
			}
			instructions := filepath.Join(gen, "guide.md")
			if err := os.WriteFile(instructions, []byte("old instructions\n"), 0444); err != nil {
				t.Fatal(err)
			}
			if tc.metadata != "" {
				if err := os.WriteFile(filepath.Join(gen, SnapshotFile), []byte(tc.metadata), 0444); err != nil {
					t.Fatal(err)
				}
			}
			if err := os.Chmod(gen, 0555); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(filepath.Base(gen), filepath.Join(root, "latest")); err != nil {
				t.Fatal(err)
			}
			if _, err := registry.DiscoverGenerations(root); err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("discover error = %v, want %q", err, tc.want)
			}
			if err := registry.Prune(root); err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("prune error = %v, want %q", err, tc.want)
			}
			if _, err := os.Stat(instructions); err != nil {
				t.Fatalf("invalid generation was removed: %v", err)
			}
		})
	}
}

func TestRemoveGenerationRestoresModeOnFailure(t *testing.T) {
	root := testRenderRoot(t)
	gen := filepath.Join(root, "generation")
	if err := os.Mkdir(gen, 0555); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(root, 0555); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(root, "permission-probe"), 0700); err == nil {
		t.Skip("current user can write a read-only directory")
	} else if !errors.Is(err, fs.ErrPermission) {
		t.Fatal(err)
	}
	if err := removeGeneration(gen); !errors.Is(err, fs.ErrPermission) {
		t.Fatalf("remove error = %v, want permission error", err)
	}
	info, err := os.Stat(gen)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0555 {
		t.Fatalf("generation mode after failed removal = %o, want 555", info.Mode().Perm())
	}
}
