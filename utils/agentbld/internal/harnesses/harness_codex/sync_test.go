package harness_codex

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResetPreservesUnrelatedResourceLink(t *testing.T) {
	root := t.TempDir()
	codex := filepath.Join(root, "codex")
	resources := filepath.Join(codex, "resources")
	if err := os.MkdirAll(resources, 0755); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(resources, "external")
	target := filepath.Join(root, "other-state", "resource")
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
	if err := Reset(codex, filepath.Join(root, "render")); err != nil {
		t.Fatal(err)
	}
	if got, err := os.Readlink(link); err != nil || got != target {
		t.Fatalf("unrelated resource link = %q, %v, want %s", got, err, target)
	}
}
