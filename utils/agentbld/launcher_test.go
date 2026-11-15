package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

const fakeGo = `#!/bin/sh
set -eu
printf '%s\n' "$PWD" "$@" > "$TEST_BUILD_ARGS"
if [ "${TEST_STAGE:-}" = build ]; then
  trap 'printf signal > "$TEST_SIGNAL"; exit 129' HUP
  trap 'printf signal > "$TEST_SIGNAL"; exit 130' INT
  trap 'printf signal > "$TEST_SIGNAL"; exit 143' TERM
  : > "$TEST_READY"
  while :; do sleep 0.05; done
fi
if [ "${TEST_BUILD_EXIT:-0}" != 0 ]; then exit "$TEST_BUILD_EXIT"; fi
if [ "$1" != -C ] || [ "$3" != build ] || [ "$4" != -o ] || [ "$6" != . ]; then exit 99; fi
cat > "$5" <<'BINARY'
#!/bin/sh
set -eu
printf '%s\n' "$PWD" "$@" > "$TEST_RUN_ARGS"
if [ "${TEST_STAGE:-}" = runtime ]; then
  trap 'printf signal > "$TEST_SIGNAL"; exit 129' HUP
  trap 'printf signal > "$TEST_SIGNAL"; exit 130' INT
  trap 'printf signal > "$TEST_SIGNAL"; exit 143' TERM
  : > "$TEST_READY"
  while :; do sleep 0.05; done
fi
exit "${TEST_RUN_EXIT:-0}"
BINARY
chmod +x "$5"
`

type launcherFixture struct {
	t                                                                     *testing.T
	root, wrapper, module, caller, tmp, buildArgs, runArgs, ready, signal string
	env                                                                   []string
}

