package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func marketplaceFixture(t *testing.T) (*renderPlan, string) {
	t.Helper()
	root := t.TempDir()
	src := filepath.Join(root, "SKILL.md")
	if err := os.WriteFile(src, []byte("---\nname: Example\ndescription: Good skill\n---\nbody\n"), 0644); err != nil {
		t.Fatal(err)
	}
	return &renderPlan{skills: map[string]string{"Example": src}, files: map[string]string{"skills/Example/SKILL.md": src}}, filepath.Join(root, "home")
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
	if _, err := populateMarketplace(p, home); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(home, "plugins/skillset-uoleg/skills/Example/SKILL.md")
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p.files["skills/Example/SKILL.md"], []byte("---\nname: Example\ndescription: [bad\n---\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := populateMarketplace(p, home); err == nil {
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
	path := filepath.Join(home, ".agents/plugins/marketplace.json")
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	original := `{"name":"custom","plugins":[{"name":"other","source":{"source":"npm","package":"example"}},{"name":"skillset-uoleg","source":"./plugins/skillset-uoleg","policy":{"installation":"NOT_AVAILABLE","authentication":"MANUAL"}}]}`
	if err := os.WriteFile(path, []byte(original), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := populateMarketplace(p, home); err != nil {
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
