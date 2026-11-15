package internal

import (
	"encoding/json"
	"fmt"
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

func TestGenerationStoresModeAndTracksConfigIdentity(t *testing.T) {
	f := modeFixture(t)
	lifecycleRun(t, f, 0, lifecycleRenderArgs(f, "--mode", `codex:{"agents":true,"skills":true}`)...)
	first := f.latest()
	data := f.read(filepath.Join(first, "config.json"))
	var metadata map[string]json.RawMessage
	if err := json.Unmarshal([]byte(data), &metadata); err != nil {
		t.Fatalf("invalid config.json: %v", err)
	}
	if _, ok := metadata["codex"]; !ok {
		t.Fatalf("missing codex mode in %s", data)
	}
	if info, err := os.Stat(filepath.Join(first, "config.json")); err != nil || info.Mode().Perm()&0222 != 0 {
		t.Fatalf("generation config not immutable: %v %v", info, err)
	}
	lifecycleRun(t, f, 0, lifecycleRenderArgs(f, "--mode", `codex:{"agents":true,"skills":true}`)...)
	if f.latest() != first {
		t.Fatal("unchanged render created generation")
	}
	lifecycleRun(t, f, 0, lifecycleRenderArgs(f, "--mode", `codex:{"agents":true,"skills":true,"plugin":"full"}`)...)
	if f.latest() == first {
		t.Fatal("changed resolved mode reused prior generation")
	}
}
func TestPluginOnlyCreatesGeneration(t *testing.T) {
	f := modeFixture(t)
	lifecycleRun(t, f, 0, lifecycleRenderArgs(f, "--mode", `codex:{"plugin":"full"}`)...)
	gen := f.latest()
	if _, err := os.Stat(filepath.Join(gen, "config.json")); err != nil {
		t.Fatal(err)
	}
	f.absent(f.codex)
	f.absent(f.home)
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
func TestGenerationMalformedMetadataRejectedWithoutActivation(t *testing.T) {
	for _, kind := range []string{"missing", "malformed"} {
		t.Run(kind, func(t *testing.T) {
			f := modeFixture(t)
			lifecycleRun(t, f, 0, lifecycleRenderArgs(f)...)
			gen := f.latest()
			number := lifecycleOrdinal(t, gen)
			if err := os.Chmod(gen, 0755); err != nil {
				t.Fatal(err)
			}
			path := filepath.Join(gen, "config.json")
			if kind == "missing" {
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
			} else {
				if err := os.Chmod(path, 0644); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, []byte("{bad"), 0444); err != nil {
					t.Fatal(err)
				}
				if err := os.Chmod(path, 0444); err != nil {
					t.Fatal(err)
				}
			}
			if err := os.Chmod(gen, 0555); err != nil {
				t.Fatal(err)
			}
			lifecycleRun(t, f, 2, "sync", number, "--render-dir", f.render)
			f.absent(f.codex)
		})
	}
}
func TestLegacyGenerationRemainsDiscoverableAndPrunable(t *testing.T) {
	f := modeFixture(t)
	gen := filepath.Join(f.render, "2025-01-01__01")
	if err := os.MkdirAll(gen, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(gen, "AGENTS.md"), []byte("<!-- generated: agent-config; do not edit -->\n"), 0444); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(gen, 0555); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Base(gen), filepath.Join(f.render, "latest")); err != nil {
		t.Fatal(err)
	}
	lifecycleRun(t, f, 2, "sync", "1", "--render-dir", f.render)
	f.absent(f.codex)
	lifecycleRun(t, f, 0, "prune", "--render-dir", f.render, "--codex-dir", f.codex, "--sections-dir", filepath.Join(f.config, "sections"))
	if _, err := os.Stat(gen); err != nil {
		t.Fatal(fmt.Errorf("legacy generation lost: %w", err))
	}
}