func newLauncherFixture(t *testing.T) *launcherFixture {
	t.Helper()
	root := t.TempDir()
	f := &launcherFixture{t: t, root: root}
	f.wrapper = filepath.Join(root, "repo/scripts/reconcile-config")
	f.module = filepath.Join(root, "repo/utils/agentbld")
	f.caller = filepath.Join(root, "caller")
	f.tmp = filepath.Join(root, "tmp")
	bin := filepath.Join(root, "bin")
	for _, p := range []string{filepath.Dir(f.wrapper), f.module, f.caller, f.tmp, bin} {
		if err := os.MkdirAll(p, 0755); err != nil {
			t.Fatal(err)
		}
	}
	source := os.Getenv("AGENTBLD_LAUNCHER")
	if source == "" {
		source = filepath.Join("..", "..", "scripts", "reconcile-config")
	}
	data, err := os.ReadFile(source)
	if err != nil {
		t.Fatalf("read compatibility wrapper %s: %v", source, err)
	}
	if err = os.WriteFile(f.wrapper, data, 0755); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(bin, "go"), []byte(fakeGo), 0755); err != nil {
		t.Fatal(err)
	}
	f.buildArgs = filepath.Join(root, "build-args")
	f.runArgs = filepath.Join(root, "run-args")
	f.ready = filepath.Join(root, "ready")
	f.signal = filepath.Join(root, "signal")
	f.env = append(os.Environ(), "PATH="+bin+string(os.PathListSeparator)+"/usr/bin:/bin", "TMPDIR="+f.tmp, "TEST_BUILD_ARGS="+f.buildArgs, "TEST_RUN_ARGS="+f.runArgs, "TEST_READY="+f.ready, "TEST_SIGNAL="+f.signal)
	return f
}
func (f *launcherFixture) command(args ...string) *exec.Cmd {
	f.t.Helper()
	cmd := exec.Command(f.wrapper, args...)
	cmd.Dir = f.caller
	cmd.Env = f.env
	return cmd
}
func (f *launcherFixture) checkEmptyTmp() {
	f.t.Helper()
	es, e := os.ReadDir(f.tmp)
	if e != nil || len(es) != 0 {
		f.t.Fatalf("temporary files remain: %v, %v", es, e)
	}
}
func (f *launcherFixture) contents(p string) string {
	f.t.Helper()
	b, e := os.ReadFile(p)
	if e != nil {
		f.t.Fatal(e)
	}
	return string(b)
}
func launcherExit(t *testing.T, cmd *exec.Cmd, want int) {
	t.Helper()
	out, err := cmd.CombinedOutput()
	code := 0
	if err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			code = e.ExitCode()
		} else {
			t.Fatal(err)
		}
	}
	if code != want {
		t.Fatalf("exit %d, want %d: %s", code, want, out)
	}
}
func TestLauncherPreservesWorkingDirectoryArgumentsAndCleansUp(t *testing.T) {
	f := newLauncherFixture(t)
	launcherExit(t, f.command("sync", "--config-dir", "relative path", "--rules-only"), 0)
	build := strings.Split(strings.TrimSuffix(f.contents(f.buildArgs), "\n"), "\n")
	if len(build) != 7 || build[0] != f.caller || build[1] != "-C" || filepath.Clean(build[2]) != f.module || build[3] != "build" || build[4] != "-o" || build[6] != "." {
		t.Fatalf("build args: %#v", build)
	}
	if !strings.HasPrefix(build[5], f.tmp+string(os.PathSeparator)) {
		t.Fatalf("binary outside temporary root: %s", build[5])
	}
	want := f.caller + "\nsync\n--config-dir\nrelative path\n--rules-only\n"
	if got := f.contents(f.runArgs); got != want {
		t.Fatalf("args: %q, want %q", got, want)
	}
	f.checkEmptyTmp()
}
func TestLauncherPropagatesExitCodesAndCleansUp(t *testing.T) {
	for _, tc := range []struct {
		name, setting string
		code          int
	}{{"build", "TEST_BUILD_EXIT=42", 42}, {"runtime", "TEST_RUN_EXIT=23", 23}} {
		t.Run(tc.name, func(t *testing.T) {
			f := newLauncherFixture(t)
			f.env = append(f.env, tc.setting)
			launcherExit(t, f.command("render"), tc.code)
			if tc.name == "build" {
				if _, err := os.Stat(f.runArgs); !os.IsNotExist(err) {
					t.Fatalf("binary ran after failed build: %v", err)
				}
			}
			f.checkEmptyTmp()
		})
	}
}
func TestLauncherForwardsSignalsAndCleansUp(t *testing.T) {
	for _, tc := range []struct {
		name   string
		signal syscall.Signal
		code   int
	}{{"hup", syscall.SIGHUP, 129}, {"int", syscall.SIGINT, 130}, {"term", syscall.SIGTERM, 143}} {
		for _, stage := range []string{"build", "runtime"} {
			t.Run(tc.name+"/"+stage, func(t *testing.T) {
				f := newLauncherFixture(t)
				f.env = append(f.env, "TEST_STAGE="+stage)
				cmd := f.command("render")
				if err := cmd.Start(); err != nil {
					t.Fatal(err)
				}
				done := make(chan error, 1)
				go func() { done <- cmd.Wait() }()
				deadline := time.Now().Add(5 * time.Second)
				for {
					if _, err := os.Stat(f.ready); err == nil {
						break
					}
					select {
					case err := <-done:
						t.Fatalf("launcher exited before child ready: %v", err)
					default:
					}
					if time.Now().After(deadline) {
						cmd.Process.Kill()
						<-done
						t.Fatal("child did not become ready")
					}
					time.Sleep(10 * time.Millisecond)
				}
				if err := cmd.Process.Signal(tc.signal); err != nil {
					t.Fatal(err)
				}
				select {
				case err := <-done:
					e, ok := err.(*exec.ExitError)
					if !ok || e.ExitCode() != tc.code {
						t.Fatalf("exit %v, want %d", err, tc.code)
					}
				case <-time.After(5 * time.Second):
					cmd.Process.Kill()
					<-done
					t.Fatal("launcher did not terminate")
				}
				if got := f.contents(f.signal); got != "signal" {
					t.Fatal(fmt.Sprintf("child did not receive signal: %q", got))
				}
				f.checkEmptyTmp()
			})
		}
	}
}
