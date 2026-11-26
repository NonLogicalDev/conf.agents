package integration_tests

import (
	"encoding/json"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func lifecycleRun(t *testing.T, f *contractFixture, want int, args ...string) string {
	t.Helper()
	cmd := exec.Command(contractBinary, args...)
	cmd.Env = append(os.Environ(), "HOME="+filepath.Join(f.root, "isolated-home"))
	out, err := cmd.CombinedOutput()
	status := 0
	if err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			status = e.ExitCode()
		} else {
			t.Fatal(err)
		}
	}
	if status != want {
		t.Fatalf("%v: exit %d, want %d\n%s", args, status, want, out)
	}
	if status == 0 {
		first, _, _ := strings.Cut(string(out), "\n")
		if !strings.HasPrefix(first, "config: ") {
			t.Fatalf("first line %q is not config", first)
		}
	}
	return string(out)
}
func lifecycleRenderArgs(f *contractFixture, extra ...string) []string {
	args := []string{"render", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render}
	for _, flag := range extra {
		if flag == "--sync" {
			args = append(args, "--codex-dir", f.codex, "--codex-marketplace-home", filepath.Join(f.home, ".agents/plugins"))
			break
		}
	}
	return append(args, extra...)
}
func lifecycleOrdinal(t *testing.T, path string) string {
	t.Helper()
	_, n, ok := strings.Cut(filepath.Base(path), "__")
	if !ok || n == "" {
		t.Fatalf("invalid generation name %s", path)
	}
	return n
}

