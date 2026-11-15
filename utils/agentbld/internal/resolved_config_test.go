package internal

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

type reportedConfig struct {
	Command              string   `json:"command"`
	SectionsDirs         []string `json:"sections_dirs"`
	RenderDir            string   `json:"render_dir"`
	CodexDir             string   `json:"codex_dir"`
	CodexMarketplaceHome string   `json:"codex_marketplace_home"`
	Agents               bool     `json:"agents"`
	Skills               bool     `json:"skills"`
	PluginName           string   `json:"plugin_name"`
	Plugin               string   `json:"plugin"`
	Legacy               struct {
		SkillMarketplace bool `json:"skill_marketplace"`
		InstructionsOnly bool `json:"instructions_only"`
		RulesOnly        bool `json:"rules_only"`
	} `json:"legacy"`
}

func TestResolvedConfigFirstLine(t *testing.T) {
	for _, tc := range []struct {
		name, command                              string
		extra                                      []string
		status                                     int
		agents, skills                             bool
		plugin, legacy                             string
		envOverrides, flagsOverride, missingSource bool
	}{
		{name: "defaults", command: "render", agents: true, skills: true, plugin: "none"},
		{name: "environment", command: "render", agents: true, skills: true, plugin: "none", envOverrides: true},
		{name: "flags-override-environment", command: "render", agents: true, skills: true, plugin: "none", envOverrides: true, flagsOverride: true},
		{name: "plugin-default", command: "render", extra: []string{"--mode", `codex:{"plugin":"full"}`}, agents: false, skills: false, plugin: "full"},
		{name: "full-agents-first", command: "render", extra: []string{"--mode", `codex:{"agents":true,"plugin":"full"}`}, agents: true, skills: false, plugin: "full"},
		{name: "full-agents-last", command: "render", extra: []string{"--mode", `codex:{"plugin":"full","agents":true}`}, agents: true, skills: false, plugin: "full"},
		{name: "skills-default", command: "render", extra: []string{"--mode", `codex:{"plugin":"skills"}`}, agents: true, skills: false, plugin: "skills"},
		{name: "plugin-skills-first", command: "render", extra: []string{"--mode", `codex:{"skills":true,"plugin":"skills"}`}, agents: true, skills: true, plugin: "skills"},
		{name: "plugin-skills-last", command: "render", extra: []string{"--mode", `codex:{"plugin":"skills","skills":true}`}, agents: true, skills: true, plugin: "skills"},
		{name: "plugin-only", command: "render", extra: []string{"--mode", `codex:{"agents":false,"plugin":"skills"}`}, agents: false, skills: false, plugin: "skills"},
		{name: "legacy-marketplace", command: "render", extra: []string{"--skill-marketplace"}, agents: true, skills: false, plugin: "full", legacy: "marketplace"},
		{name: "legacy-instructions", command: "render", extra: []string{"--instructions-only"}, agents: true, skills: false, plugin: "none", legacy: "instructions"},
		{name: "legacy-rules", command: "render", extra: []string{"--rules-only"}, agents: false, skills: false, plugin: "none", legacy: "rules"},
		{name: "check", command: "check", status: 1, agents: true, skills: true, plugin: "none"},
		{name: "prune", command: "prune", agents: true, skills: true, plugin: "none"},
		{name: "reset", command: "reset", agents: true, skills: true, plugin: "none"},
		{name: "source-error", command: "render", status: 2, agents: true, skills: true, plugin: "none", missingSource: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := modeFixture(t)
			home := filepath.Join(f.root, "isolated-home")
			if err := os.MkdirAll(home, 0755); err != nil {
				t.Fatal(err)
			}
			env := []string{}
			for _, item := range os.Environ() {
				key, _, _ := strings.Cut(item, "=")
				if key != "HOME" && key != "XDG_STATE_HOME" && key != "CODEX_HOME" {
					env = append(env, item)
				}
			}
			env = append(env, "HOME="+home)
			render := filepath.Join(home, ".local/state/agentbld")
			codex := filepath.Join(home, ".codex")
			if tc.envOverrides {
				state := filepath.Join(f.root, "state-env")
				codex = filepath.Join(f.root, "codex-env")
				render = filepath.Join(state, "agentbld")
				env = append(env, "XDG_STATE_HOME="+state, "CODEX_HOME="+codex)
			}
			sections := filepath.Join(f.config, "sections")
			if tc.missingSource {
				sections = filepath.Join(f.root, "missing-sections")
			}
			args := []string{tc.command, "--sections-dir", sections}
			if tc.command == "render" && tc.name != "defaults" {
				args = append(args, "--sync")
			}
			if tc.flagsOverride {
				render = filepath.Join(f.root, "render-flag")
				codex = filepath.Join(f.root, "codex-flag")
				args = append(args, "--render-dir", render, "--codex-dir", codex)
			}
			args = append(args, tc.extra...)
			cmd := exec.Command(contractBinary, args...)
			cmd.Env = env
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err := cmd.Run()
			status := 0
			if err != nil {
				if exit, ok := err.(*exec.ExitError); ok {
					status = exit.ExitCode()
				} else {
					t.Fatal(err)
				}
			}
			if status != tc.status {
				t.Fatalf("status %d, want %d\nstdout: %s\nstderr: %s", status, tc.status, &stdout, &stderr)
			}
			first, _, _ := strings.Cut(stdout.String(), "\n")
			if !strings.HasPrefix(first, "config: ") {
				t.Fatalf("first stdout line is not config: %q", first)
			}
			var got reportedConfig
			if err := json.Unmarshal([]byte(strings.TrimPrefix(first, "config: ")), &got); err != nil {
				t.Fatalf("invalid config JSON: %v: %s", err, first)
			}
			want := reportedConfig{Command: tc.command, SectionsDirs: []string{sections}, RenderDir: render, CodexDir: codex, CodexMarketplaceHome: filepath.Join(home, ".agents/plugins"), Agents: tc.agents, Skills: tc.skills, Plugin: tc.plugin, PluginName: "skillset"}
			want.Legacy.SkillMarketplace = tc.legacy == "marketplace"
			want.Legacy.InstructionsOnly = tc.legacy == "instructions"
			want.Legacy.RulesOnly = tc.legacy == "rules"
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("config = %#v, want %#v", got, want)
			}
			if tc.missingSource {
				f.absent(render)
				f.absent(codex)
			}
		})
	}
}

