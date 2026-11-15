package internal

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func modesEnv(home string) []string {
	env := []string{}
	for _, v := range os.Environ() {
		k, _, _ := strings.Cut(v, "=")
		if k != "HOME" && k != "XDG_STATE_HOME" && k != "CODEX_HOME" {
			env = append(env, v)
		}
	}
	if home != "" {
		env = append(env, "HOME="+home)
	}
	return env
}
func TestModesDiscoverabilityAndHelp(t *testing.T) {
	for _, args := range [][]string{{"--help"}, {"modes", "--help"}} {
		out, err := exec.Command(contractBinary, args...).CombinedOutput()
		if err != nil {
			t.Fatalf("%v: %v\n%s", args, err, out)
		}
		if !strings.Contains(string(out), "modes") {
			t.Fatalf("missing modes in %v: %s", args, out)
		}
	}
}
func TestModesRequiresNoHomeOrSourcesAndWritesNothing(t *testing.T) {
	for _, homeSet := range []bool{false, true} {
		t.Run(map[bool]string{true: "isolated-home", false: "unset-home"}[homeSet], func(t *testing.T) {
			root := t.TempDir()
			home := ""
			if homeSet {
				home = filepath.Join(root, "home")
			}
			cmd := exec.Command(contractBinary, "modes")
			cmd.Dir = root
			cmd.Env = modesEnv(home)
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("modes: %v\n%s", err, out)
			}
			first, _, _ := strings.Cut(string(out), "\n")
			if first != `config: {"command":"modes"}` {
				t.Fatalf("first line: %q", first)
			}
			for _, want := range []string{"Mode syntax:", "--mode=<harness>[:<JSON object>]", "Available harnesses:", "Codex JSON options:", "OpenCode JSON options:", "Claude JSON options:", "plugin_name", "default"} {
				if !strings.Contains(strings.ToLower(string(out)), strings.ToLower(want)) {
					t.Errorf("modes output missing %q: %s", want, out)
				}
			}
			examples := regexp.MustCompile(`(?:codex|opencode|claude):\{[^\n}]*\}`).FindAllString(string(out), -1)
			if len(examples) < 2 {
				t.Fatalf("want concrete JSON mode examples: %s", out)
			}
			for _, example := range examples {
				_, value, _ := strings.Cut(example, ":")
				if !json.Valid([]byte(value)) {
					t.Fatalf("invalid JSON example %q", example)
				}
				probe := exec.Command(contractBinary, "render", "--sync", "--mode", example, "--sections-dir", filepath.Join(root, "missing"), "--render-dir", filepath.Join(root, "render"), "--codex-dir", filepath.Join(root, "codex"), "--codex-marketplace-home", filepath.Join(root, "marketplace"), "--opencode-dir", filepath.Join(root, "opencode"), "--claude-dir", filepath.Join(root, "claude"), "--claude-marketplace-dir", filepath.Join(root, "claude-market"))
				probe.Dir = root
				probe.Env = modesEnv(home)
				result, err := probe.CombinedOutput()
				exit, ok := err.(*exec.ExitError)
				if !ok || exit.ExitCode() != 2 || !strings.HasPrefix(string(result), "config: ") {
					t.Fatalf("example not accepted by parser %q: %v %s", example, err, result)
				}
			}
			files, err := os.ReadDir(root)
			if err != nil || len(files) != 0 {
				t.Fatalf("modes wrote files: %v, %v", files, err)
			}
		})
	}
}
func TestModesRejectsUnexpectedArgumentsWithoutWrites(t *testing.T) {
	root := t.TempDir()
	cmd := exec.Command(contractBinary, "modes", "--unknown")
	cmd.Dir = root
	cmd.Env = modesEnv("")
	out, err := cmd.CombinedOutput()
	exit, ok := err.(*exec.ExitError)
	if !ok || exit.ExitCode() != 2 {
		t.Fatalf("expected exit 2, got %v: %s", err, out)
	}
	files, e := os.ReadDir(root)
	if e != nil || len(files) != 0 {
		t.Fatalf("invalid modes command wrote files: %v, %v", files, e)
	}
}