func TestHistoricalGenerationSyncToIndependentDestinations(t *testing.T) {
	f := modeFixture(t)
	mode := `codex:{"agents":true,"skills":true,"plugin":"full"}`
	f.put(f.layer, "skills/Alpha/SKILL.md", "---\nname: Alpha\ndescription: Example\n---\nversion one\n")
	f.put(f.layer, "AGENTS.md.d/10-base.md", "instruction one {{%_resources_%}} {{%_skills_%}}\n")
	f.put(f.layer, "resources/guide/first.md", "resource one\n")
	f.put(f.layer, "rules/default.rules", "rule one\n")
	render := func() {
		lifecycleRun(t, f, 0, "render", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render, "--mode", mode)
	}
	render()
	first := f.latest()
	number := lifecycleOrdinal(t, first)
	f.put(f.layer, "skills/Alpha/SKILL.md", "---\nname: Alpha\ndescription: Example\n---\nversion two\n")
	f.put(f.layer, "AGENTS.md.d/10-base.md", "instruction two {{%_resources_%}} {{%_skills_%}}\n")
	f.put(f.layer, "resources/guide/first.md", "resource two\n")
	f.put(f.layer, "rules/default.rules", "rule two\n")
	render()
	second := f.latest()
	if first == second {
		t.Fatal("changed source did not create generation")
	}
	if err := os.RemoveAll(f.config); err != nil {
		t.Fatal(err)
	}
	homes := []string{filepath.Join(f.root, "codex-a"), filepath.Join(f.root, "codex-b")}
	marketplaces := []string{filepath.Join(f.root, "a/.agents/plugins"), filepath.Join(f.root, "b/.agents/plugins")}
	before := map[string][]byte{}
	err := filepath.WalkDir(first, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Type().IsRegular() {
			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			before[path] = b
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	var firstInstructions, firstGuidance string
	for i := range homes {
		lifecycleRun(t, f, 0, "sync", number, "--render-dir", f.render, "--codex-dir", homes[i], "--codex-marketplace-home", marketplaces[i])
		instructions := f.read(filepath.Join(homes[i], "AGENTS.md"))
		guidance := f.read(filepath.Join(homes[i], "skills/Agent_Guidance/SKILL.md"))
		for _, part := range []string{"instruction one", filepath.Join(homes[i], "skills"), filepath.Join(homes[i], "resources")} {
			if !strings.Contains(instructions, part) {
				t.Fatalf("instructions for %s missing %q: %s", homes[i], part, instructions)
			}
		}
		if !strings.Contains(guidance, filepath.Join(homes[i], "skills")) {
			t.Fatalf("guidance for %s missing target path: %s", homes[i], guidance)
		}
		for path, want := range map[string]string{"skills/Alpha/SKILL.md": "version one", "resources/guide/first.md": "resource one", "rules/default.rules": "rule one"} {
			if !strings.Contains(f.read(filepath.Join(homes[i], path)), want) {
				t.Fatalf("installed %s missing %s", path, want)
			}
		}
		if got := f.read(filepath.Join(marketplaces[i], "skillset/skills/Alpha/SKILL.md")); !strings.Contains(got, "version one") || strings.Contains(got, "version two") {
			t.Fatalf("plugin used live/new source: %s", got)
		}
		if i == 0 {
			firstInstructions = instructions
			firstGuidance = guidance
		} else if f.read(filepath.Join(homes[0], "AGENTS.md")) != firstInstructions || f.read(filepath.Join(homes[0], "skills/Agent_Guidance/SKILL.md")) != firstGuidance {
			t.Fatal("second sync changed first destination")
		}
	}
	if f.latest() != second {
		t.Fatal("historical sync changed latest")
	}
	for path, old := range before {
		now, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if string(now) != string(old) {
			t.Fatalf("generation changed: %s", path)
		}
		for _, dest := range append(homes, marketplaces...) {
			if strings.Contains(string(now), dest) {
				t.Fatalf("generation %s contains destination %s", path, dest)
			}
		}
	}
}
func TestGenerationRepeatedModeValidationBeforeWrites(t *testing.T) {
	for _, extra := range [][]string{{"--mode", "codex", "--mode", `codex:{}`}, {"--mode", "unknown"}, {"--mode", "codex", "--mode", "unknown"}} {
		t.Run(strings.Join(extra, "_"), func(t *testing.T) {
			f := modeFixture(t)
			lifecycleRun(t, f, 2, lifecycleRenderArgs(f, extra...)...)
			f.absent(f.render)
			f.absent(f.codex)
			f.absent(f.home)
		})
	}
}
func TestGenerationRejectsInvalidSelection(t *testing.T) {
	f := modeFixture(t)
	lifecycleRun(t, f, 0, lifecycleRenderArgs(f)...)
	_ = f.latest()
	for _, selector := range []string{"abc", "0", "999999", "../1"} {
		lifecycleRun(t, f, 2, "sync", selector, "--render-dir", f.render)
		f.absent(f.codex)
	}
}

func TestPluginOnlyRenderThenCheckDoesNotActivate(t *testing.T) {
	for _, plugin := range []string{"full", "skills"} {
		t.Run(plugin, func(t *testing.T) {
			f := modeFixture(t)
			mode := `codex:{"agents":false,"plugin":"` + plugin + `"}`
			lifecycleRun(t, f, 0, lifecycleRenderArgs(f, "--mode", mode)...)
			gen := f.latest()
			f.absent(f.home)
			f.absent(f.codex)
			out := lifecycleRun(t, f, 0, "check", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render, "--mode", mode)
			if !strings.Contains(out, "status=up-to-date") {
				t.Fatalf("plugin-only check not up to date: %s", out)
			}
			if f.latest() != gen {
				t.Fatal("check changed generation")
			}
			f.absent(f.home)
			f.absent(f.codex)
		})
	}
}

func TestGenerationIdentityExcludesSyncDestinations(t *testing.T) {
	f := modeFixture(t)
	mode := `codex:{"agents":true,"skills":true,"plugin":"full"}`
	homes := []string{filepath.Join(f.root, "codex-a"), filepath.Join(f.root, "codex-b")}
	marketplaces := []string{filepath.Join(f.root, "a/.agents/plugins"), filepath.Join(f.root, "b/.agents/plugins")}
	var first string
	for i := range homes {
		lifecycleRun(t, f, 0, "render", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render, "--codex-dir", homes[i], "--codex-marketplace-home", marketplaces[i], "--mode", mode)
		gen := f.latest()
		if i == 0 {
			first = gen
		} else if gen != first {
			t.Fatalf("destination-only change created generation: %s -> %s", first, gen)
		}
		f.absent(homes[i])
		f.absent(marketplaces[i])
	}
	data := f.read(filepath.Join(first, "config.json"))
	for _, destination := range append(homes, marketplaces...) {
		if strings.Contains(data, destination) {
			t.Fatalf("generation config contains destination %s", destination)
		}
	}
	var metadata any
	if err := json.Unmarshal([]byte(data), &metadata); err != nil {
		t.Fatal(err)
	}
	var inspect func(any)
	inspect = func(value any) {
		switch v := value.(type) {
		case map[string]any:
			for key, item := range v {
				if key == "codex_dir" || key == "codex_marketplace_home" {
					t.Errorf("generation config contains destination key %q", key)
				}
				inspect(item)
			}
		case []any:
			for _, item := range v {
				inspect(item)
			}
		}
	}
	inspect(metadata)
}

func TestUnsupportedHarnessModeFailsBeforeWrites(t *testing.T) {
	f := modeFixture(t)
	out, err := exec.Command(contractBinary, "render", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render, "--mode", "test-harness:{}").CombinedOutput()
	exit, ok := err.(*exec.ExitError)
	if !ok || exit.ExitCode() != 2 {
		t.Fatalf("expected exit 2: %v\n%s", err, out)
	}
	f.absent(f.render)
	f.absent(f.codex)
	f.absent(f.home)
}