func TestNonDefaultSelectionRejectedOutsideSync(t *testing.T) {
	for _, op := range []string{"prune", "reset"} {
		for _, mode := range []string{`codex:{"agents":false}`, `codex:{"skills":false}`} {
			t.Run(op+mode, func(t *testing.T) {
				f := modeFixture(t)
				cmd := exec.Command(contractBinary, op, "--mode", mode, "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render, "--codex-dir", f.codex, "--home", f.home)
				out, err := cmd.CombinedOutput()
				exit, ok := err.(*exec.ExitError)
				if !ok || exit.ExitCode() != 2 {
					t.Fatalf("expected exit 2, got %v: %s", err, out)
				}
				f.absent(f.render)
				f.absent(f.codex)
				f.absent(f.home)
			})
		}
	}
}

func TestCheckPluginGenerationWithSameMode(t *testing.T) {
	for _, plugin := range []string{"full", "skills"} {
		t.Run(plugin, func(t *testing.T) {
			f := modeFixture(t)
			mode := `codex:{"agents":true,"plugin":"` + plugin + `"}`
			modeRun(t, f, 0, "--mode", mode)
			manifest := filepath.Join(modePlugin(f), ".codex-plugin/plugin.json")
			catalog := filepath.Join(f.home, ".agents/plugins/marketplace.json")
			beforeManifest, beforeCatalog := f.read(manifest), f.read(catalog)
			args := []string{"check", "--mode", mode, "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render, "--codex-dir", f.codex, "--home", f.home}
			out, err := exec.Command(contractBinary, args...).CombinedOutput()
			if err != nil {
				t.Fatalf("check after %s sync: %v\n%s", plugin, err, out)
			}
			if !strings.Contains(string(out), "status=up-to-date") {
				t.Fatalf("check did not report up-to-date: %s", out)
			}
			if f.read(manifest) != beforeManifest || f.read(catalog) != beforeCatalog {
				t.Fatal("check changed plugin manifest or catalog")
			}
		})
	}
}
