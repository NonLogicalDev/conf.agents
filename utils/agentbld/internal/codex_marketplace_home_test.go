package internal

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCodexMarketplaceHomeCLI(t *testing.T) {
	for _, tc := range []struct{ name, flag string }{{"default-home", ""}, {"explicit-dir", "--codex-marketplace-home"}, {"legacy-home", "--home"}} {
		t.Run(tc.name, func(t *testing.T) {
			f := modeFixture(t)
			home := filepath.Join(f.root, "home")
			if err := os.MkdirAll(home, 0755); err != nil {
				t.Fatal(err)
			}
			dir := filepath.Join(home, ".agents/plugins")
			args := []string{"render", "--sync", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render, "--codex-dir", f.codex, "--mode", `codex:{"agents":false,"plugin":"skills"}`}
			if tc.flag != "" {
				customHome := filepath.Join(f.root, "custom")
				dir = filepath.Join(customHome, ".agents/plugins")
				value := dir
				if tc.flag == "--home" {
					value = customHome
				}
				args = append(args, tc.flag, value)
			}
			cmd := exec.Command(contractBinary, args...)
			env := []string{}
			for _, v := range os.Environ() {
				k, _, _ := strings.Cut(v, "=")
				if k != "HOME" {
					env = append(env, v)
				}
			}
			cmd.Env = append(env, "HOME="+home)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			if err := cmd.Run(); err != nil {
				t.Fatalf("sync: %v\nstdout: %s\nstderr: %s", err, &stdout, &stderr)
			}
			line, _, _ := strings.Cut(stdout.String(), "\n")
			if !strings.HasPrefix(line, "config: ") {
				t.Fatalf("unexpected first line: %s", line)
			}
			var config reportedConfig
			if err := json.Unmarshal([]byte(strings.TrimPrefix(line, "config: ")), &config); err != nil {
				t.Fatal(err)
			}
			if config.CodexMarketplaceHome != dir {
				t.Fatalf("marketplace dir %q, want %q", config.CodexMarketplaceHome, dir)
			}
			plugin := filepath.Join(dir, "skillset")
			if _, err := os.Stat(filepath.Join(plugin, "skills/Alpha/SKILL.md")); err != nil {
				t.Fatal(err)
			}
			catalog := f.json(filepath.Join(dir, "marketplace.json"))
			entries := catalog["plugins"].([]any)
			var source string
			for _, item := range entries {
				entry := item.(map[string]any)
				if entry["name"] == "skillset" {
					source = entry["source"].(map[string]any)["path"].(string)
				}
			}
			discoveryRoot := filepath.Dir(filepath.Dir(dir))
			if source == "" || filepath.Clean(filepath.Join(discoveryRoot, source)) != plugin {
				t.Fatalf("source %q does not resolve to %s from %s", source, plugin, discoveryRoot)
			}
			_ = f.latest()
			f.absent(f.codex)
			if tc.flag != "" {
				f.absent(filepath.Join(home, ".agents"))
				f.absent(filepath.Join(home, "plugins"))
			}
		})
	}
}
func TestCodexMarketplaceHomeRejectsLegacyAliasCombination(t *testing.T) {
	f := modeFixture(t)
	newDir := filepath.Join(f.root, "new/.agents/plugins")
	oldHome := filepath.Join(f.root, "old")
	args := []string{"render", "--sync", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render, "--codex-dir", f.codex, "--mode", `codex:{"agents":false,"plugin":"full"}`, "--codex-marketplace-home", newDir, "--home", oldHome}
	out, err := exec.Command(contractBinary, args...).CombinedOutput()
	exit, ok := err.(*exec.ExitError)
	if !ok || exit.ExitCode() != 2 {
		t.Fatalf("expected exit 2: %v\n%s", err, out)
	}
	for _, p := range []string{f.render, f.codex, filepath.Join(f.root, "new"), oldHome} {
		f.absent(p)
	}
}

