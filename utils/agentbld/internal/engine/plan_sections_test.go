package engine

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSectionsResourceFileDirectoryConflict(t *testing.T) {
	f := newSkillPlanFixture(t)
	f.put(f.layer, "resources/replaced", "old file")
	second := filepath.Join(f.root, "later")
	layer := filepath.Join(second, "00-later")
	f.put(layer, "resources/replaced/child", "new child")

	registry := Registry{"alpha": {}}
	plan, err := registry.BuildPlan([]string{f.sections, second}, true, "alpha")
	if err != nil {
		t.Fatal(err)
	}
	if _, exists := plan.Files["resources/replaced"]; exists {
		t.Fatal("old resource file was kept alongside its replacement directory")
	}
	want := filepath.Join(layer, "resources/replaced/child")
	if got := plan.Files["resources/replaced/child"]; got != want {
		t.Fatalf("replacement resource = %q, want %q", got, want)
	}
}

func TestSectionsRejectMarketplaceInput(t *testing.T) {
	f := newSkillPlanFixture(t)
	f.put(f.layer, "marketplace/.agents/plugins/marketplace.json", "{}")
	registry := Registry{"alpha": {}}
	if _, err := registry.BuildPlan([]string{f.sections}, true, "alpha"); err == nil ||
		!strings.Contains(err.Error(), filepath.Join(f.layer, "marketplace")) ||
		!strings.Contains(err.Error(), "skills/") {
		t.Fatalf("wanted source path and replacement directory in error, got %v", err)
	}
}

func TestBundleResourcesDecodesPercentEscapesBeforeInvalidPercent(t *testing.T) {
	f := newSkillPlanFixture(t)
	f.put(f.layer, "skills/Alpha/SKILL.md", "---\nname: Alpha\ndescription: Example\n---\n{{%_resources_%}}/guide/a%20b.md %\n")
	f.put(f.layer, "resources/guide/a b.md", "decoded\n")

	plan := f.plan()
	if err := BundleResources(plan, plan.Agents); err != nil {
		t.Fatal(err)
	}
	path := plan.Files["skills/Alpha/resources/guide/a b.md"]
	if data, err := os.ReadFile(path); err != nil || string(data) != "decoded\n" {
		t.Fatalf("bundled resource = %q, %v", data, err)
	}
}
