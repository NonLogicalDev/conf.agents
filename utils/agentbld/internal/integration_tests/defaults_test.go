package integration_tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Each subprocess receives its own home and explicit environment. This avoids
// consulting the user's configuration while checking real CLI behavior.
func TestCLIDefaultDestinations(t *testing.T) {
	for _, tc := range []struct {
		name, xdg, codex                      string
		setEnv, explicitRender, explicitCodex bool
	}{
		{name: "unset"},
		{name: "empty", setEnv: true},
		{name: "environment", xdg: "env-state", codex: "env-codex", setEnv: true},
		{name: "explicit-flags", xdg: "env-state", codex: "env-codex", setEnv: true, explicitRender: true, explicitCodex: true},
		{name: "render-flag-only", xdg: "env-state", codex: "env-codex", setEnv: true, explicitRender: true},
		{name: "codex-flag-only", xdg: "env-state", codex: "env-codex", setEnv: true, explicitCodex: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newContractFixture(t)
			home := filepath.Join(f.root, "isolated-home")
			if err := os.MkdirAll(home, 0755); err != nil {
				t.Fatal(err)
			}
			env := []string{}
			for _, entry := range os.Environ() {
				key, _, _ := strings.Cut(entry, "=")
				if key != "HOME" && key != "XDG_STATE_HOME" && key != "CODEX_HOME" {
					env = append(env, entry)
				}
			}
			env = append(env, "HOME="+home)
			if tc.setEnv {
				xdg, codex := tc.xdg, tc.codex
				if xdg != "" {
					xdg = filepath.Join(f.root, xdg)
				}
				if codex != "" {
					codex = filepath.Join(f.root, codex)
				}
				env = append(env, "XDG_STATE_HOME="+xdg, "CODEX_HOME="+codex)
			}
			state := filepath.Join(home, ".local/state")
			if tc.xdg != "" {
				state = filepath.Join(f.root, tc.xdg)
			}
			render := filepath.Join(state, "agentbld")
			codex := filepath.Join(home, ".codex")
			if tc.codex != "" {
				codex = filepath.Join(f.root, tc.codex)
			}
			flags := []string{"--sections-dir", filepath.Join(f.config, "sections")}
			if tc.explicitRender {
				render = filepath.Join(f.root, "explicit-render")
				flags = append(flags, "--render-dir", render)
			}
			if tc.explicitCodex {
				codex = filepath.Join(f.root, "explicit-codex")
				flags = append(flags, "--codex-dir", codex)
			}
			run := func(op string, extra ...string) {
				t.Helper()
				cmd := exec.Command(contractBinary, append(append([]string{op}, flags...), extra...)...)
				cmd.Env = env
				if out, err := cmd.CombinedOutput(); err != nil {
					t.Fatalf("%s: %v\n%s", op, err, out)
				}
			}
			run("render")
			latest := f.resolved(filepath.Join(render, "latest"))
			agents := f.read(filepath.Join(latest, "codex", "AGENTS.md"))
			if !strings.Contains(agents, "{{%_skills_%}}") || !strings.Contains(agents, "{{%_resources_%}}") {
				t.Fatalf("render did not preserve destination-independent template: %s", agents)
			}
			f.absent(filepath.Join(codex, "AGENTS.md"))
			run("check")
			run("render", "--sync")
			f.symlink(filepath.Join(codex, "AGENTS.md"))
			agents = f.read(filepath.Join(codex, "AGENTS.md"))
			if !strings.Contains(agents, filepath.Join(codex, "skills")) || !strings.Contains(agents, filepath.Join(codex, "resources")) {
				t.Fatalf("activated instructions use wrong Codex path: %s", agents)
			}
			f.symlink(filepath.Join(codex, "skills/Alpha"))
			if got := f.read(filepath.Join(codex, "skills/Alpha/SKILL.md")); got != f.read(filepath.Join(latest, "codex", "skills/Alpha/SKILL.md")) {
				t.Fatal("activated skill differs from generation")
			}

			if got := f.resolved(filepath.Join(render, "latest")); got != latest {
				t.Fatalf("unchanged source created another generation: %s", got)
			}
			if tc.explicitRender && tc.xdg != "" {
				f.absent(filepath.Join(f.root, tc.xdg, "agentbld"))
			}
			if tc.explicitCodex && tc.codex != "" {
				f.absent(filepath.Join(f.root, tc.codex))
			}
		})
	}
}

func TestDefaultHomeSymlinkUsesCanonicalDestinations(t *testing.T) {
	f := newContractFixture(t)
	home := filepath.Join(f.root, "home")
	alias := filepath.Join(f.root, "home-alias")
	if err := os.MkdirAll(home, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(home, alias); err != nil {
		t.Fatal(err)
	}
	var env []string
	for _, item := range os.Environ() {
		key, _, _ := strings.Cut(item, "=")
		if key != "HOME" && key != "XDG_STATE_HOME" && key != "CODEX_HOME" {
			env = append(env, item)
		}
	}
	env = append(env, "HOME="+alias)
	args := []string{"render", "--sync", "--sections-dir", filepath.Join(f.config, "sections"), "--mode", `codex:{"plugin":"skills"}`}
	cmd := exec.Command(contractBinary, args...)
	cmd.Env = env
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("render with symlinked home: %v\n%s", err, out)
	}
	f.symlink(filepath.Join(home, ".codex", "AGENTS.md"))
	if _, err := os.Stat(filepath.Join(home, ".agents/plugins/skillset/.codex-plugin/plugin.json")); err != nil {
		t.Fatal(err)
	}
	cmd = exec.Command(contractBinary, "sync", "1")
	cmd.Env = env
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("repeat default sync with symlinked home: %v\n%s", err, out)
	}
}