func TestMarketplaceLegacyCatalogMigrationPreservesUnrelatedData(t *testing.T) {
	f := modeFixture(t)
	old := f.put(f.home, "plugins/skillset-uoleg/keep.txt", "old package")
	oldRoot := filepath.Join(f.home, "plugins/skillset-uoleg")
	f.put(oldRoot, ".codex-plugin/plugin.json", `{"version":"0.1.0"}`)
	inventory, err := marketplaceInventory(oldRoot)
	if err != nil {
		t.Fatal(err)
	}
	marker, err := json.Marshal(map[string]any{"files": inventory})
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(oldRoot, marketplaceMarker), marker, 0644); err != nil {
		t.Fatal(err)
	}
	catalog := f.put(f.home, ".agents/plugins/marketplace.json", `{"name":"custom","plugins":[{"name":"other","source":{"source":"npm","package":"example"}},{"name":"skillset-uoleg","source":"./plugins/skillset-uoleg","policy":{"installation":"NOT_AVAILABLE","authentication":"MANUAL"}}]}`)
	modeRun(t, f, 0, "--mode", `codex:{"agents":false,"plugin":"full","plugin_name":"skillset-uoleg"}`)
	if f.read(old) != "old package" {
		t.Fatal("old package changed")
	}
	data := f.json(catalog)
	plugins := data["plugins"].([]any)
	if len(plugins) != 2 || plugins[0].(map[string]any)["name"] != "other" {
		t.Fatal(plugins)
	}
	entry := plugins[1].(map[string]any)
	if entry["policy"].(map[string]any)["installation"] != "NOT_AVAILABLE" {
		t.Fatal("policy changed")
	}
	source := entry["source"]
	var path string
	switch v := source.(type) {
	case string:
		path = v
	case map[string]any:
		path, _ = v["path"].(string)
	}
	want := filepath.Join(f.home, ".agents/plugins/skillset-uoleg")
	if filepath.Clean(filepath.Join(f.home, path)) != want {
		t.Fatalf("source %q does not resolve to %s", path, want)
	}
	if _, err := os.Stat(filepath.Join(want, "skills/Alpha/SKILL.md")); err != nil {
		t.Fatal(err)
	}
}

func TestMarketplaceUnmanagedDestinationCollisionPreserved(t *testing.T) {
	f := modeFixture(t)
	marker := f.put(f.home, ".agents/plugins/skillset/keep.txt", "unmanaged")
	modeRun(t, f, 2, "--mode", `codex:{"agents":false,"plugin":"full"}`)
	if f.read(marker) != "unmanaged" {
		t.Fatal("unmanaged plugin changed")
	}
	f.absent(filepath.Join(f.home, ".agents/plugins/marketplace.json"))
	_ = f.latest()
	f.absent(f.codex)
}

func TestMarketplaceUnmanagedLegacyPackagePreserved(t *testing.T) {
	f := modeFixture(t)
	old := f.put(f.home, "plugins/skillset-uoleg/keep.txt", "unmanaged old package")
	original := `{"name":"custom","plugins":[{"name":"skillset-uoleg","source":"./plugins/skillset-uoleg"}]}`
	catalog := f.put(f.home, ".agents/plugins/marketplace.json", original)
	modeRun(t, f, 2, "--mode", `codex:{"agents":false,"plugin":"full","plugin_name":"skillset-uoleg"}`)
	if f.read(old) != "unmanaged old package" || f.read(catalog) != original {
		t.Fatal("unmanaged legacy data changed")
	}
	f.absent(filepath.Join(f.home, ".agents/plugins/skillset-uoleg"))
}

