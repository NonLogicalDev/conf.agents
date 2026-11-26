package integration_tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestProfileListAndShow(t *testing.T) {
	f := newContractFixture(t)
	config := filepath.Join(f.root, "config-home")
	t.Setenv("XDG_CONFIG_HOME", config)
	run := func(want int, args ...string) string {
		t.Helper()
		out, err := exec.Command(contractBinary, append([]string{"profile"}, args...)...).CombinedOutput()
		status := 0
		if err != nil {
			exit, ok := err.(*exec.ExitError)
			if !ok {
				t.Fatal(err)
			}
			status = exit.ExitCode()
		}
		if status != want {
			t.Fatalf("%v: exit %d, want %d: %s", args, status, want, out)
		}
		return string(out)
	}
	if got := run(0, "list"); got != "" {
		t.Fatalf("list with no profiles: %q", got)
	}
	dir := filepath.Join(config, "agentbld")
	f.absent(dir)
	body := "{\"modes\":[\"codex\"]}\n"
	f.put(dir, "work.json", "{}\n")
	f.put(dir, "default.json", body)
	f.put(dir, "notes.txt", "notes\n")
	f.put(dir, ".in-progress.json", "{}\n")
	f.put(dir, "directory.json/contents", "")
	if got := run(0, "list"); got != "default\nwork\n" {
		t.Fatalf("listed profiles: %q", got)
	}
	if got := run(0, "show", "default"); got != body {
		t.Fatalf("show default: %q", got)
	}
	run(2, "show")
	run(2, "show", "missing")
	run(2, "show", "../default")
}

func TestSyncProfileReplay(t *testing.T) {
	f := newContractFixture(t)
	config := filepath.Join(f.root, "xdg-config")
	t.Setenv("XDG_CONFIG_HOME", config)
	home := filepath.Join(f.root, "isolated-home")
	f.render = filepath.Join(home, "render")
	openCode := filepath.Join(home, "opencode")
	f.codex = filepath.Join(f.root, "isolated-home-backup", "codex")
	render := func() string {
		lifecycleRun(t, f, 0, lifecycleRenderArgs(f, "--mode=codex", "--mode=opencode")...)
		return lifecycleOrdinal(t, f.latest())
	}
	first := render()

	lifecycleRun(t, f, 0, "sync", first, "--render-dir", f.render,
		"--mode=opencode", "--opencode-dir", openCode, "--save-profile", "work")
	work := f.json(filepath.Join(config, "agentbld", "work.json"))
	if work["render_dir"] != "~/render" || work["opencode_dir"] != "~/opencode" || work["generation"] != nil {
		t.Fatalf("saved profile: %#v", work)
	}
	f.absent(filepath.Join(f.codex, "AGENTS.md"))

	lifecycleRun(t, f, 0, "sync", first, "--render-dir", f.render,
		"--mode=codex", "--codex-dir", f.codex, "--save-profile", "default")
	defaultProfile := f.json(filepath.Join(config, "agentbld", "default.json"))
	if defaultProfile["codex_dir"] != f.codex || defaultProfile["codex_marketplace_home"] != "~/.agents/plugins" {
		t.Fatalf("saved default profile: %#v", defaultProfile)
	}
	f.put(f.layer, "AGENTS.md.d/10-base.md", "Second generation\n")
	second := render()

	plainCodex := filepath.Join(f.root, "plain-codex")
	plainOpenCode := filepath.Join(f.root, "plain-opencode")
	lifecycleRun(t, f, 0, "sync", second, "--render-dir", f.render,
		"--codex-dir", plainCodex, "--opencode-dir", plainOpenCode)
	f.symlink(filepath.Join(plainCodex, "AGENTS.md"))
	f.symlink(filepath.Join(plainOpenCode, "AGENTS.md"))
	if got := f.read(filepath.Join(f.codex, "AGENTS.md")); strings.Contains(got, "Second generation") {
		t.Fatalf("sync selected default.json without --profile: %s", got)
	}

	lifecycleRun(t, f, 0, "sync", second, "--profile", "default")
	if got := f.read(filepath.Join(f.codex, "AGENTS.md")); !strings.Contains(got, "Second generation") {
		t.Fatalf("default profile did not sync the requested generation: %s", got)
	}
	if got := f.read(filepath.Join(openCode, "AGENTS.md")); strings.Contains(got, "Second generation") {
		t.Fatalf("default profile synced an unselected harness: %s", got)
	}
	lifecycleRun(t, f, 0, "sync", second, "--profile", "work")
	if got := f.read(filepath.Join(openCode, "AGENTS.md")); !strings.Contains(got, "Second generation") {
		t.Fatalf("named profile did not sync the requested generation: %s", got)
	}

	override := filepath.Join(f.root, "overridden-codex")
	lifecycleRun(t, f, 0, "sync", second, "--profile", "work", "--mode=codex", "--codex-dir", override)
	f.symlink(filepath.Join(override, "AGENTS.md"))
	if got := f.read(filepath.Join(override, "AGENTS.md")); !strings.Contains(got, "Second generation") {
		t.Fatalf("explicit mode and destination did not override the profile: %s", got)
	}

	overrideOpenCode := filepath.Join(f.root, "overridden-opencode")
	lifecycleRun(t, f, 0, "sync", second, "--profile", "work", "--opencode-dir", overrideOpenCode)
	f.symlink(filepath.Join(overrideOpenCode, "AGENTS.md"))
}

func TestSyncProfileFlagsRequireNames(t *testing.T) {
	f := newContractFixture(t)
	for _, flag := range []string{"--profile", "--save-profile"} {
		out := lifecycleRun(t, f, 2, "sync", "1", flag)
		if !strings.Contains(out, flag) || !strings.Contains(out, "expected") {
			t.Fatalf("expected a missing argument error for %s: %s", flag, out)
		}
	}
}

func TestSyncProfilePreservedOnActivationFailure(t *testing.T) {
	f := newContractFixture(t)
	t.Setenv("XDG_CONFIG_HOME", "")
	lifecycleRun(t, f, 0, lifecycleRenderArgs(f)...)
	number := lifecycleOrdinal(t, f.latest())
	args := []string{"sync", number, "--render-dir", f.render, "--codex-dir", f.codex}
	lifecycleRun(t, f, 0, append(args, "--save-profile", "work")...)
	path := filepath.Join(f.root, "isolated-home", ".config", "agentbld", "work.json")
	before := f.read(path)

	conflict := filepath.Join(f.root, "conflicting-codex")
	f.put(conflict, "AGENTS.md", "unmanaged instructions\n")
	lifecycleRun(t, f, 2, "sync", number, "--profile", "work", "--codex-dir", conflict, "--save-profile", "work")
	if got := f.read(path); got != before {
		t.Fatalf("failed activation changed profile: %s", got)
	}
	target := filepath.Join(f.root, "unused-codex")
	lifecycleRun(t, f, 2, "sync", number, "--profile", "work", "--codex-dir", target, "--save-profile", "../other")
	f.absent(filepath.Join(target, "AGENTS.md"))
	if _, err := os.Stat(filepath.Join(filepath.Dir(filepath.Dir(path)), "other.json")); !os.IsNotExist(err) {
		t.Fatalf("unexpected profile outside the config directory: %v", err)
	}
	lifecycleRun(t, f, 2, "sync", number, "--profile", "missing", "--codex-dir", target)
	f.absent(filepath.Join(target, "AGENTS.md"))
}
