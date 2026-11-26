package harness_claude

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/NonLogicalDev/conf.agents/utils/agentbld/internal/engine"
)

func marketplaceFixture(t *testing.T) (*engine.RenderPlan, string) {
	t.Helper()
	root := t.TempDir()
	source := filepath.Join(root, "SKILL.md")
	if err := os.WriteFile(source, []byte("---\nname: Alpha\ndescription: Example\n---\nbody\n"), 0644); err != nil {
		t.Fatal(err)
	}
	plan := &engine.RenderPlan{
		Files:  map[string]string{"skills/Alpha/SKILL.md": source},
		Skills: map[string]string{"Alpha": source},
	}
	return plan, filepath.Join(root, "market")
}

func writeMarketplaceCatalog(t *testing.T, market, contents string) string {
	t.Helper()
	path := filepath.Join(market, ".claude-plugin/marketplace.json")
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(contents), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestClaudePluginRejectsMalformedExistingCatalog(t *testing.T) {
	for _, tc := range []struct{ name, content, want string }{
		{"no-name", `{"owner":{"name":"Owner"},"plugins":[]}`, "requires a name"},
		{"no-owner", `{"name":"local","plugins":[]}`, "requires an owner name"},
		{"no-plugins", `{"name":"local","owner":{"name":"Owner"}}`, "plugins must be a list"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			plan, market := marketplaceFixture(t)
			catalog := writeMarketplaceCatalog(t, market, tc.content)
			if _, err := inspectClaudeMarketplace(plan, market, "skillset"); err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("catalog error = %v, want %q", err, tc.want)
			}
			if data, err := os.ReadFile(catalog); err != nil || string(data) != tc.content {
				t.Fatalf("catalog changed: %q, %v", data, err)
			}
			if _, err := os.Lstat(filepath.Join(market, "plugins")); !errors.Is(err, fs.ErrNotExist) {
				t.Fatalf("plugin directory created: %v", err)
			}
		})
	}
}

func TestClaudeMarketplaceRejectsSymlinkedPlugin(t *testing.T) {
	plan, market := marketplaceFixture(t)
	external := t.TempDir()
	plugins := filepath.Join(market, "plugins")
	if err := os.MkdirAll(plugins, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(external, filepath.Join(plugins, "skillset")); err != nil {
		t.Fatal(err)
	}
	if _, err := inspectClaudeMarketplace(plan, market, "skillset"); err == nil || !strings.Contains(err.Error(), "is a symlink") {
		t.Fatalf("expected symlink rejection, got %v", err)
	}
	if entries, err := os.ReadDir(external); err != nil || len(entries) != 0 {
		t.Fatalf("symlink target changed: %v, %v", entries, err)
	}
	if _, err := os.Lstat(filepath.Join(market, ".claude-plugin")); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("catalog directory created: %v", err)
	}
}
func TestClaudeMarketplacePreservesCatalogAndRejectsEditedPlugin(t *testing.T) {

	plan, market := marketplaceFixture(t)
	original := `{"name":"private","owner":{"name":"owner"},"plugins":[{"name":"unrelated","source":"./plugins/unrelated","description":"retain me"}],"metadata":{"keep":true}}`
	catalog := writeMarketplaceCatalog(t, market, original)
	install, err := inspectClaudeMarketplace(plan, market, "skillset")
	if err != nil {
		t.Fatal(err)
	}
	if err := populateClaudeMarketplace(install); err != nil {
		t.Fatal(err)
	}
	beforeCatalog, err := os.ReadFile(catalog)
	if err != nil {
		t.Fatal(err)
	}
	var parsed map[string]any
	if err := json.Unmarshal(beforeCatalog, &parsed); err != nil {
		t.Fatal(err)
	}
	if parsed["name"] != "private" {
		t.Fatalf("marketplace identity changed: %#v", parsed)
	}
	entries, ok := parsed["plugins"].([]any)
	if !ok || len(entries) != 2 {
		t.Fatalf("marketplace entries: %#v", parsed)
	}
	other, ok := entries[0].(map[string]any)
	if !ok || other["name"] != "unrelated" || other["description"] != "retain me" {
		t.Fatalf("unrelated plugin changed: %#v", entries)
	}

	skill := filepath.Join(market, "plugins/skillset/skills/Alpha/SKILL.md")
	if err := os.WriteFile(skill, []byte("edited skill\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := populateClaudeMarketplace(install); err == nil || !strings.Contains(err.Error(), "locally modified") {
		t.Fatalf("updated a locally modified plugin: %v", err)
	}
	if data, err := os.ReadFile(skill); err != nil || string(data) != "edited skill\n" {
		t.Fatalf("edited skill changed: %q, %v", data, err)
	}
	if data, err := os.ReadFile(catalog); err != nil || !bytes.Equal(data, beforeCatalog) {
		t.Fatalf("catalog changed: %q, %v", data, err)
	}
}