func TestMarketplaceCustomPluginName(t *testing.T) {
	f := modeFixture(t)
	out := modeRun(t, f, 0, "--mode", `codex:{"agents":false,"plugin":"skills","plugin_name":"custom-skill_1"}`)
	line, _, _ := strings.Cut(out, "\n")
	var config map[string]any
	if !strings.HasPrefix(line, "config: ") {
		t.Fatalf("first line %q", line)
	}
	if err := json.Unmarshal([]byte(strings.TrimPrefix(line, "config: ")), &config); err != nil {
		t.Fatal(err)
	}
	if config["plugin_name"] != "custom-skill_1" {
		t.Fatalf("plugin name in config: %v", config)
	}
	dir := filepath.Join(f.home, ".agents/plugins/custom-skill_1")
	manifest := f.json(filepath.Join(dir, ".codex-plugin/plugin.json"))
	if manifest["name"] != "custom-skill_1" {
		t.Fatalf("manifest name %v", manifest)
	}
	if _, err := os.Stat(filepath.Join(dir, "skills/Alpha/SKILL.md")); err != nil {
		t.Fatal(err)
	}
	f.absent(filepath.Join(dir, "skills/Agent_Guidance"))
	catalog := f.json(filepath.Join(f.home, ".agents/plugins/marketplace.json"))
	var found bool
	for _, value := range catalog["plugins"].([]any) {
		entry := value.(map[string]any)
		if entry["name"] == "custom-skill_1" {
			found = true
			source := entry["source"].(map[string]any)["path"].(string)
			if filepath.Clean(filepath.Join(f.home, source)) != dir {
				t.Fatalf("invalid source %q", source)
			}
		}
	}
	if !found {
		t.Fatal("custom plugin missing from catalog")
	}
	_ = f.latest()
	f.absent(f.codex)
}

func TestMarketplaceDefaultPreservesOldNamedPlugin(t *testing.T) {
	f := modeFixture(t)
	old := f.put(f.home, "plugins/skillset-uoleg/keep.txt", "old package")
	catalog := f.put(f.home, ".agents/plugins/marketplace.json", `{"name":"personal","plugins":[{"name":"skillset-uoleg","source":"./plugins/skillset-uoleg"}]}`)
	modeRun(t, f, 0, "--mode", `codex:{"agents":false,"plugin":"full"}`)
	if f.read(old) != "old package" {
		t.Fatal("old plugin changed")
	}
	plugins := f.json(catalog)["plugins"].([]any)
	names := map[string]bool{}
	for _, v := range plugins {
		e := v.(map[string]any)
		names[e["name"].(string)] = true
		if e["name"] == "skillset-uoleg" && e["source"] != "./plugins/skillset-uoleg" {
			t.Fatalf("old entry changed: %v", e)
		}
	}
	if !names["skillset"] || !names["skillset-uoleg"] {
		t.Fatalf("unexpected entries: %v", names)
	}
	if _, err := os.Stat(filepath.Join(f.home, ".agents/plugins/skillset/skills/Alpha/SKILL.md")); err != nil {
		t.Fatal(err)
	}
}

func TestMarketplaceRejectsInvalidPluginNamesBeforeWrites(t *testing.T) {
	for _, name := range []string{"", "..", "../outside", "/tmp/escape", "Uppercase", "bad name", "bad.name"} {
		t.Run(name, func(t *testing.T) {
			f := modeFixture(t)
			value, _ := json.Marshal(name)
			mode := `codex:{"agents":false,"plugin":"full","plugin_name":` + string(value) + `}`
			modeRun(t, f, 2, "--mode", mode)
			f.absent(f.home)
			f.absent(f.codex)
			f.absent(f.render)
		})
	}
	for _, mode := range []string{`codex:{"plugin_name":"valid"}`, `codex:{"agents":false,"plugin_name":"valid"}`, `codex:{"plugin":"full","plugin_name":null}`, `codex:{"plugin":"full","plugin_name":true}`} {
		t.Run(mode, func(t *testing.T) {
			f := modeFixture(t)
			modeRun(t, f, 2, "--mode", mode)
			f.absent(f.home)
			f.absent(f.codex)
			f.absent(f.render)
		})
	}
}

func TestMarketplaceMissingLegacyPackagePreservesCatalog(t *testing.T) {
	f := modeFixture(t)
	original := `{"name":"personal","plugins":[{"name":"skillset-uoleg","source":"./plugins/skillset-uoleg"}]}`
	catalog := f.put(f.home, ".agents/plugins/marketplace.json", original)
	modeRun(t, f, 2, "--mode", `codex:{"agents":false,"plugin":"full","plugin_name":"skillset-uoleg"}`)
	if f.read(catalog) != original {
		t.Fatal("catalog changed")
	}
	f.absent(filepath.Join(f.home, ".agents/plugins/skillset-uoleg"))
	_ = f.latest()
	f.absent(f.codex)
}
