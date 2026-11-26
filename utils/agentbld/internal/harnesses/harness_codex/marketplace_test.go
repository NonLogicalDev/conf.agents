package harness_codex

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/NonLogicalDev/conf.agents/utils/agentbld/internal/engine"
)

func marketplaceFixture(t *testing.T) (*engine.RenderPlan, string) {
	t.Helper()
	root := t.TempDir()
	src := filepath.Join(root, "SKILL.md")
	if err := os.WriteFile(src, []byte("---\nname: Example\ndescription: Good skill\n---\nbody\n"), 0644); err != nil {
		t.Fatal(err)
	}
	return &engine.RenderPlan{Skills: map[string]string{"Example": src}, Files: map[string]string{"skills/Example/SKILL.md": src}}, filepath.Join(root, "home/.agents/plugins")
}

func TestMarketplaceRejectsCatalogPathTraversal(t *testing.T) {
	path := filepath.Join(t.TempDir(), "marketplace.json")
	data := `{"name":"x","plugins":[{"name":"a","source":{"source":"local","path":"./x/../y"}}]}`
	if err := os.WriteFile(path, []byte(data), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := loadMarketplace(path); err == nil {
		t.Fatal("wanted path rejection")
	}
}

func TestMarketplaceValidationPreservesInstalledPackage(t *testing.T) {
	p, home := marketplaceFixture(t)
	if _, err := populateMarketplace(p, home, "skillset"); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(home, "skillset/skills/Example/SKILL.md")
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p.Files["skills/Example/SKILL.md"], []byte("---\nname: Example\ndescription: [bad\n---\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := populateMarketplace(p, home, "skillset"); err == nil {
		t.Fatal("wanted invalid YAML rejection")
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != string(before) {
		t.Fatal("validation changed the installed package")
	}
}

func TestMarketplacePreservesExistingCatalogEntriesAndPolicy(t *testing.T) {
	p, home := marketplaceFixture(t)
	path := filepath.Join(home, "marketplace.json")
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	original := `{"name":"custom","plugins":[{"name":"other","source":{"source":"npm","package":"example"}},{"name":"skillset","source":"./.agents/plugins/skillset","policy":{"installation":"NOT_AVAILABLE","authentication":"MANUAL"}}]}`
	if err := os.WriteFile(path, []byte(original), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := populateMarketplace(p, home, "skillset"); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var before, after any
	if err := json.Unmarshal([]byte(original), &before); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &after); err != nil {
		t.Fatal(err)
	}
	b, _ := json.Marshal(before)
	a, _ := json.Marshal(after)
	if string(a) != string(b) {
		t.Fatalf("catalog changed: %s", data)
	}
}

func TestMarketplaceRejectsSymlinkDestination(t *testing.T) {
	plan, home := marketplaceFixture(t)
	outside := t.TempDir()
	if err := os.MkdirAll(filepath.Dir(home), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, home); err != nil {
		t.Fatal(err)
	}
	if _, err := populateMarketplace(plan, home, "skillset"); err == nil {
		t.Fatal("marketplace accepted a symlinked destination")
	}
	if entries, err := os.ReadDir(outside); err != nil || len(entries) != 0 {
		t.Fatalf("symlinked destination changed: %v, %v", entries, err)
	}
}