func TestPruneProtectsGenerationActiveInOtherDestinations(t *testing.T) {
	f := modeFixture(t)
	render := func() {
		lifecycleRun(t, f, 0, "render", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render)
	}
	render()
	first := f.latest()
	number := lifecycleOrdinal(t, first)
	homes := []string{filepath.Join(f.root, "codex-a"), filepath.Join(f.root, "codex-b")}
	for _, home := range homes {
		lifecycleRun(t, f, 0, "sync", number, "--render-dir", f.render, "--codex-dir", home)
	}
	for i := 0; i < 5; i++ {
		f.put(f.layer, "AGENTS.md.d/10-base.md", fmt.Sprintf("new generation %d\n", i))
		render()
	}
	args := []string{"prune", "--render-dir", f.render, "--codex-dir", filepath.Join(f.root, "unrelated-codex"), "--sections-dir", filepath.Join(f.config, "sections")}
	cmd := exec.Command(contractBinary, args...)
	cmd.Env = append(os.Environ(), "HOME="+filepath.Join(f.root, "isolated-home"))
	out, err := cmd.CombinedOutput()
	if err != nil {
		exit, ok := err.(*exec.ExitError)
		if !ok || exit.ExitCode() != 2 || (!strings.Contains(strings.ToLower(string(out)), "active") && !strings.Contains(strings.ToLower(string(out)), "referenc")) {
			t.Fatalf("prune failed for unexpected reason: %v\n%s", err, out)
		}
	}
	if _, err := os.Stat(first); err != nil {
		t.Fatalf("pruned generation used elsewhere: %v\n%s", err, out)
	}
	for _, home := range homes {
		if got := f.read(filepath.Join(home, "AGENTS.md")); !strings.Contains(got, "Base ") {
			t.Fatalf("active target broken in %s: %s", home, got)
		}
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

func TestPruneFailsClosedWhenRecordedDestinationUnavailable(t *testing.T) {
	f := modeFixture(t)
	render := func() {
		lifecycleRun(t, f, 0, "render", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render)
	}
	render()
	first := f.latest()
	number := lifecycleOrdinal(t, first)
	codexB := filepath.Join(f.root, "codex-b")
	offline := filepath.Join(f.root, "codex-b-offline")
	lifecycleRun(t, f, 0, "sync", number, "--render-dir", f.render, "--codex-dir", codexB)
	if err := os.Rename(codexB, offline); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if _, err := os.Lstat(offline); err == nil {
			_ = os.Rename(offline, codexB)
		}
	})
	for i := 0; i < 5; i++ {
		f.put(f.layer, "AGENTS.md.d/10-base.md", fmt.Sprintf("new generation %d\n", i))
		render()
	}
	out := lifecycleRun(t, f, 2, "prune", "--render-dir", f.render, "--codex-dir", filepath.Join(f.root, "codex-a"), "--sections-dir", filepath.Join(f.config, "sections"))
	if _, err := os.Stat(first); err != nil {
		t.Fatalf("pruned generation while destination was unavailable: %v\n%s", err, out)
	}
	if err := os.Rename(offline, codexB); err != nil {
		t.Fatal(err)
	}
	if got := f.read(filepath.Join(codexB, "AGENTS.md")); !strings.Contains(got, "Base ") {
		t.Fatalf("restored instructions broken: %s", got)
	}
	if _, err := os.ReadFile(filepath.Join(codexB, "skills/Alpha/SKILL.md")); err != nil {
		t.Fatalf("restored skill link broken: %v", err)
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

func TestPruneFailsClosedWhenDestinationIsReplacedByEmptyDirectory(t *testing.T) {
	f := modeFixture(t)
	render := func() {
		lifecycleRun(t, f, 0, "render", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render)
	}
	render()
	first := f.latest()
	number := lifecycleOrdinal(t, first)
	codexB := filepath.Join(f.root, "codex-b")
	offline := filepath.Join(f.root, "codex-b-offline")
	lifecycleRun(t, f, 0, "sync", number, "--render-dir", f.render, "--codex-dir", codexB)
	if err := os.Rename(codexB, offline); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(codexB, 0755); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if _, err := os.Lstat(offline); err == nil {
			_ = os.Remove(codexB)
			_ = os.Rename(offline, codexB)
		}
	})
	for i := 0; i < 5; i++ {
		f.put(f.layer, "AGENTS.md.d/10-base.md", fmt.Sprintf("new generation %d\n", i))
		render()
	}
	out := lifecycleRun(t, f, 2, "prune", "--render-dir", f.render, "--codex-dir", filepath.Join(f.root, "codex-a"), "--sections-dir", filepath.Join(f.config, "sections"))
	if _, err := os.Stat(first); err != nil {
		t.Fatalf("pruned generation while destination was replaced: %v\n%s", err, out)
	}
	if err := os.Remove(codexB); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(offline, codexB); err != nil {
		t.Fatal(err)
	}
	if got := f.read(filepath.Join(codexB, "AGENTS.md")); !strings.Contains(got, "Base ") {
		t.Fatalf("restored instructions broken: %s", got)
	}
	if _, err := os.ReadFile(filepath.Join(codexB, "skills/Alpha/SKILL.md")); err != nil {
		t.Fatalf("restored skill link broken: %v", err)
	}
}

func TestSyncRejectsReplacedRecordedDestination(t *testing.T) {
	f := modeFixture(t)
	lifecycleRun(t, f, 0, "render", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render)
	first := f.latest()
	number := lifecycleOrdinal(t, first)
	codexB := filepath.Join(f.root, "codex-b")
	offline := filepath.Join(f.root, "codex-b-offline")
	syncArgs := []string{"sync", number, "--render-dir", f.render, "--codex-dir", codexB}
	lifecycleRun(t, f, 0, syncArgs...)
	if err := os.Rename(codexB, offline); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(codexB, 0755); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if _, err := os.Lstat(offline); err == nil {
			_ = os.Remove(codexB)
			_ = os.Rename(offline, codexB)
		}
	})
	lifecycleRun(t, f, 2, syncArgs...)
	entries, err := os.ReadDir(codexB)
	if err != nil || len(entries) != 0 {
		t.Fatalf("replaced destination changed: %v, %v", entries, err)
	}
	if _, err := os.Stat(first); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(codexB); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(offline, codexB); err != nil {
		t.Fatal(err)
	}
	lifecycleRun(t, f, 0, syncArgs...)
	if got := f.read(filepath.Join(codexB, "AGENTS.md")); !strings.Contains(got, "Base ") {
		t.Fatalf("restored instructions broken: %s", got)
	}
}
